package logic

import (
	// "sort"
	"math"
)


const MaxDepth = 64
const MaxMovesPerPosition = 218 // максимум возможных ходов в шахматах

type MoveList struct {
	moves [MaxDepth][MaxMovesPerPosition]Move
	sizes [MaxDepth]int
}


func (b *Board) alphaBeta(ml *MoveList, depth int, alpha, beta int, team int) int {
	if depth == 0 {
		return b.Evaluate()
	}

	b.getMoves(ml, depth, team)

	if ml.sizes[depth] == 0 {
		if b.isChecked(b.getKing(team), team) {
			return -1000000
		}
		return 0
	}

	for i := 0; i < ml.sizes[depth]; i++ {
		m := ml.moves[depth][i]

		b.move(m.from, m.to, m.team, m.movingPiece)

		if b.isChecked(b.getKing(team), team) {
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


func (b *Board)GetBestMove(depth int, team int) Move {
	ml := &MoveList{}

	b.getMoves(ml, depth, team)

	alpha := math.MinInt + 1
	beta := math.MaxInt

	var res Move
	best := math.MinInt + 1

	for i := 0; i < ml.sizes[depth]; i++ {
		m := ml.moves[depth][i]

		b.move(m.from, m.to, m.team, m.movingPiece)

		if b.isChecked(b.getKing(team), team) {
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


func (b *Board)IsTerminated() bool {
	 s := b.checkGameStatus(b.turn)

	return s == Checkmate || s == Draw
}


func (b *Board) getMoves(ml *MoveList, depth int, team int) {
	ml.sizes[depth] = 0
	bb := b.getBitboard(team)

	fromPieceBoards := [6]uint64{bb.pawns, bb.bishops, bb.knights, bb.rooks, bb.queens, bb.king}
	pieceTypes := [6]int{Pawn, Bishop, Knight, Rook, Queen, King}

	for i := 0; i < 6; i++ {
		fromBb := fromPieceBoards[i]

		for fromBb != 0 {
			from := fromBb & -fromBb

			toBb := b.getPossibleMoves(from, team, pieceTypes[i])

			for toBb != 0 {
				to := toBb & -toBb

				m := Move{
					from: from, to: to,
					team: team,
					movingPiece: pieceTypes[i],
					eatenPiece: b.GetPieceType(to),

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
