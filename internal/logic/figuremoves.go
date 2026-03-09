package logic

import (
	"math/bits"
)


func (b *Board) getPawnSquaresAttacked(p uint64) uint64 {
	var moves uint64

	t := b.getPieceTeam(p)
	if t == -1 { panic(t) }

	if (t == White) {
		moves |= (p << 7) & b.blackFigures.all & NotH
	 	moves |= (p << 9) & b.blackFigures.all & NotA

		step := p << 8
		if step & b.occupied == 0 {
			moves |= step
			if p & Rank2 != 0 {
				step = p << 16
				if step & b.occupied == 0 { moves |= step }
			}
		}
	} else {
		moves |= (p >> 7) & b.whiteFigures.all & NotA
	 	moves |= (p >> 9) & b.whiteFigures.all & NotH

		step := p >> 8
		if step & b.occupied == 0 {
			moves |= step

			if p & Rank7 != 0 {
				step = p >> 16
				if step & b.occupied == 0 { moves |= step }
			}
		}
	}

	return moves
}


func (b *Board) getBishopSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= b.ray(p, 9, 0, NotH)
	moves |= b.ray(p, 7, 0, NotA)
	moves |= b.ray(p, 0, 7, NotH)
	moves |= b.ray(p, 0, 9, NotA)

	return moves
}


func (b *Board) getKnightSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= (p << 17) & NotA
	moves |= (p << 15) & NotH
	moves |= (p << 10) & NotAB
	moves |= (p << 6) & NotHG

	moves |= (p >> 15) & NotA
	moves |= (p >> 17) & NotH
	moves |= (p >> 6) & NotAB
	moves |= (p >> 10) & NotHG

	return moves
}


func (b *Board) getRookSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= b.ray(p, 8, 0, 0)
	moves |= b.ray(p, 0, 8, 0)
	moves |= b.ray(p, 0, 1, NotA)
	moves |= b.ray(p, 1, 0, NotH)

	return moves
}


func (b *Board) getQueenSquaresAttacked(p uint64) uint64 {
	return b.getBishopSquaresAttacked(p) | b.getRookSquaresAttacked(p)
}


func (b *Board) getKingSquaresAttacked(p uint64) uint64 {
	var moves uint64

	moves |= (p << 1) & NotA
	moves |= (p >> 1) & NotH

	moves |= (p << 8)
	moves |= (p >> 8)

	moves |= (p << 7) & NotH
	moves |= (p << 9) & NotA

	moves |= (p >> 7) & NotA
	moves |= (p >> 9) & NotH

	return moves
}


func (b *Board) getCastling(p uint64) uint64 {
	var res uint64
	team := b.getPieceTeam(p)

	if team == White {
		if b.whiteShortCastling && (b.occupied & WhiteShortEmptyCells) == 0 && b.isPathSafe(WhiteShortSafeCells, White) {
			res |= G1
		}
		if b.whiteLongCastling && (b.occupied&WhiteLongEmptyCells) == 0 && b.isPathSafe(WhiteLongSafeCells, White) {
				res |= C1
		}
	} else {
		if b.blackShortCastling && (b.occupied&BlackShortEmptyCells) == 0 && b.isPathSafe(BlackShortSafeCells, Black) {
			res |= G8
		}
		if b.blackLongCastling && (b.occupied&BlackLongEmptyCells) == 0 && b.isPathSafe(BlackLongSafeCells, Black) {
			res |= C8
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

		if curr & b.occupied != 0 { break }
	}
	return res
}
