package logic


func (b *Board) GetPlayerEnemyFigures(p uint64) uint64 {
	if p & b.allOccupied == 0 {
		return 0
	}


	return b.teamOccupied[getOppositeTeam(b.getPieceTeam(p))]
}


func (b *Board) MakePlayerMove(from, to uint64) (int, int) {
	b.move(from, to, b.getPieceType(from), b.getPieceType(to), b.turn)
	b.turn = getOppositeTeam(b.turn)

	// panic(fmt.Sprintf("%d %d", b.turn, b.turn))

	// bestMove := b.GetBestMove(botDepth, b.turn)
	// b.move(bestMove.from, bestMove.to, bestMove.movingPiece, bestMove.eatenPiece, bestMove.team)
	// b.turn = getOppositeTeam(b.turn)

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
	if p & b.allOccupied == 0 {
		return 0
	}

	team := b.getPieceTeam(p)

	if team != b.turn {
		return 0
	}

	return b.getLegalMoves(p, b.getPieceType(p), team)
}
