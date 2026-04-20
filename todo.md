# flexgo — TODO

Snapshot of outstanding work. Everything previously tracked as A1–A7, T1, T2,
T3, T4, T5, B1, B2, P3, P4 has landed.

Legend: 🐞 bug · 🧪 test gap · 🧹 polish · 🎨 feature · 🏛️ architectural

---

## Priority 1 — Safe polish (non-breaking)

### P2. 🧹 ~~`surround(dir, lead, mid, trail)` helper~~ — shipped.
Lives in `render_helpers.go`; used by `join()` for `JustifyCenter`
and `JustifyEnd`.

### P5. 🐞 ~~`layouts.Tabs` byte-length underline~~ — shipped.

### P6. 🧹 ~~`layouts.Form` spacer → `Gap: 1`~~ — shipped.

### P7. 🧹 ~~`layouts.Dashboard` sidebar construction~~ — shipped.

### P8. 🧪 ~~GoDoc `Example*` stubs~~ — shipped; deleted in favour of
golden coverage via `example/layouts/*`.

### P9. 🧹 ~~Uniform `// Customize:` block on each recipe~~ — shipped.

---

## Priority 1.5 — Architectural review findings

### AR1. 🐞 ~~Silent conflicts between co-located fields~~ — shipped.
Field-level GoDoc on `Node` spells out each precedence
(`View`/`Children`, `Padding`/`Paddings`, `Margin`/`Margins`,
`Flex`/`Width`/`Height`, `ShowBorder`/`Debug`). See `node.go`.

### AR2. 🧹 ~~Duplicate main-axis clamp in `renderChildren`~~ — shipped.
Second clamp is now cross-axis only; main-axis clamping stays in
`resolveMainAxisSizes`.

### AR3. 🧹 ~~Extract auto-margin application from `render`~~ — shipped.
`applyAutoMargins` + `splitAuto` helpers in `render.go`.

### AR4. 🐞 ~~`Debug` + explicit `Border` double-reserves a row~~ — shipped.
Policy: explicit border wins; debug wrapper is suppressed when
`ShowBorder` is set. Pinned by `TestDebugPlusExplicitBorderPrefersBorder`.

### AR5. 🧹 ~~`hasExplicitBorder` relies on struct zero-value equality~~ — shipped.
Added `Node.ShowBorder bool` as the explicit opt-in.
`builder.Border()` keeps the old behaviour by auto-setting
`ShowBorder = true`; struct-literal users must set it explicitly.

### AR6. 🎨 ~~`Validate(root) error` + `Inspect(root) string`~~ — shipped.
New file `inspect.go`, tested in `inspect_test.go`.

### AR7. 🏛️ ~~Split `layouts/layouts.go` into file-per-recipe~~ — shipped.
One file per recipe (`dashboard.go`, `form.go`, `grid.go`,
`headerbodyfooter.go`, `modal.go`, `splitpane.go`, `tabs.go`),
shared helpers in `internal.go`, package doc in `doc.go`.

### AR8. 🧹 ~~`NodeBuilder` drift guard~~ — shipped.
Kept the builder (public API), added `TestBuilderCoversAllNodeFields`
so new `Node` fields fail CI until the builder catches up.

### AR9. 🧪 ~~Test gaps~~ — shipped.
`render_edge_test.go` covers ambient-bg propagation, impossible
min constraints (no-panic), auto-margin + fixed size, and the
Debug + Border interaction.

### AR10. 🧹 ~~Docs drift~~ — shipped.
CLAUDE.md architecture table rewritten to match the seven-ish
files that actually exist; README examples index refreshed to
match the real `example/` tree.

### AR11. 🎨 ~~`SplitPaneFlex` variant~~ — shipped.
`layouts/splitpane.go`; percentage `SplitPane` now delegates to it.

### AR12. 🧹 ~~Repeated "full-box centered label" closure~~ — shipped.
`centeredView` helper in `layouts/internal.go`; `Tabs` now uses it.

---

## Priority 2 — Recipes (`flexgo/layouts`)

API foundation is in place (per-side spacing, typed colours,
HBox/VBox/Leaf/Spacer, Border, AlignSelf, Name). Ready to start.

### R1. Package layout: `flexgo/layouts`

Sub-package. Each recipe is a function returning `*flexgo.Node`, with
all dimensions/styles overridable by the caller.

### R2. Initial recipe set (in order of usefulness)

1. ~~**`layouts.Dashboard`**~~ — shipped.
   The htop/k9s/lazygit shape.
2. ~~**`layouts.HeaderBodyFooter`**~~ — shipped.
3. ~~**`layouts.Modal`**~~ — shipped (no dimmed backdrop yet — the
   caller overlays it by swapping trees. A proper backdrop needs
   composite rendering; revisit after X1).
4. ~~**`layouts.Form`**~~ — shipped.
   stacked `Label: Input` rows with aligned
   labels.
5. ~~**`layouts.Tabs`**~~ — shipped.
   tab bar + active panel. Uses `AlignSelf` for
   the active-tab underline.
6. ~~**`layouts.SplitPane`**~~ — shipped.
   fixed-ratio two-pane split.
7. ~~**`layouts.Grid`**~~ — shipped.
   uniform N×M grid.

### R3. Per-recipe deliverables

- GoDoc `Example*` test that renders and golden-checks.
- Matching `example/<recipe>/main.go` with BubbleTea wrapping.
- Short "customization" section in the recipe doc (top 3 overrides).

### R4. Testing

Every recipe lands with a golden file. Non-negotiable — recipes drift
invisibly otherwise.

---

## Priority 3 — Architectural work

### X1. 🏛️ Two-phase layout (measure → arrange → paint)

Splitting `Render` into `Layout(w,h) → laidOut` + `Paint() → string`
unlocks X2 (Wrap), X3 (natural-size `View`), layout caching, and
`Inspect(root)` for debug. Significant refactor — defer until one of
its dependants becomes necessary.

### X2. 🏛️ `Wrap` support

CSS `flex-wrap: wrap`. Needs X1.

### X3. 🏛️ Richer `View` signature

Current `func(w, h int) string` has no way to say "I only need 3
rows." A `Measure(w, h) (int, int)` interface would let `distribute`
honour natural sizes. Needs X1.

### X4. 🏛️ `FlexGrow` / `FlexShrink` / `FlexBasis`

Single `Flex` conflates grow with basis; shrink is undefined. Either
document "no shrink, fixed widths win, overflow is silent" or extend
the model. Decide explicitly; low priority until someone hits it.

---

## Priority 4 — Nice-to-haves / later

- 🎨 **`flexgo.Place()`** — absolute positioning (tooltips, floating
  panels), port of lipgloss `Place`.
- 🎨 **`Overflow`** — `OverflowClip` / `OverflowEllipsis` for text
  that exceeds its box. Currently overflow is silent.
- 🎨 **Percentage sizes** — `Width: "50%"` alt to flex. Requires
  parsing; maybe overkill.
- 🎨 **Theme** — named colours a recipe references; palettes
  swappable without rewriting recipes.
- 🧪 **Benchmark suite** — baseline before X1 so the refactor can
  demonstrate no regression.
- 🧹 ~~**`example/basics/debug_tree`**~~ — shipped. Also demos
  `Inspect(root)`.
- 🧹 ~~**`example/basics/gap_vs_spacebetween`**~~ — shipped.
- 🧹 **`example/bubble_components`** — wrap `textinput`, `viewport`,
  `spinner` as `View` callbacks. Makes the BubbleTea value prop
  concrete.
- 🧹 ~~**`example/basics/min_max`**~~ — shipped.
- 🧹 **`example/responsive`** — breakpoint switching.

---

## Reference app

`demo/scanner/` is a larger skeleton showing how to structure an app
with multiple screens, tabs, scrollable panels, a modal, and a
centralised keymap. Its `README.md` is the architectural tour; copy
the patterns (root-model routing, one-way screen data flow, central
`KeyMap`, viewport state living on the screen) when building apps on
flexgo.

---

## Suggested order of work

1. **X1** — two-phase layout, only when Wrap / richer min-max /
   caching becomes genuinely needed. Priority-3 items (X2/X3/X4)
   unblock behind it.
2. Remaining §4 items are opportunistic — `example/bubble_components`,
   `example/responsive`, the Place/Overflow/Percentage/Theme features,
   and the benchmark suite. Each warrants its own design discussion
   before coding.
