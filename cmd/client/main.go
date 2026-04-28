package main

import (
	"flag"
	"fmt"
)

func main() {
	serverURL := flag.String("server", "ws://localhost:8080/ws", "WebSocket URL des Servers")
	flag.Parse()

	fmt.Printf("TODO: lokaler Client verbindet sich mit %s\n", *serverURL)
}
