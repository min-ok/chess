package logic

// Здесь функции, которые напрямую работают с Board, это низкие функции, достаточно тупые функции
// Написаны просто перебором значений для скорости


type Board struct {
	bitboard [2][6]uint64

	teamOccupied [2]uint64
	allOccupied uint64

	flags uint8

	turn int

	history [2048]Move
	historyLen int
}


type Move struct {
	from uint64
	to uint64

	team int

	oldFlags uint8

	eatenPiece int
	movingPiece int
}


func (b *Board) arrangeFigures() {
	b.bitboard[White][Pawn] = 0x000000000000FF00
	b.bitboard[White][Bishop] = 0x0000000000000024
	b.bitboard[White][Knight] = 0x0000000000000042
	b.bitboard[White][Rook] = 0x0000000000000081
	b.bitboard[White][Queen] = 0x0000000000000008
	b.bitboard[White][King] = 0x00000000000000010

	b.bitboard[Black][Pawn] = 0x00FF000000000000
	b.bitboard[Black][Bishop] = 0x2400000000000000
	b.bitboard[Black][Knight] = 0x4200000000000000
	b.bitboard[Black][Rook] = 0x8100000000000000
	b.bitboard[Black][Queen] = 0x0800000000000000
	b.bitboard[Black][King] = 0x1000000000000000

	b.flags = 0x0F
}


func (b *Board) getPieceType(p uint64) int {
	for i := 0; i < 6; i += 1 {
		if (b.bitboard[White][i] | b.bitboard[Black][i]) & p != 0 { return i }
	}

	return Empty
}


func (b *Board) getPieceTeam(p uint64) int {
	if b.teamOccupied[White] & p != 0 {
		return White
	}

	if b.teamOccupied[Black] & p != 0 {
		return Black
	}

	panic("not such team (getPieceTeam)")
}


func getOppositeTeam(team int) int {
	return 1 - team
}


func (b *Board) updateAll() {
	b.teamOccupied[White] = 0
	b.teamOccupied[Black] = 0

	for i := 0; i < 6; i += 1 {
		b.teamOccupied[White] |= b.bitboard[White][i]
		b.teamOccupied[Black] |= b.bitboard[Black][i]
	}
	b.allOccupied = b.teamOccupied[White] | b.teamOccupied[Black]
}


func (b *Board) GetFigures() [2][6]uint64 {
	return b.bitboard
}
