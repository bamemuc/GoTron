package game

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Position struct {
	X int
	Y int
}

type Player struct {
	Position    Position
	TrailRender []Position        // used to render the trail
	TrailMap    map[Position]bool // used to check if the player crossed a trail
	Direction   Direction
	Alive       bool
	Id          int
}

type State struct {
	Players [2]*Player
	Length  int
	Width   int
	Tick    int
}
