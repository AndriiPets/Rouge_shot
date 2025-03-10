package main

import (
	"fmt"
	"image"
	"log"
	"math"
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
	if entityID == gameGlobal.player.ID {
		enemy, ok := gameGlobal.enemies[targetID]
		if ok {
			enemy.health -= int(gameGlobal.player.MeleeDamage)
			fmt.Println("Player melee hit enemy ", "health remaining ", enemy.health)
		}
	}

	if enemy, ok := gameGlobal.enemies[entityID]; ok {
		if targetID == gameGlobal.player.ID {
			gameGlobal.player.health -= enemy.meleeDamage
			fmt.Println("Enemy melee hit player ", "health remaining :", gameGlobal.player.health)
		}
	}
}

func CellsInLine(x0, x1, y0, y1 float64) []Vec2 {
	x0, x1 = math.Floor(x0/16), math.Floor(x1/16)
	y0, y1 = math.Floor(y0/16), math.Floor(y1/16)
	var cells []Vec2
	dx := math.Abs(x0 - x1)
	dy := math.Abs(y0 - y1)
	var sx, sy float64
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	er := dx - dy

	for {
		cells = append(cells, Vec2{x0, y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * er
		if e2 > -dy {
			er -= dy
			x0 += sx
		}
		if e2 < dx {
			er += dx
			y0 += sy
		}
	}
	return cells
}
