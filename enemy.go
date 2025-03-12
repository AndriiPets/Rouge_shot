package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Enemy struct {
	ID       EntityID
	sprite   *Sprite
	collider *Collider

	//COMBAT
	health      int
	weapon      *Weapon
	meleeDamage int

	//AI
	Behaviors map[BehaviorState]BehaviorFunc
	State     BehaviorState

	Path PathList

	DetectedPlayer bool
	IsMoving       bool
	HasFiered      bool
}

func NewEnemyShooter(x, y float64) {
	enemy := NewEnemy(x, y)
	enemy.Behaviors[EnemyAttack] = AttackBehaviorShooter
}

func NewEnemyBomber(x, y float64) {
	enemy := NewEnemy(x, y)
	enemy.Behaviors[EnemyAttack] = AttackBehaviorBomber
}
func NewEnemy(x, y float64) *Enemy {
	id := NewID()
	enemySprite := NewSprite(x, y, id)
	enemyCollider := NewCollider(int(x), int(y), 16, 16, id)

	var health, melee, wRange, cooldown int

	enemySprite.LoadImageFromFile("assets/images/skeleton.png")
	health = 2
	melee = 2
	wRange = 4
	cooldown = 7

	enemy := &Enemy{
		ID:          id,
		sprite:      enemySprite,
		collider:    enemyCollider,
		health:      health,
		weapon:      &Weapon{fireRange: float32(wRange), damage: 1, AimDir: None, Cooldown: cooldown, CooldownCount: cooldown},
		meleeDamage: melee,
		Behaviors: map[BehaviorState]BehaviorFunc{
			EnemyIdle:   IdleBehavior,
			EnemyChase:  ChaseBehavior,
			EnemyAttack: AttackBehavior,
		},
		State: EnemyIdle,
	}
	gameGlobal.enemies[id] = enemy

	enemy.Move(int(x), int(y))

	return enemy
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

func (e *Enemy) WeaponRangeCheck(p *Player) bool {
	dist := math.Max(
		math.Abs(e.sprite.X-p.sprite.X),
		math.Abs(e.sprite.Y-p.sprite.Y),
	)

	return dist <= float64(e.weapon.fireRange*16)

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

	if e.State == EnemyChase {
		//fmt.Println("Detected by enemy")
		e.BuildPath()
	}
	//fmt.Println(e.State)

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

	e.weapon.UpdateCooldown()
	//e.UpdateAI()
	if behavior, ok := e.Behaviors[e.State]; ok {
		e.State = behavior(e, gameGlobal.player)
	}

	if e.HasFiered {
		e.HasFiered = false
		e.IsMoving = false
	} else {
		e.IsMoving = true
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
	if e.State == EnemyChase {
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
