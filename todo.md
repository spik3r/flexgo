# flexgo — TODO

Snapshot of outstanding work. Everything previously tracked as A1–A7, T1, T2,
T3, T4, T5, B1, B2, P3, P4 has landed.

Legend: 🐞 bug · 🧪 test gap · 🧹 polish · 🎨 feature · 🏛️ architectural

---

## Priority 1 — Safe polish (non-breaking)

### P2. 🧹 `surround(dir, lead, mid, trail)` helper

**File:** `render_helpers.go`

Several call sites build a 3-element slice just to pass to `concat`.
A `surround` helper reads cleaner. Micro-win on the next touch of
that file.

---

## Priority 2 — Recipes (`flexgo/layouts`)

API foundation is in place (per-side spacing, typed colours,
HBox/VBox/Leaf/Spacer, Border, AlignSelf, Name). Ready to start.

### R1. Package layout: `flexgo/layouts`

Sub-package. Each recipe is a function returning `*flexgo.Node`, with
all dimensions/styles overridable by the caller.

### R2. Initial recipe set (in order of usefulness)

1. **`layouts.Dashboard`** — sidebar + header + main + status bar.
   The htop/k9s/lazygit shape.
2. ~~**`layouts.HeaderBodyFooter`**~~ — shipped.
3. ~~**`layouts.Modal`**~~ — shipped (no dimmed backdrop yet — the
   caller overlays it by swapping trees. A proper backdrop needs
   composite rendering; revisit after X1).
4. **`layouts.Form`** — stacked `Label: Input` rows with aligned
   labels.
5. **`layouts.Tabs`** — tab bar + active panel. Uses `AlignSelf` for
   the active-tab underline.
6. **`layouts.SplitPane`** — fixed-ratio two-pane split.
7. **`layouts.Grid`** — uniform N×M grid.

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
- 🧹 **`example/debug_tree`** — deep tree with `Debug: true` on every
  container; visual reference for the debug mode.
- 🧹 **`example/gap_vs_spacebetween`** — side-by-side. Most common
  flexbox confusion; worth a dedicated demo.
- 🧹 **`example/bubble_components`** — wrap `textinput`, `viewport`,
  `spinner` as `View` callbacks. Makes the BubbleTea value prop
  concrete.
- 🧹 **`example/min_max`** — exercises `MinWidth`/`MaxWidth`.
- 🧹 **`example/responsive`** — breakpoint switching.

---

## Suggested order of work

1. **P2** — `surround` helper, low-risk polish.
2. **§2 Recipes** — start with `Dashboard` and `HeaderBodyFooter`.
3. **X1** — two-phase layout, only when Wrap / richer min-max /
   caching becomes genuinely needed.

§4 is opportunistic — pick up as time allows.
