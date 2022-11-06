package gencorpus

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/subcommands"
	"github.com/nelhage/taktician/ai"
	"github.com/nelhage/taktician/prove"
	"github.com/nelhage/taktician/ptn"
	"github.com/nelhage/taktician/tak"
	"golang.org/x/sync/errgroup"
)

type Command struct {
	seed int64
	size int

	games int

	epsilon float64
	depth   int
	threads int

	limit    time.Duration
	analysis string

	stats  bool
	output string
}

func (*Command) Name() string     { return "gencorpus" }
func (*Command) Synopsis() string { return "Generate a corpus of 3x3 positions" }
func (*Command) Usage() string {
	return `gencorpus [flags]
`
}

func (c *Command) SetFlags(flags *flag.FlagSet) {
	flags.IntVar(&c.size, "size", 3, "what size to analyze")
	flags.IntVar(&c.games, "games", 100, "games to generate")
	flags.Int64Var(&c.seed, "seed", 0, "Random seed")
	flags.IntVar(&c.threads, "threads", runtime.NumCPU(), "Number of threads")

	flags.StringVar(&c.analysis, "analysis", "minimax", "Analysis engine to run: minimax,dfpn,winning,none")
	flags.DurationVar(&c.limit, "limit", 5*time.Second, "Minimax time limit when scoring")

	flags.BoolVar(&c.stats, "stats", false, "compute and print stats")
	flags.IntVar(&c.depth, "depth", 2, "minimax depth")
	flags.Float64Var(&c.epsilon, "epsilon", 0.95, "epsilon for epsilon-greedy generation")

	flags.StringVar(&c.output, "output", "positions.txt", "output file")

}

type game struct {
	positions []*tak.Position
	moves     []tak.Move
}

func growslice[T any](sl []T, newlen int) []T {
	if len(sl) >= newlen {
		return sl
	}
	newsl := make([]T, newlen)
	copy(newsl, sl)
	return newsl
}

type entry struct {
	pos        *tak.Position
	move       tak.Move
	value      float64
	otherMoves []tak.Move
}

func fmtMoves(moves []tak.Move) string {
	var sb strings.Builder
	for i, m := range moves {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(ptn.FormatMove(m))
	}
	return sb.String()
}

func (c *Command) Execute(ctx context.Context, flag *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var byLength []int
	var posCount []map[uint64]int

	grp, ctx := errgroup.WithContext(ctx)

	games := make(chan *game)
	positions := make(chan *tak.Position)

	grp.Go(func() error {
		c.generateGames(ctx, games)
		return nil
	})
	grp.Go(func() error {
		defer close(positions)
		rng := rand.New(rand.NewSource(c.seed))

		seen := make(map[uint64]struct{})

		for g := range games {
			if c.stats {
				moves := len(g.positions)
				byLength = growslice(byLength, moves)
				byLength[moves-1] += 1
				posCount = growslice(posCount, moves+1)
				for i, p := range g.positions {
					if posCount[i] == nil {
						posCount[i] = make(map[uint64]int)
					}
					posCount[i][p.Hash()] += 1
				}
			}
			// select position
			var idx int
			r := rng.Float64()
			if r < 0.01 {
				idx = int(rng.Int31n(4))
			} else if r < 0.25 {
				idx = 4 + int(rng.Int31n(5))
			} else if r < 0.95 {
				npos := len(g.positions)
				if npos <= 9 {
					continue
				}
				idx = 9 + int(rng.Int31n(int32(npos)-9))
			}
			if idx >= len(g.positions)-1 {
				continue
			}
			pos := g.positions[idx]
			if _, ok := seen[pos.Hash()]; ok {
				continue
			}
			seen[pos.Hash()] = struct{}{}
			select {
			case positions <- pos:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if c.stats {
			for i := range byLength {
				log.Printf("ply=%3d games=%3d uniq=%4d", i, byLength[i], len(posCount[i]))
			}
		}
		return nil
	})

	results := make(chan entry)
	grp.Go(func() error {
		c.evaluate(ctx, positions, results)
		return nil
	})
	grp.Go(func() error {
		fh, err := os.Create(c.output)
		if err != nil {
			return fmt.Errorf("open %q: %w", c.output, err)
		}
		defer fh.Close()
		wr := csv.NewWriter(fh)
		defer wr.Flush()

		for e := range results {
			wr.Write([]string{
				ptn.FormatTPS(e.pos),
				ptn.FormatMove(e.move),
				fmt.Sprintf("%+f", e.value),
				fmtMoves(e.otherMoves),
			})
		}
		return nil
	})

	if err := grp.Wait(); err != nil {
		log.Println(err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (c *Command) analyzeWinning(mm *ai.MinimaxAI, e *entry) {
	ctx, cancel := context.WithTimeout(context.Background(), c.limit)
	defer cancel()
	pv, val, _ := mm.Analyze(
		ctx,
		e.pos,
	)
	if val <= ai.WinThreshold {
		return
	}

	e.move = pv[0]

	var buf [100]tak.Move
	moves := e.pos.AllMoves(buf[:0])
	var ch *tak.Position
	var err error
	for _, m := range moves {
		ch, err = e.pos.MovePreallocated(m, ch)
		if err != nil {
			continue
		}
		if ok, winner := ch.GameOver(); ok {
			if winner == e.pos.ToMove() {
				e.otherMoves = append(e.otherMoves, m)
			}
			continue
		}
		_, val, _ = mm.Analyze(ctx, ch)
		if val < -ai.WinThreshold {
			e.otherMoves = append(e.otherMoves, m)
		}

		if ctx.Err() != nil {
			return
		}
	}

	e.value = 1.0
}

func (c *Command) evaluate(ctx context.Context, positions <-chan *tak.Position, results chan<- entry) {
	defer close(results)
	grp, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.threads; i++ {
		grp.Go(func() error {
			var analyze func(p *entry)
			if c.analysis == "dfpn" {
				prover := prove.NewDFPN(&prove.DFPNConfig{
					// Attacker: tak.White,
					TableMem: 100 * 1 << 20,
				})
				analyze = func(e *entry) {
					res, _ := prover.Prove(e.pos)
					e.move = res.Move

					if res.Result == prove.EvalUnknown {
						log.Printf("unprovable! %q bounds=%d,%d",
							ptn.FormatTPS(e.pos),
							res.Proof,
							res.Disproof,
						)
					}
					if res.Result == prove.EvalTrue {
						e.value = 1.0
					} else {
						e.value = -1.0
					}
				}
			} else if c.analysis == "minimax" {
				mm := ai.NewMinimax(ai.MinimaxConfig{
					Size:     c.size,
					TableMem: 100 * 1 << 20,
				})

				analyze = func(e *entry) {
					ctx, cancel := context.WithTimeout(context.Background(), c.limit)
					defer cancel()
					pv, val, _ := mm.Analyze(
						ctx,
						e.pos,
					)
					e.move = pv[0]
					if val > ai.WinThreshold {
						e.value = 1.0
					} else if val < -ai.WinThreshold {
						e.value = -1.0
					} else if val > 0 {
						e.value = 0.5
					} else if val < 0 {
						e.value = 0.5
					}
				}
			} else if c.analysis == "winning" {
				weights := ai.DefaultWeights[c.size]
				weights[ai.Terminal_Flats] = 0
				weights[ai.Terminal_Plies] = 0
				weights[ai.Terminal_OpponentReserves] = 0
				weights[ai.Terminal_Reserves] = 0
				mm := ai.NewMinimax(ai.MinimaxConfig{
					Size:     c.size,
					TableMem: 100 * 1 << 20,
					Evaluate: ai.MakeEvaluator(c.size, &weights),
				})

				analyze = func(e *entry) {
					c.analyzeWinning(mm, e)
				}
			} else if c.analysis == "none" {
				analyze = func(e *entry) {}
			} else {
				log.Fatalf("unknown analysis: %q", c.analysis)
			}

			for p := range positions {
				ent := entry{pos: p}
				analyze(&ent)
				if c.analysis == "winning" && ent.value != 1.0 {
					continue
				}

				select {
				case results <- ent:
				case <-ctx.Done():
					return ctx.Err()
				}

			}
			return nil
		})
	}
	grp.Wait()
}

func (c *Command) generateGames(ctx context.Context, games chan<- *game) {
	defer close(games)
	todo := int64(c.games)

	grp, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.threads; i++ {
		grp.Go(func() error {
			c.generateWorker(ctx, games, &todo, i)
			return nil
		})
	}
	grp.Wait()
}

const prime = 1099511628211

func (c *Command) generateWorker(ctx context.Context, games chan<- *game, todo *int64, id int) {
	rng := rand.New(rand.NewSource(prime*c.seed + int64(id)))
	mm := ai.NewMinimax(ai.MinimaxConfig{
		Size:     c.size,
		Seed:     rng.Int63(),
		Depth:    c.depth,
		TableMem: -1,
	})
	rnd := ai.NewRandom(rng.Int63())
	for {
		gid := atomic.AddInt64(todo, -1)
		if gid < 0 {
			return
		}
		pos := tak.New(tak.Config{Size: c.size})
		g := game{positions: []*tak.Position{pos}}
		for {
			if done, _ := pos.GameOver(); done {
				break
			}
			var player ai.TakPlayer
			if rng.Float64() < c.epsilon {
				player = rnd
			} else {
				player = mm
			}

			for {
				m := player.GetMove(ctx, pos)
				child, err := pos.Move(m)
				if err != nil {
					continue
				}
				g.positions = append(g.positions, child)
				g.moves = append(g.moves, m)
				pos = child
				break
			}
		}
		select {
		case games <- &g:
		case <-ctx.Done():
			return
		}
	}
}
