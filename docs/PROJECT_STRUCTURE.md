# Projektstruktur (MVP)

Diese Struktur ist fuer ein simples 2-3 Wochen Tron-MVP ausgelegt: **Server in Docker, Client lokal**.

## Verzeichnisbaum

```text
GoTron/
├── cmd/
│   ├── server/main.go
│   └── client/main.go
├── internal/
│   ├── config/config.go
│   ├── game/
│   │   ├── state.go
│   │   ├── loop.go
│   │   └── collision.go
│   ├── protocol/messages.go
│   ├── server/
│   │   ├── ws_handler.go
│   │   └── room.go
│   └── session/session.go
├── deployments/docker/
│   ├── Dockerfile.server
│   └── docker-compose.yml
├── configs/server.env.example
├── docs/PROJECT_STRUCTURE.md
├── test/integration/.gitkeep
├── go.mod
├── .gitignore
└── README.md
```

## Was kommt wohin?

- `cmd/server/main.go`: Startpunkt des Servers; HTTP-Routen, WebSocket-Route und Serverstart.
- `cmd/client/main.go`: Startpunkt fuer lokalen Client; verbindet sich auf `ws://localhost:8080/ws`.
- `internal/config/config.go`: Einlesen und Halten der Konfiguration (Port, Tickrate, Max-Spieler).
- `internal/game/state.go`: Datenmodell fuer Spielfeld, Spieler, Positionen und Trails.
- `internal/game/loop.go`: Autoritative Tick-Schleife (Bewegung berechnen, Regeln anwenden).
- `internal/game/collision.go`: Wand- und Trail-Kollisionen pruefen.
- `internal/protocol/messages.go`: JSON-Nachrichten fuer Client/Server (z. B. `join`, `input`, `state_snapshot`).
- `internal/server/ws_handler.go`: WebSocket-Verbindungen annehmen, lesen, schreiben, disconnect behandeln.
- `internal/server/room.go`: Match-Verwaltung fuer ein laufendes 1v1-Spiel (spaeter erweiterbar).
- `internal/session/session.go`: Session-Metadaten pro Spieler (ID, Verbindungsstatus, letzter Ping).
- `deployments/docker/Dockerfile.server`: Image-Build nur fuer den Server.
- `deployments/docker/docker-compose.yml`: Startet den Server-Container mit Port-Mapping.
- `configs/server.env.example`: Beispielwerte fuer lokale Konfiguration.
- `test/integration/`: Spaeter End-to-End-Tests fuer Join -> Matchstart -> Matchende.

## Empfohlene Reihenfolge der Implementierung

1. `cmd/server/main.go` + `internal/server/ws_handler.go`
2. `internal/protocol/messages.go`
3. `internal/game/state.go` + `internal/game/loop.go` + `internal/game/collision.go`
4. `internal/server/room.go` + `internal/session/session.go`
5. `cmd/client/main.go`
6. `deployments/docker/*` final pruefen

## MVP-Regeln (knapp)

- Erstziel: **2 Spieler, 1 Match, 1 Serverinstanz**.
- Kein Matchmaking, kein Reconnect, keine Persistenz vor MVP-Abschluss.
- Server bleibt authoritative: Client sendet nur Inputs, Server berechnet den State.

