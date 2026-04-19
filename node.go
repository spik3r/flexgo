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

type Node struct {
	Dir      Direction
	Children []*Node

	Flex   int
	Width  int
	Height int

	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int

	Gap int

	Padding int
	Margin  int

	Paddings Spacing
	Margins  Spacing

	MarginTopAuto    bool
	MarginBottomAuto bool
	MarginLeftAuto   bool
	MarginRightAuto  bool

	Justify Justify
	Align   Align

	AlignSelf *Align

	Debug bool
	Name  string

	Border           lipgloss.Border
	BorderForeground color.Color
	BorderBackground color.Color

	Background color.Color

	View func(w, h int) string
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
