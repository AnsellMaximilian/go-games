package tictactoe

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	width  int = 600
	height int = 600
	padding int = 10
	tileSize int = 75
	outerPadding int = (width - (tileSize * 3 + padding * 2)) / 2
	
)

var (
	mplusNormalFont font.Face
	fontSize = 24
)

type Game struct {
	grid [9]string
	turn string
	winner string
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				tileY := (row*tileSize + (padding * row)) + outerPadding
				tileX := (col*tileSize + (padding * col)) + outerPadding

				isMouseOver := mouseX >= tileX && mouseX <= tileX+tileSize && mouseY >= tileY && mouseY <= tileY+tileSize

				if isMouseOver {
					index := col + (row * 3)
					// fmt.Printf("(%v, %v); index; %v\n", col, row, index)
					if g.grid[index] == "" {
						g.grid[index] = g.turn
						g.checkWin()
						if g.turn == "X" {
							g.turn = "O"
						}else {
							g.turn = "X"
						}
						
					}
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.White)

	if g.winner != "" {
		drawInfo(fmt.Sprintf("%v Won!", g.winner), screen)
	}else {
		drawInfo(fmt.Sprintf("%v's Turn", g.turn), screen)
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// draw board
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			tile := ebiten.NewImage(tileSize, tileSize)
			tile.Fill(color.Black)

			tileY := (row*tileSize + (padding * row)) + outerPadding
			tileX := (col*tileSize + (padding * col)) + outerPadding

			isMouseOver := mouseX >= tileX && mouseX <= tileX+tileSize && mouseY >= tileY && mouseY <= tileY+tileSize

			if isMouseOver {
				tile.Fill(color.White)
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tileX), float64(tileY))
			screen.DrawImage(tile, op)
			// inner tile
			innerTile := ebiten.NewImage(tileSize-2, tileSize-2)
			innerTile.Fill(color.Black)
			innerOp := &ebiten.DrawImageOptions{}
			innerOp.GeoM.Translate(float64(tileX)+1, float64(tileY)+1)
			screen.DrawImage(innerTile, innerOp)

			// draw x's and o's
			index := col + (row * 3)

			if g.grid[index] != "" {
				text.Draw(screen, g.grid[index], mplusNormalFont, tileX + 27, tileY + 50, color.White)
			}

		}
	}


}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func RunGame() {
	g := &Game{
		turn: "X",
	}
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Tic-Tac-Toe")
	

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) checkWin() bool {
	var winner string
	// Winning combinations
	winningCombinations := [8][3]int{
		{0, 1, 2}, // First row
		{3, 4, 5}, // Second row
		{6, 7, 8}, // Third row
		{0, 3, 6}, // First column
		{1, 4, 7}, // Second column
		{2, 5, 8}, // Third column
		{0, 4, 8}, // Diagonal from top-left
		{2, 4, 6}, // Diagonal from top-right
	}

	// Check each winning combination
	for _, combo := range winningCombinations {
		if g.grid[combo[0]] != "" && g.grid[combo[0]] == g.grid[combo[1]] && g.grid[combo[1]] == g.grid[combo[2]] {
			winner = g.grid[combo[0]]
		}
	}

	g.winner = winner
	return winner != ""

}

func drawInfo(info string, screen *ebiten.Image) {	
	outerPadding := (width - (len(info) * fontSize)) / 2
	for i := 0; i < len(info); i++ {
		char := string([]rune(info)[i])
		text.Draw(screen, char, mplusNormalFont, i*fontSize + outerPadding, 100, color.Black)
	}

	
}