package logic

// Здесь функции, которые напрямую работают с Board, это низкие функции, достаточно тупые функции
// Написаны просто перебором значений для скорости

type Bitboard struct {
	pawns uint64
	bishops uint64
	knights uint64
	rooks uint64
	queens uint64
	king uint64
	all uint64
}


type Board struct {
	WhiteFigures Bitboard
	BlackFigures Bitboard
	occupied uint64

	turn int
}


func (b *Board) arrangeFigures() {
	b.WhiteFigures.pawns = 0x000000000000FF00
	b.WhiteFigures.knights = 0x0000000000000042
	b.WhiteFigures.bishops = 0x0000000000000024
	b.WhiteFigures.rooks = 0x0000000000000081
	b.WhiteFigures.queens = 0x0000000000000008
	b.WhiteFigures.king = 0x00000000000000010

	b.BlackFigures.pawns = 0x00FF000000000000
	b.BlackFigures.knights = 0x4200000000000000
	b.BlackFigures.bishops = 0x2400000000000000
	b.BlackFigures.rooks = 0x8100000000000000
	b.BlackFigures.queens = 0x0800000000000000
	b.BlackFigures.king = 0x1000000000000000
}


func (b *Board) getPieceType(p uint64) int {
	if (b.WhiteFigures.pawns | b.BlackFigures.pawns) & p != 0 { return pawn }
	if (b.WhiteFigures.bishops | b.BlackFigures.bishops) & p != 0 { return bishop }
	if (b.WhiteFigures.knights | b.BlackFigures.knights) & p != 0 { return knight }
	if (b.WhiteFigures.rooks | b.BlackFigures.rooks) & p != 0 { return rook }
	if (b.WhiteFigures.queens | b.BlackFigures.queens) & p != 0 { return queen }
	if (b.WhiteFigures.king | b.BlackFigures.king) & p != 0 { return king }

	panic("not such type (getPieceType)")
}


func (b *Board) getPieceTeam(p uint64) int {
	if b.WhiteFigures.all & p != 0 {
		return White
	}

	if b.BlackFigures.all & p != 0 {
		return Black
	}

	panic("not such team (getPieceTeam)")
}


func (b *Board) getBitboard(team int) *Bitboard {
	 if team == White {
		return &b.WhiteFigures
	 }

	if team == Black {
		return &b.BlackFigures
	}

	panic("not right team (getFriendlyBitboard)")
}


func (b *Board) getKing(team int) uint64 {
	if team == White {
		return b.WhiteFigures.king
	}

	if team == Black {
		return b.BlackFigures.king
	}

	panic("not such team (getKing)")
}


func getOppositeTeam(team int) int {
	return 1 - team
}


func (b *Board) removePiece(p uint64) {
	b.WhiteFigures.pawns &= ^p
	b.WhiteFigures.bishops &= ^p
	b.WhiteFigures.knights &= ^p
	b.WhiteFigures.rooks &= ^p
	b.WhiteFigures.queens &= ^p
	b.WhiteFigures.king &= ^p

	b.BlackFigures.pawns &= ^p
	b.BlackFigures.bishops &= ^p
	b.BlackFigures.knights &= ^p
	b.BlackFigures.rooks &= ^p
	b.BlackFigures.queens &= ^p
	b.BlackFigures.king &= ^p
}


func (b *Board) updateAll() {
	b.WhiteFigures.all = b.WhiteFigures.pawns | b.WhiteFigures.bishops | b.WhiteFigures.knights | b.WhiteFigures.rooks | b.WhiteFigures.queens | b.WhiteFigures.king
	b.BlackFigures.all = b.BlackFigures.pawns | b.BlackFigures.bishops | b.BlackFigures.knights | b.BlackFigures.rooks | b.BlackFigures.queens | b.BlackFigures.king
	b.occupied = b.WhiteFigures.all | b.BlackFigures.all
}


func (b *Board) GetFigures() [2][6]uint64 {
	var res [2][6]uint64

	res[White][pawn] = b.WhiteFigures.pawns
	res[White][bishop] = b.WhiteFigures.bishops
	res[White][knight] = b.WhiteFigures.knights
	res[White][rook] = b.WhiteFigures.rooks
	res[White][queen] = b.WhiteFigures.queens
	res[White][king] = b.WhiteFigures.king

	res[Black][pawn] = b.BlackFigures.pawns
	res[Black][bishop] = b.BlackFigures.bishops
	res[Black][knight] = b.BlackFigures.knights
	res[Black][rook] = b.BlackFigures.rooks
	res[Black][queen] = b.BlackFigures.queens
	res[Black][king] = b.BlackFigures.king

	return res
}
