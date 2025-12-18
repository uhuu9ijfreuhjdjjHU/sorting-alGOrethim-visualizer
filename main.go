package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	data []int

	i int // outer loop index
	j int // inner loop index

	sorted bool
}


var (
	visualizerBar *ebiten.Image
	visualizerPosition = float64(0)
	tickCount = int(0)
)

type barStats struct {
	height int
}

func init() {
	var err error

	visualizerBar, _, err = ebitenutil.NewImageFromFile("assets/Sprite-0001.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error { //game logic
	// Initialize once
	if g.data == nil {
		size := randInt(340, 340)
		g.data = make([]int, size)
		for i := range g.data {
			g.data[i] = randInt(1, 340)
		}
		return nil
	}

	if g.sorted {
		return nil
	}

	// ---- ONE comparison per tick ----
	if g.data[g.j] > g.data[g.j+1] {
		g.data[g.j], g.data[g.j+1] =
			g.data[g.j+1], g.data[g.j]
	}

	// Advance inner index
	g.j++

	// End of inner pass
	if g.j >= len(g.data)-g.i-1 {
		g.j = 0
		g.i++
	}

	// Fully sorted
	if g.i >= len(g.data)-1 {
		g.sorted = true
	}

	return nil
}



func (g *Game) Draw(screen *ebiten.Image) {
	
	visualizerPosition = 0
	
	for i := range g.data {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Scale(0.01, float64(g.data[i])*0.01)
		op.GeoM.Translate(float64(100 * (visualizerPosition * 0.01)), 0)

		screen.DrawImage(visualizerBar, op)

		visualizerPosition++
	}
}



func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
