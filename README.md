# GoTron

A real-time multiplayer 1v1 Tron game built in Go as a learning project.
Two players connect via browser, control their light cycles, and the last one alive wins.

![Go](https://img.shields.io/badge/Go-1.22-blue) ![WebSockets](https://img.shields.io/badge/WebSockets-gorilla-green)

---

## What I built

- Authoritative game server in Go — clients send only directional inputs, the server owns all game state
- Real-time WebSocket communication at 25 ticks/second
- Browser-based client with Canvas rendering, a camera that follows your player, and trail/collision visuals
- Rematch system — no server restart needed between games

---

## What I learned

### Go fundamentals
- **Structs and methods** — modelling game entities (`Player`, `State`, `Position`) with methods attached to types
- **Exported vs unexported identifiers** — understanding why `TrailMap` vs `trailMap` matters across packages
- **JSON tags** — serializing structs to JSON with `json:"field_name"` for WebSocket messages
- **Interfaces** — using `http.Handler` and understanding Go's implicit interface model

### Concurrency
- **Goroutines** — spawning lightweight threads for the game loop, per-player input readers, and the output writer
- **Channels** — communicating safely between goroutines (`inputs chan PlayerInput`, `states chan State`)
- **Mutexes** — protecting shared state (`sync.Mutex` in `Room.Join`) against race conditions
- **`time.Ticker`** — driving a fixed-rate game loop without busy-waiting

### Networking
- **HTTP server** — setting up routes with `http.NewServeMux`, serving static files
- **WebSocket upgrade** — turning an HTTP connection into a persistent WebSocket with `gorilla/websocket`
- **JSON protocol design** — a typed message envelope (`Message.Type` + `Message.Payload`) so both sides know how to decode each message
- **Authoritative server model** — why the server must own state to prevent cheating and keep clients in sync

### Game logic
- **Collision detection** — checking wall bounds and trail maps per tick
- **Fixed timestep loop** — moving players at a predictable rate regardless of processing time
- **Input handling** — ignoring invalid inputs (180° turns) on the server side so clients can't cheat

### Testing
- **Table-free unit tests** in Go — testing `checkCollision` across wall, trail, and head-on scenarios
- **Test helpers** — `makeState()` to reduce boilerplate and keep tests readable

---

## How to run

**Requirements:** Go 1.22+

```bash
# Start the server
go run ./cmd/server

# Open in browser
http://localhost:8080
```

Open two browser tabs — player 1 gets cyan, player 2 gets magenta. Game starts automatically when both are connected.

**Controls:** Arrow keys or WASD

---

## Play online

To play across different networks, run a Cloudflare Tunnel alongside the server:

```bash
cloudflared tunnel --url http://localhost:8080
```

Share the printed URL with your opponent — no port forwarding or cloud server needed.

---

## Project structure

```
internal/
  protocol/   — JSON message types (join, input, state, end)
  game/       — pure game logic: state, collision, loop
  server/     — WebSocket handler and room/match management
  session/    — per-player session metadata
cmd/server/   — server entrypoint
web/          — browser client (HTML + Canvas + JS)
```

---

## Docker — WIP

Container setup is in progress. The goal is to run the server in Docker so it can be deployed to any VPS without installing Go.

```bash
# Not ready yet
docker compose -f deployments/docker/docker-compose.yml up --build server
```

---

## Tech stack

| Layer | Technology |
|---|---|
| Server | Go 1.22 |
| WebSockets | gorilla/websocket |
| Client | HTML + Canvas + vanilla JS |
| Tunneling | Cloudflare Tunnel |
| Containerization | Docker *(WIP)* |
