package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	ViewPort   Vec2
	Position   Vec2
	ZoomFactor int
	Rotation   int
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, Zoom: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() Vec2 {
	return Vec2{
		c.ViewPort.X * 0.5,
		c.ViewPort.Y * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position.X, -c.Position.Y)

	//Scale and rotate around the center of the screen
	//m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	/*m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)*/
	m.Translate(c.viewportCenter().X, c.viewportCenter().Y)
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) Update(pos Vec2) {
	c.Position = pos

	if ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		if c.ZoomFactor > -2400 {
			c.ZoomFactor -= 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		if c.ZoomFactor < 2400 {
			c.ZoomFactor += 1
		}
	}
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		return math.NaN(), math.NaN()
	}
}
