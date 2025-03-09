package main

import (
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/dngn"
)

type GenerationMode int

const (
	BSP GenerationMode = iota
	DrunkWalk
	RandomRooms
)

type MapManager struct {
	Map            *dngn.Layout
	Tileset        *ebiten.Image
	GenerationMode GenerationMode
}

func NewMapManager(w, h int, tilesetPath string) MapManager {
	mapM := MapManager{}
	mapM.Map = dngn.NewLayout(w, h)

	mapM.Tileset = LoadImageFromFile(tilesetPath)
	mapM.GenerationMode = BSP

	return mapM
}

func LoadImageFromFile(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return img

}

func (m *MapManager) GenerateMap() {

	mapSelection := m.Map.Select()

	switch m.GenerationMode {
	case BSP:

		bspOptions := dngn.NewDefaultBSPOptions()
		bspOptions.SplitCount = 40
		bspOptions.MinimumRoomSize = 5

		bspRooms := m.Map.GenerateBSP(bspOptions)

		start := bspRooms[0]

		for _, subroom := range bspRooms {

			subroomCenter := subroom.Center()
			center := m.Map.Center()

			margin := 10

			if subroomCenter.X > center.X-margin &&
				subroomCenter.X < center.X+margin &&
				subroomCenter.Y > center.Y-margin &&
				subroomCenter.Y < center.Y+margin {
				start = subroom
				break
			}
		}

		for _, room := range bspRooms {
			hoops := room.CountHopsTo(start)

			if hoops < 0 || hoops > 4 {
				mapSelection.FilterByArea(room.X, room.Y, room.W+1, room.H+1).Fill('x')
				room.Disconnect()
			}
		}

		player_pos := start.Center()
		m.Map.Set(player_pos.X, player_pos.Y, 'p')

	case DrunkWalk:

		m.Map.GenerateDrunkWalk(' ', 'x', 0.5)

	case RandomRooms:

		m.Map.GenerateRandomRooms(' ', 'x', 6, 3, 3, 5, 5, true)

		mapSelection.FilterByRune(' ').FilterBy(func(x, y int) bool {
			return (m.Map.Get(x+1, y) == 'x' &&
				m.Map.Get(x-1, y) == 'x') ||
				(m.Map.Get(x, y-1) == 'x' && m.Map.Get(x, y+1) == 'x')
		}).FilterByPercentage(0.25).Fill('#')
	}

	// Fill the outer walls
	mapSelection.Remove(mapSelection.FilterByArea(1, 1, m.Map.Width-2, m.Map.Height-2)).Fill('x')

	// Add a different tile for an alternate floor
	mapSelection.FilterByRune(' ').FilterByPercentage(0.1).Fill('.')
	mapSelection.FilterByRune(' ').FilterByPercentage(0.005).Fill('e')
}

func (m *MapManager) ToString() string {
	return m.Map.DataToString()
}

func (m *MapManager) SetGenerationMode(genMode GenerationMode) {
	m.GenerationMode = genMode
}

func (m *MapManager) DrawTiles(screen *ebiten.Image) {

	roomSelect := m.Map.Select()

	for cell := range roomSelect.Cells {

		v := m.Map.Get(cell.X, cell.Y)

		left := m.Map.Get(cell.X-1, cell.Y) == v
		right := m.Map.Get(cell.X+1, cell.Y) == v
		up := m.Map.Get(cell.X, cell.Y-1) == v
		down := m.Map.Get(cell.X, cell.Y+1) == v
		rotation := 0.0

		// Tile graphic defaults to plain ground
		srcOffsetX := 0
		srcOffsetY := 16

		if v == ' ' || v == '#' {
			if m.Map.Get(cell.X, cell.Y-1) == 'x' {
				srcOffsetY = 0
			}
		}

		// Minor scratches on the ground
		if v == '.' {
			if m.Map.Get(cell.X, cell.Y-1) == 'x' {
				srcOffsetY = 0
			} else {
				srcOffsetY = 32
			}
		}

		// Wall
		if v == 'x' {

			//add to collide map maybe move to its own function later

			num := 0
			if left {
				num++
			}
			if right {
				num++
			}
			if up {
				num++
			}
			if down {
				num++
			}

			if num == 0 {

				srcOffsetX = 48
				srcOffsetY = 16

			} else if num == 1 {

				srcOffsetX = 48
				srcOffsetY = 32

				if right {
					rotation = math.Pi
				} else if up {
					rotation = math.Pi / 2
				} else if down {
					rotation = -math.Pi / 2
				}

			} else if num == 2 {

				if left && right {
					srcOffsetX = 32
					srcOffsetY = 32
				} else if up && down {
					srcOffsetX = 32
					srcOffsetY = 32
					rotation = math.Pi / 2
				} else {

					srcOffsetX = 48
					srcOffsetY = 0

					if up && right {
						rotation = math.Pi / 2
					} else if right && down {
						rotation = math.Pi
					} else if down && left {
						rotation = -math.Pi / 2
					}

				}

			} else if num == 3 {
				srcOffsetX = 32
				srcOffsetY = 0

				if up && right && down {
					rotation = math.Pi / 2
				} else if right && down && left {
					rotation = math.Pi
				} else if down && left && up {
					rotation = -math.Pi / 2
				}

			} else if num == 4 {
				srcOffsetX = 32
				srcOffsetY = 16
			}

		}

		src := image.Rect(0, 0, 16, 16)
		src = src.Add(image.Point{srcOffsetX, srcOffsetY})

		tile := m.Tileset.SubImage(src).(*ebiten.Image)
		geoM := ebiten.GeoM{}
		geoM.Translate(-float64(src.Dx()/2), -float64(src.Dy()/2))

		geoM.Rotate(rotation)

		geoM.Translate(float64(src.Dx()/2), float64(src.Dy()/2))

		geoM.Translate(float64(cell.X*src.Dx()), float64(cell.Y*src.Dy()))
		screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geoM})

	}

	doors := roomSelect.FilterByRune('#')

	for d := range doors.Cells {

		src := image.Rect(16, 0, 32, 16)

		dstX, dstY := float64(d.X*src.Dx()), float64(d.Y*src.Dy())

		if m.Map.Get(d.X-1, d.Y) != 'x' && m.Map.Get(d.X+1, d.Y) != 'x' { // Horizontal door
			src = src.Add(image.Point{0, 16})
		} else {
			dstY += 4
		}

		tile := m.Tileset.SubImage(src).(*ebiten.Image)
		geoM := ebiten.GeoM{}
		geoM.Translate(dstX, dstY)
		screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geoM})

	}
}
