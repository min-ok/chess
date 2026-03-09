package game

import (
	"fmt"

	"chess/internal/logic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)


func (g *Game) Input() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if xx, yy := x - g.boardMargin * 3, y - g.boardMargin * 3; xx >= g.cellSizeX * 3 * 8 || xx < 0 || yy >= g.cellSizeY * 3 * 8 || yy < 0 {
			return
		}

		tx := (x - g.boardMargin * 3) / g.cellSizeX / 3
		ty := (y - g.boardMargin * 3) / g.cellSizeY / 3

		p := uint64(1) << ((7 - ty) * 8 + tx)

		if g.possibleMoves & p != 0 {
			s, t := g.board.MakePlayerMove(g.pointer, p)

			switch s {
				case logic.Check: g.drawTextScreen(fmt.Sprintf("Check to %d", t))
				case logic.Checkmate: g.drawTextScreen(fmt.Sprintf("Checkmate to %d", t))
				case logic.Draw: g.drawTextScreen("Draw")
			}

			g.resetPointer()
			g.resetPossibleMoves()

		} else {
			g.pointer = p

			g.possibleMoves = g.board.GetLegalMoves(g.pointer)
		}
	}
}
