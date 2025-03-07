package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EntityID uuid.UUID

func NewID() EntityID {
	return EntityID(uuid.New())
}

type Entity interface {
	GetID() EntityID
	GetPos() Vec2
}

type StaticObstacle struct {
	ID  EntityID
	Pos Vec2
}

func NewStaticObstacle(x, y float64) *StaticObstacle {
	s := StaticObstacle{
		ID:  NewID(),
		Pos: Vec2{x, y},
	}
	return &s
}

func (o *StaticObstacle) GetID() EntityID {
	return o.ID
}

func (o *StaticObstacle) GetPos() Vec2 {
	return o.Pos
}

func (s *Sprite) LoadImageFromFile(path string) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	s.Img = img
}
