package ai

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"github.com/nelhage/taktician/bitboard"
	"github.com/nelhage/taktician/ptn"
	"github.com/nelhage/taktician/tak"
)

const (
	maxEval      int64 = 1 << 30
	minEval            = -maxEval
	winThreshold       = 1 << 29
)

type MinimaxAI struct {
	depth int
	size  int
	rand  *rand.Rand

	Debug bool

	st stats
	c  bitboard.Constants
}

type stats struct {
	generated uint64
	evaluated uint64
	cutoffs   uint64
}

func formatpv(ms []tak.Move) string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, m := range ms {
		if i != 0 {
			out.WriteString(" ")
		}
		out.WriteString(ptn.FormatMove(&m))
	}
	out.WriteString("]")
	return out.String()
}

func (m *MinimaxAI) GetMove(p *tak.Position) tak.Move {
	ms, _ := m.Analyze(p)
	return ms[0]
}

func (m *MinimaxAI) Analyze(p *tak.Position) ([]tak.Move, int64) {
	if m.size != p.Size() {
		panic("Analyze: wrong size")
	}
	seed := time.Now().Unix()
	m.rand = rand.New(rand.NewSource(seed))
	if m.Debug {
		log.Printf("seed=%d", seed)
	}

	var ms []tak.Move
	var v int64
	top := time.Now()
	var prevEval uint64
	for i := 1; i <= m.depth; i++ {
		m.st = stats{}
		start := time.Now()
		ms, v = m.minimax(p, i, ms, minEval-1, maxEval+1)
		if m.Debug {
			log.Printf("[minimax] depth=%d val=%d pv=%s time=%s total=%s evaluated=%d branch=%d",
				i, v, formatpv(ms),
				time.Now().Sub(start),
				time.Now().Sub(top),
				m.st.evaluated,
				m.st.evaluated/(prevEval+1),
			)
		}
		prevEval = m.st.evaluated
		if v > winThreshold || v < -winThreshold {
			break
		}
	}
	return ms, v
}

func (ai *MinimaxAI) minimax(
	p *tak.Position,
	depth int,
	pv []tak.Move,
	α, β int64) ([]tak.Move, int64) {
	over, _ := p.GameOver()
	if depth == 0 || over {
		ai.st.evaluated++
		return nil, ai.evaluate(p)
	}

	if p.MoveNumber() < 2 {
		for _, c := range [][]int{{0, 0}, {p.Size() - 1, 0}, {0, p.Size() - 1}, {p.Size() - 1, p.Size() - 1}} {
			x, y := c[0], c[1]
			if len(p.At(x, y)) == 0 {
				return []tak.Move{{X: x, Y: y, Type: tak.PlaceFlat}}, 0
			}
		}
	}
	moves := p.AllMoves()
	ai.st.generated += uint64(len(moves))
	if depth == ai.depth {
		for i := len(moves) - 1; i > 0; i-- {
			j := ai.rand.Int31n(int32(i))
			moves[j], moves[i] = moves[i], moves[j]
		}
	}
	if len(pv) > 0 {
		for i, m := range moves {
			if m.Equal(&pv[0]) {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	best := make([]tak.Move, 0, depth)
	max := minEval - 1
	for _, m := range moves {
		child, e := p.Move(&m)
		if e != nil {
			continue
		}
		var ms []tak.Move
		var v int64
		if len(best) == 0 {
			ms, v = ai.minimax(child, depth-1, nil, -β, -α)
		} else {
			ms, v = ai.minimax(child, depth-1, best[1:], -β, -α)
		}
		v = -v
		if v > max {
			max = v
			best = append(best[:0], m)
			best = append(best, ms...)
		}
		if v > α {
			α = v
			if α >= β {
				ai.st.cutoffs++
				break
			}
		}
	}
	return best, max
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

const (
	weightFlat       = 200
	weightCaptured   = 100
	weightControlled = 500
	weightCapstone   = -50
	weightThreat     = 150
	weightCenter     = 50
)

func (m *MinimaxAI) evaluate(p *tak.Position) int64 {
	if over, winner := p.GameOver(); over {
		switch winner {
		case tak.NoColor:
			return 0
		case p.ToMove():
			return maxEval - int64(p.MoveNumber())
		default:
			return minEval + int64(p.MoveNumber())
		}
	}
	mine, theirs := 0, 0
	me := p.ToMove()
	addw := func(c tak.Color, w int) {
		if c == me {
			mine += w
		} else {
			theirs += w
		}
	}
	analysis := p.Analysis()
	for x := 0; x < p.Size(); x++ {
		for y := 0; y < p.Size(); y++ {
			sq := p.At(x, y)
			if len(sq) == 0 {
				continue
			}
			addw(sq[0].Color(), weightControlled)
			if sq[0].Kind() == tak.Capstone {
				addw(sq[0].Color(), weightCapstone)
			}

			if x == p.Size()/2 && y == p.Size()/2 {
				// TODO(board-size)
				addw(sq[0].Color(), weightCenter)
			}
			for i, stone := range sq {
				if i > 0 && i < p.Size() {
					addw(sq[0].Color(), weightCaptured)
				}
				if stone.Kind() == tak.Flat {
					addw(stone.Color(), weightFlat)
				}
			}
		}
	}
	addw(tak.White, weightThreat*m.threats(analysis.WhiteGroups, analysis.Occupied))
	addw(tak.Black, weightThreat*m.threats(analysis.BlackGroups, analysis.Occupied))

	return int64(mine - theirs)
}

func (m *MinimaxAI) threats(groups []uint64, filled uint64) int {
	count := 0
	empty := ^filled
	s := uint(m.size)
	for _, g := range groups {
		if g&m.c.L != 0 {
			if g&(m.c.R<<1) != 0 && empty&m.c.R != 0 {
				count++
			}
		}
		if g&m.c.R != 0 {
			if g&(m.c.L>>1) != 0 && empty&m.c.L != 0 {
				count++
			}
		}
		if g&m.c.B != 0 {
			if g&(m.c.T>>s) != 0 && empty&m.c.T != 0 {
				count++
			}
		}
		if g&m.c.T != 0 {
			if g&(m.c.B<<s) != 0 && empty&m.c.B != 0 {
				count++
			}
		}
	}
	return count
}

func NewMinimax(size int, depth int) *MinimaxAI {
	m := &MinimaxAI{size: size, depth: depth}
	m.c = bitboard.Precompute(uint(size))
	return m
}
