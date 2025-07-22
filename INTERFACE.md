# Constants

empty = 0
WKING = 1
WQUEEN = 2
WROOK = 3
WBISHOP = 4
WKNIGHT = 5
WPAWN = 6
BKING = 7
BQUEEN = 8
BROOK = 9
BBISHOP = 10
BKNIGHT = 11
BPAWN = 12

a1 = 0
a2 = 1
a3 = 2
a4 = 3
a5 = 4
a6 = 5
a7 = 6
a8 = 7
b1 = 8
b2 = 9
b3 = 10
b4 = 11
b5 = 12
b6 = 13
b7 = 14
b8 = 15
c1 = 16
c2 = 17
c3 = 18
c4 = 19
c5 = 20
c6 = 21
c7 = 22
c8 = 23
d1 = 24
d2 = 25
d3 = 26
d4 = 27
d5 = 28
d6 = 29
d7 = 30
d8 = 31
e1 = 32
e2 = 33
e3 = 34
e4 = 35
e5 = 36
e6 = 37
e7 = 38
e8 = 39
f1 = 40
f2 = 41
f3 = 42
f4 = 43
f5 = 44
f6 = 45
f7 = 46
f8 = 47
g1 = 48
g2 = 49
g3 = 50
g4 = 51
g5 = 52
g6 = 53
g7 = 54
g8 = 55
h1 = 56
h2 = 57
h3 = 58
h4 = 59
h5 = 60
h6 = 61
h7 = 62
h8 = 63



# Types


type Piece uint8

type Square uint8

### type Move 
type Move struct{
    from square
    to square
    promotion piece 
    fromPiece piece
    toPiece piece

    isLongCastle bool
    isShortCastle bool
    isEnPassant bool
    
}

### type Chessboard
type Chessboard struct {
	// Configuration fields
	Variant     string
	StartingFEN string // Default to initialFEN

	// State fields
	BoardState       [13]uint64
	WhiteToMove      bool
	EnPassantSquare  square 
	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool

	// Information fields
	Moves            []Move
	HalfmoveClock   int
	FullmoveCounter int

	// PGN necessary fields
	Event string
	Site string
	Date string
	Round int
	White string
	Black string
	Result result
    
}


# Functions

## func 




/*
// Main interaction interface

// func (c *Chessboard) CheckMoveLegality() bool
// func (c *Chessboard) GetMoveList() []Move
// func (c *Chessboard) GetTripleRepetition() bool  // Check if there is a tripple repetition to claim draw
// func (c *Chessboard) InsufficientMaterial() bool // Check if there is sufficient material
// func (c *Chessboard) GetResult() Result
// func (c *Chessboard) GetFEN() string
func (c *Chessboard) GetPGN() string

// func (c *Chessboard) MakeMove(Move) bool // Returns if the move was executed
// func (c *Chessboard) DeclareResult(Result)
// func (c *Chessboard) UndoMove()

// func (c *Chessboard) PrintBoard()
// func (c *Chessboard) GameOver() bool




*/


