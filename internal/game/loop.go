package game

import "time"

// TODO: Feste Tick-Schleife implementieren und State pro Tick aktualisieren.

func Run(state *State, tickRate int) int {
	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	for range ticker.C {
		movePlayer(state.Players[0])
		movePlayer(state.Players[1])
		c0 := checkCollision(*state, 0)
		c1 := checkCollision(*state, 1)
		if c0 && c1 {
			return -1
		} else if c0 {
			return 0
		} else if c1 {
			return 1
		}
		state.Players[0].TrailMap[state.Players[0].Position] = true
		state.Players[1].TrailMap[state.Players[1].Position] = true
	}
	panic("unreachable")
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
