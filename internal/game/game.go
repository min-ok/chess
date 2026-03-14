package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"chess/internal/logic"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)


type Game struct {

	pieceImages [2][6]*ebiten.Image
	boardImage *ebiten.Image

	screenImg *ebiten.Image
	screenState *ebiten.Image

	FaceSource *text.GoTextFaceSource

	pointerImage *ebiten.Image
	pointer uint64

	board *logic.Board

	possibleMoves uint64

	boardMargin int
	boardSizeX, boardSizeY int
	cellSizeX, cellSizeY int

	scale int

	gameended bool
}



func (g *Game) resetPointer() {
	g.pointer = 0
}


func (g *Game) resetPossibleMoves() {
	g.possibleMoves = 0
}


func (g *Game) Update() error {
	g.Input()

	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	g.updateScreen()
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(2, 2)

	screen.DrawImage(g.screenImg, nil)
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.boardSizeX, g.boardSizeY
}


func setDimensions(g *Game) {
	bs := g.boardImage.Bounds().Size()

	g.boardSizeX = bs.X
	g.boardSizeY = bs.Y

	g.cellSizeX = ( g.boardSizeX - g.boardMargin * 2 ) / 8
	g.cellSizeY = ( g.boardSizeY - g.boardMargin  * 2 ) / 8
}


func Start() {
	g := &Game{}
	ebiten.SetWindowSize(176 * 5, 176 * 5)
	ebiten.SetWindowTitle("Chess")
	g.scale = 5

	g.boardMargin = 8 * g.scale


	loadImages(g)
	loadFont(g)
	setDimensions(g)

	g.board = logic.CreateBoard()
	g.screenImg = ebiten.NewImage(g.boardSizeX, g.boardSizeY)
	g.screenState = ebiten.NewImage(g.boardSizeX, g.boardSizeY)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
