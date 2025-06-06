package chessboard

import (
	"testing"
)

func TestAddPair(t *testing.T) {
	// test 1
	resultPair := addPair(
		pair{1, 2},
		pair{2, 1},
	)

	if resultPair != (pair{3, 3}) {
		t.Errorf(`addPair({1,2}, {2,1}) = %v, should equal {3,3}`, resultPair)
	}
}

func TestCreateChessgame(t *testing.T) {
	chessgame := CreateChessboard("new game")
}
