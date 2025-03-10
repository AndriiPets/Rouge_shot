package main

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
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
	img := ebiten.NewImageFromImage(readImage(path))
	s.Img = img
}

func readImage(file string) image.Image {
	b, err := assets.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("Cannot find a file %s", file))
	}
	return bytes2Image(&b)
}

func bytes2Image(raw *[]byte) image.Image {
	img, format, err := image.Decode(bytes.NewReader(*raw))
	if err != nil {
		log.Fatal("Byte2Image Failed:", format, err)
	}

	return img
}
