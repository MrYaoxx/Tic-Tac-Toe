package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
)

type Game struct {
}

var (
	//go:embed icon.png
	IconBytes []byte

	//go:embed font.ttf
	FontBytes []byte
	Font      font.Face
	SmallFont font.Face

	Width      = 320
	Height     = 320
	GridSpace  = 10
	GridWidth  = (Width - GridSpace*4) / 3
	GridHeight = (Height - GridSpace*4) / 3

	Grid          = [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	CurrentPlayer = 1
	Winner        = 0

	BGColor   = color.RGBA{139, 87, 42, 255}
	GridColor = []color.RGBA{{12, 142, 223, 255}, {223, 12, 16, 255}, {223, 211, 12, 255}}
)

func Check() int {
	for i := 0; i < 3; i++ {
		if Grid[i][0] > 0 && Grid[i][0] == Grid[i][1] && Grid[i][1] == Grid[i][2] {
			return Grid[i][0]
		}
	}
	for j := 0; j < 3; j++ {
		if Grid[0][j] > 0 && Grid[0][j] == Grid[1][j] && Grid[1][j] == Grid[2][j] {
			return Grid[0][j]
		}
	}
	if (Grid[0][0] == Grid[1][1] && Grid[1][1] == Grid[2][2]) ||
		(Grid[0][2] == Grid[1][1] && Grid[1][1] == Grid[2][0]) {
		return Grid[1][1]
	}
	return 0
}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	x, y := my/(Width/3), mx/(Width/3)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && Winner == 0 {
		if Grid[x][y] == 0 {
			Grid[x][y] = CurrentPlayer
			Winner = Check()

			if CurrentPlayer == 1 {
				CurrentPlayer = 2
			} else {
				CurrentPlayer = 1
			}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		Winner = 0
		Grid = [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		CurrentPlayer = 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if Winner > 0 {
		screen.Fill(GridColor[Winner])
		text.Draw(screen, "win", Font, 70, 170, BGColor)
		text.Draw(screen, "Right click to restart", SmallFont, 15, 250, BGColor)
	} else {
		screen.Fill(BGColor)
		for i := 0; i < 3; i++ {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(GridSpace+1), float64(i*(GridHeight+GridSpace)+GridSpace+1))
			for j := 0; j < 3; j++ {
				grid := ebiten.NewImage(GridWidth, GridHeight)
				grid.Fill(GridColor[Grid[i][j]])
				screen.DrawImage(ebiten.NewImageFromImage(grid), opt)
				opt.GeoM.Translate(float64(GridWidth+GridSpace), 0)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

func init() {
	tt, _ := opentype.Parse(FontBytes)
	const dpi = 72
	Font, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    120,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	SmallFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func main() {
	icon, _, _ := image.Decode(bytes.NewReader(IconBytes))
	ebiten.SetWindowIcon([]image.Image{icon})
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
