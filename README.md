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

## How to play

### Same network (local)

**Requirements:** Go 1.22+

```bash
go run ./cmd/server
```

Open `http://localhost:8080` in two browser tabs. Game starts when both players are connected.

**Controls:** Arrow keys or WASD

---

### With friends online

**Requirements:** [Docker Desktop](https://www.docker.com/products/docker-desktop/)

**1. Start the server**
```bash
docker run -p 8080:8080 bastiangrut/gotron:latest
```

**2. Get your public URL** (in a second terminal)
```bash
docker logs $(docker ps -q --filter ancestor=bastiangrut/gotron:latest) 2>&1 | grep "trycloudflare"
```
You'll see something like:
```
https://something-random.trycloudflare.com
```

**3. Share the URL** with your opponent — both open it in a browser and the game starts automatically.

**4. Stop when done**
```bash
docker stop $(docker ps -q --filter ancestor=bastiangrut/gotron:latest)
```

> The URL changes every time you restart — share a fresh one each session.

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

## Play with friends

The server and Cloudflare Tunnel are bundled into one Docker container. No Go, no extra tools needed — just Docker.

**Requirements:** [Docker Desktop](https://www.docker.com/products/docker-desktop/)

**1. Start**
```bash
docker run -p 8080:8080 bastiangrut/gotron:latest
```

**2. Get the public URL**
```bash
docker logs $(docker ps -q --filter ancestor=bastiangrut/gotron:latest) 2>&1 | grep "trycloudflare"
```
Look for a line like:
```
https://something-random.trycloudflare.com
```
Share that URL with your opponent — both open it in a browser and the game starts.

**3. Stop**
```bash
docker stop $(docker ps -q --filter ancestor=bastiangrut/gotron:latest)
```

> The URL changes every time you restart, so share a fresh one each session.

---

## Tech stack

| Layer | Technology |
|---|---|
| Server | Go 1.22 |
| WebSockets | gorilla/websocket |
| Client | HTML + Canvas + vanilla JS |
| Tunneling | Cloudflare Tunnel |
| Containerization | Docker + Docker Compose |
