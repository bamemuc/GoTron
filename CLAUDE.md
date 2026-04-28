# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project context

GoTron is a learning project: an authoritative multiplayer 1v1 Tron game over WebSockets. The user is a Go beginner — **act as a tutor**. Guide and review rather than writing code directly. Only implement when explicitly asked or to unblock.

See `docs/LERNPLAN.md` for the step-by-step learning plan.

## Commands

```bash
# Run server locally
go run ./cmd/server

# Build server binary
go build -o bin/server ./cmd/server

# Run all tests
go test ./...

# Run a single test
go test ./internal/game/... -run TestCollision

# Run server in Docker
docker compose -f deployments/docker/docker-compose.yml up --build server

# Add a dependency
go get github.com/some/package
```

## Architecture

**Authoritative server model**: clients send only directional inputs; the server owns all game state and broadcasts a full `StatePayload` every tick.

Message flow:
1. Client connects via WebSocket (`/ws`)
2. Client sends `join` → server responds with `joined` (includes `player_id: 0|1`)
3. Server starts the game loop once both players are connected
4. Every tick: server sends `state` to both clients
5. On collision: server sends `match_end`, loop stops

Key data flow across packages:
- `protocol/` — JSON message types shared between server and (future) client. `Message.Type` is the discriminator; `Message.Payload` is `json.RawMessage` decoded into the appropriate `*Payload` struct.
- `server/ws_handler.go` — upgrades HTTP → WebSocket, spawns per-connection read/write goroutines, hands the connection to `room`.
- `server/room.go` — holds exactly 2 sessions, starts the game loop when full, routes incoming inputs to the loop, broadcasts outgoing state.
- `game/` — pure game logic with no network dependencies: `state.go` holds `GameState`/`Player`/`Position`, `loop.go` runs the ticker, `collision.go` checks walls and trails.
- `session/session.go` — per-player metadata (ID, connection reference).

## Configuration

Environment variables (see `configs/server.env.example`):

| Variable | Default | Meaning |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP/WS listen port |
| `TICK_RATE` | `12` | Game ticks per second |
| `MAX_PLAYERS` | `2` | Players per match (fixed for MVP) |

## Client

The Go client stub in `cmd/client/` is **not used**. The frontend is a browser-based client (HTML + JS + Canvas) served as static files from the Go server under `/`. See `docs/LERNPLAN.md` Schritt 9.

## MVP boundaries

No matchmaking, reconnect, persistence, or accounts before the MVP is complete. One match, one server instance, two players.
