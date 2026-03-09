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
	case Pawn: return b.getPawnSquaresAttacked(p)
	case Bishop: return b.getBishopSquaresAttacked(p) & ^friendlyPiece
	case Knight: return b.getKnightSquaresAttacked(p) & ^friendlyPiece
	case Rook: return b.getRookSquaresAttacked(p) & ^friendlyPiece
	case Queen: return b.getQueenSquaresAttacked(p) & ^friendlyPiece
	case King: return ( b.getKingSquaresAttacked(p) | b.getCastling(p) ) & ^friendlyPiece
	}

	return 0
}


func (b *Board) isMoveSafe(p1, p2 uint64) bool {
	team := b.getPieceTeam(p1)

	tempBoard := *b

	tempBoard.move(p1, p2)

	var kingPos uint64
	if team == White {
		kingPos = tempBoard.whiteFigures.king
	} else {
		kingPos = tempBoard.blackFigures.king
	}

	return !tempBoard.isChecked(kingPos, team)
}


func (b *Board) GetLegalMoves(p uint64) uint64 {
	if p & b.occupied == 0 /*|| b.getPieceTeam(p) != b.turn*/ {
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

	if from == A1 || to == A1 { b.whiteLongCastling = false }
	if from == H1 || to == H1 { b.whiteShortCastling = false }
	if from == A8 || to == A8 { b.blackLongCastling = false }
	if from == H8 || to == H8 { b.blackShortCastling = false }


	switch t {
	case Pawn:

	if team == White && to & Rank8 != 0 {
		figs.queens |= to
	} else if team == Black && to & Rank1 != 0 {
		figs.queens |= to
	} else {
		figs.pawns |= to
	}

	case Bishop: figs.bishops |= to
	case Knight: figs.knights |= to
	case Rook: figs.rooks |= to
	case Queen: figs.queens |= to
	case King:


	if team == White && from == E1 {
		switch to {
		case G1:
			b.removePiece(H1)
			figs.rooks |= F1
		case C1:
			b.removePiece(A1)
			figs.rooks |= D1
		}
		b.whiteShortCastling, b.whiteLongCastling = false, false
	} else if from == E8 {
		switch to {
		case G8:
			b.removePiece(H8)
			figs.rooks |= F8
		case C8:
			b.removePiece(A8)
			figs.rooks |= D8
		}
		b.blackShortCastling, b.blackLongCastling = false, false
	}

	figs.king |= to
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
		if (p << 7 & b.blackFigures.pawns & NotH) != 0 { return true }
		if (p << 9 & b.blackFigures.pawns & NotA) != 0 { return true }
	} else {
		if (p >> 7 & b.whiteFigures.pawns & NotA) != 0 { return true }
		if (p >> 9 & b.whiteFigures.pawns & NotH) != 0 { return true }
	}

	if (b.getBishopSquaresAttacked(p) & (enemyBitboard.bishops | enemyBitboard.queens)) != 0 { return true }
	if (b.getKnightSquaresAttacked(p) & enemyBitboard.knights) != 0 { return true }
	if (b.getRookSquaresAttacked(p) & (enemyBitboard.rooks | enemyBitboard.queens)) != 0 { return true }
	if (b.getKingSquaresAttacked(p) & enemyBitboard.king) != 0 { return true }

	return false
}
