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

func (s *Sprite) LoadImageFromFile(path string) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	s.Img = img
}
