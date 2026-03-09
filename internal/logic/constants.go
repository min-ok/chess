package logic


const (
	White = 0
	Black = 1
)


const (
	pawn = 0
	bishop = 1
	knight = 2
	rook = 3
	queen = 4
	king = 5
)


const (
	None = 0
	Check = 1
	Checkmate = 2
	Draw = 3
)


const (
	rank1 = uint64(0x00000000000000FF)
	rank2 = uint64(0x000000000000FF00)
	rank7 = uint64(0x00FF000000000000)
	rank8 = uint64(0xFF00000000000000)
	notA  = uint64(0xFEFEFEFEFEFEFEFE)
	notH  = uint64(0x7F7F7F7F7F7F7F7F)
	notAB = uint64(0xfcfcfcfcfcfcfcfc)
	notHG = uint64(0x3f3f3f3f3f3f3f3f)
)
