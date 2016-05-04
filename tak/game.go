package tak

import "errors"

type Config struct {
	Size      int
	Pieces    int
	Capstones int

	l, r, t, b uint64
	mask       uint64
}

var defaultPieces = []int{0, 0, 0, 10, 15, 21, 30, 40, 50}
var defaultCaps = []int{0, 0, 0, 0, 0, 1, 1, 1, 2}

func (c *Config) precompute() {
	s := uint(c.Size)
	for i := uint(0); i < s; i++ {
		c.l |= 1 << (s*(i+1) - 1)
		c.r |= 1 << (s * i)
		c.t |= 1 << (s*(s-1) + i)
		c.b |= 1 << i
	}
	c.mask = 1<<uint(c.Size*c.Size) - 1
}

func New(g Config) *Position {
	if g.Pieces == 0 {
		g.Pieces = defaultPieces[g.Size]
	}
	if g.Capstones == 0 {
		g.Capstones = defaultCaps[g.Size]
	}
	g.precompute()
	p := &Position{
		cfg:         &g,
		whiteStones: byte(g.Pieces),
		whiteCaps:   byte(g.Capstones),
		blackStones: byte(g.Pieces),
		blackCaps:   byte(g.Capstones),
		move:        0,
		board:       make([]Square, g.Size*g.Size),
	}
	return p
}

type Square []Piece

type Position struct {
	cfg         *Config
	whiteStones byte
	whiteCaps   byte
	blackStones byte
	blackCaps   byte

	move     int
	board    []Square
	analysis Analysis
}

type Analysis struct {
	WhiteRoad   uint64
	BlackRoad   uint64
	WhiteGroups []uint64
	BlackGroups []uint64
}

// FromSquares initializes a Position with the specified squares and
// move number. `board` is a slice of rows, numbered from low to high,
// each of which is a slice of positions.
func FromSquares(cfg Config, board [][]Square, move int) (*Position, error) {
	p := New(cfg)
	p.move = move
	for x := 0; x < p.Size(); x++ {
		for y := 0; y < p.Size(); y++ {
			p.set(x, y, board[y][x])
			for _, piece := range board[y][x] {
				switch piece {
				case MakePiece(White, Capstone):
					p.whiteCaps--
				case MakePiece(Black, Capstone):
					p.blackCaps--
				case MakePiece(White, Flat), MakePiece(White, Standing):
					p.whiteStones--
				case MakePiece(Black, Flat), MakePiece(Black, Standing):
					p.blackStones--
				default:
					return nil, errors.New("bad stone")
				}
			}
		}
	}
	p.analyze()
	return p, nil
}

func (p *Position) Size() int {
	return p.cfg.Size
}

func (p *Position) At(x, y int) Square {
	return p.board[y*p.cfg.Size+x]
}

func (p *Position) set(x, y int, s Square) {
	p.board[y*p.cfg.Size+x] = s
}

func (p *Position) ToMove() Color {
	if p.move%2 == 0 {
		return White
	}
	return Black
}

func (p *Position) MoveNumber() int {
	return p.move
}

func (p *Position) WhiteStones() int {
	return int(p.whiteStones)
}

func (p *Position) BlackStones() int {
	return int(p.blackStones)
}

func (p *Position) GameOver() (over bool, winner Color) {
	if p, ok := p.hasRoad(); ok {
		return true, p
	}

	if (p.whiteStones+p.whiteCaps) != 0 && (p.blackStones+p.blackCaps) != 0 {
		return false, White
	}

	return true, p.flatsWinner()
}

func (p *Position) roadAt(x, y int) (Color, bool) {
	sq := p.At(x, y)
	if len(sq) == 0 {
		return White, false
	}
	return sq[0].Color(), sq[0].IsRoad()
}

func (p *Position) hasRoad() (Color, bool) {
	white, black := false, false

	for _, g := range p.analysis.WhiteGroups {
		if ((g&p.cfg.t) != 0 && (g&p.cfg.b) != 0) ||
			((g&p.cfg.l) != 0 && (g&p.cfg.r) != 0) {
			white = true
			break
		}
	}
	for _, g := range p.analysis.BlackGroups {
		if ((g&p.cfg.t) != 0 && (g&p.cfg.b) != 0) ||
			((g&p.cfg.l) != 0 && (g&p.cfg.r) != 0) {
			black = true
			break
		}
	}

	switch {
	case white && black:
		if p.ToMove() == White {
			return Black, true
		}
		return White, true
	case white:
		return White, true
	case black:
		return Black, true
	default:
		return White, false
	}

}

func (p *Position) analyze() {
	var bb uint64
	var bw uint64
	for _, sq := range p.board {
		bw <<= 1
		bb <<= 1
		if len(sq) > 0 && sq[0].IsRoad() {
			if sq[0].Color() == White {
				bw |= 1
			} else {
				bb |= 1
			}
		}
	}
	p.analysis.WhiteRoad = bw
	p.analysis.BlackRoad = bb

	alloc := make([]uint64, 0, 2*p.Size())
	p.analysis.WhiteGroups = p.floodone(bw, alloc)
	alloc = p.analysis.WhiteGroups
	p.analysis.BlackGroups = p.floodone(bb, alloc[len(alloc):cap(alloc)])
}

func (p *Position) floodone(bits uint64, out []uint64) []uint64 {
	var seen uint64
	for bits != 0 {
		next := bits & (bits - 1)
		bit := bits &^ next

		if seen&bit == 0 {
			g := p.flood(bits, bit)
			if g != bit && popcount(g) > 2 {
				out = append(out, g)
			}
			seen |= g
		}

		bits = next
	}
	return out
}

func (p *Position) flood(all uint64, seed uint64) uint64 {
	for {
		next := seed
		next |= (seed << 1) &^ p.cfg.r
		next |= (seed >> 1) &^ p.cfg.l
		next |= (seed >> uint(p.cfg.Size))
		next |= (seed << uint(p.cfg.Size))
		next &= all & p.cfg.mask
		if next == seed {
			return next
		}
		seed = next
	}
}

func popcount(x uint64) (n int) {
	// bit population count, see
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return int(x >> 56)
}

func (p *Position) bitroad(bits uint64) bool {
	s := uint(p.cfg.Size)
	var mask uint64 = (1 << s) - 1
	row := bits & mask
	for i := uint(1); i < s; i++ {
		if row == 0 {
			return false
		}
		next := (bits >> (i * s)) & mask
		row &= next
		for {
			last := row
			row |= ((row >> 1) & next) |
				((row << 1) & next)
			row &= mask
			if row == last {
				break
			}
		}
	}
	return row != 0

}

func (p *Position) countFlats() (w int, b int) {
	cw, cb := 0, 0
	for i := 0; i < p.cfg.Size*p.cfg.Size; i++ {
		stack := p.board[i]
		if len(stack) > 0 {
			if stack[0].Kind() == Flat {
				if stack[0].Color() == White {
					cw++
				} else {
					cb++
				}
			}
		}
	}
	return cw, cb
}

func (p *Position) flatsWinner() Color {
	cw, cb := p.countFlats()
	if cw > cb {
		return White
	}
	if cb > cw {
		return Black
	}
	return NoColor
}

type WinReason int

const (
	RoadWin WinReason = iota
	FlatsWin
	Resignation
)

type WinDetails struct {
	Reason     WinReason
	Winner     Color
	WhiteFlats int
	BlackFlats int
}

func (p *Position) WinDetails() WinDetails {
	over, c := p.GameOver()
	if !over {
		panic("WinDetails on a game not over")
	}
	var d WinDetails
	d.Winner = c
	d.WhiteFlats, d.BlackFlats = p.countFlats()
	if _, ok := p.hasRoad(); ok {
		d.Reason = RoadWin
	} else {
		d.Reason = FlatsWin
	}
	return d
}
