package game

import (
	"fmt"

	"chess/internal/logic"
	// "math/bits"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)


func intToTeamString(team int) string {
	if team == logic.White {
		return "White"
	}

	if team == logic.Black {
		return "Black"
	}

	return ""
}







func (g *Game) Input() {
	if g.gameended {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if xx, yy := x - g.boardMargin, y - g.boardMargin; xx >= g.cellSizeX * 8 || xx < 0 || yy >= g.cellSizeY * 8 || yy < 0 {
			return
		}

		tx := (x - g.boardMargin) / g.cellSizeX
		ty := (y - g.boardMargin) / g.cellSizeY

		p := uint64(1) << ((7 - ty) * 8 + tx)

		if g.possibleMoves & p != 0 {
			s, t := g.board.MakePlayerMove(g.pointer, p)

			switch s {
				case logic.None: g.drawTextScreen("")
				case logic.Check: g.drawTextScreen(fmt.Sprintf("Check to %s", intToTeamString(t)))
				case logic.Checkmate:
				g.drawTextScreen(fmt.Sprintf("Checkmate to %s", intToTeamString(t)))
				g.gameended = true
				case logic.Draw:
				g.drawTextScreen("Draw")
				g.gameended = true
			}

			// bestMove := g.board.GetBestMove(7, false)

			// g.board.Move(bestMove)

			g.resetPointer()
			g.resetPossibleMoves()

		} else {
			g.pointer = p

			g.possibleMoves = g.board.GetLegalMoves(g.pointer, g.board.GetPieceTeam(g.pointer), g.board.GetPieceType(g.pointer))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.board.MakePlayerUndo()
	}
}
