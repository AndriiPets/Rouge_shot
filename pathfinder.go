package main

import (
	path "github.com/quasilyte/pathing"
)

type PathAlgo int
type PathList path.BuildPathResult
type Step path.Direction

const (
	BFS PathAlgo = iota
	AStar
)

const (
	TileWall = iota
	TileFloor
)

type PathFinder struct {
	Grid   *path.Grid
	Layers path.GridLayer
	BFS    *path.GreedyBFS
	AStar  *path.AStar
}

func NewPathFinder(w, h, cell_size int) PathFinder {
	grid := path.NewGrid(path.GridConfig{
		WorldWidth:  uint(w) * uint(cell_size),
		WorldHeight: uint(h) * uint(cell_size),
		CellWidth:   uint(cell_size),
		CellHeight:  uint(cell_size),
	})

	bfs := path.NewGreedyBFS(path.GreedyBFSConfig{
		NumCols: uint(grid.NumCols()),
		NumRows: uint(grid.NumRows()),
	})

	aStar := path.NewAStar(path.AStarConfig{
		NumCols: uint(grid.NumCols()),
		NumRows: uint(grid.NumRows()),
	})

	return PathFinder{Grid: grid, BFS: bfs, AStar: aStar}
}

func (p *PathFinder) GenerateLayout(data [][]rune, wall rune) {
	for y, row := range data {
		for x, val := range row {

			if val == wall {

				p.Grid.SetCellTile(path.GridCoord{X: x, Y: y}, TileWall)

			} else {

				p.Grid.SetCellTile(path.GridCoord{X: x, Y: y}, TileFloor)

			}
		}
	}

	groundltLayer := path.MakeGridLayer([4]uint8{
		TileFloor: 1,
		TileWall:  0,
	})

	//airLayer := path.MakeGridLayer([4]uint8{
	//	TileFloor: 1,
	//	TileWall:  1,
	//})

	p.Layers = groundltLayer
}

func (p *PathFinder) MakePath(fromX, fromY, toX, toY float64, algo PathAlgo) PathList {
	startPos := p.Grid.PosToCoord(fromX, fromY)
	endPos := p.Grid.PosToCoord(toX, toY)

	var steps PathList

	switch algo {
	case BFS:
		steps = PathList(p.BFS.BuildPath(p.Grid, startPos, endPos, p.Layers))
	case AStar:
		steps = PathList(p.AStar.BuildPath(p.Grid, startPos, endPos, p.Layers))
	}

	return steps
}
