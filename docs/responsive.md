# Responsive master/detail screens

Run `go run ./example/responsive`. It uses synthetic local records: no scanner,
network access, files or background commands. Resize your terminal while a
record is selected and its details are scrolled.

| Area | Behaviour |
|---|---|
| At least 72 columns and 6 rows | 24-column list, one-column gap, flexible details |
| 24–71 columns and at least 6 rows | One pane at a time: list or details |
| Under 24 columns or 6 rows | Bounded resize hint; navigation pauses, quit still works |
| Zero-sized terminal | Empty output until a usable size arrives |

The thresholds belong to this example, **not** the library. A different app
should choose them from its content and minimum useful controls.

## Try this two-minute walkthrough

1. At 80×24, move to record 10 with Down, then Enter (or Tab).
2. Scroll its details with Down. Resize to 40×10: the same record and detail
   position remain selected, now as a full-width pane.
3. Press Esc to return to the list. Record 10 is still selected and visible.
   Enter reopens the detail at its remembered position.
4. Shrink below 24×6, then grow back. The resize hint disappears and navigation
   resumes without losing selection or scroll intent.
5. Select another record: its detail starts at the top. Press `q` to quit.

`Up`/`Down` (or `k`/`j`) move the selection when the list is active, or scroll
when details are active. `Tab` toggles active pane, `Enter` selects details,
`Esc` selects the list, and `q`/`ctrl+c` quit. The `*` marks the active pane.

## State survives; the layout changes

```text
WindowSizeMsg ──► dimensions only ──► choose wide / compact / tiny
                                             │
selection + active pane + scroll intent ──────┤
                                             ▼
                                  rebuild a flexgo Node tree
                                             │
                                  Render(terminal dimensions)
```

Think of the terminal as a desk. On a large desk, keep the index and the open
book side by side. On a small desk, show one at a time. Moving the index off the
desk should not erase your bookmark.

The example stores the bookmark (`selected`, `detail`, `offset`) in the Bubble
Tea model. The `WindowSizeMsg` handler only updates dimensions. `View` builds
new nodes and callbacks from that state; it does not initialise a new model,
move selection, or write back a viewport offset.

The visible scroll offset is clamped against the current viewport **during
rendering**. A large viewport may show earlier content so it has no blank tail;
shrinking again restores the previous scroll intent. The next scroll key starts
from the currently visible offset. Changing records intentionally resets it.
The list's visible start is derived from selection and available height, so
selection stays on screen after a shrink without a second persistent list offset.

## Applying this to generated or hand-written apps

The same pattern works for a Foundry-generated record browser or a Monstera
screen without moving application state into flexgo:

1. Keep stable record IDs, selected record, focus and scroll state in the app.
   This static demo uses indices because its records never reorder; live results
   should reconcile IDs when records change and handle removed selections.
2. Receive the latest terminal size in the app's update loop.
3. Choose a composition using that size. Include mandatory navigation at every
   usable breakpoint. Do not just hide a panel that owns the only Back control.
4. Build the tree from callbacks. The wide composition reserves sidebar width
   only when enough space remains; the compact composition has no hidden fixed
   sidebar consuming width.
5. Bound leaf output to the callback's supplied width and height. Flexgo does
   **not** automatically wrap, scroll or clip overconstrained content. `fit`
   joins logical records and uses `flexgo.ClipView` to bound terminal-cell width
   and line count, including embedded newlines.
   The example clips horizontally rather than wrapping; apps needing wrapped
   prose should use a viewport whose line mapping is recomputed on resize.
6. Keep IO in `tea.Cmd`, not `View`. Store component models outside rebuilt
   layout nodes so interrupted or repeated navigation does not reset them.

The example owns its breakpoints and state; `ClipView` handles leaf text only.
It does not introduce a shrink algorithm or global breakpoint policy.

## Verification

```sh
go test -race -count=1 ./example/responsive
go test -run 'TestExampleGolden/responsive$' .
FLEXGO_GOLDEN=1 go run ./example/responsive
```

The tests cover breakpoint boundaries, zero/negative/tiny sizes, output bounds,
Unicode cell clipping, wide → compact → tiny → wide state retention, repeated
Back/Enter/Tab flows, empty data and scroll/selection limits. The deterministic
80×12 golden covers the wide screen without opening a TTY.

The synthetic data and allocation limits are deliberate. This is not an input
sanitiser for arbitrary control sequences, a mouse/hit-test demo, or a promise
that an arbitrary Node tree will fit every terminal. Use existing application
sanitisation and interaction components when adapting it to external data.
