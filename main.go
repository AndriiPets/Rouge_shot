package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var gameGlobal *Game

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

func NewSprite(x, y float64, id EntityID) *Sprite {
	s := &Sprite{
		X: x,
		Y: y,
	}

	if gameGlobal == nil {
		log.Fatal("Global game object is not initialized")
	}

	gameGlobal.sprites[id] = s

	return s
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(s.X, s.Y)

	screen.DrawImage(
		s.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

}

type Game struct {
	player    *Player
	enemies   map[EntityID]*Enemy
	sprites   map[EntityID]*Sprite
	colliders map[EntityID]*Collider

	//camera stuff
	camera Camera
	world  *ebiten.Image

	//grid stuff
	grid SparseGrid

	//map gen stuff
	Level      MapManager
	PathFinder PathFinder
}

func NewGame() *Game {
	game := &Game{}
	game.enemies = make(map[EntityID]*Enemy)
	game.sprites = make(map[EntityID]*Sprite)
	game.colliders = make(map[EntityID]*Collider)

	game.grid = make(SparseGrid)

	game.camera = Camera{
		ViewPort:   Vec2{VIRTUAL_WIDTH, VIRTUAL_HEIGHT},
		ZoomFactor: 48,
		Position:   Vec2{10 * CELL_SIZE, 5 * CELL_SIZE},
	}
	game.world = ebiten.NewImage(GRID_WIDTH*CELL_SIZE, GRID_HEIGHT*CELL_SIZE)

	game.Level = NewMapManager(GRID_WIDTH, GRID_HEIGHT, "assets/images/Tileset.png")
	game.Level.GenerateMap()

	game.PathFinder = NewPathFinder(GRID_WIDTH, GRID_HEIGHT, CELL_SIZE)
	game.PathFinder.GenerateLayout(game.Level.Map.Data, 'x')

	return game
}

func (g *Game) GameTick() {
	for _, enemy := range g.enemies {
		if enemy != nil {
			enemy.Update()
		}
	}
}

func (g *Game) ParceMap() {
	for y, row := range g.Level.Map.Data {
		for x, val := range row {

			posX, posY := (x * CELL_SIZE), (y * CELL_SIZE)

			if val == 'x' {
				AddToGrid(NewStaticObstacle(float64(posX), float64(posY)), posX, posY)
			}
			if val == 'p' {
				if g.player != nil {
					g.player.Move(posX, posY)
				}
			}
			if val == 'e' {
				NewEnemy(float64(posX), float64(posY))
			}

		}
	}
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	g.player.Update()

	playerPos := Vec2{g.player.sprite.X, g.player.sprite.Y}
	g.camera.Update(playerPos)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.Fill(color.RGBA{120, 180, 255, 255})
	g.world.Clear()

	//draw map
	g.Level.DrawTiles(g.world)
	//draw grid
	for x := 0; x < GRID_WIDTH*CELL_SIZE; x += CELL_SIZE {
		vector.StrokeLine(g.world, float32(x), 0, float32(x), GRID_HEIGHT*CELL_SIZE, 1, color.RGBA{255, 255, 255, 155}, false)
	}

	for y := 0; y < GRID_HEIGHT*CELL_SIZE; y += CELL_SIZE {
		vector.StrokeLine(g.world, 0, float32(y), GRID_WIDTH*CELL_SIZE, float32(y), 1, color.RGBA{255, 255, 255, 155}, false)
	}

	g.player.Draw(g.world)

	for _, sprite := range g.sprites {
		if sprite != nil {
			sprite.Draw(g.world)
		}
	}

	for _, coll := range g.colliders {
		if coll != nil {
			vector.StrokeRect(
				g.world,
				float32(coll.BB.Min.X),
				float32(coll.BB.Min.Y),
				float32(coll.BB.Dx()),
				float32(coll.BB.Dy()),
				1,
				color.RGBA{0, 255, 0, 255},
				false,
			)
		}
	}

	g.camera.Render(g.world, screen)

	ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'f', 1, 64))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return VIRTUAL_WIDTH, VIRTUAL_HEIGHT
}

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Rouge!")
	game := NewGame()
	gameGlobal = game
	game.player = NewPlayer()
	game.ParceMap()
	//NewEnemy(game.player.sprite.X+CELL_SIZE, game.player.sprite.Y-CELL_SIZE)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
