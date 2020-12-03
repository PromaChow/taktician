package selfplay

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"runtime/pprof"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/google/subcommands"
	"github.com/nelhage/taktician/ptn"
	"github.com/nelhage/taktician/tak"
)

type Command struct {
	size int
	zero bool
	p1   string
	p2   string
	seed int64

	games  int
	cutoff int
	swap   bool

	prefix   string
	openings string

	debug       int
	limit       time.Duration
	timeControl time.Duration

	threads int

	out     string
	summary string
	verbose bool

	merge bool

	memProfile string
}

func (*Command) Name() string     { return "selfplay" }
func (*Command) Synopsis() string { return "Play two AIs against each other and report results" }
func (*Command) Usage() string {
	return `selfplay [flags]
`
}

func (c *Command) SetFlags(flags *flag.FlagSet) {
	flags.IntVar(&c.size, "size", 5, "board size")
	flags.StringVar(&c.p1, "p1", "taktician tei", "player1 TIE driver")
	flags.StringVar(&c.p2, "p2", "taktician tei", "player2 TIE driver")

	flags.Int64Var(&c.seed, "seed", 0, "starting random seed")
	flags.IntVar(&c.games, "games", 10, "number of games to play per opening/color")
	flags.IntVar(&c.cutoff, "cutoff", 80, "cut games off after how many plies")
	flags.BoolVar(&c.swap, "swap", true, "swap colors each game")
	flags.StringVar(&c.prefix, "prefix", "", "ptn file to start games at the end of")
	flags.StringVar(&c.openings, "openings", "", "File of openings, 1/line in TPS")
	flags.IntVar(&c.debug, "debug", 0, "debug level")
	flags.DurationVar(&c.limit, "limit", 0, "amount of time to search each move")
	flags.DurationVar(&c.timeControl, "tc", 0, "Time control for each side")
	flags.IntVar(&c.threads, "threads", 4, "number of parallel threads")
	flags.StringVar(&c.out, "out", "", "directory to write ptns to")
	flags.StringVar(&c.summary, "summary", "", "write summary JSON file")
	flags.BoolVar(&c.verbose, "v", false, "verbose output")
	flags.StringVar(&c.memProfile, "mem-profile", "", "write memory profile")

	flags.BoolVar(&c.merge, "merge", false, "merge+analyze multiple summary files")
}

func readOpenings(path string) ([]*tak.Position, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []*tak.Position
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Text()
		pos, err := ptn.ParseTPS(line)
		if err != nil {
			return nil, fmt.Errorf("parse TPS: %q: %w", line, err)
		}
		out = append(out, pos)
	}
	return out, nil
}

func (c *Command) Execute(ctx context.Context, flag *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.memProfile != "" {
		defer func() {
			f, e := os.OpenFile(c.memProfile,
				os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			if e != nil {
				log.Printf("open memory profile: %v", e)
				return
			}
			pprof.Lookup("heap").WriteTo(f, 0)
		}()
	}

	if c.merge {
		var st Stats
		for _, arg := range flag.Args() {
			if err := mergeStats(&st, arg); err != nil {
				log.Fatalf("merge: %q: %s", arg, err.Error())
			}
		}
		printSummary(&st)
		return subcommands.ExitSuccess
	}

	if c.seed == 0 {
		c.seed = time.Now().Unix()
	}

	var openings []*tak.Position
	if c.prefix != "" {
		pt, e := ptn.ParseFile(c.prefix)
		if e != nil {
			log.Fatalf("Parse PTN: %v", e)
		}
		p, e := pt.PositionAtMove(0, tak.NoColor)
		if e != nil {
			log.Fatalf("PTN: %v", e)
		}
		openings = []*tak.Position{p}
	}
	if c.openings != "" {
		var e error
		openings, e = readOpenings(c.openings)
		if e != nil {
			log.Fatalf("-openings: %v", e)
		}
	}
	if len(openings) == 0 {
		openings = []*tak.Position{tak.New(tak.Config{Size: c.size})}
	}

	cfg := &Config{
		Zero:        c.zero,
		Size:        c.size,
		Debug:       c.debug,
		Swap:        c.swap,
		Games:       c.games,
		Threads:     c.threads,
		Seed:        c.seed,
		Cutoff:      c.cutoff,
		Limit:       c.limit,
		TimeControl: c.timeControl,
		Initial:     openings,
		Verbose:     c.verbose,
		P1:          strings.Split(c.p1, " "),
		P2:          strings.Split(c.p2, " "),
	}

	st := Simulate(cfg)

	if c.out != "" {
		if c.summary == "" {
			c.summary = path.Join(c.out, "summary.json")
		}
		for _, r := range st.Games {
			writeGame(c.out, &r)
		}
	}
	if c.summary != "" {
		if err := c.writeSummary(c.summary, &st); err != nil {
			log.Println("writing summary: ", err.Error())
		}
	}

	log.Printf("done games=%d seed=%d ties=%d cutoff=%d white=%d black=%d limit=%s",
		len(st.Games), c.seed, st.Ties, st.Cutoff, st.White, st.Black, c.limit)

	printSummary(&st)

	return subcommands.ExitSuccess
}

func printSummary(st *Stats) {
	log.Printf("p1.wins=%d (%d road/%d flat) p2.wins=%d (%d road/%d flat)",
		st.Players[0].Wins, st.Players[0].RoadWins, st.Players[0].FlatWins,
		st.Players[1].Wins, st.Players[1].RoadWins, st.Players[1].FlatWins)
	tw := tabwriter.NewWriter(os.Stderr, 2, 4, 2, ' ', 0)

	fmt.Fprintf(tw, "\twhite\tblack\tsum\n")
	fmt.Fprintf(tw, "p1\t%d\t%d\t%d\n", st.Players[0].WhiteWins, st.Players[0].BlackWins, st.Players[0].Wins)
	fmt.Fprintf(tw, "p2\t%d\t%d\t%d\n", st.Players[1].WhiteWins, st.Players[1].BlackWins, st.Players[1].Wins)
	fmt.Fprintf(tw, "sum\t%d\t%d\t%d\n",
		st.Players[0].WhiteWins+st.Players[1].WhiteWins,
		st.Players[0].BlackWins+st.Players[1].BlackWins,
		st.Players[0].Wins+st.Players[1].Wins,
	)
	tw.Flush()

	score := (float64(st.Players[0].Wins) + float64(st.Ties+st.Cutoff)/2) / float64(st.Count())
	if score > 0 && score < 1 {
		elo := -400 * math.Log10(1/score-1)
		log.Printf("ΔELO=%.0f\n", elo)
	}

	a, b := int64(st.Players[0].Wins), int64(st.Players[1].Wins)
	if a < b {
		a, b = b, a
	}
	log.Printf("p[one-sided]=%f", binomTest(a, b, 0.5))
}

func writeGame(d string, r *Result) {
	os.MkdirAll(d, 0755)
	p := &ptn.PTN{}
	p.Tags = []ptn.Tag{
		{Name: "Size", Value: fmt.Sprintf("%d", r.Position.Size())},
		{Name: "Player1", Value: r.spec.p1color.String()},
	}
	if r.Initial.MoveNumber() != 0 {
		p.Tags = append(p.Tags, ptn.Tag{
			Name: "TPS", Value: ptn.FormatTPS(r.Initial)})
	}
	for i, m := range r.Moves {
		if i%2 == 0 {
			p.Ops = append(p.Ops, &ptn.MoveNumber{Number: i/2 + 1})
		}
		p.Ops = append(p.Ops, &ptn.Move{Move: m})
	}
	ptnPath := path.Join(d, fmt.Sprintf("%d-%d.ptn", r.spec.oi, r.spec.i))
	ioutil.WriteFile(ptnPath, []byte(p.Render()), 0644)
}

type Summary struct {
	Cmdline []string
	Player1 string
	Player2 string
	Limit   time.Duration
	Stats   *Stats
}

func mergeStats(st *Stats, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var summary Summary
	if err := json.Unmarshal(data, &summary); err != nil {
		return err
	}

	*st = st.Merge(summary.Stats)
	return nil
}

func (c *Command) writeSummary(path string, stats *Stats) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	summary := Summary{
		Cmdline: os.Args,
		Player1: c.p1,
		Player2: c.p2,
		Limit:   c.limit,
		Stats:   stats,
	}

	bs, err := json.MarshalIndent(&summary, "", "  ")
	if err != nil {
		return err
	}
	f.Write(bs)
	fmt.Fprintln(f)
	return nil
}
