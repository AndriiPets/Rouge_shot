package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ExplosiveType uint8

const (
	Barrel ExplosiveType = iota
	Dynamite
)

const (
	None AimDirection = iota
	Down
	Up
	Left
	Right
	AllDirections
)

func AttackArea(cells []Vec2, weapon *Weapon) {
	for _, cell := range cells {
		if IsCellOcupied(int(cell.X), int(cell.Y)) {
			Id := gameGlobal.grid[int(cell.Y)][int(cell.X)]
			weapon.DoDamage(Id)
		}
	}
}

type Explosive struct {
	ID         EntityID
	Sprite     *Sprite
	Weapon     *Weapon
	TimeToBoom int
	health     int
	IsExploded bool
	Velocity   Vec2
}

func NewExplosive(x, y float64, ExType ExplosiveType) *Explosive {

	id := NewID()
	expSprite := NewSprite(x, y, id)
	var boom int

	switch ExType {
	case Dynamite:
		expSprite.LoadImageFromFile("assets/images/dynamite.png")
		boom = 8

	case Barrel:
		expSprite.LoadImageFromFile("assets/images/barrel.png")
		boom = 100
	}

	wp := Weapon{
		fireRange: 2,
		damage:    2,
		AimDir:    AllDirections,
	}

	exp := Explosive{
		ID:         id,
		Sprite:     expSprite,
		Weapon:     &wp,
		TimeToBoom: boom,
		health:     1,
		Velocity:   Vec2{0, 0},
	}

	gameGlobal.explosives[id] = &exp
	exp.Move(int(x), int(y))

	return &exp
}

func (ex *Explosive) SetVelocity(vel Vec2) {
	ex.Velocity = vel
}

func (ex *Explosive) Update() {
	if ex.TimeToBoom < 20 {
		ex.TimeToBoom -= 1
		fmt.Println("time to boom", ex.TimeToBoom)
	}

	if ex.Velocity.X != 0 || ex.Velocity.Y != 0 {
		ex.Move(int(ex.Sprite.X+ex.Velocity.X), int(ex.Sprite.Y+ex.Velocity.Y))
	}

	if ex.IsExploded {
		AttackArea(ex.Weapon.DamageArea, ex.Weapon)
		ex.Kill()
	}

	if ex.TimeToBoom <= 0 || ex.health <= 0 {
		ex.Weapon.UpdateAim(ex.Sprite.X, ex.Sprite.Y)
		ex.IsExploded = true
	}

}

func (ex *Explosive) Kill() {
	RemoveFromGrid(ex)
	delete(gameGlobal.sprites, ex.ID)
	delete(gameGlobal.explosives, ex.ID)
}

func (ex *Explosive) Move(x, y int) {
	HandleMovement(ex, x, y, ex.Sprite)
}

func (ex *Explosive) GetID() EntityID {
	return ex.ID
}

func (ex *Explosive) GetPos() Vec2 {
	return Vec2{ex.Sprite.X, ex.Sprite.Y}
}

type Weapon struct {
	fireRange  float32
	damage     int
	BB         image.Rectangle
	AimDir     AimDirection
	DamageArea []Vec2

	//Cooldown
	OnCooldown    bool
	Cooldown      int
	CooldownCount int
}

func (w *Weapon) UpdateCooldown() {
	if w.OnCooldown {
		w.CooldownCount -= 1
		if w.CooldownCount <= 0 {
			w.OnCooldown = false
			w.CooldownCount = w.Cooldown
		}
	}
}

func (w *Weapon) UpdateAim(X, Y float64) {
	area := make([]Vec2, int(w.fireRange))
	switch w.AimDir {
	case Up:
		for i := range int(w.fireRange) {
			cell := Vec2{X, Y - float64(16*(i+1))}
			area = append(area, cell)
		}

	case Down:
		for i := range int(w.fireRange) {
			cell := Vec2{X, Y + float64(16*(i+1))}
			area = append(area, cell)
		}

	case Left:
		for i := range int(w.fireRange) {
			cell := Vec2{X - float64(16*(i+1)), Y}
			area = append(area, cell)
		}

	case Right:
		for i := range int(w.fireRange) {
			cell := Vec2{X + float64(16*(i+1)), Y}
			area = append(area, cell)
		}

	case AllDirections:
		x, y := math.Floor(X/16), math.Floor(Y/16)
		for i := x - float64(w.fireRange); i <= x+float64(w.fireRange); i++ {
			for j := y - float64(w.fireRange); j <= y+float64(w.fireRange); j++ {
				area = append(area, Vec2{i * 16, j * 16})

			}
		}

	}

	w.DamageArea = area
}

func (w *Weapon) DoDamage(id EntityID) {
	if id == gameGlobal.player.ID {
		gameGlobal.player.health -= w.damage
		fmt.Println("Enemy ranged hit player ", "health remaining ", gameGlobal.player.health)
	}

	if enemy, ok := gameGlobal.enemies[id]; ok {
		fmt.Println("Player ranged hit enemy ", "health remaining ", enemy.health)
		enemy.health -= w.damage
	}

	if ex, ok := gameGlobal.explosives[id]; ok {
		fmt.Println("Player ranged hit explosive ", "it goes boom ")
		ex.health -= w.damage
	}
}

func (ex *Explosive) DrawDebug(screen *ebiten.Image) {
	for _, v := range ex.Weapon.DamageArea {
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
