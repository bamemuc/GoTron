package server

// TODO: WebSocket-Verbindungen annehmen und Read/Write-Loops starten.
import (
	"encoding/json"
	"gotron/internal/protocol"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsHandler(w http.ResponseWriter, r *http.Request, room *Room) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	room.Join(conn)
	<-room.done
}

func write(conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte("hello"))
}

func read(conn *websocket.Conn) {
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

		case "join":
			var p protocol.JoinPayload
			err = json.Unmarshal(msg.Payload, &p)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(p)

		case "input":
			var p protocol.InputPayload
			err = json.Unmarshal(msg.Payload, &p)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(p)
		}
	}
}
