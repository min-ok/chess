package logic

import (
	// "fmt"
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


func (b *Board) getPossibleMoves(p uint64, piece int, team int) uint64 {
	friendlyPiece := b.teamOccupied[team]

	switch piece {
	case Pawn: return b.getPawnSquaresAttacked(p, team)
	case Bishop: return b.getBishopSquaresAttacked(p) & ^friendlyPiece
	case Knight: return b.getKnightSquaresAttacked(p) & ^friendlyPiece
	case Rook: return b.getRookSquaresAttacked(p) & ^friendlyPiece
	case Queen: return b.getQueenSquaresAttacked(p) & ^friendlyPiece
	case King: return ( b.getKingSquaresAttacked(p) | b.getCastling(team) ) & ^friendlyPiece
	}

	return 0
}


func (b *Board) isMoveSafe(from, to uint64, fromPieceType int, toPieceType int, team int) bool {
	b.move(from, to, fromPieceType, toPieceType, team)

	res := !b.isChecked(b.bitboard[team][King], team)

	b.undo()

	return res
}


func (b *Board) getLegalMoves(p uint64, piece int, team int) uint64 {
	var legalMoves uint64

	pm := b.getPossibleMoves(p, piece, team)

	for pm != 0 {
		inx := bits.TrailingZeros64(pm)
		to := uint64(1) << inx

		if b.isMoveSafe(p, to, piece, b.getPieceType(to), team) {
			legalMoves |= to
		}

		pm &= (pm - 1)
	}

	return legalMoves
}



func (b *Board) move(from uint64, to uint64, fromPieceType, toPieceType int, team int) {
	flip := from | to

	m := Move {
		from: from, to: to,
		team: team,
		movingPiece: fromPieceType,
		eatenPiece: toPieceType,
		oldFlags: b.flags,
	}

	b.history[b.historyLen] = m
	b.historyLen += 1

	if toPieceType != Empty {
		b.bitboard[getOppositeTeam(team)][toPieceType] &= ^to
	}

	bbt := &b.bitboard[team]

	if from == h1 || to == h1 { b.flags &= ^whiteShortCastling }
	if from == a1 || to == a1 { b.flags &= ^whiteLongCastling }
	if from == h8 || to == h8 { b.flags &= ^blackShortCastling }
	if from == a8 || to == a8 { b.flags &= ^blackLongCastling }

	switch fromPieceType {
	case Pawn:
		if to & ( rank1 | rank8 ) != 0 {
			bbt[Pawn] &= ^from
			bbt[Queen] |= to
		} else {
			bbt[Pawn] ^= flip
		}
	case King:
		if team == White && from == e1 {
			switch to {
			case g1: bbt[Rook] ^= h1 | f1
			case c1: bbt[Rook] ^= a1 | d1
			}
			b.flags &= notWhiteCastling
		} else if team == Black && from == e8 {
			switch to {
			case g8: bbt[Rook] ^= h8 | f8
			case c8: bbt[Rook] ^= a8 | d8
			}
			b.flags &= notBlackCastling
		}

		bbt[King] ^= flip
	default:
		bbt[fromPieceType] ^= flip
	}

	b.updateAll()
}


func (b *Board) undo() {
	b.historyLen -= 1
	m := b.history[b.historyLen]

	flip := m.from | m.to

	bbt := &b.bitboard[m.team]

	switch m.movingPiece {
	case Pawn:
		if m.to & (rank1 | rank8) != 0 {
			bbt[Queen] &= ^m.to
			bbt[Pawn] |= m.from
		} else {
			bbt[Pawn] ^= flip
		}

	case King:
		if m.team == White && m.from == e1 {
			switch m.to {
			case g1: bbt[Rook] ^= h1 | f1
			case c1: bbt[Rook] ^= a1 | d1
			}
		} else if m.team == Black && m.from == e8 {
			switch m.to {
			case g8: bbt[Rook] ^= h8 | f8
			case c8: bbt[Rook] ^= a8 | d8
			}
		}

		bbt[King] ^= flip

		default:
			bbt[m.movingPiece] ^= flip
	}

	if m.eatenPiece != Empty {
		b.bitboard[getOppositeTeam(m.team)][m.eatenPiece] |= m.to
	}

	b.flags = m.oldFlags

	b.updateAll()
}


func (b *Board) checkGameStatus(team int) int {
	hasMoves := false

	for i := 0; i < 6 && !hasMoves; i++ {
		bitsLeft := b.bitboard[team][i]
		for bitsLeft != 0 {
			p := bitsLeft & -bitsLeft
			if b.getLegalMoves(p, pieceTypes[i], team) != 0 {
				hasMoves = true
				break
			}
			bitsLeft &= bitsLeft - 1
		}
	}

	if b.isChecked(b.bitboard[team][King], team) {
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
	enemyBitboard := b.bitboard[getOppositeTeam(team)]

	if team == White {
		if (p << 7 & b.bitboard[Black][Pawn] & notH) != 0 { return true }
		if (p << 9 & b.bitboard[Black][Pawn] & notA) != 0 { return true }
	} else {
		if (p >> 7 & b.bitboard[White][Pawn] & notA) != 0 { return true }
		if (p >> 9 & b.bitboard[White][Pawn] & notH) != 0 { return true }
	}

	if (b.getBishopSquaresAttacked(p) & (enemyBitboard[Bishop] | enemyBitboard[Queen])) != 0 { return true }
	if (b.getKnightSquaresAttacked(p) & enemyBitboard[Knight]) != 0 { return true }
	if (b.getRookSquaresAttacked(p) & (enemyBitboard[Rook] | enemyBitboard[Queen])) != 0 { return true }
	if (b.getKingSquaresAttacked(p) & enemyBitboard[King]) != 0 { return true }

	return false
}
