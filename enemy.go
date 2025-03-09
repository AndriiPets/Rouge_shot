package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AiType uint8

const (
	Melee AiType = iota
	Shooter
)

type Enemy struct {
	ID       EntityID
	sprite   *Sprite
	collider *Collider

	//COMBAT
	health int
	weapon *Weapon
	AI     AiType

	Path PathList

	DetectedPlayer bool
	IsMoving       bool
}

func NewEnemy(x, y float64, aiType AiType) {
	id := NewID()
	enemySprite := NewSprite(x, y, id)
	enemyCollider := NewCollider(int(x), int(y), 16, 16, id)
	enemySprite.LoadImageFromFile("assets/images/skeleton.png")
	enemy := &Enemy{
		ID:       id,
		sprite:   enemySprite,
		collider: enemyCollider,
		health:   3,
		weapon:   &Weapon{fireRange: 6, damage: 1, AimDir: None},
		AI:       aiType,
	}
	gameGlobal.enemies[id] = enemy

	enemy.Move(int(x), int(y))

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

func PlayerInRange(x, y, dRange int) bool {
	detectionRange := dRange * CELL_SIZE
	playerX, playerY := gameGlobal.player.sprite.X, gameGlobal.player.sprite.Y
	return math.Abs(float64(x)-playerX) < float64(detectionRange) || math.Abs(float64(y)-playerY) < float64(detectionRange)
}

func (e *Enemy) UpdateAI() {
	mov := true
	detect := true
	playerX, playerY := gameGlobal.player.sprite.X, gameGlobal.player.sprite.Y
	line := CellsInLine(e.sprite.X, gameGlobal.player.sprite.X, e.sprite.Y, gameGlobal.player.sprite.Y)

	for _, v := range line {
		if gameGlobal.Level.Map.Get(int(v.X), int(v.Y)) == 'x' {
			detect = false
		}
	}

	switch e.AI {
	case Melee:

	case Shooter:

		if e.sprite.X == playerX || e.sprite.Y == playerY {

			if len(line) <= int(e.weapon.fireRange) && len(line) > 2 {

				//TODO: turn weapon towards the player before shooting

				fmt.Println("enemy ready to fire", len(line))
				e.weapon.UpdateAim(e.sprite.X, e.sprite.Y)

				AttackArea(e.weapon.DamageArea, e.weapon)

				mov = false
			}
		}

	}

	e.IsMoving = mov

	e.DetectedPlayer = detect
}

func (e *Enemy) Update() {

	dx, dy := int(e.sprite.X), int(e.sprite.Y)
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
	e.UpdateAI()

	if e.DetectedPlayer {
		//fmt.Println("Detected by enemy")
		e.BuildPath()
	}

	if e.Path.Steps.HasNext() && e.IsMoving {
		switch e.Path.Steps.Next().String() {
		case "Right":
			dx += CELL_SIZE
			e.weapon.AimDir = Right
		case "Left":
			dx -= CELL_SIZE
			e.weapon.AimDir = Left
		case "Up":
			dy -= CELL_SIZE
			e.weapon.AimDir = Up
		case "Down":
			e.weapon.AimDir = Down
			dy += CELL_SIZE

		}

		e.Move(dx, dy)

	}

	if e.health <= 0 {
		fmt.Println("enemy Dead")
		e.Kill()
	}
}

func (e *Enemy) Move(x, y int) {
	HandleMovement(e, x, y, e.sprite)
	e.collider.Move(int(e.sprite.X), int(e.sprite.Y))
}

func (e *Enemy) DrawDebug(screen *ebiten.Image) {
	if e.DetectedPlayer {
		line := CellsInLine(e.sprite.X, gameGlobal.player.sprite.X, e.sprite.Y, gameGlobal.player.sprite.Y)
		for _, v := range line {
			vector.StrokeRect(
				screen,
				float32(v.X*16),
				float32(v.Y*16),
				16,
				16,
				1,
				color.RGBA{155, 155, 0, 255},
				false,
			)
		}
	}

	for _, v := range e.weapon.DamageArea {
		vector.StrokeRect(
			screen,
			float32(v.X),
			float32(v.Y),
			16,
			16,
			1,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	}
}
