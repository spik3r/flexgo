# Cards that reflow inside their allocated width

`layouts.AutoGrid` complements the [responsive master/detail pattern](responsive.md)
when the content is a collection of peer cards rather than two panes. It chooses
columns from its **content allocation**, including its parent's sidebar, gutters,
padding and border. A nested grid therefore does not need the terminal width or
an application-wide breakpoint.

```go
cells := []func(int, int) string{
    func(w, h int) string { return renderCard("One", w, h) },
    func(w, h int) string { return renderCard("Two", w, h) },
    func(w, h int) string { return renderCard("Three", w, h) },
}
grid := layouts.AutoGrid(14, 1, cells...)
grid.Padding = 1 // column selection uses the width remaining after this inset
```

`renderCard` is your application callback. It must return content bounded to the
supplied terminal-cell width and row count. The runnable examples include a
bounded callback using Lipgloss, with synthetic labels and no external data.

For minimum cell width 14 and gap 1:

| Content width | Columns | Result |
|---|---|---|
| 1–13 | 1 | Narrow fallback; the requested minimum cannot be met |
| 14–28 | 1 | One full-width slot |
| 29–43 | 2 | Two slots, sharing any extra column |
| 44–58 | 3 | Three slots |
| 59+ | Up to the number of cells | More columns as space permits |

In general, `columns × minimum + (columns − 1) × gap` must fit. The implementation
avoids adding `width + gap` so extreme integer arguments cannot overflow column
selection. The column count never exceeds the number of cells. Nonpositive
minimum widths become 1; negative gaps become 0. Empty collections render no
content; a nil callback reserves a blank slot. No cell callback runs when its
assigned width or height is zero.

The existing `Grid` recipe distributes row heights and column widths. Cells stay
in input order from left to right, then top to bottom. A short final row reserves
its missing column slots instead of stretching its remaining cards. For example,
five cards in three columns are arranged as:

```text
One    Two    Three
Four   Five   [blank]
```

## State and builder integration

Create the callbacks once or rebuild them from stable component models. Keep the
selected card ID, input models and focus in the Bubble Tea model, outside the
layout tree. Resize changes allocations; AutoGrid does not select, reorder or
reset components. Live collections should reconcile selection by stable ID when
items are added or removed, as in the master/detail guide.

The callback slice is copied at construction, so replacing an entry in the
caller's slice does not change an existing grid. Captured state still belongs to
the caller and must be synchronized if it is changed while rendering. Pure cell
callbacks may be rendered concurrently on the same grid.

The returned node is a leaf whose `View` creates a temporary Grid for each render.
`Inspect` and `Validate` see this leaf, not the generated rows. Keep its `Children`
empty, and set its ordinary frame fields (`Padding`, `Background`, borders, etc.).
The constructor's gap applies inside the generated grid; setting the returned
leaf's `Gap` does not change it. To use NodeBuilder, reuse the same callback:

```go
node := flexgo.NewNode().
    View(layouts.AutoGrid(14, 1, cells...).View).
    Padding(1).
    Build()
```

This keeps the raw and builder paths on the same reflow algorithm.

## Height and overflow remain explicit

AutoGrid reflows by width. It does not wrap arbitrary Node children, introduce a
measure pass, scroll, paginate, or clip callback output. Rows divide the supplied
height after vertical gutters. If that height cannot hold all gutters, output
can overflow, just as with Grid. If callbacks ignore their supplied dimensions,
their output can also overflow. A minimum cell width is a column-selection target,
not a Node minimum that forces a narrow terminal wider.

Apps should choose an explicit short-terminal policy: a resize hint, pagination,
or a scrolling component that owns its logical content area. The examples show
a bounded hint below 18×17, sufficient for seven one-row cards at the narrowest
usable width. This threshold belongs to the demo, not the library. The grid
reflows inside the same root tree as the terminal grows or shrinks.

## Try and verify

```sh
go run ./example/layouts/autogrid
go run ./example/builder/layouts/autogrid
# Resize the terminal; q or ctrl+c quits, including at tiny sizes.

go test -race -count=1 ./layouts ./example/layouts/autogrid ./example/builder/layouts/autogrid
go test -run '^TestExampleGolden$/^layouts$/^autogrid$' .
go test -run '^TestExampleGolden$/^builder$/^layouts$/^autogrid$' .
```

Tests cover gap-aware column boundaries, nested parent budgets and decoration,
row allocations, final-row blanks, callback ordering, nil/empty collections,
normalization, zero-area callbacks, integer-limit column arithmetic, concurrent
rendering, repeated shrink/grow cycles, tiny-terminal bounds and quitting. Both
example paths have deterministic 80×20 goldens and reuse the same public recipe.
