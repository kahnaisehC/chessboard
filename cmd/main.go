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
	fmt.Println(c.GetFEN())
	c.PrintBoard()
	for {
		var move string
		fmt.Scanf("%s", &move)
		err := c.MakeMove("0" + move + "_")
		if err != nil {
			fmt.Println(err.Error())
		}
		c.PrintBoard()
		fmt.Println(c.GetFEN())
	}
}
