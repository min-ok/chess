package logic

import (
	"math"
	"math/bits"
)

const MaxDepth = 64
const MaxMovesPerPosition = 218

type MoveList struct {
	moves [MaxDepth][MaxMovesPerPosition]Move
	sizes [MaxDepth]int
}


func (b *Board)getBestMove(depth int, team int) Move {
	ml := &MoveList{}

	b.getMoves(ml, depth, team)

	alpha := math.MinInt + 1
	beta := math.MaxInt

	var res Move
	best := math.MinInt + 1

	for i := 0; i < ml.sizes[depth]; i++ {
		m := ml.moves[depth][i]

		b.move(m.from, m.to, m.movingPiece, m.eatenPiece, m.team)

		if b.isChecked(b.bitboard[team][King], team) {
			b.undo()
			continue
		}

		score := -b.alphaBeta(ml, depth - 1, -beta, -alpha, getOppositeTeam(team))
		b.undo()

		if score > best {
			best = score
			res = m
			if score > alpha {
				alpha = score
			}
		}
	}

	return res
}


func (b *Board) alphaBeta(ml *MoveList, depth int, alpha, beta int, team int) int {
	if depth == 0 {
		return b.evaluate()
	}

	b.getMoves(ml, depth, team)

	if ml.sizes[depth] == 0 {
		if b.isChecked(b.bitboard[team][King], team) {
			return -1000000
		}
		return 0
	}

	for i := 0; i < ml.sizes[depth]; i++ {
		m := ml.moves[depth][i]

		b.move(m.from, m.to, m.movingPiece, m.eatenPiece, m.team)

		if b.isChecked(b.bitboard[team][King], team) {
			b.undo()
			continue
		}

		score := -b.alphaBeta(ml, depth - 1, -beta, -alpha, getOppositeTeam(team))
		b.undo()

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
		}
	}

	return alpha
}


func (b *Board) getMoves(ml *MoveList, depth int, team int) {
	ml.sizes[depth] = 0
	bb := b.bitboard[team]

	for i := 0; i < 6; i++ {
		fromBb := bb[i]

		for fromBb != 0 {
			from := fromBb & -fromBb

			toBb := b.getPossibleMoves(from, pieceTypes[i], team)

			for toBb != 0 {
				to := toBb & -toBb

				m := Move{
					from: from, to: to,
					team: team,
					movingPiece: pieceTypes[i],
					eatenPiece: b.getPieceType(to),

					oldFlags: b.flags,
				}

				ml.moves[depth][ml.sizes[depth]] = m
				ml.sizes[depth]++

				toBb &= toBb - 1
			}

			fromBb &= fromBb - 1
		}
	}

}


func (b *Board) evaluate() int {
	sum := 0

	for i := 0; i < 6; i += 1 {
		sum += ( bits.OnesCount64(b.bitboard[White][i]) - bits.OnesCount64(b.bitboard[Black][i]) ) * pieceCost[i]
	}


	if b.turn == White {
		return -sum
	} else {
		return sum
	}
}
