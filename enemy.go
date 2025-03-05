package main

type Enemy struct {
	ID       EntityID
	sprite   *Sprite
	collider *Collider
	health   int
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

func (e *Enemy) Update() {
	dx, dy := int(e.sprite.X), int(e.sprite.Y)
	if e.sprite.X < gameGlobal.player.sprite.X {
		dx += CELL_SIZE
	} else if e.sprite.X > gameGlobal.player.sprite.X {
		dx -= CELL_SIZE
	} else if e.sprite.Y < gameGlobal.player.sprite.Y {
		dy += CELL_SIZE
	} else if e.sprite.Y > gameGlobal.player.sprite.Y {
		dy -= CELL_SIZE
	}

	HandleMovement(e, dx, dy, e.sprite)
	e.collider.Move(int(e.sprite.X), int(e.sprite.Y))

}
