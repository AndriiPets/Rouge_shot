package main

const (
	SCREEN_WIDTH   = 640
	SCREEN_HEIGHT  = 480
	VIRTUAL_WIDTH  = 320
	VIRTUAL_HEIGHT = 240

	//Grid
	CELL_SIZE   = 16
	GRID_HEIGHT = 40
	GRID_WIDTH  = 60
)

type Vec2 struct {
	X, Y float64
}

type SparseGrid map[int]map[int]EntityID
