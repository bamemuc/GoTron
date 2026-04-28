package protocol

// TODO: JSON-Nachrichten (join, input, state_snapshot, match_end) definieren.
import (
	"encoding/json"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type JoinPayload struct {
	Username string `json:"username"`
	Color    string `json:"color"`
}

type InputPayload struct {
	Direction string `json:"direction"`
}

type PlayerState struct {
	Position Position   `json:"position"`
	Trail    []Position `json:"trail"`
	Alive    bool       `json:"alive"`
}

type StatePayload struct {
	Players [2]PlayerState `json:"players"`ok,
	Status  string         `json:"status"`
	Tick    int            `json:"tick"`
}