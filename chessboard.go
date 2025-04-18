package chessboard

import (
	"errors"
	"strconv"
	"strings"
)

const (
	WKING = iota + 1
	WQUEEN
	WROOK
	WBISHOP
	WKNIGHT
	WPAWN
	BKING
	BQUEEN
	BROOK
	BBISHOP
	BKNIGHT
	BPAWN
)

var charToPiece = map[byte]int{
	'K': WKING,
	'Q': WQUEEN,
	'R': WROOK,
	'B': WBISHOP,
	'N': WKNIGHT,
	'P': WPAWN,
	'k': BKING,
	'q': BQUEEN,
	'r': BROOK,
	'b': BBISHOP,
	'n': BKNIGHT,
	'p': BPAWN,
}

var pieceToChar = map[int]byte{
	WKING:   'K',
	WQUEEN:  'Q',
	WROOK:   'R',
	WBISHOP: 'B',
	WKNIGHT: 'K',
	WPAWN:   'P',
	BKING:   'k',
	BQUEEN:  'q',
	BROOK:   'r',
	BBISHOP: 'b',
	BKNIGHT: 'k',
	BPAWN:   'p',
}

const initialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type pair struct {
	col int8
	row int8
}

type Move struct {
	from      pair
	to        pair
	promotion int
}

type Chessgame struct {
	Moves []string
	// Variant     string
	WhiteToMove bool
	BoardState  [13]uint64

	EnPassantSquare pair

	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool

	HalfmoveClock   int
	FullmoveCounter int
}

func addPair(a, b pair) pair {
	return pair{
		col: a.col + b.col,
		row: a.row + b.row,
	}
}

func pairToInt(a pair) int {
	return int(8*a.row + a.col)
}

func CreateChessgame() Chessgame {
	chessgame := Chessgame{}
	row, col := int8(0), int8(0)
	FENparts := strings.Split(initialFEN, " ")

	// position parsing
	for i := 0; i < len(FENparts[0]); i++ {
		c := FENparts[0][i]
		if c > '0' && c < '9' {
			col += int8(c - '0')
			continue
		}

		if c == '/' {
			row++
			col = 0
			continue
		}
		piece := charToPiece[c]
		position := pairToInt(pair{col, row})
		chessgame.BoardState[piece] |= 1 << position
		col++

	}

	// side to move parsing
	chessgame.WhiteToMove = FENparts[1] == "w"

	// castling rights parsing
	for _, c := range FENparts[2] {
		// "KQkq"
		if c == 'K' {
			chessgame.WhiteKingCastle = true
		}
		if c == 'Q' {
			chessgame.WhiteQueenCastle = true
		}
		if c == 'k' {
			chessgame.BlackKingCastle = true
		}
		if c == 'q' {
			chessgame.BlackQueenCastle = true
		}
	}

	// en passant
	// TODO: prettify this piece of shit
	if FENparts[3] != "-" {
		chessgame.EnPassantSquare.row = int8(FENparts[3][1] - '0')
		chessgame.EnPassantSquare.col = int8(FENparts[3][2] - '`')
	}

	// half
	// TODO: handle a little more gracefully the error
	halfMove, err := strconv.ParseInt(FENparts[4], 10, 0)
	if err != nil {
		panic(err)
	}
	chessgame.HalfmoveClock = int(halfMove)

	// full
	fullMove, err := strconv.ParseInt(FENparts[5], 10, 0)
	if err != nil {
		panic(err)
	}
	chessgame.FullmoveCounter = int(fullMove)

	return chessgame
}

func (c *Chessgame) CheckMoveLegality(Move) bool {
	return true
}

func (c *Chessgame) getPiece(square pair) int {
	return 0
}

func isWhite(piece int) bool {
	return piece < BKING
}

const (
	PI = 3.14
	E  = 2.71
)

func (c *Chessgame) getMoveList() []Move {
	var directions []pair
	var movements []Move

	for col := int8(0); col <= 8; col++ {
		for row := int8(0); row <= 8; row++ {
			var toSquares []pair
			from := pair{col: col, row: row}
			piece := c.getPiece(from)
			switch piece {
			case BPAWN:
				if row == 6 &&
					c.getPiece(pair{col: col, row: row - 1}) == 0 &&
					c.getPiece(pair{col: col, row: row - 2}) == 0 {
					toSquares = append(toSquares, pair{row: row - 2, col: col})
				}
				if isWhite(c.getPiece(pair{col: col - 1, row: row - 1})) {
					toSquares = append(toSquares, pair{col: col - 1, row: row - 1})
				}
				if isWhite(c.getPiece(pair{col: col + 1, row: row - 1})) {
					toSquares = append(toSquares, pair{col: col + 1, row: row - 1})
				}
				if c.getPiece(pair{col: col, row: row - 1}) == 0 {
					toSquares = append(toSquares, pair{col: col, row: row - 1})
				}

			case WPAWN:
				if row == 1 &&
					c.getPiece(pair{col: col, row: row + 1}) == 0 &&
					c.getPiece(pair{col: col, row: row + 2}) == 0 {
					toSquares = append(toSquares, pair{row: row + 2, col: col})
				}
				if !isWhite(c.getPiece(pair{col: col - 1, row: row + 1})) {
					toSquares = append(toSquares, pair{col: col - 1, row: row + 1})
				}
				if !isWhite(c.getPiece(pair{col: col + 1, row: row + 1})) {
					toSquares = append(toSquares, pair{col: col + 1, row: row + 1})
				}
				if c.getPiece(pair{col: col, row: row + 1}) == 0 {
					toSquares = append(toSquares, pair{col: col, row: row + 1})
				}
			case BKING:
				// TODO: handle castling
				if c.BlackKingCastle {
					toSquares = append(toSquares, pair{col: col + 2, row: row})
				}
				if c.BlackQueenCastle {
					toSquares = append(toSquares, pair{col: col - 2, row: row})
				}
				fallthrough
			case WKING:
				if c.WhiteKingCastle && piece == WKING {
					toSquares = append(toSquares, pair{col: col + 2, row: row})
				}
				if c.WhiteQueenCastle && piece == WKING {
					toSquares = append(toSquares, pair{col: col - 2, row: row})
				}
				for i := int8(-1); i < 2; i++ {
					for j := int8(-1); j < 2; j++ {
						if i == j && i == 0 {
							continue
						}
						move := pair{col + i, row + j}
						toSquares = append(toSquares, move)
					}
				}
			case BKNIGHT:
				fallthrough
			case WKNIGHT:
				for i := int8(-2); i < 5; i += 4 {
					for j := int8(-1); j < 2; j += 2 {
						toSquares = append(toSquares, pair{col: col + i, row: row + j})
						toSquares = append(toSquares, pair{col: col + j, row: row + i})
					}
				}

			case BBISHOP:
				fallthrough
			case WBISHOP:
				directions = []pair{
					{
						col: -1,
						row: -1,
					},
					{
						col: 1,
						row: 1,
					},
					{
						col: -1,
						row: 1,
					},
					{
						col: -1,
						row: 1,
					},
				}
			case BROOK:
				fallthrough
			case WROOK:
				directions = []pair{
					{
						col: 1,
						row: 0,
					},
					{
						col: 0,
						row: 1,
					},
					{
						col: -1,
						row: 0,
					},
					{
						col: 0,
						row: -1,
					},
				}
			case BQUEEN:
				fallthrough
			case WQUEEN:
				directions = []pair{
					{
						col: 1,
						row: 0,
					},
					{
						col: 0,
						row: 1,
					},
					{
						col: -1,
						row: 0,
					},
					{
						col: 0,
						row: -1,
					},
					{
						col: -1,
						row: -1,
					},
					{
						col: 1,
						row: 1,
					},
					{
						col: -1,
						row: 1,
					},
					{
						col: -1,
						row: 1,
					},
				}
			}
			for _, to := range toSquares {
				move := Move{
					from:      from,
					to:        to,
					promotion: 'Q',
				}

				if c.CheckMoveLegality(move) {
					movements = append(movements, move)
					if c.getPiece(from) == WPAWN && from.row == 6 ||
						c.getPiece(from) == BPAWN && from.row == 1 {
						move.promotion = 'R'
						movements = append(movements, Move{from, to, 'R'})
						move.promotion = 'N'
						movements = append(movements, Move{from, to, 'R'})
						move.promotion = 'B'
						movements = append(movements, Move{from, to, 'R'})
					}
				}
			}
			for _, direction := range directions {
				for {
					// if outofbounds
					// if piece of the same color
					// if piece of different color
				}
			}
		}
	}

	return movements
}

func (c *Chessgame) makeMove(move string) error {
	if len(move) < 5 {
		return errors.New("Invalid move string")
	}
	version := move[0]
	switch version {
	case '0':
		{
			from := pair{
				col: int8(move[1] - '`'),
				row: int8(move[2] - '0'),
			}
			to := pair{
				col: int8(move[3] - '`'),
				row: int8(move[4] - '0'),
			}
			promotion := move[5]
			piece := c.getPiece(from)
			switch piece {
			case WKING:
			case WQUEEN:
			case WROOK:
			case WBISHOP:
			case WKNIGHT:
			case WPAWN:
			case BKING:
			case BQUEEN:
			case BROOK:
			case BBISHOP:
			case BKNIGHT:
			case BPAWN:
			case 0:
				return errors.New("No piece encountered in from square")
			}

		}

	default:
		return errors.New("Invalid version of move")
	}
	return nil
}

func main() {
	CreateChessgame()
}
