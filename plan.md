# Web Client Plan

A two-phase plan for the browser client. Phase 1 = a playable game end-to-end.
Phase 2 = visual polish without changing the game logic.

---

## Constants from the server (so the client matches)

- World size: **3000 x 3000** (from `room.go`)
- Tick rate: **25 Hz**
- Player 0 spawn: `(1500, 10)`, default direction `Up` (0)
- Player 1 spawn: `(1500, 2990)`, default direction `Down` (1)
- Direction enum (server side): `0=Up, 1=Down, 2=Left, 3=Right`

The server broadcasts `StatePayload` every tick:
```json
{ "players": [{ "position": {x,y}, "trail": [...], "alive": true }, ...],
  "status": "running", "tick": 42 }
```

---

## Server gap to fix first (small)

The client needs to know **which player it is** (0 or 1) so the camera can follow
the right one. Currently `room.Join()` accepts the connection but never sends a
`joined` message back. The `JoinedPayload` struct already exists in
`protocol/messages.go`.

**Fix**: in `room.Join()`, after assigning `r.ws[0]` / `r.ws[1]`, send a
`Message{Type:"joined", Payload: {id: 0|1}}` to that connection.

Without this, the client cannot reliably know which player to follow.

---

## Phase 1 — Make it work

Goal: two browser tabs can play a full match end to end. Visuals can be ugly.

### 1. Wire up static file serving
In `cmd/server/main.go`, add:
```go
mux.Handle("/", http.FileServer(http.Dir("web")))
```
*(must come last so `/ws` and `/health` are not shadowed)*

### 2. Create `web/index.html`
Single file, no build step. Structure:
```
<canvas id="game" width="800" height="800"></canvas>
<div id="status">Connecting...</div>
<script>
  // 1. WebSocket connection
  // 2. Send join → wait for joined → store myId
  // 3. Capture arrow keys → send input messages
  // 4. On state message: update local cache, request animation frame
  // 5. On end message: show winner overlay
</script>
```

### 3. WebSocket message handling
- `onopen`: send `{"type":"join","payload":{"username":"p","color":"#0ff"}}`
- `onmessage`: `switch msg.type`
  - `joined`  → store `myId` (0 or 1)
  - `state`   → store latest snapshot, mark dirty for redraw
  - `end`     → show overlay (winner, draw)

### 4. Input handling
- `keydown` listener for ArrowUp/Down/Left/Right
- Map to `"up"|"down"|"left"|"right"` strings (server `parseDirection` expects these)
- Send `{"type":"input","payload":{"direction":"up"}}`
- Throttle: only send when direction actually changed

### 5. Minimal rendering (no camera yet)
- Render everything at 1:1 scale, fit world into canvas (scale = canvas.width / 3000)
- Clear canvas, draw world border rect, draw both trails as polylines, draw both heads as small squares
- This is intentionally crude — it just proves the game works

### Phase 1 done when:
- Two tabs can join, see each other move, and one wins on collision

---

## Phase 2 — Camera, borders, direction

Goal: player-centered camera with clear visibility of borders and heading.

### 1. Camera that follows the player
The world is 3000x3000 — far too big to show on screen at once. We render in
**world coordinates** then translate so the player is centered.

```js
const VIEW_W = 800, VIEW_H = 800;   // canvas size
const ZOOM   = 0.5;                 // 1 world unit = 0.5 px (tweakable)

// Camera target = my player's position
const me = state.players[myId];
const camX = me.position.x;
const camY = me.position.y;

ctx.save();
ctx.translate(VIEW_W/2, VIEW_H/2);
ctx.scale(ZOOM, ZOOM);
ctx.translate(-camX, -camY);
// ... draw world here in world coordinates ...
ctx.restore();
```

Optional: add a tiny smoothing factor (lerp camera toward player by ~0.2 per
frame) so it doesn't snap.

### 2. World border clearly visible
Draw a thick rectangle from `(0,0)` to `(3000,3000)` in world coords, e.g.
`strokeStyle = "#f55"; lineWidth = 8`. Because it's drawn in world space, when
the player is near the edge, the border becomes visible naturally.

Optional: subtle grid every 100 units (helps depth perception when moving).

### 3. Direction indicator
Three options, ordered cheapest → fanciest:

a) **Triangle pointing forward** (recommended). Draw the player as a triangle
   rotated by direction:
   ```
   Up    → 0°
   Right → 90°
   Down  → 180°
   Left  → 270°
   ```
   The trail is still a fat polyline behind it.

b) Square head with a small line poking out the front.

c) Glow/light cone in the heading direction (shader-ish, more work).

Direction comes from… the server doesn't send it in `StatePayload` today (only
position). Two ways to derive it:
- **Infer client-side**: keep the previous position; direction = `current - prev`
  normalized. Free, works fine.
- **Add it to `StatePayload`**: server change. Cleaner but extra work.

Recommend: infer client-side first.

### 4. Trail styling
- Thick line (`lineWidth = 8` in world units), rounded line caps
- Different color per player (player 0 = cyan, player 1 = magenta is classic Tron)
- Slight glow via `shadowBlur` if cheap enough

### 5. UI overlays
- Status text top-left in screen coords (after `ctx.restore()`):
  - "Waiting for opponent…" before game start
  - "Tick: 142" while running (debug, optional)
- Centered overlay on `end`:
  - "You win!" / "You lose." / "Draw."
  - "Refresh to play again."

### Phase 2 done when:
- Camera smoothly follows your player
- Border is unmistakable when near the edge
- You can tell which way you're heading at a glance
- It looks recognizably like Tron

---

## Order of work next session

1. Add `joined` message to `room.Join()` (5 min, server side)
2. Add static file serving line to `main.go` (1 min)
3. Build Phase 1 `web/index.html` — get a playable ugly version (30–60 min)
4. Test with two browser tabs locally
5. Iterate on Phase 2 visuals (camera → border → direction → trail)

---

## Things explicitly NOT in this plan

- Mobile / touch input
- Reconnect after disconnect
- Lobby / play-again button (refresh both tabs to restart)
- Sound
- Sprites / images (canvas primitives are enough for Tron)
- Build tools (no webpack, no TypeScript — plain HTML/JS only)
