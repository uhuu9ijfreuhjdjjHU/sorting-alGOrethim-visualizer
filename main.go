package main

import (
	"log"
	"os"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 44100

type Game struct {
  data []int
  i, j int
  sorted bool

  audioContext *audio.Context
  soundFile    string
}

var (
	gameSpeed = int(240)
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

func NewGame() *Game {
  ctx := audio.NewContext(sampleRate)

  return &Game{
  	audioContext: ctx,
  	soundFile:    "assets/boop.mp3",
  }
}

func (g *Game) Update() error { //game logic
	
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		gameSpeed = 3000
	}
	
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		gameSpeed = 120
	}

	if g.data == nil { // Initialize once
		size := randInt(340, 340)
		g.data = make([]int, size)
		for i := range g.data {
			g.data[i] = randInt(1, 220)
		}
		return nil
	}

	if g.sorted {
		return nil
	}

	if g.data[g.j] < g.data[g.j+1] { //one comparison per tick
    g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]	
    // Create a new player for overlapping sound
    f, err := os.Open(g.soundFile)
    if err != nil {
      log.Println("failed to open sound:", err)
    } else {
      stream, err := mp3.DecodeWithoutResampling(f)
      if err != nil {
        log.Println("failed to decode sound:", err)
      } else {
        player, err := g.audioContext.NewPlayer(stream)
        if err != nil {
          log.Println("failed to create player:", err)
        } else {
          player.Play() // independent player, overlaps allowed
        }
      }
  	}
	}

	g.j++ // Advance inner index

	if g.j >= len(g.data)-g.i-1 {	// End of inner pass
		g.j = 0
		g.i++
	}

	if g.i >= len(g.data)-1 { // Fully sorted
		g.sorted = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	
	visualizerPosition = 0

	for i := range g.data {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.01, float64(g.data[i])*0.01)
		op.GeoM.Translate(float64(100*(visualizerPosition*0.01)), 0)

		// Tint bars green if sorted
		if g.sorted {
			op.ColorM.Scale(0, 1, 0, 1)
		}

		screen.DrawImage(visualizerBar, op)
		visualizerPosition++
	}
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetTPS(gameSpeed)
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Hello, World!")

	game := NewGame() // <-- use the constructor
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
