package logic

import "math/bits"


// SquaresAttacked - это именно клетки под атакой, там и свои и чужие фигуры могут быть
// PossibleMoves - это возможные ходы: куда может походить или кого съесть (SquaresAttacked без клеток своей команды)
// LegalMoves - легальные ходы, которые получаются после обрабоки PossibleMoves, LegalMoves учитывают шах своему королю
// (если после PossibleMove свой король под шахом, то такой ход запрещен)


func CreateBoard() *Board {
	b := &Board{}

	b.arrangeFigures()
	b.updateAll()

	return b
}



func (b *Board) getPossibleMoves(p uint64) uint64 {
	t := b.getPieceType(p)

	friendlyPiece := b.getBitboard(b.getPieceTeam(p)).all

	switch t {
	case pawn: return b.getPawnSquaresAttacked(p)
	case bishop: return b.getBishopSquaresAttacked(p) & ^friendlyPiece
	case knight: return b.getKnightSquaresAttacked(p) & ^friendlyPiece
	case rook: return b.getRookSquaresAttacked(p) & ^friendlyPiece
	case queen: return b.getQueenSquaresAttacked(p) & ^friendlyPiece
	case king: return b.getKingSquaresAttacked(p) & ^friendlyPiece
	}

	return 0
}


func (b *Board) isMoveSafe(p1, p2 uint64) bool {
	team := b.getPieceTeam(p1)

	tempBoard := *b

	tempBoard.move(p1, p2)

	var kingPos uint64
	if team == White {
		kingPos = tempBoard.WhiteFigures.king
	} else {
		kingPos = tempBoard.BlackFigures.king
	}

	return !tempBoard.isChecked(kingPos, team)
}


func (b *Board) GetLegalMoves(p uint64) uint64 {
	if p & b.occupied == 0 || b.getPieceTeam(p) != b.turn {
		return 0
	}

	var legalMoves uint64

	pm := b.getPossibleMoves(p)

	for pm != 0 {
		inx := bits.TrailingZeros64(pm)
		to := uint64(1) << inx

		if b.isMoveSafe(p, to) {
			legalMoves |= to
		}

		pm &= (pm - 1)
	}

	return legalMoves
}



func (b *Board) GetEnemyFigures(p uint64) uint64 {
	if p & b.occupied == 0 {
		return 0
	}
	return b.getBitboard(getOppositeTeam(b.getPieceTeam(p))).all
}


func (b *Board) MakePlayerMove(from, to uint64) (int, int) {
	b.move(from, to)
	b.turn = getOppositeTeam(b.turn)
	return b.checkGameStatus(b.turn), b.turn
}


func (b *Board) move(from uint64, to uint64) {
	t := b.getPieceType(from)
	team := b.getPieceTeam(from)

	b.removePiece(from)
	b.removePiece(to)

	figs := b.getBitboard(team)

	switch t {
	case pawn:

	if team == White && to & rank8 != 0 {
		figs.queens |= to
	} else if team == Black && to & rank1 != 0 {
		figs.queens |= to
	} else {
		figs.pawns |= to
	}

	case bishop: figs.bishops |= to
	case knight: figs.knights |= to
	case rook: figs.rooks |= to
	case queen: figs.queens |= to
	case king: figs.king |= to
	}

	b.updateAll()
}


func (b *Board) checkGameStatus(team int) int {
	kingPos := b.getKing(team)
	inCheck := b.isChecked(kingPos, team)

	hasMoves := false

	allPieces := b.getBitboard(team).all

	for allPieces != 0 {
		idx := bits.TrailingZeros64(allPieces)
		p := uint64(1) << idx

		if b.GetLegalMoves(p) != 0 {
			hasMoves = true
			break
		}

		allPieces &= allPieces - 1
	}

	if inCheck {
		if !hasMoves {
			return Checkmate
		} else {
			return Check
		}
	} else if !hasMoves {
		return Draw
	}

	return None
}


func (b *Board) isChecked(p uint64, team int) bool {
	enemyBitboard := b.getBitboard(getOppositeTeam(team))

	if team == White {
		if (p << 7 & b.BlackFigures.pawns & notH) != 0 { return true }
		if (p << 9 & b.BlackFigures.pawns & notA) != 0 { return true }
	} else {
		if (p >> 7 & b.WhiteFigures.pawns & notA) != 0 { return true }
		if (p >> 9 & b.WhiteFigures.pawns & notH) != 0 { return true }
	}

	if (b.getBishopSquaresAttacked(p) & (enemyBitboard.bishops | enemyBitboard.queens)) != 0 { return true }
	if (b.getKnightSquaresAttacked(p) & enemyBitboard.knights) != 0 { return true }
	if (b.getRookSquaresAttacked(p) & (enemyBitboard.rooks | enemyBitboard.queens)) != 0 { return true }
	if (b.getKingSquaresAttacked(p) & enemyBitboard.king) != 0 { return true }

	return false
}


// func (b *Board) castling()
