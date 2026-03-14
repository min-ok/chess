package game

import (
	"math"
	"image/color"
	"math/bits"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)




func (g *Game) drawTextScreen(t string) {
	g.screenState.Clear()

	op := &text.DrawOptions{}

	face := &text.GoTextFace{
		Source: g.FaceSource,
		Size: 24 * float64(g.scale),
	}

	w, h := text.Measure(t, face, op.LineSpacing)
	x := (float64(g.boardSizeX) - w) / 2
	y := (float64(g.boardSizeY) - h) / 2

	op.GeoM.Translate(math.Floor(x), math.Floor(y))
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(g.screenState, t, face, op)
}


func (g *Game) drawFigures() {
	figures := g.board.GetFigures()

	op := &ebiten.DrawImageOptions{}
	for teamInx := 0; teamInx < len(figures); teamInx += 1 {
		for pieceInx := 0; pieceInx < len(figures[teamInx]); pieceInx += 1 {
			f := figures[teamInx][pieceInx]

			for f != 0 {
				idx := bits.TrailingZeros64(f)
				i, j := idx % 8, 7 - idx / 8

				op.GeoM.Reset()
				op.GeoM.Translate(float64(i * g.cellSizeX + g.boardMargin), float64(j * g.cellSizeY + g.boardMargin))

				g.screenImg.DrawImage(g.pieceImages[teamInx][pieceInx], op)

				f &= (f - 1)
			}
		}
	}
}


func (g *Game) drawPointers() {
	pm := g.possibleMoves

	enemy := g.board.GetPlayerEnemyFigures(g.pointer)

	op := &ebiten.DrawImageOptions{}
	for pm != 0 {
		idx := bits.TrailingZeros64(pm)
		i, j := idx % 8, 7 - idx / 8

		op.GeoM.Reset()
		op.ColorScale.Reset()

		op.GeoM.Translate(float64(i * g.cellSizeX + g.boardMargin), float64(j * g.cellSizeY + g.boardMargin))
		if enemy & (uint64(1) << idx) != 0 { op.ColorScale.Scale(0, 1, 0, 0.5) }
		g.screenImg.DrawImage(g.pointerImage, op)

		pm &= pm - 1
	}
}


func (g *Game) updateScreen() {
	g.screenImg.Clear()

	g.screenImg.DrawImage(g.boardImage, nil)


	g.drawFigures()
	g.drawPointers()

	g.screenImg.DrawImage(g.screenState, nil)
}
