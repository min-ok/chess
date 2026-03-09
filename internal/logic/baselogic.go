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
	whiteFigures Bitboard
	blackFigures Bitboard
	occupied uint64

	whiteLongCastling bool
	whiteShortCastling bool
	blackLongCastling bool
	blackShortCastling bool

	turn int
}


func (b *Board) arrangeFigures() {
	b.whiteFigures.pawns = 0x000000000000FF00
	b.whiteFigures.knights = 0x0000000000000042
	b.whiteFigures.bishops = 0x0000000000000024
	b.whiteFigures.rooks = 0x0000000000000081
	b.whiteFigures.queens = 0x0000000000000008
	b.whiteFigures.king = 0x00000000000000010

	b.blackFigures.pawns = 0x00FF000000000000
	b.blackFigures.knights = 0x4200000000000000
	b.blackFigures.bishops = 0x2400000000000000
	b.blackFigures.rooks = 0x8100000000000000
	b.blackFigures.queens = 0x0800000000000000
	b.blackFigures.king = 0x1000000000000000

	b.whiteLongCastling = true
	b.whiteShortCastling = true
	b.blackLongCastling = true
	b.blackShortCastling = true
}


func (b *Board) getPieceType(p uint64) int {
	if (b.whiteFigures.pawns | b.blackFigures.pawns) & p != 0 { return Pawn }
	if (b.whiteFigures.bishops | b.blackFigures.bishops) & p != 0 { return Bishop }
	if (b.whiteFigures.knights | b.blackFigures.knights) & p != 0 { return Knight }
	if (b.whiteFigures.rooks | b.blackFigures.rooks) & p != 0 { return Rook }
	if (b.whiteFigures.queens | b.blackFigures.queens) & p != 0 { return Queen }
	if (b.whiteFigures.king | b.blackFigures.king) & p != 0 { return King }

	panic("not such type (getPieceType)")
}


func (b *Board) getPieceTeam(p uint64) int {
	if b.whiteFigures.all & p != 0 {
		return White
	}

	if b.blackFigures.all & p != 0 {
		return Black
	}

	panic("not such team (getPieceTeam)")
}


func (b *Board) getBitboard(team int) *Bitboard {
	 if team == White {
		return &b.whiteFigures
	 }

	if team == Black {
		return &b.blackFigures
	}

	panic("not right team (getFriendlyBitboard)")
}


func (b *Board) getKing(team int) uint64 {
	if team == White {
		return b.whiteFigures.king
	}

	if team == Black {
		return b.blackFigures.king
	}

	panic("not such team (getKing)")
}


func getOppositeTeam(team int) int {
	return 1 - team
}


func (b *Board) removePiece(p uint64) {
	b.whiteFigures.pawns &= ^p
	b.whiteFigures.bishops &= ^p
	b.whiteFigures.knights &= ^p
	b.whiteFigures.rooks &= ^p
	b.whiteFigures.queens &= ^p
	b.whiteFigures.king &= ^p

	b.blackFigures.pawns &= ^p
	b.blackFigures.bishops &= ^p
	b.blackFigures.knights &= ^p
	b.blackFigures.rooks &= ^p
	b.blackFigures.queens &= ^p
	b.blackFigures.king &= ^p
}


func (b *Board) updateAll() {
	b.whiteFigures.all = b.whiteFigures.pawns | b.whiteFigures.bishops | b.whiteFigures.knights | b.whiteFigures.rooks | b.whiteFigures.queens | b.whiteFigures.king
	b.blackFigures.all = b.blackFigures.pawns | b.blackFigures.bishops | b.blackFigures.knights | b.blackFigures.rooks | b.blackFigures.queens | b.blackFigures.king
	b.occupied = b.whiteFigures.all | b.blackFigures.all
}


func (b *Board) GetFigures() [2][6]uint64 {
	var res [2][6]uint64

	res[White][Pawn] = b.whiteFigures.pawns
	res[White][Bishop] = b.whiteFigures.bishops
	res[White][Knight] = b.whiteFigures.knights
	res[White][Rook] = b.whiteFigures.rooks
	res[White][Queen] = b.whiteFigures.queens
	res[White][King] = b.whiteFigures.king

	res[Black][Pawn] = b.blackFigures.pawns
	res[Black][Bishop] = b.blackFigures.bishops
	res[Black][Knight] = b.blackFigures.knights
	res[Black][Rook] = b.blackFigures.rooks
	res[Black][Queen] = b.blackFigures.queens
	res[Black][King] = b.blackFigures.king

	return res
}
