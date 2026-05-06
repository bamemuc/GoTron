package game

import (
	"time"
)

const StepsPerTick = 10

type PlayerInput struct {
	PlayerId  int
	Direction Direction
}

func Run(state *State, tickRate int, inputs <-chan PlayerInput, states chan<- State) int {
	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	for range ticker.C {
		select {
		case val := <-inputs:
			cur := state.Players[val.PlayerId].Direction
			if !isOpposite(cur, val.Direction) {
				state.Players[val.PlayerId].Direction = val.Direction
			}
		default:
		}

		loser := -2 // sentinel: no collision yet

		for s := 0; s < StepsPerTick; s++ {
			movePlayer(state.Players[0])
			movePlayer(state.Players[1])

			c0 := checkCollision(*state, 0)
			c1 := checkCollision(*state, 1)

			if c0 && c1 {
				loser = -1
				break
			} else if c0 {
				loser = 0
				break
			} else if c1 {
				loser = 1
				break
			}

			state.Players[0].TrailMap[state.Players[0].Position] = true
			state.Players[1].TrailMap[state.Players[1].Position] = true
			state.Players[0].TrailRender = append(state.Players[0].TrailRender, state.Players[0].Position)
			state.Players[1].TrailRender = append(state.Players[1].TrailRender, state.Players[1].Position)
		}

		state.Tick++

		p0copy := *state.Players[0]
		p1copy := *state.Players[1]
		snapshot := *state
		snapshot.Players = [2]*Player{&p0copy, &p1copy}

		states <- snapshot

		if loser != -2 {
			return loser
		}
	}
	return -1
}

func isOpposite(a, b Direction) bool {
	return (a == Up && b == Down) || (a == Down && b == Up) ||
		(a == Left && b == Right) || (a == Right && b == Left)
}

func movePlayer(player *Player) {
	switch player.Direction {
	case Up:
		player.Position.Y--
	case Down:
		player.Position.Y++
	case Left:
		player.Position.X--
	case Right:
		player.Position.X++
	}
}
