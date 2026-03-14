package logic

import (
	"math/bits"
)

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


func (b *Board) getPossibleMoves(p uint64, team int, piece int) uint64 {
	friendlyPiece := b.getBitboard(team).all

	switch piece {
	case Pawn: return b.getPawnSquaresAttacked(p, team)
	case Bishop: return b.getBishopSquaresAttacked(p) & ^friendlyPiece
	case Knight: return b.getKnightSquaresAttacked(p) & ^friendlyPiece
	case Rook: return b.getRookSquaresAttacked(p) & ^friendlyPiece
	case Queen: return b.getQueenSquaresAttacked(p) & ^friendlyPiece
	case King: return ( b.getKingSquaresAttacked(p) | b.getCastling(p, team) ) & ^friendlyPiece
	}

	return 0
}


func (b *Board) isMoveSafe(from, to uint64, team int, piece int) bool {
	b.move(from, to, team, piece)

	var kingPos uint64
	if team == White {
		kingPos = b.whiteFigures.king
	} else {
		kingPos = b.blackFigures.king
	}

	res := !b.isChecked(kingPos, team)

	b.undo()

	return res
}


func (b *Board) getLegalMoves(p uint64, team int, piece int) uint64 {
	var legalMoves uint64

	pm := b.getPossibleMoves(p, team, piece)

	for pm != 0 {
		inx := bits.TrailingZeros64(pm)
		to := uint64(1) << inx

		if b.isMoveSafe(p, to, team, piece) {
			legalMoves |= to
		}

		pm &= (pm - 1)
	}

	return legalMoves
}



func (b *Board) move(from uint64, to uint64, team int, piece int) {
	eatenPieceType := b.GetPieceType(to)

	flip := from | to

	m := Move {
		from: from, to: to,
		team: team,
		eatenPiece: eatenPieceType,
		movingPiece: piece,
		oldFlags: b.flags,
	}

	b.history[b.historyLen] = m
	b.historyLen += 1

	if eatenPieceType != Empty {
		b.removePiece(to)
	}

	figs := b.getBitboard(team)

	if from == h1 || to == h1 { b.flags &= ^whiteShortCastling }
	if from == a1 || to == a1 { b.flags &= ^whiteLongCastling }
	if from == h8 || to == h8 { b.flags &= ^blackShortCastling }
	if from == a8 || to == a8 { b.flags &= ^blackLongCastling }

	switch piece {
	case Pawn:
		if to & ( rank1 | rank8 ) != 0 {
			figs.pawns &= ^from
			figs.queens |= to
		} else {
			figs.pawns ^= flip
		}

	case Bishop: figs.bishops ^= flip
	case Knight: figs.knights ^= flip
	case Rook: figs.rooks ^= flip
	case Queen: figs.queens ^= flip
	case King:

		if team == White && from == e1 {
			switch to {
			case g1: figs.rooks ^= h1 | f1
			case c1: figs.rooks ^= a1 | d1
			}
			b.flags &= notWhiteCastling
		} else if from == e8 {
			switch to {
			case g8: figs.rooks ^= h8 | f8
			case c8: figs.rooks ^= a8 | d8
			}
			b.flags &= notBlackCastling
		}

		figs.king ^= flip
	}

	b.updateAll()
}


func (b *Board) undo() {
	b.historyLen -= 1
	m := b.history[b.historyLen]

	from := m.from
	to := m.to

	flip := from | to

	t := m.movingPiece
	team := m.team

	figs := b.getBitboard(team)
	figs2 := b.getBitboard(getOppositeTeam(team))

	switch t {
	case Pawn:
		if to & (rank1 | rank8) != 0 {
			figs.queens &= ^to
			figs.pawns |= from
		} else {
			figs.pawns ^= flip
		}

	case Bishop: figs.bishops ^= flip
	case Knight: figs.knights ^= flip
	case Rook: figs.rooks ^= flip
	case Queen: figs.queens ^= flip

	case King:
		if team == White && from == e1 {
			switch to {
			case g1: figs.rooks ^= h1 | f1
			case c1: figs.rooks ^= a1 | d1
			}
		} else if from == e8 {
			switch to {
			case g8: figs.rooks ^= h8 | f8
			case c8: figs.rooks ^= a8 | d8
			}
		}

		figs.king ^= flip
	}

	switch m.eatenPiece {
		case Empty: break
		case Pawn: figs2.pawns |= to
		case Bishop: figs2.bishops |= to
		case Knight: figs2.knights |= to
		case Rook: figs2.rooks |= to
		case Queen: figs2.queens |= to
	}

	b.flags = m.oldFlags

	b.updateAll()
}


func (b *Board) checkGameStatus(team int) int {
	kingPos := b.getKing(team)
	inCheck := b.isChecked(kingPos, team)

	hasMoves := false

	bb := b.getBitboard(team)

	pieceBoards := [6]uint64{bb.pawns, bb.bishops, bb.knights, bb.rooks, bb.queens, bb.king}
	pieceTypes := [6]int{Pawn, Bishop, Knight, Rook, Queen, King}

	for i := 0; i < 6 && !hasMoves; i++ {
		bitsLeft := pieceBoards[i]
		for bitsLeft != 0 {
			p := bitsLeft & -bitsLeft
			if b.getLegalMoves(p, team, pieceTypes[i]) != 0 {
				hasMoves = true
				break
			}
			bitsLeft &= bitsLeft - 1
		}
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
		if (p << 7 & b.blackFigures.pawns & notH) != 0 { return true }
		if (p << 9 & b.blackFigures.pawns & notA) != 0 { return true }
	} else {
		if (p >> 7 & b.whiteFigures.pawns & notA) != 0 { return true }
		if (p >> 9 & b.whiteFigures.pawns & notH) != 0 { return true }
	}

	if (b.getBishopSquaresAttacked(p) & (enemyBitboard.bishops | enemyBitboard.queens)) != 0 { return true }
	if (b.getKnightSquaresAttacked(p) & enemyBitboard.knights) != 0 { return true }
	if (b.getRookSquaresAttacked(p) & (enemyBitboard.rooks | enemyBitboard.queens)) != 0 { return true }
	if (b.getKingSquaresAttacked(p) & enemyBitboard.king) != 0 { return true }

	return false
}
