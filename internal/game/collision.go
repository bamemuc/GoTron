package game

func checkCollision(state State, pId int) bool {
	var p = state.Players[pId]
	var e = state.Players[1-pId]
	var ppos Position = p.Position
	var epos Position = e.Position

	// check if in map bounds
	if ppos.X < 0 || ppos.X >= state.Length || ppos.Y < 0 || ppos.Y >= state.Width {
		return true
	}

	// check for trail collisions
	if p.TrailMap[ppos] || e.TrailMap[ppos] {
		return true
	}

	// check for player collisions
	if ppos == epos {
		return true
	}

	return false
}
