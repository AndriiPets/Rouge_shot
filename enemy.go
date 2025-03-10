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
	Bomber
)

type Enemy struct {
	ID       EntityID
	sprite   *Sprite
	collider *Collider

	//COMBAT
	health      int
	weapon      *Weapon
	meleeDamage int
	AI          AiType

	Path PathList

	DetectedPlayer bool
	IsMoving       bool
	HasFiered      bool
}

func NewEnemy(x, y float64, aiType AiType) {
	id := NewID()
	enemySprite := NewSprite(x, y, id)
	enemyCollider := NewCollider(int(x), int(y), 16, 16, id)

	var health, melee, wRange, cooldown int

	switch aiType {
	case Melee:
		enemySprite.LoadImageFromFile("assets/images/skeleton.png")
		health = 2
		melee = 2
		wRange = 2
		cooldown = 1

	case Shooter:
		enemySprite.LoadImageFromFile("assets/images/skeleton2.png")
		health = 3
		melee = 1
		wRange = 4
		cooldown = 1

	case Bomber:
		enemySprite.LoadImageFromFile("assets/images/skeleton2.png")
		health = 4
		melee = 2
		wRange = 8
		cooldown = 7
	}
	enemy := &Enemy{
		ID:          id,
		sprite:      enemySprite,
		collider:    enemyCollider,
		health:      health,
		weapon:      &Weapon{fireRange: float32(wRange), damage: 1, AimDir: None, Cooldown: cooldown, CooldownCount: cooldown},
		AI:          aiType,
		meleeDamage: melee,
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

	if e.HasFiered {
		mov = false
		e.HasFiered = false
	}

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

		if (e.sprite.X == playerX || e.sprite.Y == playerY) && detect {

			if len(line) <= int(e.weapon.fireRange) && len(line) > 2 {

				//turn weapon towards the player before shooting
				if e.sprite.X > playerX {
					e.weapon.AimDir = Left
				}
				if e.sprite.X < playerX {
					e.weapon.AimDir = Right
				}
				if e.sprite.Y < playerY {
					e.weapon.AimDir = Down
				}
				if e.sprite.Y > playerY {
					e.weapon.AimDir = Up
				}

				fmt.Println("enemy ready to fire", len(line))
				e.weapon.UpdateAim(e.sprite.X, e.sprite.Y)

				AttackArea(e.weapon.DamageArea, e.weapon)
				e.HasFiered = true

				mov = false
			}
		}

	case Bomber:
		e.weapon.UpdateCooldown()
		fmt.Println("bomber cooldown", e.weapon.CooldownCount)
		if detect {

			if len(line) <= int(e.weapon.fireRange) && len(line) > 2 && !e.weapon.OnCooldown {
				fmt.Println("Bomber ready to bomb", len(line))
				vel := Vec2{}
				//turn weapon towards the player before shooting
				if e.sprite.X > playerX {
					e.weapon.AimDir = Left
					if e.sprite.Y < playerY {
						vel = Vec2{-16, 16}
					}
					if e.sprite.Y > playerY {
						vel = Vec2{-16, -16}
					}
					if e.sprite.Y == playerY {
						vel = Vec2{-16, 0}
					}
				}
				if e.sprite.X < playerX {
					e.weapon.AimDir = Right
					if e.sprite.Y < playerY {
						vel = Vec2{16, 16}
					}
					if e.sprite.Y > playerY {
						vel = Vec2{16, -16}
					}
					if e.sprite.Y == playerY {
						vel = Vec2{16, 0}
					}
				}
				if e.sprite.X == playerX {
					if e.sprite.Y < playerY {
						e.weapon.AimDir = Down
						vel = Vec2{0, 16}
					}
					if e.sprite.Y > playerY {
						e.weapon.AimDir = Up
						vel = Vec2{0, -16}
					}
				}

				exp := NewExplosive(e.sprite.X+vel.X, e.sprite.Y+vel.Y, Dynamite)
				exp.SetVelocity(vel)
				e.weapon.OnCooldown = true
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

	if e.DetectedPlayer {
		//fmt.Println("Detected by enemy")
		e.BuildPath()
	}

	if e.Path.Steps.HasNext() && e.IsMoving {
		switch e.Path.Steps.Next().String() {
		case "Right":
			dx += CELL_SIZE
			//e.weapon.AimDir = Right
		case "Left":
			dx -= CELL_SIZE
			//e.weapon.AimDir = Left
		case "Up":
			dy -= CELL_SIZE
			//e.weapon.AimDir = Up
		case "Down":
			//e.weapon.AimDir = Down
			dy += CELL_SIZE

		}

		e.Move(dx, dy)

	}

	e.UpdateAI()

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
