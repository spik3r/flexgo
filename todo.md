# flexgo — TODO

Snapshot of outstanding work.

Legend: 🐞 bug · 🧪 test gap · 🧹 polish · 🎨 feature · 🏛️ architectural

---

## Priority 1 — Architectural work

### X1. 🏛️ Two-phase layout (measure → arrange → paint)

Splitting `Render` into `Layout(w,h) → laidOut` + `Paint() → string`
unlocks X2 (Wrap), X3 (natural-size `View`), layout caching, and a
fuller `Inspect(root)` that shows actual laid-out sizes. Significant
refactor — defer until one of its dependants becomes necessary.

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

## Priority 2 — Features

- 🎨 **`flexgo.Place()`** — absolute positioning (tooltips, floating
  panels), port of lipgloss `Place`.
- 🎨 **`Overflow`** — `OverflowClip` / `OverflowEllipsis` for text
  that exceeds its box. Currently overflow is silent.
- 🎨 **Percentage sizes** — `Width: "50%"` alt to flex. Requires
  parsing; maybe overkill.
- 🎨 **Theme** — named colours a recipe references; palettes
  swappable without rewriting recipes. The `demo/scanner/widgets.go`
  palette block is a working sketch of what this would look like.

Each of these warrants its own design discussion before coding.

---

## Priority 3 — Examples & infrastructure

- 🧹 **`example/bubble_components`** — wrap `textinput`, `viewport`,
  `spinner` as `View` callbacks. Makes the BubbleTea value prop
  concrete.
- 🧹 **`example/responsive`** — breakpoint switching based on
  terminal size.
- 🧪 **Benchmark suite** — baseline before X1 so the refactor can
  demonstrate no regression.

---

## Reference app

`demo/scanner/` is a larger skeleton showing how to structure an app
with multiple screens, tabs, scrollable panels, a modal, and a
centralised keymap. Its `README.md` is the architectural tour — copy
the patterns (root-model routing, one-way screen data flow, central
`KeyMap`, viewport state living on the screen, `centered` helper for
fixed-size cards) when building apps on flexgo.

---

## Recipe deliverables checklist

Standing checklist for any new `flexgo/layouts/*` recipe added:

- Recipe function returning `*flexgo.Node`, all dimensions/styles
  overridable on the returned node.
- Doc comment with a `// Customize (top 3 overrides):` block.
- `example/layouts/<recipe>/main.go` with BubbleTea wrapping.
- `example/builder/layouts/<recipe>/main.go` with the equivalent
  built via `NodeBuilder`.
- Golden file via `golden_test.go` registration.

---

## Suggested order of work

1. **X1** — two-phase layout, only when Wrap / richer min-max /
   caching becomes genuinely needed. X2/X3 unblock behind it.
2. The Priority 2 features and Priority 3 examples are opportunistic
   — pick up as time and need allow.
