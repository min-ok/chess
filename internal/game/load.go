package game

import (
	"fmt"
	"os"
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)


func loadSprite(path string, g *Game) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}

	res := ebiten.NewImage(img.Bounds().Dx() * g.scale, img.Bounds().Dy() * g.scale)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.scale), float64(g.scale))

	res.DrawImage(img, op)
	return res
}


func loadImages(g *Game) {
	g.boardImage = loadSprite("internal/assets/textures/boards/board_1.png", g)
	g.pointerImage = loadSprite("internal/assets/textures/pointer.png", g)

	for teamNum, pieceTeam := range [2]string{"white", "black"} {
		for typeNum, pieceType := range [6]string{"pawn", "bishop", "knight", "rook", "queen", "king"} {
			g.pieceImages[teamNum][typeNum] = loadSprite(fmt.Sprintf("internal/assets/textures/pieces/%s_%s.png", pieceType, pieceTeam), g)
		}
	}
}


func loadFont(g *Game)  {
	fontData, err := os.ReadFile("internal/assets/Jersey10-Regular.ttf")
	if err != nil {
		panic("cant load font")
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic("cant load font")
	}

	g.FaceSource = s
}
