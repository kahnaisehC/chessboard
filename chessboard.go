package chessboard

import (
	"errors"
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

type pair struct {
	col int8
	row int8
}

type Chessgame struct {
	Moves       []string
	Variant     string
	WhiteToMove bool
	BoardState  []uint64

	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool
	EnPassantSquare  pair
}

func addPair(a, b pair) pair {
	return pair{
		col: a.col + b.col,
		row: a.row + b.row,
	}
}

func CreateChessgame() Chessgame {
	chessgame := Chessgame{}
	return chessgame
}
