package main

import (
	"fmt"

	"github.com/kahnaisehC/chessboard"
)

func main() {
	c := chessboard.CreateChessboard("")
	for _, row := range c.BoardState {
		fmt.Printf("%b\n", row)
	}
	c.PrintBoard()
}
