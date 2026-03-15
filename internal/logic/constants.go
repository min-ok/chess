package logic


const (
	botDepth = 1
)


const (
	White = 0
	Black = 1
)


const (
	Pawn = 0
	Bishop = 1
	Knight = 2
	Rook = 3
	Queen = 4
	King = 5
	Empty = -1
)


var pieceTypes = [6]int{Pawn, Bishop, Knight, Rook, Queen, King}


const (
	None = 0
	Check = 1
	Checkmate = 2
	Draw = 3
)


const (
	pawnCost = 1000
	bishopCost = 3000
	knightCost = 4000
	rookCost = 5000
	queenCost = 9000
	kingCost = 1000000
)


const (
	whiteShortCastling uint8 = 0x01
	whiteLongCastling uint8 = 0x02
	blackShortCastling uint8 = 0x04
	blackLongCastling uint8 = 0x08

	notWhiteCastling uint8 = 0x03
	notBlackCastling uint8 = 0x0C
)


const (
	rank1 uint64 = 0x00000000000000FF
	rank2 uint64 = 0x000000000000FF00
	rank7 uint64 = 0x00FF000000000000
	rank8 uint64 = 0xFF00000000000000
)


const (
	notA uint64 = 0xFEFEFEFEFEFEFEFE
	notH uint64 = 0x7F7F7F7F7F7F7F7F
	notAB uint64 = 0xfcfcfcfcfcfcfcfc
	notHG uint64 = 0x3f3f3f3f3f3f3f3f
)


const (
	a1 uint64 = 1 << iota;
		b1; c1; d1; e1; f1; g1; h1
	a2; b2; c2; d2; e2; f2; g2; h2
	a3; b3; c3; d3; e3; f3; g3; h3
	a4; b4; c4; d4; e4; f4; g4; h4
	a5; b5; c5; d5; e5; f5; g5; h5
	a6; b6; c6; d6; e6; f6; g6; h6
	a7; b7; c7; d7; e7; f7; g7; h7
	a8; b8; c8; d8; e8; f8; g8; h8
)


const (
	whiteShortEmptyCells uint64 = 0x60
	whiteShortSafeCells uint64 = 0x70
	whiteLongEmptyCells uint64 = 0x0E
	whiteLongSafeCells uint64 = 0x1C

	blackShortEmptyCells uint64 = 0x6000000000000000
	blackShortSafeCells uint64 = 0x7000000000000000
	blackLongEmptyCells uint64 = 0x0E00000000000000
	blackLongSafeCells uint64 = 0x1C00000000000000
)
