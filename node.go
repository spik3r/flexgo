package flexgo

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Direction int

const (
	Row Direction = iota
	Col
)

type Justify int

const (
	JustifyStart Justify = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
)

type Align int

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
)

type Spacing struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func (s Spacing) TotalWidth() int {
	return s.Left + s.Right
}

func (s Spacing) TotalHeight() int {
	return s.Top + s.Bottom
}

// Node is the single public type in flexgo. Every element in the UI —
// leaf or container — is a Node. Fields are grouped below by concern;
// several fields interact via silent precedence rules documented inline.
//
// Precedence summary (callers should avoid setting both at once):
//   - View vs Children: a node with Children is a container and View is
//     ignored. A node with neither is a blank-space Spacer.
//   - Flex vs Width/Height: the fixed dimension wins on its axis.
//   - Padding vs Paddings: if Paddings is non-zero, it is used and the
//     Padding shorthand is ignored. Same for Margin vs Margins.
//   - ShowBorder + Border: ShowBorder=true applies Border (or the
//     lipgloss zero Border, effectively no visible frame). Debug adds
//     a debug wrapper only when ShowBorder is false.
type Node struct {
	Dir      Direction
	Children []*Node // container: children render inside; View is ignored.

	Flex   int // flex weight on the main axis; 0 = size by Width/Height.
	Width  int // fixed; takes precedence over Flex on the horizontal axis.
	Height int // fixed; takes precedence over Flex on the vertical axis.

	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int

	Gap int // main-axis space inserted between siblings.

	Padding int     // uniform inner padding shorthand; shadowed by Paddings.
	Margin  int     // uniform outer margin shorthand; shadowed by Margins.
	Paddings Spacing // per-side padding; non-zero wins over Padding.
	Margins  Spacing // per-side margin; non-zero wins over Margin.

	MarginTopAuto    bool
	MarginBottomAuto bool
	MarginLeftAuto   bool
	MarginRightAuto  bool

	Justify Justify // main-axis distribution.
	Align   Align   // cross-axis alignment applied to each child.

	AlignSelf *Align // per-child override of the parent's Align.

	Debug bool   // draw a debug wrapper if no explicit border is set.
	Name  string // label used inside the debug wrapper.

	ShowBorder       bool // true = render Border around this node.
	Border           lipgloss.Border
	BorderForeground color.Color
	BorderBackground color.Color

	Background color.Color

	View func(w, h int) string // leaf renderer; ignored when Children is non-empty.
}

func HBox(children ...*Node) *Node {
	return &Node{Dir: Row, Children: children}
}

func VBox(children ...*Node) *Node {
	return &Node{Dir: Col, Children: children}
}

func Leaf(view func(w, h int) string) *Node {
	return &Node{View: view}
}

func Spacer(flex int) *Node {
	return &Node{Flex: flex}
}

func DebugAll(root *Node) {
	if root == nil {
		return
	}
	root.Debug = true
	for _, child := range root.Children {
		DebugAll(child)
	}
}

func DebugOff(root *Node) {
	if root == nil {
		return
	}
	root.Debug = false
	for _, child := range root.Children {
		DebugOff(child)
	}
}
