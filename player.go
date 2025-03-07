package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayerState uint8

const (
	Idle PlayerState = iota
	AimedDown
	AimedUp
	AimedRight
	AimedLeft
)

type Player struct {
	ID       EntityID
	sprite   *Sprite
	collider *Collider
	weapon   *Weapon
	health   int
	aimed    bool
	state    PlayerState
}

func NewPlayer() *Player {
	id := NewID()
	posX, posY := 10*CELL_SIZE, 5*CELL_SIZE
	palyerSpite := NewSprite(float64(posX), float64(posY), id)
	playerCollider := NewCollider(posX, posY, 16, 16, id)
	palyerSpite.LoadImageFromFile("assets/images/ninja.png")
	player := Player{
		ID:       id,
		sprite:   palyerSpite,
		collider: playerCollider,
		health:   10,
		state:    Idle,
		weapon:   &Weapon{fireRange: 4, damage: 1},
	}

	return &player

}

type Weapon struct {
	fireRange float32
	damage    int
	BB        image.Rectangle
}

func AttackEnemy(rect image.Rectangle, damage int) {
	for _, enemy := range gameGlobal.enemies {
		if enemy != nil {
			if rect.Overlaps(*enemy.collider.BB) {
				enemy.health -= damage
				fmt.Printf("enemy damaged for %d health %d health left\n", damage, enemy.health)
			}

			if enemy.health <= 0 {
				enemy.Kill()
				fmt.Println("Enemy just died")
			}
		}
	}
}

func (p *Player) GetID() EntityID {
	return p.ID
}

func (p *Player) GetPos() Vec2 {
	return Vec2{p.sprite.X, p.sprite.Y}
}

func (p *Player) Update() {
	tick := false
	dx, dy := int(p.sprite.X), int(p.sprite.Y)

	var rect image.Rectangle
	switch p.state {
	case AimedDown:
		rect = image.Rect(int(p.sprite.X), int(p.sprite.Y+16), int(p.sprite.X+16), int(p.sprite.Y+float64(16*p.weapon.fireRange+16)))
	case AimedUp:
		rect = image.Rect(int(p.sprite.X), int(p.sprite.Y-float64(16*p.weapon.fireRange)), int(p.sprite.X+16), int(p.sprite.Y))
	case AimedLeft:
		rect = image.Rect(int(p.sprite.X-float64(16*p.weapon.fireRange)), int(p.sprite.Y+16), int(p.sprite.X), int(p.sprite.Y))
	case AimedRight:
		rect = image.Rect(int(p.sprite.X+16), int(p.sprite.Y), int(p.sprite.X+float64(16*p.weapon.fireRange+16)), int(p.sprite.Y+16))
	}
	p.weapon.BB = rect

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if !p.aimed {
			dx += CELL_SIZE
			tick = true
		} else {
			fmt.Println("aimed right")
			p.state = AimedRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if !p.aimed {
			dx -= CELL_SIZE
			tick = true
		} else {
			fmt.Println("aimed left")
			p.state = AimedLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if !p.aimed {
			dy -= CELL_SIZE
			tick = true
		} else {
			fmt.Println("aimed up")
			p.state = AimedUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if !p.aimed {
			dy += CELL_SIZE
			tick = true
		} else {
			fmt.Println("aimed down")
			p.state = AimedDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		p.aimed = !p.aimed
		if p.aimed {
			fmt.Println("gun ready")

		} else {
			fmt.Println("gun holsterd")
			p.state = Idle
		}

		tick = true
	}

	//fire weapon
	if inpututil.IsKeyJustPressed(ebiten.KeyX) && p.aimed {
		AttackEnemy(p.weapon.BB, p.weapon.damage)
		tick = true
	}

	if tick {
		p.Move(dx, dy)
		gameGlobal.GameTick()
	}
}

func (p *Player) Move(x, y int) {
	HandleMovement(p, x, y, p.sprite)
	p.collider.Move(int(p.sprite.X), int(p.sprite.Y))

}

func (p *Player) Draw(screen *ebiten.Image) {
	vector.StrokeRect(
		screen,
		float32(p.weapon.BB.Min.X),
		float32(p.weapon.BB.Min.Y),
		float32(p.weapon.BB.Dx()),
		float32(p.weapon.BB.Dy()),
		1,
		color.RGBA{255, 0, 0, 255},
		false,
	)

}
