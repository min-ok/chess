package logic


func (b *Board) GetPlayerEnemyFigures(p uint64) uint64 {
	if p & b.occupied == 0 {
		return 0
	}
	return b.getBitboard(getOppositeTeam(b.GetPieceTeam(p))).all
}


func (b *Board) MakePlayerMove(from, to uint64) (int, int) {
	b.move(from, to, b.turn, b.GetPieceType(from))
	b.turn = getOppositeTeam(b.turn)

	bestMove := b.GetBestMove(5, b.turn)
	b.move(bestMove.from, bestMove.to, bestMove.team, bestMove.movingPiece)
	b.turn = getOppositeTeam(b.turn)

	return b.checkGameStatus(b.turn), b.turn
}


func (b *Board) MakePlayerUndo() {
	if b.historyLen == 0 {
		return
	}

	b.undo()
	b.turn = getOppositeTeam(b.turn)
}


func (b *Board) GetPlayerLegalMoves(p uint64) uint64 {
	if p & b.occupied == 0 {
		return 0
	}

	team := b.GetPieceTeam(p)

	if team != b.turn {
		return 0
	}

	return b.getLegalMoves(p, team, b.GetPieceType(p))
}
