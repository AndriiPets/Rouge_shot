package main

import (
	"fmt"
	"math"
)

type Enemy struct {
	ID       EntityID
	sprite   *Sprite
	collider *Collider
	health   int

	Path PathList
}

func NewEnemy(x, y float64) {
	id := NewID()
	enemySprite := NewSprite(x, y, id)
	enemyCollider := NewCollider(int(x), int(y), 16, 16, id)
	enemySprite.LoadImageFromFile("assets/images/skeleton.png")
	enemy := &Enemy{
		ID:       id,
		sprite:   enemySprite,
		collider: enemyCollider,
		health:   1,
	}
	gameGlobal.enemies[id] = enemy

}

func (e *Enemy) GetID() EntityID {
	return e.ID
}

func (e *Enemy) GetPos() Vec2 {
	return Vec2{e.sprite.X, e.sprite.Y}
}

func (e *Enemy) Kill() {
	RemoveFromGrid(e)
	delete(gameGlobal.colliders, e.ID)
	delete(gameGlobal.sprites, e.ID)
	delete(gameGlobal.enemies, e.ID)
}

//TODO: line of sight check for the enemies

func (e *Enemy) BuildPath() {
	if gameGlobal.PathFinder.Grid != nil {

		steps := gameGlobal.PathFinder.MakePath(
			e.sprite.X,
			e.sprite.Y,
			gameGlobal.player.sprite.X,
			gameGlobal.player.sprite.Y,
			AStar,
		)

		e.Path = steps
	}
}

func (e *Enemy) Update() {
	dx, dy := int(e.sprite.X), int(e.sprite.Y)
	detectionRange := 5 * CELL_SIZE
	playerX, playerY := gameGlobal.player.sprite.X, gameGlobal.player.sprite.Y
	/*
		if e.sprite.X < gameGlobal.player.sprite.X {
			dx += CELL_SIZE
		} else if e.sprite.X > gameGlobal.player.sprite.X {
			dx -= CELL_SIZE
		} else if e.sprite.Y < gameGlobal.player.sprite.Y {
			dy += CELL_SIZE
		} else if e.sprite.Y > gameGlobal.player.sprite.Y {
			dy -= CELL_SIZE
		}

		e.Move(dx, dy)*/
	if math.Abs(playerX-float64(dx)) < float64(detectionRange) || math.Abs(playerY-float64(dy)) < float64(detectionRange) {
		fmt.Println("Detected by enemy")
		e.BuildPath()
	}

	if e.Path.Steps.HasNext() {
		switch e.Path.Steps.Next().String() {
		case "Right":
			dx += CELL_SIZE
		case "Left":
			dx -= CELL_SIZE
		case "Up":
			dy -= CELL_SIZE
		case "Down":
			dy += CELL_SIZE

		}

		e.Move(dx, dy)
	}
}

func (e *Enemy) Move(x, y int) {
	HandleMovement(e, x, y, e.sprite)
	e.collider.Move(int(e.sprite.X), int(e.sprite.Y))
}
