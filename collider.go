package main

import (
	"image"
	"log"
)

type Collider struct {
	BB *image.Rectangle
}

func NewCollider(x, y, w, h int, id EntityID) *Collider {
	rec := image.Rect(x, y, x+w, y+h)
	c := &Collider{
		BB: &rec,
	}

	if gameGlobal == nil {
		log.Fatal("Global game object is not initialized")
	}

	gameGlobal.colliders[id] = c

	return c
}

func (c *Collider) Move(x, y int) {
	c.BB.Min.X = x
	c.BB.Min.Y = y
	c.BB.Max.X = x + 16
	c.BB.Max.Y = y + 16
}

func AddToGrid(entity Entity, X, Y int) {
	id := entity.GetID()
	if _, exists := gameGlobal.grid[Y]; !exists {
		gameGlobal.grid[Y] = make(map[int]EntityID)
	}
	gameGlobal.grid[Y][X] = id
}

func RemoveFromGrid(entity Entity) {
	pos := entity.GetPos()
	if row, exists := gameGlobal.grid[int(pos.Y)]; exists {
		delete(row, int(pos.X))
		if len(row) == 0 {
			delete(gameGlobal.grid, int(pos.Y))
		}
	}
}

func IsCellOcupied(X, Y int) bool {
	if row, exists := gameGlobal.grid[Y]; exists {
		_, occupied := row[X]
		return occupied
	}
	return false
}

func MoveEntity(entity Entity, newX, newY int) {
	RemoveFromGrid(entity)
	AddToGrid(entity, newX, newY)
}

func HandleMovement(entity Entity, newX, newY int, sprite *Sprite) {
	if IsCellOcupied(newX, newY) {
		targetEntityID := gameGlobal.grid[newY][newX]
		ResolveCollision(entity.GetID(), targetEntityID)
	} else {
		MoveEntity(entity, newX, newY)
		sprite.X = float64(newX)
		sprite.Y = float64(newY)
	}
}

func ResolveCollision(entityID, targetID EntityID) {
	//fmt.Println("handle collision")
	//fmt.Println(gameGlobal.grid)
}
