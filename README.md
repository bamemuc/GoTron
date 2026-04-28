# GoTron: Minimal Multiplayer MVP

Ein kleines Lernprojekt in Go: ein **autoritativer Multiplayer-Server** für ein simples Tron/Snake-ähnliches 1v1.
Der Fokus liegt auf Networking und Server-Architektur, nicht auf komplexem Gameplay.

## MVP-Umfang (2-3 Wochen)

- 2 Spieler
- 1 Match gleichzeitig
- Feste Tickrate
- Client sendet nur Input, Server berechnet den State
- Server läuft in Docker, Client lokal

## Out of Scope (vor MVP-Abschluss)

- Matchmaking/Lobby-System
- Reconnect
- Persistenz, Accounts, Ranking
- Performance-Optimierung und horizontales Scaling

## Tech-Stack

- Backend: Go
- Kommunikation: WebSockets (JSON)
- Architektur: Authoritative Server Model
- Infrastruktur: Docker + Docker Compose (nur Server)

## Projektstruktur

Eine genaue Übersicht findest du in `docs/PROJECT_STRUCTURE.md`.

## Start (Server in Docker, Client lokal)

Voraussetzung: Docker ist installiert, Go ist lokal installiert.

```bash
docker compose -f deployments/docker/docker-compose.yml up --build server
```

```bash
go run ./cmd/client --server ws://localhost:8080/ws
```

## Nächster Schritt

Implementiere zuerst in dieser Reihenfolge:

1. `internal/protocol/messages.go`
2. `internal/server/ws_handler.go`
3. `internal/game/state.go` und `internal/game/loop.go`
4. `internal/game/collision.go`
5. `internal/server/room.go`
