# flexgo

A Go library for building flexible, responsive terminal UIs (TUIs) using a flexbox-inspired layout system. flexgo wraps the [Charm](https://charm.sh/) ecosystem ([BubbleTea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss)) and exposes a tree-based layout API modelled after CSS Flexbox.

## Installation

```bash
go get github.com/spik3r/flexgo
```

## Quick Start

```go
package main

import (
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

func main() {
	page := flexgo.VBox(
		&flexgo.Node{Height: 3, Justify: flexgo.JustifyCenter, View: box("HEADER")},
		&flexgo.Node{
			Dir:     flexgo.Row,
			Flex:    1,
			Justify: flexgo.JustifySpaceBetween,
			Children: []*flexgo.Node{
				{Flex: 3, View: box("LEFT")},
				{Flex: 7, View: box("RIGHT")},
			},
		},
		&flexgo.Node{Height: 3, Justify: flexgo.JustifyCenter, View: box("FOOTER")},
	)

	out := page.Render(80, 24)
	fmt.Println(out)
}

func box(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}
```

## Core Concepts

### Node

The `Node` is the fundamental building block. Every element in your UI is a `Node`, including leaf nodes (with `View` functions) and container nodes (with `Children`).

### Direction

`Row` lays out children horizontally; `Col` lays out children vertically.

```go
node.Dir = flexgo.Row  // horizontal layout
node.Dir = flexgo.Col   // vertical layout
```

### Sizing

Three ways to size a node:

```go
// Fixed size (takes precedence over flex)
node.Width = 80
node.Height = 24

// Flex (fills remaining space, 0 = don't stretch)
node.Flex = 1

// Constraints
node.MinWidth = 10
node.MaxWidth = 100
```

When a container has `Flex > 0`, it expands to fill available space from its parent. Children within use their `Flex`, `Width`, or `Height` to determine their size.

### Justify (Main Axis)

Controls where children are placed along the main axis when there is remaining space:

```go
node.Justify = flexgo.JustifyStart       // children packed to start (default)
node.Justify = flexgo.JustifyCenter      // children centered
node.Justify = flexgo.JustifyEnd         // children packed to end
node.Justify = flexgo.JustifySpaceBetween // space distributed between children
```

### Align (Cross Axis)

Controls where children sit on the cross axis (vertical for Row containers, horizontal for Col containers):

```go
node.Align = flexgo.AlignStart   // children at top/left (default)
node.Align = flexgo.AlignCenter  // children centered
node.Align = flexgo.AlignEnd     // children at bottom/right
```

### Spacing

**Gap** inserts space between siblings:

```go
node.Gap = 3  // 3 characters of space between children
```

**Padding** creates inner space between a container's edge and its children:

```go
node.Padding = 2  // 2 characters of padding on each side

// Per-side padding
node.Paddings = flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}
```

**Margin** shrinks a node within its allocated slot:

```go
node.Margin = 1  // 1 character of margin on each side

// Per-side margin
node.Margins = flexgo.Spacing{Top: 1, Bottom: 1}
```

### Colors and Borders

Colors are typed as `color.Color` (use `lipgloss.Color(...)` for ANSI/hex values):

```go
node.Background = lipgloss.Color("237")

node.Border = lipgloss.NormalBorder()
node.BorderForeground = lipgloss.Color("240")
node.BorderBackground = lipgloss.Color("237")
```

### Auto Margins

For centering or right-aligning a fixed-width node within its parent:

```go
node.Width = 80
node.MarginLeftAuto = true
node.MarginRightAuto = true // centered

node.MarginLeftAuto = true
node.MarginRightAuto = false // right-aligned
```

`MarginLeftAuto`/`MarginRightAuto` and `MarginTopAuto`/`MarginBottomAuto`
only work when there is spare space in the allocated axis. They are effective
for roots and in containers that do not exact-size children on that axis.

### AlignSelf

Use `AlignSelf` on a child to override the parent container's `Align` setting:

```go
alignEnd := flexgo.AlignEnd
child.AlignSelf = &alignEnd
```

### Debug Mode

Set `Debug: true` to draw a border around the node for layout debugging:

```go
node.Debug = true
```

Set `Name` to include a label in debug borders, and use helpers for whole trees:

```go
node.Name = "sidebar"
flexgo.DebugAll(root)
flexgo.DebugOff(root)
```

### Convenience Constructors

```go
root := flexgo.VBox(
	flexgo.HBox(left, right),
	flexgo.Spacer(1),
	flexgo.Leaf(viewFn),
)
```

## Builder Pattern

flexgo provides a fluent builder API for constructing node trees:

```go
page := flexgo.NewNode().
    Dir(flexgo.Col).
    Flex(1).
    Children(
        flexgo.NewNode().
            Dir(flexgo.Row).
            Height(3).
            Justify(flexgo.JustifyCenter).
            View(headerView),
        flexgo.NewNode().
            Dir(flexgo.Row).
            Flex(1).
            Justify(flexgo.JustifySpaceBetween).
            Children(
                flexgo.NewNode().Flex(3).View(leftView),
                flexgo.NewNode().Flex(7).View(rightView),
            ),
    ).
    Build()
```

## Integration with BubbleTea

flexgo integrates seamlessly with [BubbleTea](https://github.com/charmbracelet/bubbletea):

```go
type model struct {
    width  int
    height int
    page   *flexgo.Node
    ready  bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.ready = true
    }
    return m, nil
}

func (m model) View() string {
    if !m.ready {
        return "Loading..."
    }
    return m.page.Render(m.width, m.height)
}
```

## Layout Algorithm

The rendering flow:

```
Render(w, h)
  → subtract margin → subtract padding
  → if leaf: call View(w, h)
  → if container: renderChildren()
      → distribute() assigns widths or heights based on Flex
      → each child.Render() called recursively
      → applyAlign() for cross-axis alignment
      → join() for main-axis (Justify) arrangement
  → reapply padding → reapply margin
```

## Examples

Run the examples to see flexgo in action. Full index: `example/README.md`.

```bash
# Static layout example
cd example/basic && go run main.go

# Builder API (basic dashboard)
cd example/builder_basic && go run main.go

# Builder API (AlignSelf + Border)
cd example/builder_alignself && go run main.go

# Dynamic BubbleTea example
cd example/dynamic && go run main.go

# Centered layout example
cd example/centered && go run main.go

# Justify modes demo
cd example/justify && go run main.go

# Align modes demo
cd example/align && go run main.go

# Spacing controls demo
cd example/spacing && go run main.go

# Horizontal auto-margin centering demo
cd example/hautocenter && go run main.go

# Vertical auto-margin centering demo
cd example/vautocenter && go run main.go

# Full center (horizontal + vertical) demo
cd example/centeredLayout && go run main.go
```

## CSS Flexbox Comparison

| CSS Flexbox      | flexgo           | Description                                    |
|-----------------|------------------|------------------------------------------------|
| `flex-direction` | `Dir`            | `Row` = row, `Col` = column                     |
| `flex`          | `Flex`           | Growth factor (0 = fixed size)                  |
| `width/height`  | `Width/Height`   | Fixed dimensions                               |
| `justify-content` | `Justify`       | Main axis alignment                            |
| `align-items`    | `Align`          | Cross axis alignment                           |
| `gap`           | `Gap`            | Space between children                         |
| `padding`       | `Padding`        | Inner spacing                                  |
| `margin`        | `Margin`         | Outer spacing                                  |
| `margin-auto`    | `MarginLeftAuto`  | Horizontal centering/alignment                 |

## Testing

```bash
go test ./...

# Regenerate example golden snapshots
go test -run TestExampleGolden -update
```

## Concurrency

`Render` is safe to call concurrently on the same node tree. Mutating node
fields is not safe concurrently with rendering; synchronize external writes if
your update loop changes nodes while another goroutine calls `Render`.

## License

MIT
