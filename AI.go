package main

import "fmt"

type BehaviorState uint8
type BehaviorFunc func(*Enemy, *Player) BehaviorState

const (
	EnemyIdle BehaviorState = iota
	EnemyChase
	EnemyAttack
	EnemyFlee
	EnemyDead
)

func LineOfSightCheck(from, to Vec2) bool {
	line := CellsInLine(from.X, to.X, from.Y, to.Y)

	for _, v := range line {
		if gameGlobal.Level.Map.Get(int(v.X), int(v.Y)) == 'x' {
			return false
		}
	}

	return true
}

func IdleBehavior(enenmy *Enemy, player *Player) BehaviorState {
	playerPos, enemyPos := Vec2{player.sprite.X, player.sprite.Y}, Vec2{enenmy.sprite.X, enenmy.sprite.Y}
	if LineOfSightCheck(enemyPos, playerPos) {
		return EnemyChase
	}
	return EnemyIdle
}

func ChaseBehavior(enemy *Enemy, player *Player) BehaviorState {
	if enemy.WeaponRangeCheck(player) {
		fmt.Println("Enemy attacikg!")
		return EnemyAttack
	}

	return EnemyChase
}

func AttackBehavior(enemy *Enemy, player *Player) BehaviorState {
	if !enemy.WeaponRangeCheck(player) {
		return EnemyChase
	}

	return EnemyAttack
}

func AttackBehaviorShooter(enemy *Enemy, player *Player) BehaviorState {
	if enemy.WeaponRangeCheck(player) {

		//check axis aligment
		if enemy.sprite.X == player.sprite.X || enemy.sprite.Y == player.sprite.Y {

			//turn weapon towards the player before shooting
			if enemy.sprite.X > player.sprite.X {
				enemy.weapon.AimDir = Left
			}
			if enemy.sprite.X < player.sprite.X {
				enemy.weapon.AimDir = Right
			}
			if enemy.sprite.Y < player.sprite.Y {
				enemy.weapon.AimDir = Down
			}
			if enemy.sprite.Y > player.sprite.Y {
				enemy.weapon.AimDir = Up
			}

			fmt.Println("enemy ready to fire")
			enemy.weapon.UpdateAim(enemy.sprite.X, enemy.sprite.Y)

			AttackArea(enemy.weapon.DamageArea, enemy.weapon)

			enemy.HasFiered = true
			enemy.IsMoving = false

			return EnemyAttack

		}
	}

	return EnemyChase
}

func AttackBehaviorBomber(enemy *Enemy, player *Player) BehaviorState {
	if enemy.WeaponRangeCheck(player) && !enemy.weapon.OnCooldown {
		fmt.Println("Bomber ready to bomb")
		vel := Vec2{}
		//turn weapon towards the player before shooting
		if enemy.sprite.X > player.sprite.X {
			enemy.weapon.AimDir = Left
			if enemy.sprite.Y < player.sprite.Y {
				vel = Vec2{-16, 16}
			}
			if enemy.sprite.Y > player.sprite.Y {
				vel = Vec2{-16, -16}
			}
			if enemy.sprite.Y == player.sprite.Y {
				vel = Vec2{-16, 0}
			}
		}
		if enemy.sprite.X < player.sprite.X {
			enemy.weapon.AimDir = Right
			if enemy.sprite.Y < player.sprite.Y {
				vel = Vec2{16, 16}
			}
			if enemy.sprite.Y > player.sprite.Y {
				vel = Vec2{16, -16}
			}
			if enemy.sprite.Y == player.sprite.Y {
				vel = Vec2{16, 0}
			}
		}
		if enemy.sprite.X == player.sprite.X {
			if enemy.sprite.Y < player.sprite.Y {
				enemy.weapon.AimDir = Down
				vel = Vec2{0, 16}
			}
			if enemy.sprite.Y > player.sprite.Y {
				enemy.weapon.AimDir = Up
				vel = Vec2{0, -16}
			}
		}

		exp := NewExplosive(enemy.sprite.X+vel.X, enemy.sprite.Y+vel.Y, Dynamite)
		exp.SetVelocity(vel)
		enemy.weapon.OnCooldown = true

		enemy.HasFiered = true

		enemy.IsMoving = false

		return EnemyAttack
	}

	return EnemyChase
}
