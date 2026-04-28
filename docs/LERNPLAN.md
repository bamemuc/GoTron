# GoTron Lernplan

Ziel: Du programmierst so viel wie möglich selbst. Claude erklärt, gibt Hinweise, und hilft wenn du steckst.

---

## Wie wir zusammenarbeiten

- **Du machst**: Tippen, nachdenken, entscheiden.
- **Claude macht**: Erklären, reviewen, unblockieren wenn du fragst.
- Wenn du nicht weiterkommst: Frag! "Ich verstehe X nicht" ist immer die richtige Antwort.

---

## Wochenende 1 — Grundlagen & Verbindung

### Lernziele
- Go Basics: Structs, Packages, exportierte Felder
- HTTP-Server starten
- Erste WebSocket-Verbindung aufbauen

### Schritt 1: Protocol fertigstellen (`internal/protocol/messages.go`)
Was du lernen wirst: Structs in Go, JSON-Tags, exportierte vs. unexportierte Felder.

**Aufgabe:** Korrigiere und vervollständige die Nachrichten-Typen.

Checkliste:
- [ ] Alle Struct-Felder sind **großgeschrieben** (exportiert)
- [ ] Jedes Feld hat einen JSON-Tag (`json:"field_name"`)
- [ ] `JoinPayload.Username` ist ein `string`, nicht `int`
- [ ] `StatePayload` enthält Trail-Daten beider Spieler (nicht nur aktuelle Position)
- [ ] `MatchEndPayload` existiert (wer hat gewonnen?)

Hinweis zu JSON-Tags:
```go
type Example struct {
    Name string `json:"name"`  // ← so sieht ein JSON-Tag aus
}
```

### Schritt 2: WebSocket-Dependency hinzufügen
```bash
go get github.com/gorilla/websocket
```
Danach: schau dir `go.mod` und `go.sum` an — was hat Go hinzugefügt?

### Schritt 3: WebSocket-Handler (`internal/server/ws_handler.go`)
Was du lernen wirst: HTTP-Handler in Go, WebSocket-Upgrade, Goroutinen, Channels.

**Aufgabe:** Schreibe einen Handler, der:
1. Eine HTTP-Verbindung zu WebSocket upgraded
2. In einer Goroutine Nachrichten liest und sie (als Echo) zurückschickt

Das ist dein erstes echtes Networking in Go. Fang klein an — Echo zuerst, dann weiter.

### Schritt 4: Server verdrahten (`cmd/server/main.go`)
**Aufgabe:** Registriere deinen ws_handler unter `/ws`.

Test: Verbinde dich mit einem WebSocket-Testclient (z.B. `wscat` oder einem Browser-Tool) und schau ob der Echo funktioniert.

---

## Wochenende 2 — Spiellogik

### Lernziele
- Go Interfaces und Methoden auf Structs
- Ticker für feste Tickrate
- Mutex für parallelen Zugriff (Goroutinen + shared state)

### Schritt 5: Game State (`internal/game/state.go`)
Was du lernen wirst: Structs mit Methoden, Slices für Trails.

**Aufgabe:** Definiere:
- `Direction` (Up/Down/Left/Right)
- `Player` mit Position, Trail (Slice von Punkten), Farbe, Status (alive/dead)
- `GameState` mit Spielfeld-Größe, beiden Spielern, aktuellem Tick

### Schritt 6: Kollisionserkennung (`internal/game/collision.go`)
**Aufgabe:** Schreibe eine Funktion `CheckCollision(state *GameState) (bool, PlayerID)` die prüft:
- Spieler trifft Wand
- Spieler trifft seinen eigenen Trail
- Spieler trifft Trail des Gegners

### Schritt 7: Game Loop (`internal/game/loop.go`)
Was du lernen wirst: `time.Ticker`, Goroutinen, wie man einen Game Loop baut.

**Aufgabe:** Schreibe eine Schleife die:
1. Jeden Tick: Spieler bewegt, Trail erweitert, Kollision prüft
2. Bei Kollision: Match beendet
3. Nach jedem Tick: State an alle Spieler broadcastet

### Schritt 8: Room (`internal/server/room.go`)
**Aufgabe:** Verbinde WebSocket-Verbindungen mit dem Game Loop. Room hält:
- 2 Spieler-Sessions
- Den GameState
- Startet den Loop wenn beide verbunden sind

---

## Wochenende 3 — Client & Abschluss

### Lernziele
- Entweder: Terminal-Rendering (Go + tcell)
- Oder: Browser-Client (HTML + JavaScript) — empfohlen für CV

### Schritt 9: Browser-Client
**Entschieden: Browser-Client** (HTML + JavaScript + Canvas)

Der Go-Client unter `cmd/client/` entfällt. Stattdessen liefert der Go-Server eine statische HTML-Seite aus.

Aufgabe:
- `web/index.html` mit einem `<canvas>` Element
- JavaScript das sich per WebSocket verbindet, Inputs sendet und den State zeichnet
- Go-Server liefert `web/` als statische Dateien unter `/` aus

### Schritt 10: Docker & End-to-End
**Aufgabe:** `deployments/docker/` finalisieren, Server containerisieren, End-to-End testen.

---

## Go-Konzepte die du in diesem Projekt lernst

| Konzept | Wo du es anwendest |
|---|---|
| Structs & Methoden | game/state.go |
| Interfaces | session.go, ws_handler.go |
| Goroutinen | game loop, ws reader/writer |
| Channels | Kommunikation Room ↔ Loop |
| Mutex | Shared state absichern |
| JSON marshal/unmarshal | protocol/messages.go |
| HTTP/WebSocket | server/ws_handler.go |
| time.Ticker | game/loop.go |
| Docker | deployments/ |

---

## Aktueller Stand

- [x] Projektstruktur erstellt
- [ ] **JETZT: `protocol/messages.go` fertigstellen** ← Start hier
- [ ] WebSocket-Dependency hinzufügen
- [ ] WebSocket-Handler mit Echo
- [ ] Game State definieren
- [ ] Kollisionserkennung
- [ ] Game Loop
- [ ] Room
- [ ] Client
- [ ] Docker End-to-End
