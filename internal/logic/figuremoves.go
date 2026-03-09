package logic


func (b *Board) getPawnSquaresAttacked(p uint64) uint64 {
	var moves uint64

	t := b.getPieceTeam(p)
	if t == -1 { panic(t) }

	if (t == White) {
		moves |= (p << 7) & b.BlackFigures.all & notH
	 	moves |= (p << 9) & b.BlackFigures.all & notA

		step := p << 8
		if step & b.occupied == 0 {
			moves |= step
			if p & rank2 != 0 {
				step = p << 16
				if step & b.occupied == 0 { moves |= step }
			}
		}
	} else {
		moves |= (p >> 7) & b.WhiteFigures.all & notA
	 	moves |= (p >> 9) & b.WhiteFigures.all & notH

		step := p >> 8
		if step & b.occupied == 0 {
			moves |= step

			if p & rank7 != 0 {
				step = p >> 16
				if step & b.occupied == 0 { moves |= step }
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
