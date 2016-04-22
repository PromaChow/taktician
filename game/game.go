package tak

type Game struct {
	size           int
	startingPieces int
}

type Color byte
type Kind byte
type Piece byte

const (
	White Color = 1 << 7
	Black Color = 0 << 7

	colorMask byte = 1 << 7

	Flat     Kind = 0
	Standing Kind = 1
	Capstone Kind = 2

	typeMask byte = 1<<2 - 1
)

func makePiece(color Color, kind Kind) Piece {
	return Piece(byte(color) | byte(kind))
}

func pieceColor(p Piece) Color {
	return Color(byte(p) & colorMask)
}

func pieceKind(p Piece) Kind {
	return Kind(byte(p) & typeMask)
}

func isRoad(p Piece) bool {
	return pieceKind(p) == Flat || pieceKind(p) == Capstone
}

type Square []Piece

type Position struct {
	game       *Game
	whiteFlats int
	blackFlats int
	move       int
	board      []Square
}

func (p *Position) At(x, y int) Square {
	return p.board[y*p.game.size+x]
}

func (p *Position) ToMove() Color {
	if p.move%2 == 0 {
		return White
	}
	return Black
}

func (p *Position) GameOver() (over bool, winner Color) {
	if p, ok := p.hasRoad(); ok {
		return true, p
	}

	if p.whiteFlats != 0 && p.blackFlats != 0 {
		return false, White
	}

	return true, p.flatsWinner()
}

func (p *Position) roadAt(x, y int) (Color, bool) {
	sq := p.At(x, y)
	if len(sq) == 0 {
		return White, false
	}
	return pieceColor(sq[0]), isRoad(sq[0])
}

func (p *Position) hasRoad() (Color, bool) {
	s := p.game.size
	white, black := false, false
	reachable := make([]Piece, s*s)
	for x := 0; x < s; x++ {
		if c, ok := p.roadAt(x, 0); ok {
			reachable[x] = makePiece(c, Flat)
		}
	}
	for y := 1; y < s; y++ {
		for x := 0; x < s; x++ {
			c, ok := p.roadAt(x, y)
			if !ok {
				continue
			}
			if reachable[x+(y-1)*s] == makePiece(c, Flat) {
				reachable[x+y*s] = makePiece(c, Flat)
			}
		}
		for x := 0; x < s; x++ {
			c, ok := p.roadAt(x, y)
			if !ok {
				continue
			}
			if x > 0 && reachable[x-1+y*s] == makePiece(c, Flat) {
				reachable[x+y*s] = makePiece(c, Flat)
			}
			if x < s-1 && reachable[x+1+y*s] == makePiece(c, Flat) {
				reachable[x+y*s] = makePiece(c, Flat)
			}
		}
	}
	for x := 0; x < s; x++ {
		r := reachable[x+(s-1)*s]
		if r == makePiece(White, Flat) {
			white = true
		}
		if r == makePiece(Black, Flat) {
			black = true
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

func (p *Position) flatsWinner() Color {
	cw, cb := 0, 0
	for i := 0; i < p.game.size*p.game.size; i++ {
		stack := p.board[i]
		if len(stack) > 0 {
			if pieceKind(stack[0]) == Flat {
				if pieceColor(stack[0]) == White {
					cw++
				} else {
					cb++
				}
			}
		}
	}
	if cw > cb {
		return White
	}
	return Black
}
