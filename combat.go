package main

func AttackArea(cells []Vec2, weapon *Weapon) {
	for _, cell := range cells {
		if IsCellOcupied(int(cell.X), int(cell.Y)) {
			Id := gameGlobal.grid[int(cell.Y)][int(cell.X)]
			weapon.DoDamage(Id)
		}
	}
}
