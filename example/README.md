# Examples

Organised by what the example teaches. Each subdirectory is a runnable
`main` package.

## `basics/` — static layout primitives

Small one-screen demos exercising a single Node property.

| Path | Shows |
|------|-------|
| `basics/basic` | Canonical dashboard-style Row/Col composition |
| `basics/centered` | Fixed-width page, centred horizontally |
| `basics/justify` | All four `Justify*` modes side by side |
| `basics/align` | All three `Align*` modes |
| `basics/spacing` | `Gap`, `Padding`, `Margin` interactions |

## `margins/` — auto-margin centring

Dedicated demos of the `MarginTopAuto` / `MarginBottomAuto` /
`MarginLeftAuto` / `MarginRightAuto` flags.

| Path | Shows |
|------|-------|
| `margins/hautocenter` | Horizontal auto-margins |
| `margins/vautocenter` | Vertical auto-margins |
| `margins/centered_layout` | Full 2-axis centring |

## `builder/` — fluent builder API

| Path | Shows |
|------|-------|
| `builder/basic` | Same shape as `basics/basic`, via `NewNode()` |
| `builder/alignself` | Per-child `AlignSelf` + first-class `Border` |

## `layouts/` — recipes from `github.com/spik3r/flexgo/layouts`

| Path | Shows |
|------|-------|
| `layouts/headerbodyfooter` | `layouts.HeaderBodyFooter` 3-row shape |
| `layouts/modal` | Scrollable body (j/k) + `layouts.Modal` overlay (space opens, esc closes) |

## `dynamic/`

Flagship interactive demo — reactive BubbleTea app using the raw API.

## Running

```bash
go run ./example/basics/basic
go run ./example/layouts/modal
go run ./example/dynamic
```

## Golden snapshot mode

Every example supports a non-interactive mode used by `golden_test.go`:

```bash
FLEXGO_GOLDEN=1 go run ./example/basics/basic
```
