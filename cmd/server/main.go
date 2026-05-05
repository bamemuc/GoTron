package main

import (
	"gotron/internal/server"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	room := server.NewRoom()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.WsHandler(w, r, room)
	})

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
