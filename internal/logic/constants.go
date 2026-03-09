package logic


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
)


const (
	None = 0
	Check = 1
	Checkmate = 2
	Draw = 3
)


const (
	ShortSide = iota
	LongSide
)


const (
	Rank1 uint64 = 0x00000000000000FF
	Rank2 uint64 = 0x000000000000FF00
	Rank7 uint64 = 0x00FF000000000000
	Rank8 uint64 = 0xFF00000000000000
)


const (
	NotA uint64 = 0xFEFEFEFEFEFEFEFE
	NotH uint64 = 0x7F7F7F7F7F7F7F7F
	NotAB uint64 = 0xfcfcfcfcfcfcfcfc
	NotHG uint64 = 0x3f3f3f3f3f3f3f3f
)


const (
	A1 uint64 = 1 << iota;
		B1; C1; D1; E1; F1; G1; H1
	A2; B2; C2; D2; E2; F2; G2; H2
	A3; B3; C3; D3; E3; F3; G3; H3
	A4; B4; C4; D4; E4; F4; G4; H4
	A5; B5; C5; D5; E5; F5; G5; H5
	A6; B6; C6; D6; E6; F6; G6; H6
	A7; B7; C7; D7; E7; F7; G7; H7
	A8; B8; C8; D8; E8; F8; G8; H8
)


const (
	WhiteShortEmptyCells uint64 = 0x60
	WhiteShortSafeCells uint64 = 0x70
	WhiteLongEmptyCells uint64 = 0x0E
	WhiteLongSafeCells uint64 = 0x1C

	BlackShortEmptyCells uint64 = 0x6000000000000000
	BlackShortSafeCells uint64 = 0x7000000000000000
	BlackLongEmptyCells uint64 = 0x0E00000000000000
	BlackLongSafeCells uint64 = 0x1C00000000000000
)
