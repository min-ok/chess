package logic

import (
	"math/bits"
)


func (b *Board) getPawnSquaresAttacked(p uint64, team int) uint64 {
	var moves uint64

	if (team == White) {
		moves |= (p << 7) & b.teamOccupied[Black] & notH
	 	moves |= (p << 9) & b.teamOccupied[Black] & notA

		step := p << 8
		if step & b.allOccupied == 0 {
			moves |= step
			if p & rank2 != 0 {
				step = p << 16
				if step & b.allOccupied == 0 { moves |= step }
			}
		}
	} else {
		moves |= (p >> 7) & b.teamOccupied[White] & notA
	 	moves |= (p >> 9) & b.teamOccupied[White] & notH

		step := p >> 8
		if step & b.allOccupied == 0 {
			moves |= step

			if p & rank7 != 0 {
				step = p >> 16
				if step & b.allOccupied == 0 { moves |= step }
			}
		}
	}

	return moves
}


func (b *Board) getBishopSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= b.ray(p, 9, 0, notH)
	moves |= b.ray(p, 7, 0, notA)
	moves |= b.ray(p, 0, 7, notH)
	moves |= b.ray(p, 0, 9, notA)

	return moves
}


func (b *Board) getKnightSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= (p << 17) & notA
	moves |= (p << 15) & notH
	moves |= (p << 10) & notAB
	moves |= (p << 6) & notHG

	moves |= (p >> 15) & notA
	moves |= (p >> 17) & notH
	moves |= (p >> 6) & notAB
	moves |= (p >> 10) & notHG

	return moves
}


func (b *Board) getRookSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= b.ray(p, 8, 0, 0)
	moves |= b.ray(p, 0, 8, 0)
	moves |= b.ray(p, 0, 1, notA)
	moves |= b.ray(p, 1, 0, notH)

	return moves
}


func (b *Board) getQueenSquaresAttacked(p uint64) uint64 {
	return b.getBishopSquaresAttacked(p) | b.getRookSquaresAttacked(p)
}


func (b *Board) getKingSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= (p << 1) & notA
	moves |= (p >> 1) & notH

	moves |= (p << 8)
	moves |= (p >> 8)

	moves |= (p << 7) & notH
	moves |= (p << 9) & notA

	moves |= (p >> 7) & notA
	moves |= (p >> 9) & notH

	return moves
}


func (b *Board) getCastling(team int) uint64 {
	var res uint64

	kingPos := b.bitboard[team][King]

	if team == White && kingPos == e1 {
		if b.flags & whiteShortCastling != 0 && (b.allOccupied & whiteShortEmptyCells) == 0 && b.isPathSafe(whiteShortSafeCells, White) {
			res |= g1
		}
		if b.flags & whiteLongCastling != 0 && (b.allOccupied & whiteLongEmptyCells) == 0 && b.isPathSafe(whiteLongSafeCells, White) {
			res |= c1
		}
	} else if team == Black && kingPos == e8 {
		if b.flags & blackShortCastling != 0 && (b.allOccupied & blackShortEmptyCells) == 0 && b.isPathSafe(blackShortSafeCells, Black) {
			res |= g8
		}
		if b.flags & blackLongCastling != 0 && (b.allOccupied & blackLongEmptyCells) == 0 && b.isPathSafe(blackLongSafeCells, Black) {
			res |= c8
		}
	}

	return res
}


func (b *Board) isPathSafe(cells uint64, team int) bool {
	for cells != 0 {
		idx := bits.TrailingZeros64(cells)
		if b.isChecked(uint64(1)<<idx, team) {
			return false
		}
		cells &= cells - 1
	}
	return true
}



func (b *Board) ray(p uint64, shift1 uint64, shift2 uint64, mask uint64) uint64 {
	var res uint64
	curr := p

	for {
		if mask != 0 && (curr & mask == 0) { break }

		curr <<= shift1
		curr >>= shift2

		if curr == 0 { break }

		res |= curr

		if curr & b.allOccupied != 0 { break }
	}

	return res
}
