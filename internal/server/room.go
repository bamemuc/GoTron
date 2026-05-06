package server

import (
	"encoding/json"
	"gotron/internal/game"
	"gotron/internal/protocol"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	mu     sync.Mutex
	inputs chan<- game.PlayerInput
	states <-chan game.State
	ws     [2]*websocket.Conn
	done   chan struct{}
	loser  int
}

func (r *Room) Start(ws0 *websocket.Conn, ws1 *websocket.Conn) {
	inputs := make(chan game.PlayerInput)
	states := make(chan game.State)

	var pos0 = game.Position{
		X: 1500,
		Y: 500,
	}
	var pos1 = game.Position{
		X: 1500,
		Y: 2500,
	}

	var p0 = game.Player{
		Position:    pos0,
		TrailRender: nil,
		TrailMap:    make(map[game.Position]bool),
		Direction:   game.Down,
		Alive:       true,
		Id:          0,
	}

	var p1 = game.Player{
		Position:    pos1,
		TrailRender: nil,
		TrailMap:    make(map[game.Position]bool),
		Direction:   game.Up,
		Alive:       true,
		Id:          1,
	}

	var state = game.State{
		Players: [2]*game.Player{&p0, &p1},
		Length:  3000,
		Width:   3000,
		Tick:    0,
	}

	r.inputs = inputs
	r.states = states
	r.ws[0] = ws0
	r.ws[1] = ws1

	go func() {
		r.loser = game.Run(&state, 25, inputs, states)
		close(states)
	}()
	go r.inputReader(ws0, 0, inputs)
	go r.inputReader(ws1, 1, inputs)
	go r.outputWriter()
}

func (r *Room) inputReader(conn *websocket.Conn, playerId int, inputs chan<- game.PlayerInput) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg protocol.Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		switch msg.Type {

		case "input":
			var p protocol.InputPayload
			err = json.Unmarshal(msg.Payload, &p)
			if err != nil {
				log.Println(err)
				return
			}
			inputs <- game.PlayerInput{
				PlayerId:  playerId,
				Direction: parseDirection(p.Direction),
			}
		}
	}
}

func (r *Room) outputWriter() {
	defer close(r.done)
	for state := range r.states {
		statePayload, err := json.Marshal(r.toStatePayload(state))
		if err != nil {
			log.Println(err)
			return
		}
		msg := protocol.Message{
			Type:    "state",
			Payload: statePayload,
		}
		msgJson, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}

		err = r.ws[0].WriteMessage(websocket.TextMessage, msgJson)
		if err != nil {
			return
		}

		err = r.ws[1].WriteMessage(websocket.TextMessage, msgJson)
		if err != nil {
			return
		}
	}
	r.sendMatchEnd(r.loser)
}

func parseDirection(s string) game.Direction {
	switch s {
	case "up":
		return game.Up
	case "down":
		return game.Down
	case "left":
		return game.Left
	default:
		return game.Right
	}
}

func (r *Room) Join(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ws[0] == nil {
		r.ws[0] = conn
		sendJoined(conn, 0)
	} else if r.ws[1] == nil {
		r.ws[1] = conn
		sendJoined(conn, 1)
		r.Start(r.ws[0], r.ws[1])
	}
}

func sendJoined(conn *websocket.Conn, id int) {
	payload, err := json.Marshal(protocol.JoinedPayload{Id: id})
	if err != nil {
		log.Println(err)
		return
	}
	msg := protocol.Message{Type: "joined", Payload: payload}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, msgJson); err != nil {
		log.Println(err)
	}
}

func (r *Room) toStatePayload(state game.State) protocol.StatePayload {

	var p0 = state.Players[0]
	var p1 = state.Players[1]

	trail0 := make([]protocol.Position, len(p0.TrailRender))
	for i, pos := range p0.TrailRender {
		trail0[i] = protocol.Position{X: pos.X, Y: pos.Y}
	}

	trail1 := make([]protocol.Position, len(p1.TrailRender))
	for i, pos := range p1.TrailRender {
		trail1[i] = protocol.Position{X: pos.X, Y: pos.Y}
	}

	var ps0 = protocol.PlayerState{
		Position: protocol.Position{
			X: state.Players[0].Position.X,
			Y: state.Players[0].Position.Y,
		},
		Trail: trail0,
		Alive: state.Players[0].Alive,
	}

	var ps1 = protocol.PlayerState{
		Position: protocol.Position{
			X: state.Players[1].Position.X,
			Y: state.Players[1].Position.Y,
		},
		Trail: trail1,
		Alive: state.Players[1].Alive,
	}

	return protocol.StatePayload{
		Players: [2]protocol.PlayerState{ps0, ps1},
		Status:  "running",
		Tick:    state.Tick,
	}

}

func (r *Room) sendMatchEnd(loser int) {

	var endPayload protocol.EndPayload

	if loser == -1 {
		endPayload = protocol.EndPayload{
			Draw:     true,
			WinnerID: -1,
		}
	} else {
		endPayload = protocol.EndPayload{
			Draw:     false,
			WinnerID: 1 - loser,
		}
	}

	EndPayload, err := json.Marshal(endPayload)
	if err != nil {
		log.Println(err)
		return
	}

	msg := protocol.Message{
		Type:    "end",
		Payload: EndPayload,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	err = r.ws[0].WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		return
	}

	err = r.ws[1].WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		return
	}
}

func (r *Room) IsDone() bool {
	select {
	case <-r.done:
		return true
	default:
		return false
	}
}

func NewRoom() *Room {
	return &Room{done: make(chan struct{})}
}
