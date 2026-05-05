package game

import (
	"time"
)

// TODO: Feste Tick-Schleife implementieren und State pro Tick aktualisieren.
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
			state.Players[val.PlayerId].Direction = val.Direction
		default:
		}

		movePlayer(state.Players[0])
		movePlayer(state.Players[1])

		c0 := checkCollision(*state, 0)
		c1 := checkCollision(*state, 1)

		p0copy := *state.Players[0]
		p1copy := *state.Players[1]
		snapshot := *state
		snapshot.Players = [2]*Player{&p0copy, &p1copy}

		if c0 && c1 {
			states <- snapshot
			return -1
		} else if c0 {
			states <- snapshot
			return 0
		} else if c1 {
			states <- snapshot
			return 1
		}
		state.Players[0].TrailMap[state.Players[0].Position] = true
		state.Players[1].TrailMap[state.Players[1].Position] = true
		state.Players[0].TrailRender = append(state.Players[0].TrailRender, state.Players[0].Position)
		state.Players[1].TrailRender = append(state.Players[1].TrailRender, state.Players[1].Position)

		states <- snapshot
	}
	return -1
}

func movePlayer(player *Player) {
	var direction = player.Direction
	switch direction {
	case Up:
		player.Position.Y -= 1
	case Down:
		player.Position.Y += 1
	case Left:
		player.Position.X -= 1
	case Right:
		player.Position.X += 1
	}
}
