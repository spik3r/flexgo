# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`flexgo` is a Go library for building flexible, responsive terminal UIs (TUIs) using a flexbox-inspired layout system. It wraps the Charm ecosystem (BubbleTea, Lipgloss, Bubbles) and exposes a tree-based layout API modelled after CSS Flexbox.

Module path: `github.com/spik3r/flexgo`

## Commands

```bash
# Run tests (none exist yet)
go test ./...

# Build the module
go build ./...

# Run the static layout example
cd example/basic && go run main.go

# Run the interactive BubbleTea example
cd example/dynamic && go run main.go
```

## Architecture

The library is implemented in four files:

| File | Role |
|------|------|
| `node.go` | Defines the `Node` struct — the single public type. All layout properties live here. |
| `layout.go` | `distribute()` — partitions available width/height among children using flex weights and fixed sizes. |
| `render.go` | `Node.Render(w, h int) string` — the rendering entry point. Recursively renders the node tree, applying margin → padding → content/children → alignment. |
| `align.go` | Cross-axis alignment helpers built on top of lipgloss. |

### Node tree rendering flow

```
Render(w, h)
  → subtract margin → subtract padding
  → if leaf: call View(w, h)
  → if container: renderChildren()
      → distribute() assigns widths or heights
      → each child.Render() called recursively
      → applyAlign() for cross-axis alignment
      → join() for main-axis (Justify) arrangement
  → reapply padding → reapply margin
```

### Key `Node` fields

- `Dir` — `Row` or `Col`
- `Flex` — flex weight (0 = fixed size)
- `Width`, `Height` — fixed dimensions (used when `Flex == 0`)
- `MinWidth`, `MaxWidth`, `MinHeight`, `MaxHeight` — constraints applied after distribution
- `Gap` — space between children
- `Padding`, `Margin` — lipgloss `Style`-based spacing
- `MarginTopAuto`, `MarginBottomAuto`, `MarginLeftAuto`, `MarginRightAuto` — auto-margin support
- `Justify` — main-axis: `JustifyStart`, `JustifyCenter`, `JustifyEnd`, `JustifySpaceBetween`
- `Align` — cross-axis: `AlignStart`, `AlignCenter`, `AlignEnd`
- `View func(w, h int) string` — set on leaf nodes to render content
- `Children []*Node` — set on container nodes
- `Debug bool` — draws a lipgloss border around the node for layout debugging

### Integration with BubbleTea

Leaf nodes use a `View` callback that receives the allocated `(w, h)` at render time. In the dynamic example, component models implement `View(w, h int) string` and the outer BubbleTea model builds the `Node` tree inside its own `View()` method, calling `root.Render(msg.Width, msg.Height)`.
