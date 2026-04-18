package flexgo

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

type Node struct {
	Dir      Direction
	Children []*Node

	// sizing
	Flex   int
	Width  int
	Height int

	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int

	// spacing
	Gap int

	Padding int
	Margin  int

	// Margin
	MarginTopAuto    bool
	MarginBottomAuto bool
	MarginLeftAuto   bool
	MarginRightAuto  bool

	// alignment
	Justify Justify
	Align   Align

	// debug
	Debug bool

	// content
	View func(w, h int) string
}
