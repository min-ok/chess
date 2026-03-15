package logic

import (
	"math"
	"math/bits"
)


func (ml *MoveList) sortMoves(depth int) {
    size := ml.sizes[depth]
    for i := 1; i < size; i++ {
        key := ml.moves[depth][i]
        keyScore := moveScore(key)
        j := i - 1
        for j >= 0 && moveScore(ml.moves[depth][j]) < keyScore {
            ml.moves[depth][j + 1] = ml.moves[depth][j]
            j--
        }
        ml.moves[depth][j+1] = key
    }
}

var pieceValue = [7]int{0, 100, 330, 320, 500, 900, 20000}

func moveScore(m Move) int {
    if m.eatenPiece == Empty {
        return 0
    }

    return pieceValue[m.eatenPiece] - pieceValue[m.movingPiece]
}


const MaxDepth = 64
const MaxMovesPerPosition = 218

type MoveList struct {
	moves [MaxDepth][MaxMovesPerPosition]Move
	sizes [MaxDepth]int
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

		b.move(m.from, m.to, m.movingPiece, m.eatenPiece, m.team)

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


func (b *Board) alphaBeta(ml *MoveList, depth int, alpha, beta int, team int) int {
	if depth == 0 {
		return b.evaluate()
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

		b.move(m.from, m.to, m.movingPiece, m.eatenPiece, m.team)

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


func (b *Board) getMoves(ml *MoveList, depth int, team int) {
	ml.sizes[depth] = 0
	bb := b.getBitboard(team)

	fromPieceBoards := [6]uint64{bb.pawns, bb.bishops, bb.knights, bb.rooks, bb.queens, bb.king}
	pieceTypes := [6]int{Pawn, Bishop, Knight, Rook, Queen, King}

	for i := 0; i < 6; i++ {
		fromBb := fromPieceBoards[i]

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

	ml.sortMoves(depth)
}


func (b *Board) evaluate() int {
	sum := ( bits.OnesCount64(b.whiteFigures.pawns) - bits.OnesCount64(b.blackFigures.pawns) ) * pawnCost
	sum += ( bits.OnesCount64(b.whiteFigures.bishops) - bits.OnesCount64(b.blackFigures.bishops) ) * bishopCost
	sum += ( bits.OnesCount64(b.whiteFigures.knights) - bits.OnesCount64(b.blackFigures.knights) ) * knightCost
	sum += ( bits.OnesCount64(b.whiteFigures.rooks) - bits.OnesCount64(b.blackFigures.rooks) ) * rookCost
	sum += ( bits.OnesCount64(b.whiteFigures.queens) - bits.OnesCount64(b.blackFigures.queens) ) * queenCost
	sum += ( bits.OnesCount64(b.whiteFigures.king) - bits.OnesCount64(b.blackFigures.king) ) * kingCost

	return sum
}
