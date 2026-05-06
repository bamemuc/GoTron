package main

import (
	"gotron/internal/server"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	var (
		roomMu      sync.Mutex
		currentRoom = server.NewRoom()
	)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		roomMu.Lock()
		if currentRoom.IsDone() {
			currentRoom = server.NewRoom()
		}
		room := currentRoom
		roomMu.Unlock()
		server.WsHandler(w, r, room)
	})

	mux.Handle("/", http.FileServer(http.Dir("web")))

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
