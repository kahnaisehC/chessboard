package chessboard

import (
	"errors"
	"fmt"
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

var bishopSlides = []pair{
	{col: 1, row: 1},
	{col: -1, row: 1},
	{col: -1, row: -1},
	{col: 1, row: -1},
}

var rookSlides = []pair{
	{col: 1, row: 0},
	{col: -1, row: 0},
	{col: 0, row: 1},
	{col: 0, row: -1},
}

var kingMoves = []pair{
	{col: 1, row: 1},
	{col: 1, row: 0},
	{col: 1, row: -1},
	{col: -1, row: 0},
	{col: -1, row: 1},
	{col: -1, row: -1},
	{col: 0, row: 1},
	{col: 0, row: -1},
}

var knightMoves = []pair{
	{col: 2, row: 1},
	{col: 2, row: -1},
	{col: -2, row: 1},
	{col: -2, row: -1},
	{col: 1, row: 2},
	{col: 1, row: -2},
	{col: -1, row: 2},
	{col: -1, row: -2},
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

const (
	BLACK = false
	WHITE = true
)

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

/*

	CreateChessboard(FEN);

	Chessboard.CheckMoveLegality(Move);
	Chessboard.GetMoveList();
	Chessboard.GetTripleRepetition();
	Chessboard.InsufficientMaterial();
	Chessboard.GetResult();
	Chessboard.GetFEN();

	Chessboard.MakeMove(Move);
	Chessboard.DeclareResult(Result);
	Chessboard.UndoMove()


	BoardState       [13]uint64
	WhiteToMove      bool
	EnPassantSquare  pair
	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool
	Moves            []string

*/

type Chessboard struct {
	// Variant     string
	BoardState       [13]uint64
	WhiteToMove      bool
	EnPassantSquare  pair
	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool
	Moves            []string

	HalfmoveClock   int
	FullmoveCounter int
}

// Main interaction interface

func (c *Chessboard) CheckMoveLegality(Move)
func (c *Chessboard) GetMoveList()
func (c *Chessboard) GetTripleRepetition()
func (c *Chessboard) InsufficientMaterial()
func (c *Chessboard) GetResult()
func (c *Chessboard) GetFEN()

func (c *Chessboard) MakeMove(Move)
func (c *Chessboard) DeclareResult()
func (c *Chessboard) UndoMove()

func (c *Chessboard) PrintBoard()
func (c *Chessboard) GameOver() bool

// WARNING: FUNCTION VERY PERIGLOSA. Use at your own risk or smth
func sq(s string) pair {
	return pair{col: int8(s[0] - 'a'), row: int8(s[1] - '0')}
}

func addPair(a, b pair) pair {
	return pair{
		col: a.col + b.col,
		row: a.row + b.row,
	}
}

func intToPair(n int) pair {
	return pair{col: int8(n / 8), row: int8(n % 8)}
}

func pairToInt(a pair) int {
	return int(8*a.row + a.col)
}

func CreateChessboard(FEN string) Chessboard {
	chessgame := Chessboard{}
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

func inBounds(p pair) bool {
	return p.col >= 0 && p.col < 8 && p.row >= 0 && p.row < 8
}

func (c *Chessboard) erasePiece(s pair) {
	// NOTE: check if range is a copy or a reference. It copies
	piecePosition := pairToInt(s)
	bitAux := ^uint64(0) ^ uint64(1<<piecePosition)

	for i := 0; i < len(c.BoardState); i++ {
		c.BoardState[i] &= bitAux
	}
}

func (c *Chessboard) putPiece(s pair, piece int) {
	c.erasePiece(s)
	piecePosition := pairToInt(s)
	bitAux := uint64(1 << piecePosition)
	c.BoardState[piece] |= bitAux
}

func (c *Chessboard) SquareIsThreatened(white bool, p pair) bool {
	if white {
		// pawn
		if c.getPiece(addPair(p, pair{col: 1, row: 1})) == BPAWN {
			return true
		}
		if c.getPiece(addPair(p, pair{col: 1, row: -1})) == BPAWN {
			return true
		}
		// knight
		for _, move := range knightMoves {
			if c.getPiece(addPair(p, move)) == BKNIGHT {
				return true
			}
		}
		// king
		for _, move := range kingMoves {
			if c.getPiece(addPair(p, move)) == BKING {
				return true
			}
		}
		// bishop
		for _, direction := range bishopSlides {
			nextSquare := addPair(p, direction)
			for inBounds(nextSquare) {
				piece := c.getPiece(nextSquare)
				if isWhite(piece) {
					break
				}
				if piece == BQUEEN || piece == BBISHOP {
					return true
				}
				if piece != 0 {
					break
				}

			}
		}
		// rook
		for _, direction := range rookSlides {
			nextSquare := addPair(p, direction)
			for inBounds(nextSquare) {
				piece := c.getPiece(nextSquare)
				if isWhite(piece) {
					break
				}
				if piece == BQUEEN || piece == BROOK {
					return true
				}
				if piece != 0 {
					break
				}

			}
		}
	} else {
		// pawn
		if c.getPiece(addPair(p, pair{col: -1, row: 1})) == WPAWN {
			return true
		}
		if c.getPiece(addPair(p, pair{col: -1, row: -1})) == WPAWN {
			return true
		}
		// knight
		for _, move := range knightMoves {
			if c.getPiece(addPair(p, move)) == WKNIGHT {
				return true
			}
		}
		// king
		for _, move := range kingMoves {
			if c.getPiece(addPair(p, move)) == WKING {
				return true
			}
		}
		// bishop
		for _, direction := range bishopSlides {
			nextSquare := addPair(p, direction)
			for inBounds(nextSquare) {
				piece := c.getPiece(nextSquare)
				if !isWhite(piece) {
					break
				}
				if piece == WQUEEN || piece == WBISHOP {
					return true
				}
				if piece != 0 {
					break
				}

			}
		}
		// rook
		for _, direction := range rookSlides {
			nextSquare := addPair(p, direction)
			for inBounds(nextSquare) {
				piece := c.getPiece(nextSquare)
				if !isWhite(piece) {
					break
				}
				if piece == WQUEEN || piece == WROOK {
					return true
				}
				if piece != 0 {
					break
				}

			}
		}
	}
	return false
}

func (c *Chessboard) GetKingPosition() pair {
	king := BKING
	if c.WhiteToMove {
		king = WKING
	}
	for i := 0; i < 64; i++ {
		if 1 == ((1 << i) | c.BoardState[king]) {
			return intToPair(i)
		}
	}
	return intToPair(64)
}

func (c *Chessboard) CheckMoveLegality(move Move) bool {
	// check inbounds
	if !(inBounds(move.from) || inBounds(move.from)) {
		return false
	}
	fromPiece := c.getPiece(move.from)
	//  check if frompiece is not empty
	if fromPiece == 0 {
		return false
	}
	toPiece := c.getPiece(move.to)

	// toPiece is not of the same color
	if toPiece != 0 && isWhite(toPiece) == isWhite(fromPiece) {
		return false
	}

	// if is the correct turn to play
	if isWhite(fromPiece) != c.WhiteToMove {
		return false
	}

	// lists of changes
	//

	var squaresToErase []pair     //
	var squaresToPutPieces []pair // squares where there are pieces
	var piecesPerSquare []int     //

	// make move
	switch {
	// castling
	case fromPiece == BKING:
		if (move.from == pair{col: 4, row: 7}) && (move.to == pair{col: 2, row: 7}) {
			if !c.BlackQueenCastle {
				return false
			}
			if c.SquareIsThreatened(WHITE, sq("e8")) || c.SquareIsThreatened(WHITE, sq("d8")) || c.SquareIsThreatened(WHITE, sq("c8")) {
				return false
			}
			if c.getPiece(sq("d8")) != 0 || c.getPiece(sq("c8")) != 0 || c.getPiece(sq("b8")) != 0 {
				return false
			}

			// WARNING: I THINK THERE IS NO OTHER EDGE CASE.
			return true
		}
		if (move.from == pair{col: 4, row: 7}) && (move.to == pair{col: 6, row: 7}) {
			if !c.BlackKingCastle {
				return false
			}
			if c.SquareIsThreatened(WHITE, sq("e8")) || c.SquareIsThreatened(WHITE, sq("f8")) || c.SquareIsThreatened(WHITE, sq("g8")) {
				return false
			}
			if c.getPiece(sq("f8")) != 0 || c.getPiece(sq("g8")) != 0 {
				return false
			}
			// WARNING: I THINK THERE IS NO OTHER EDGE CASE
			return true
		}
	case fromPiece == WKING:
		if (move.from == pair{col: 4, row: 0}) && (move.to == pair{col: 2, row: 0}) {
			if !c.WhiteQueenCastle {
				return false
			}
			if c.SquareIsThreatened(BLACK, sq("e1")) || c.SquareIsThreatened(BLACK, sq("d1")) || c.SquareIsThreatened(BLACK, sq("c1")) {
				return false
			}
			if c.getPiece(sq("d1")) != 0 || c.getPiece(sq("c1")) != 0 || c.getPiece(sq("b1")) != 0 {
				return false
			}
			// WARNING: I THINK THERE IS NO OTHER EDGE CASE
			return true
		}
		if (move.from == pair{col: 4, row: 0}) && (move.to == pair{col: 6, row: 0}) {
			if !c.WhiteKingCastle {
				return false
			}
			if c.SquareIsThreatened(BLACK, sq("e1")) || c.SquareIsThreatened(BLACK, sq("f1")) || c.SquareIsThreatened(BLACK, sq("g1")) {
				return false
			}
			if c.getPiece(sq("f1")) != 0 || c.getPiece(sq("g1")) != 0 {
				return false
			}
			// WARNING: I THINK THERE IS NO OTHER EDGE CASE
			return true
		}
	// en passant

	// promotion

	// if non special move
	default:
		c.erasePiece(move.from)
		c.putPiece(move.to, fromPiece)
	}

	// check legality of c.BoardState
	kingPosition := c.GetKingPosition()
	threat := c.SquareIsThreatened(c.WhiteToMove, kingPosition)

	// restore c.BoardState
	for _, square := range squaresToErase {
		c.erasePiece(square)
	}

	for i := 0; i < len(piecesPerSquare); i++ {
		c.putPiece(squaresToPutPieces[i], piecesPerSquare[i])
	}

	return threat
}

func (c *Chessboard) getPiece(square pair) int {
	if !inBounds(square) {
		return 0
	}
	for i := 0; i < len(c.BoardState); i++ {
		if c.BoardState[i]&uint64(1<<pairToInt(square)) == 1 {
			return i
		}
	}
	return 0
}

func isWhite(piece int) bool {
	return piece > 0 && piece < BKING
}

func (c *Chessboard) getMoveList() []Move {
	// Check all possible moves of all pieces
	// Check move legality of every move
	// Return filtered movement list
	var movements []Move
	for col := int8(0); col <= 8; col++ {
		for row := int8(0); row <= 8; row++ {
			var directions []pair
			var toSquares []pair
			var auxSquare pair
			from := pair{col: col, row: row}
			piece := c.getPiece(from)
			pieceColor := isWhite(piece)
			switch piece {
			case BPAWN:
				if row == 6 &&
					c.getPiece(pair{col: col, row: row - 1}) == 0 &&
					c.getPiece(pair{col: col, row: row - 2}) == 0 {
					toSquares = append(toSquares, pair{row: row - 2, col: col})
				}
				auxSquare = pair{col: col - 1, row: row - 1}
				if isWhite(c.getPiece(auxSquare)) || c.EnPassantSquare == auxSquare {
					toSquares = append(toSquares, auxSquare)
				}
				auxSquare = pair{col: col + 1, row: row - 1}
				if isWhite(c.getPiece(auxSquare)) || c.EnPassantSquare == auxSquare {
					toSquares = append(toSquares, auxSquare)
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
				auxSquare = pair{col: col - 1, row: row + 1}
				if !isWhite(c.getPiece(auxSquare)) || c.EnPassantSquare == auxSquare {
					toSquares = append(toSquares, pair{col: col - 1, row: row + 1})
				}
				auxSquare = pair{col: col + 1, row: row + 1}
				if !isWhite(c.getPiece(auxSquare)) || c.EnPassantSquare == auxSquare {
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
			for _, direction := range directions {
				currSquare := addPair(direction, from)
				currSquarePieceColor := isWhite(c.getPiece(currSquare))
				for {
					if !inBounds(currSquare) {
						break
					}
					if currSquarePieceColor == pieceColor {
						break
					}
					if currSquarePieceColor != pieceColor {
						toSquares = append(toSquares, currSquare)
						break
					}
					toSquares = append(toSquares, currSquare)
					currSquare = addPair(currSquare, direction)
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
						movements = append(movements, Move{from, to, 'R'})
						movements = append(movements, Move{from, to, 'N'})
						movements = append(movements, Move{from, to, 'B'})
					}
				}
			}
		}
	}
	return movements
}

func (c *Chessboard) makeMove(move string) error {
	if len(move) < 5 {
		return errors.New("Invalid move string")
	}
	version := move[0]
	var from pair
	var to pair
	var promotion int
	switch version {
	case '0':
		{
			from = pair{
				col: int8(move[1] - '`'),
				row: int8(move[2] - '0'),
			}
			to = pair{
				col: int8(move[3] - '`'),
				row: int8(move[4] - '0'),
			}
			promotion = int(move[5] - '0')
		}

	default:
		return errors.New("Invalid version of move")
	}

	if !c.CheckMoveLegality(Move{from: from, to: to, promotion: int(promotion)}) {
		return errors.New("the move is Illegal")
	}

	c.putPiece(to, c.getPiece(from))
	c.erasePiece(from)
	if promotion != 0 {
		c.putPiece(to, promotion)
	}

	return nil
}

func (c *Chessboard) DoSomething() {
	fmt.Printf("hello")
}
