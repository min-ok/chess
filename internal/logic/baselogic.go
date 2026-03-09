package logic

// Здесь функции, которые напрямую работают с Board, это низкие функции, достаточно тупфые функции
// (но из-за отсутсвия циклов и еще чего-то для лучшей читаемости, они досаточно быстрые наверно (так gemini сказал, я тут ни при чем))

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
}


func (b *Board) getPieceType(p uint64) int {
	if (b.whiteFigures.pawns | b.blackFigures.pawns) & p != 0 { return pawn }
	if (b.whiteFigures.bishops | b.blackFigures.bishops) & p != 0 { return bishop }
	if (b.whiteFigures.knights | b.blackFigures.knights) & p != 0 { return knight }
	if (b.whiteFigures.rooks | b.blackFigures.rooks) & p != 0 { return rook }
	if (b.whiteFigures.queens | b.blackFigures.queens) & p != 0 { return queen }
	if (b.whiteFigures.king | b.blackFigures.king) & p != 0 { return king }

	panic("not such type (getPieceType)")
}


func (b *Board) getPieceTeam(p uint64) int {
	if b.whiteFigures.all & p != 0 {
		return white
	}

	if b.blackFigures.all & p != 0 {
		return black
	}

	panic("not such team (getPieceTeam)")
}


func (b *Board) getBitboard(team int) *Bitboard {
	 if team == white {
		return &b.whiteFigures
	 }

	if team == black {
		return &b.blackFigures
	}

	panic("not right team (getFriendlyBitboard)")
}


func (b *Board) getKing(team int) uint64 {
	if team == white {
		return b.whiteFigures.king
	}

	if team == black {
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

	res[white][pawn] = b.whiteFigures.pawns
	res[white][bishop] = b.whiteFigures.bishops
	res[white][knight] = b.whiteFigures.knights
	res[white][rook] = b.whiteFigures.rooks
	res[white][queen] = b.whiteFigures.queens
	res[white][king] = b.whiteFigures.king

	res[black][pawn] = b.blackFigures.pawns
	res[black][bishop] = b.blackFigures.bishops
	res[black][knight] = b.blackFigures.knights
	res[black][rook] = b.blackFigures.rooks
	res[black][queen] = b.blackFigures.queens
	res[black][king] = b.blackFigures.king

	return res
}
