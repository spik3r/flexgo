# demo/scanner — structuring a larger flexgo app

This demo is a reference skeleton for a TUI that has outgrown a single
screen. It's not a working code scanner; it's a layout and state
architecture you can copy when your own app starts acquiring screens,
modals, tabs, and scrollable panels.

Run it:

```bash
cd demo/scanner && go run .
```

Keys: `1`/`2`/`3` switch screens, `?` opens the keymap modal,
`tab`/`shift+tab` cycles tabs on the Scan screen, `j`/`k` scroll,
`enter` opens a file in the Files tab, `esc` closes a modal or the
file viewer, `q` quits.

---

## What's here

| File | Role |
|------|------|
| `main.go` | Entry point — constructs `App`, hands it to BubbleTea. |
| `app.go` | Root model. Owns screen models, keymap, modal flag; does all routing. |
| `keys.go` | Central `KeyMap` — every binding lives here, nowhere else. |
| `scan.go` | Scan screen (header + sidebar + tabbed main + footer). |
| `launcher.go` | Launcher / profile-editor screen. |
| `history.go` | History screen (sidebar of runs + detail viewport). |
| `widgets.go` | Shared UI primitives: viewport, filetree, header/footer, keymap modal, formatters. |
| `state.go` | Domain types + sample data stand-ins. |

Three screens, four tabs on the Scan screen, one modal — and the app
is ~800 lines total. Most of that is formatting strings, not routing.

---

## 1. Layout & page structure

### One page shell, interchangeable bodies

Every screen exposes four things:

```go
Title() string               // "Scanner"
Subtitle() string            // "profile: default"
Footer() string              // screen-specific hint
Body(w, h int) *flexgo.Node  // the middle chunk
```

`App.View` feeds those into one `pageShell(subtitle, active, body,
screenHint)` helper in `widgets.go`. The shell is a `VBox` with seven
rows of stacked chrome and a flex body between them:

```
┌────────────── colAccent (3 rows) ─────────────┐
│                                               │
│                 ⬢  FlexGo                     │   brand title, centred both axes
│                                               │
├────────────── colPanel (1 row) ───────────────┤
│  Scanner · profile: default   1 Launcher  ·   │   flex-between: subtitle + indicator
│                               ● 2 Scan · 3 …  │
├────────────── colBg (1 row) ──────────────────┤
│  ───────────────────────────────────────────  │   divider
│                                               │
│           (body with 2-col left/right margin) │
│                                               │
│                                               │
├────────────── colBg (1 row) ──────────────────┤
│  progress  ·  j/k scroll  ·  tab next         │   screen-specific hint
├────────────── colPanel (1 row) ───────────────┤
│        ? help  ·  1/2/3 screens  ·  q quit    │   global hint, centred
└───────────────────────────────────────────────┘
```

The screen indicator on the right of the subheader mirrors CSS
`justify-content: space-between` — rendered via an explicit bg-tinted
gap between the subtitle and the 1/2/3 entries, so the active entry
stays flush with the right edge regardless of terminal width. The
active screen is bolded + in accent; the other two are muted.

All keymap hints live in the footer — there's nothing keyboard-y in
the title bar. Screens give their own hint via `Footer()`; the global
hint (`? help · 1/2/3 screens · q quit`) is hardcoded in pageShell and
the same on every screen.

The body can be anything:

- **Split view** (Scan, History) — `Row(sidebar, main)` with a `Gap`
  between them so the panels don't touch.
- **Centred card** (Launcher, keymap modal) — a fixed-size child on a
  `Flex:1, Dir:Row, Justify:JustifyCenter, Align:AlignCenter`
  container. See the centring gotcha below.
- **Full-width viewport** (Scan's opened-file sub-view) — one
  `Flex:1` node with a `View` callback.

#### Centring gotcha: explicit Flex-spacers, not auto-margins

For a card with a fixed `Width`, `MarginLeftAuto + MarginRightAuto`
*looks* like the right answer but produces nothing visible. flexgo
allocates the card a slot exactly its `Width` wide, so there's no
"spare space" for the auto-margin to consume. The card stays flush
left.

`Justify:JustifyCenter + Align:AlignCenter` on the parent is closer
to right — flexgo inserts `spacer` nodes on the sides — but there are
two subtle failure modes:

1. `JustifyCenter` splits the spare evenly (`remaining/2` each side);
   when the remaining count is odd the output is 1 cell short, and
   lipgloss's outer Width-pad doesn't always fill that last cell with
   the bg when the content carries embedded ANSI resets.
2. `Node.Paddings` goes through lipgloss padding, which has the same
   transparent-cell problem on multi-line styled content.

The approach this demo uses — in the `centered` helper in
`widgets.go` — is to build explicit `Flex:1` spacer nodes with their
own `Background` and a `solidView` that emits a fully bg-painted
rectangle. Every cell is its own painted leaf; there's nothing for
the terminal to leave unpainted.

```go
func centered(card *flexgo.Node, cardHeight int) *flexgo.Node {
    middle := &flexgo.Node{
        Height:     cardHeight,
        Dir:        flexgo.Row,
        Background: colBg,
        Children: []*flexgo.Node{
            {Flex: 1, Background: colBg, View: solidView(colBg)},
            card,
            {Flex: 1, Background: colBg, View: solidView(colBg)},
        },
    }
    return &flexgo.Node{
        Flex:       1,
        Dir:        flexgo.Col,
        Background: colBg,
        Children: []*flexgo.Node{
            {Flex: 1, Background: colBg, View: solidView(colBg)},
            middle,
            {Flex: 1, Background: colBg, View: solidView(colBg)},
        },
    }
}
```

The Launcher profile card and the keymap modal both centre via this
helper.

Modal overlay reuses this: when the keymap modal is open, `App.View`
swaps only the body for `buildKeymapCard(keys)` — the header and
footer stay visible and the user sees they're still inside the same
screen. No more "replace the whole tree" stunt.

### Palette

The palette block at the top of `widgets.go` is the Tokyo Night Storm
theme (#24283b bg, #7aa2f7 accent blue, etc) expressed as
`lipgloss.Color("#hex")` values. Every colour the demo uses is a
named constant there, so swapping themes is a one-file change.

### Every cell must have a defined background

Terminals (and tmux in particular) leave unpainted cells transparent,
so a hole in one line of output shows whatever was on screen before
— your wallpaper, previous tmux content, whatever. Two rules keep
this from happening:

1. Set `Background` on every leaf `*flexgo.Node`. flexgo pads around
   the content with that colour, so the edges are covered.
2. When your `View` callback paints styled text (with its own
   foreground/background), make sure the same bg flows through the
   gap cells between styled segments — lipgloss's `[m` resets become
   transparent. Use the `padLineBg` helper, build segments with the
   bg baked in, or let flexgo's Node.Background absorb the padding.

`widgets.go` has `padLineBg`, `panelPaint`, `textStyle`, and
`headingStyle` so screens don't re-derive these rules. When a screen
renders its own content, it uses them; when it delegates to a
viewport or filetree helper, those already do it.

### Rebuild every frame

Every frame rebuilds the whole tree. Don't try to cache nodes across
frames — the tree is cheap (a few hundred node structs); caching
invites stale-state bugs and makes transitions awkward.

### Reach for `flexgo/layouts` when the shape fits

`layouts.HeaderBodyFooter`, `layouts.Modal`, etc. are there for the
common shapes. Use them as starting points, override fields after
construction for customisation — every recipe lists the top three
useful overrides in its doc comment.

When none fit, drop to struct literals. The `Scan` screen builds its
own tab strip directly rather than using `layouts.Tabs` because it
wants a different active-tab style — recipes are a shortcut, not a
straitjacket.

### Scrolling and overflow

flexgo has no built-in scrolling (overflow is silent — see todo).
Every scrollable panel here uses a small `ViewportState` (offset
only) + the `viewportView` helper in `widgets.go`. For a production
app, swap that stub for `charm.land/bubbles/v2/viewport`, which
handles word wrap, horizontal scroll, and mouse. The state/key flow
stays identical — only the paint function changes.

### Modal overlay

flexgo can't composite today (it's on the backlog as X1: two-phase
layout). "Overlay" here means *body swap*: the page shell's header
and footer keep rendering, but the body becomes the modal card. The
screen content behind the modal is hidden, not dimmed — good enough
for a reference modal like a keymap, not good enough for a tooltip
that needs to sit over live content. That upgrade waits on X1.

---

## 2. State management

### One root model. Screen models underneath. No globals.

```
App
├── keys       KeyMap
├── screen     Screen        // which one is active
├── modalOpen  bool
├── launcher   LauncherScreen
├── scan       ScanScreen
└── history    HistoryScreen
```

`App` is the single source of truth for "which screen", "is a modal
open", "what keys exist". Each screen model owns its own state:
ScanScreen holds four viewport offsets, a file cursor, the currently
opened file; HistoryScreen holds the run cursor and detail offset.

### Screen-to-screen data flow is one-way

When you press `2` to jump from Launcher to Scan, `App.Update` reads
the profile from the launcher and constructs a **new** `ScanScreen`
with it. The scan never reaches back into the launcher; the launcher
never knows it was consumed.

```go
case Matches(key, a.keys.Scan):
    a.scan = NewScanScreen(a.launcher.Profile())
    a.screen = ScreenScan
```

This is the same pattern you want for history → scan ("re-run this
profile") or any other transition. Rebuilding a screen's state at
transition time is cheap and makes state ownership unambiguous.

### Screens are plain structs with value receivers

`func (s ScanScreen) Update(...) ScanScreen` — not pointer receivers.
This keeps updates functional: the parent stores the new value. It
also means screens are trivially snapshotable if you ever want
undo/history or time-travel debugging.

The cost is allocation on every keypress. If a screen grows to where
this matters, switch it to a pointer receiver in isolation — the root
model only cares that `Update` returns the new state.

### Side effects go through `tea.Cmd`

Everything expensive (reading a file, running a subprocess, waiting
for a scan worker) returns a `tea.Cmd`, never blocks in `Update`.
The demo stubs this — `fakeFileBody` synthesises content inline — but
the comment in `scan.go` points at the real pattern:

```go
// Production:
case Matches(key, keys.Open):
    return s, tea.Cmd(func() tea.Msg {
        body, err := os.ReadFile(path)
        return fileOpenedMsg{path: path, body: body, err: err}
    })

// ...later, in Update, handle fileOpenedMsg and set openedPath.
```

---

## 3. Key handling

### One `KeyMap`, consulted by everyone

`keys.go` defines every binding in the app. Screens never compare
against literal strings like `"j"` — they call `Matches(key,
keys.Down)`. Three wins from this:

- **Rebind once.** Swap `Down: ["j", "down"]` for `Down: ["down"]`
  and every scrollable panel adopts it.
- **Help is auto-generated.** The keymap modal reads `HelpEntries()`
  off the same struct; no drift between "what works" and "what the
  help says works".
- **Contextual shadowing.** When the file viewer is open it rebinds
  `Close` to "back out of the file", but `Up`/`Down` still scroll.
  The keys are the same; the meaning differs by context. The map
  makes that explicit.

### Dispatch order: modal → global → screen

`App.Update` hard-codes this order:

```
1. Modal open?  →  only Close/Help escape; everything else ignored.
2. Global key?  →  quit, help toggle, screen switch — always work.
3. Screen key   →  active screen's Update gets the keypress.
```

Two rules for picking which bucket a new binding lives in:

- **Global** if pressing it on any screen should do the same thing
  (quit, help, screen switch, maybe a "command palette").
- **Screen** if its meaning changes by screen (scroll, select, tab
  switch, open file).

Don't let screens re-implement global actions. If the Scan screen
had its own `q` → quit, rebinding it would leave other screens
out of sync. Everything that must stay consistent across screens
belongs in the global block.

### Screens return themselves, they don't emit commands for UI state

`ScanScreen.Update` signature is
`func (ScanScreen, tea.KeyPressMsg, KeyMap) ScanScreen` — no `tea.Cmd`
return. For UI state (scroll, cursor, tab) this is right: it's
synchronous, nothing to wait on. Reserve `tea.Cmd` for genuine side
effects (file I/O, network, timers). Screens that need both can return
`(ScanScreen, tea.Cmd)`; keep the Cmd for effects only.

---

## 4. Scaling tips

### When to split a screen's file

Split when a screen grows:

- a sub-screen with its own keymap shadowing (the file viewer on the
  Scan screen is about to cross this line — if it gains search, syntax
  highlighting, and navigation bindings, pull it into `scan_file.go`);
- a distinct background lifecycle (worker fan-out, streaming content);
- substantial non-UI logic (diffing, filtering, sorting).

Don't pre-split. `scan.go` at 200 lines is fine; split at 500-ish or
when two concerns in the file start needing different test setups.

### When to introduce a `components/` subpackage

Right now `widgets.go` is all plain functions in `main`. Promote to
a `components` package when:

- another binary (a second demo, a real app) wants the same widgets;
- the widgets need their own tests that don't pull in the whole app.

Until then, resist the abstraction tax. In-package functions are one
of Go's strengths.

### When to reach for bubbles

The stubs here (`viewportView`, `filetreeView`) cover the flow but
not the polish. Reach for [bubbles](https://github.com/charmbracelet/bubbles)
when you need:

- proper text input (`textinput`, `textarea`),
- real word-wrapping scroll (`viewport`),
- a production filetree (`list` + custom item delegate, or a
  community filetree),
- spinners, progress bars with animation.

flexgo renders the box; bubbles fills it. The `View` callback on a
flexgo `Node` is where bubbles models plug in — the outer BubbleTea
model stores the bubble component, the `View` callback calls its
`View()` with the `(w, h)` flexgo allocated. Already demonstrated in
`example/dynamic/`.

### When your app outgrows this shape

Two signals the single-root pattern is straining:

- **Screen coupling.** Two screens need to mutate the same piece of
  state (the current scan is shared between Scan and History, and
  both want to react to updates). Lift the shared state up — put it
  on `App` directly, pass read-only views to screens.
- **Command fan-out.** Updates start returning tea.Cmds that other
  screens also want to react to. Introduce a thin message bus on
  `App` that broadcasts to all screens. Don't reach for a reactive
  framework; a `switch` over typed messages handles this cleanly
  for a long time.

Only then consider splitting into packages. A TUI this complex
stays in one package longer than most people expect.
