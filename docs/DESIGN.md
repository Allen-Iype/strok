# strok — Design Document

A terminal-based typing app, built for developers.
Terminal-first, lightweight, responsive. Single Go binary. Linux / macOS / Windows.

> This is the written design produced in Phases 1–4. Implementation (Phase 5)
> follows the roadmap at the bottom, one milestone at a time, keeping the app
> compiling after each.

---

## Phase 1 — Product Understanding

### Primary goal
A terminal-first typing app that **always shows a full color-coded keyboard
with finger mapping**, gives **live per-character feedback**, computes
**WPM / accuracy / errors live**, generates **progressive lessons**, and
**persists progress as JSON**. Feels like a polished open-source project.

### MVP scope
- Bubble Tea + Lip Gloss TUI, always-visible keyboard, resize-safe.
- Live keystroke handling.
- Full keyboard, every key → one finger, 9 finger colors.
- Current key highlighted; correct → green flash; wrong → expected yellow + pressed red.
- Typing area: practice text, user input, cursor; correct=green, wrong=red, current=underlined.
- Live stats: WPM, accuracy, errors, chars typed, elapsed.
- Progressive drill-style lesson generator behind a swappable interface.
- JSON persistence: best WPM, avg WPM, accuracy, practice time, completed lessons, weak keys.

### Explicitly later (design for, do NOT build)
Adaptive data-driven lesson algorithm, heatmaps, streaks, themes, multiple layouts,
programming-symbols mode, vim mode, multiplayer, sound, plugin generators.

### Edge cases
Terminal too small; resize mid-session; stats before first keystroke (div-by-zero);
backspace semantics; non-printable/modifier keys, Ctrl+C/Esc; first run (no file);
corrupt/partial JSON; multi-key-per-finger coloring; complete vs quit mid-lesson;
Unicode/wide chars (MVP: ASCII only).

### Assumptions / decisions
1. Layout: US ANSI QWERTY, standard touch-typing finger map.
2. Practice text is ASCII lowercase letters + spaces in MVP.
3. **Backspace allowed, errors permanent:** backspace fixes visible
   text only; any wrong keystroke permanently counts against accuracy/errors.
4. WPM = (correct chars / 5) / minutes. Timer starts on first keystroke.
5. A lesson is one fixed-length practice string (~30–60 chars). Completing it
   records stats and generates the next lesson.
6. Save on lesson completion and on graceful quit. No per-keystroke autosave.
7. Weak keys = per-key error rate accumulated across sessions.
8. No global state; dependency-injected; single binary; Go 1.22+.
9. **Lesson style: drill-style** — random combos of unlocked letters, widening
   the alphabet as the user progresses (`f / j / fj / fjfj`...).
10. Module path `strok`, command at `cmd/strok`, binary `strok`.

---

## Phase 2 — Technical Design

### 1. High-level architecture
Layered, unidirectional. Bubble Tea Elm loop (Model→Update→View) in `ui`; all
domain logic kept outside the UI so it's testable headlessly.

```
cmd/strok            wiring / DI / lifecycle
   │
internal/ui          Bubble Tea Model (orchestration, owns no logic)
   │  ┌──────────┬──────────┬──────────┬──────────┐
engine  lesson   keyboard   stats     storage
   │
internal/domain      pure structs, no deps
```

Principles:
- UI owns no domain logic. `Update` routes msgs to engine/stats; `View` renders.
- Domain is pure (no Bubble Tea, no I/O) — trivially unit-testable.
- Everything pluggable is an interface (`LessonGenerator`, `Store`, `Layout`,
  `Clock`), injected from `cmd/strok`. No globals.
- One-way data flow: keystroke → engine → state → stats → View; on completion → save.

Data flow (one keystroke):
```
tea.KeyMsg → ui.Update → engine.HandleKey(rune) → TypingState
           → stats.Compute(state, now) → ui.View
on lesson end → profile.Apply(session) → store.Save → generator.Next → new state
```

### 2. Project structure
```
strok/
├── go.mod
├── README.md
├── Makefile
├── docs/DESIGN.md
├── cmd/strok/main.go            entry: flags, wire deps, run, handle errors
└── internal/
    ├── domain/                  pure data types, zero deps
    │   ├── finger.go            Finger enum + names
    │   ├── key.go               Key (rune, label, finger, row, width)
    │   ├── lesson.go            Lesson (target text + keyset)
    │   ├── stats.go             Stats snapshot
    │   ├── session.go           Session (one lesson result)
    │   └── profile.go           Profile, KeyStat, Apply, WeakKeys
    ├── engine/
    │   ├── engine.go            TypingState, HandleKey, Backspace, Done
    │   └── entry.go             per-char Entry (expected, typed, status)
    ├── stats/stats.go           Compute snapshot from state + elapsed
    ├── keyboard/
    │   ├── layout.go            Layout interface + Find
    │   └── qwerty.go            US QWERTY ANSI implementation
    ├── lesson/
    │   ├── generator.go         LessonGenerator interface
    │   └── progressive.go       drill-style progressive generator
    ├── storage/
    │   ├── store.go             Store interface
    │   └── jsonstore.go         JSON file impl + paths
    └── ui/
        ├── model.go             Model, Init, Update
        ├── view.go              top-level View, layout, resize guard
        ├── theme.go             finger colors + status styles
        ├── render_keyboard.go   keyboard widget
        ├── render_text.go       typing area
        └── render_stats.go      header/stats/footer
```

### 3. Terminal UI layout
Header (title + active keyset) · Stats bar (WPM/ACC/ERR/CHARS/TIME) ·
Lesson area (target text, per-char color, underlined cursor) ·
Keyboard (always rendered, finger-colored, current key highlit) · Footer (hints).

Resize: `WindowSizeMsg` updates width/height; `View` recomputes layout each frame.
If `width < keyboardWidth` or `height < minHeight`, show a centered
"terminal too small" message. Frame centered with Lip Gloss; sections stacked
with `JoinVertical`; typing text wraps.

### 4. Keyboard design
US ANSI QWERTY, 4 letter rows + spacebar. Each `Key` = {Rune, Label, Finger, Row, Width}.

Finger map (standard touch typing):
- L pinky: `` ` 1 q a z `` + Tab/Caps/Shift
- L ring: `2 w s x`
- L middle: `3 e d c`
- L index: `4 5 r t f g v b`
- R index: `6 7 y u h j n m`
- R middle: `8 i k ,`
- R ring: `9 o l .`
- R pinky: `0 - = p [ ] \ ; ' /` + Enter/Backspace/Shift
- Thumb: space

Colors: 9 distinct Lip Gloss colors in `Theme` (index warm, middle green,
ring blue, pinky purple/magenta, thumb gray). Key states: normal (finger fg),
current (bold inverse highlight), correct flash (green bg ~120ms), incorrect
(expected=yellow bg, pressed=red bg until next keystroke).

Animation: time-based via `flashUntil` + `tea.Tick`; no goroutines (Elm-pure).
Current key always derived from engine state — never cached.

### 5. Data models (internal/domain)
```go
type Finger int // LPinky..RPinky, Thumb
type Key struct { Rune rune; Label string; Finger Finger; Row, Width int }
type Lesson struct { Text string; Keyset []rune }
type Stats struct { WPM, Accuracy float64; Errors, Typed int; Elapsed time.Duration }
type Session struct {
    WPM, Accuracy float64; Errors int; Duration time.Duration
    Keyset []rune; KeyErrors map[rune]int
}
type Profile struct {
    Version int; BestWPM, AvgWPM, Accuracy float64
    PracticeTime time.Duration; LessonsDone, UnlockedLevel int
    KeyStats map[string]KeyStat
}
type KeyStat struct { Presses, Errors int }
func (p *Profile) Apply(s Session)
func (p *Profile) WeakKeys(n int) []rune
```
Profile.Apply owns all aggregation; storage has no business logic.

### 6. Interfaces
```go
type LessonGenerator interface { Next(p domain.Profile) domain.Lesson }
type Store interface { Load() (domain.Profile, error); Save(domain.Profile) error }
type Layout interface { Rows() [][]domain.Key; Find(r rune) (domain.Key, bool); Name() string }
type Clock interface { Now() time.Time }
```

### 7. Application flow
1. Startup: parse flags, build Clock/Store/Layout/Generator/Theme.
2. Load profile (missing→fresh; corrupt→back up + fresh, never crash).
3. generator.Next → first lesson → engine TypingState.
4. tea.NewProgram(model, WithAltScreen); View composes sections.
5. Update: WindowSizeMsg→dims; KeyMsg ctrl+c/esc→save+quit, backspace→Backspace,
   tab→restart, rune→HandleKey; tick→advance flashes.
6. HandleKey records Entry, advances cursor, updates per-key tallies, sets flash.
7. Each frame stats.Compute → Stats rendered live.
8. On Done: build Session, profile.Apply, store.Save, generator.Next → new state.
9. Quit: save, tea.Quit; main restores terminal, prints summary.

### 8. Future extensibility
Adaptive lessons → new LessonGenerator. Layouts → new Layout. Symbols → generator
+ wider keyset. Themes → injected Theme + registry. Vim mode → input strategy in
Update. Multiplayer → transport pkg + remote msgs. Heatmaps/analytics → read-only
pkg over Profile.KeyStats. Plugins → LessonGenerator registry. Streaks → Profile
fields + Apply. Four interfaces + pure domain + versioned profile make all additive.

---

## Phase 3 — Implementation roadmap
Each milestone leaves the app compiling and runnable.

1. Project setup — go.mod, dirs, main.go banner.
2. Domain models — domain package.
3. Keyboard layout — keyboard package (rows + finger map + Find).
4. Bubble Tea skeleton — ui Model/Init/Update/View, static keyboard, quit.
5. Keyboard renderer + finger colors — theme + render_keyboard.
6. Typing engine — engine package; wire into Update; render typed text.
7. Statistics engine — stats package + render_stats; flashes.
8. Lesson generator — progressive drill generator; advance on completion.
9. Persistence — storage package; save on complete/quit; aggregation; weak keys.
10. Polish — resize guard, footer, README, Makefile, tests, complete feedback.

---

## Phase 4 — Architecture review (resolved)
1. ui god-package → split renderers + theme; Model orchestration-only.
2. Time coupling → Clock interface; engine/stats take time.Time as input.
3. Logic in storage → all aggregation in Profile.Apply; storage just marshals.
4. Over-abstraction → rejected Renderer iface + event bus as premature.
5. Redundant state → current key always derived from engine, never cached.
6. Schema evolution → Profile.Version now.
7. Import cycles → strict direction: domain→nothing; others→domain; ui→all; cmd wires.
