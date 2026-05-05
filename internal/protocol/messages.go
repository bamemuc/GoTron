package protocol

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

type JoinedPayload struct {
	Id int `json:"id"`
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
	Players [2]PlayerState `json:"players"`
	Status  string         `json:"status"`
	Tick    int            `json:"tick"`
}

type EndPayload struct {
	Draw     bool `json:"draw"`
	WinnerID int  `json:"winnerID"`
}
