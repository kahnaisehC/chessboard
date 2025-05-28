package main

import (
	"fmt"

	"github.com/kahnaisehC/chessboard"
)

func main() {
	c := chessboard.CreateChessgame()
	for _, row := range c.BoardState {
		fmt.Printf("%b\n", row)
	}
	c.DoSomething()
}
