package game

import (
	"testing"
)

func makeState(p0pos, p1pos Position) State {
	return State{
		Players: [2]*Player{
			{Position: p0pos, TrailMap: make(map[Position]bool), Alive: true, Id: 0},
			{Position: p1pos, TrailMap: make(map[Position]bool), Alive: true, Id: 1},
		},
		Length: 100,
		Width:  100,
	}
}

func TestNoCollision(t *testing.T) {
	state := makeState(Position{10, 10}, Position{50, 50})
	if checkCollision(state, 0) {
		t.Error("expected no collision")
	}
}

func TestWallLeft(t *testing.T) {
	state := makeState(Position{-1, 10}, Position{50, 50})
	if !checkCollision(state, 0) {
		t.Error("expected wall collision on left")
	}
}

func TestWallRight(t *testing.T) {
	state := makeState(Position{100, 10}, Position{50, 50})
	if !checkCollision(state, 0) {
		t.Error("expected wall collision on right")
	}
}

func TestWallTop(t *testing.T) {
	state := makeState(Position{10, -1}, Position{50, 50})
	if !checkCollision(state, 0) {
		t.Error("expected wall collision on top")
	}
}

func TestWallBottom(t *testing.T) {
	state := makeState(Position{10, 100}, Position{50, 50})
	if !checkCollision(state, 0) {
		t.Error("expected wall collision on bottom")
	}
}

func TestOwnTrail(t *testing.T) {
	state := makeState(Position{10, 10}, Position{50, 50})
	state.Players[0].TrailMap[Position{10, 10}] = true
	if !checkCollision(state, 0) {
		t.Error("expected collision with own trail")
	}
}

func TestEnemyTrail(t *testing.T) {
	state := makeState(Position{10, 10}, Position{50, 50})
	state.Players[1].TrailMap[Position{10, 10}] = true
	if !checkCollision(state, 0) {
		t.Error("expected collision with enemy trail")
	}
}

func TestHeadOn(t *testing.T) {
	state := makeState(Position{10, 10}, Position{10, 10})
	if !checkCollision(state, 0) {
		t.Error("expected head-on collision")
	}
	if !checkCollision(state, 1) {
		t.Error("expected head-on collision for player 1")
	}
}
