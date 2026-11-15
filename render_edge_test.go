package flexgo

import (
	"reflect"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// When a parent sets Background, the ambient bg should propagate to a
// child's gap/margin painting even when that child has no bg of its own.
func TestAmbientBackgroundPropagatesThroughContainers(t *testing.T) {
	parentBg := lipgloss.Color("57")
	root := &Node{
		Dir:        Col,
		Background: parentBg,
		Gap:        1,
		Children: []*Node{
			{Height: 1, View: func(w, h int) string { return lipgloss.NewStyle().Width(w).Render("a") }},
			{Height: 1, View: func(w, h int) string { return lipgloss.NewStyle().Width(w).Render("b") }},
		},
	}
	out := root.Render(4, 3)
	if !containsAnsi(out, "48;5;57") && !containsAnsi(out, "48:5:57") {
		t.Fatalf("expected ambient background ANSI 57 somewhere in output, got %q", out)
	}
}

// Two children with min constraints that can't both fit — library
// should produce some output without panic and honour mins on both.
func TestImpossibleMinConstraintsDoNotPanic(t *testing.T) {
	root := &Node{
		Dir: Row,
		Children: []*Node{
			{Flex: 1, MinWidth: 30, View: func(w, h int) string { return "" }},
			{Flex: 1, MinWidth: 30, View: func(w, h int) string { return "" }},
		},
	}
	out := root.Render(20, 3)
	if out == "" {
		t.Fatalf("expected output, got empty string")
	}
}

// Auto-margins only engage when the allocated size exceeds the
// rendered size; a child that already fills its slot must not shrink.
func TestAutoMarginWithFixedSizeOnlyUsesSpareSpace(t *testing.T) {
	// Node with fixed width exactly the container width — no spare space.
	node := &Node{
		Width:           10,
		Height:          1,
		MarginLeftAuto:  true,
		MarginRightAuto: true,
		View:            func(w, h int) string { return lipgloss.NewStyle().Width(w).Render("hi") },
	}
	out := node.Render(10, 1)
	if w := lipgloss.Width(out); w != 10 {
		t.Fatalf("no spare space: expected width 10, got %d", w)
	}

	// Same node with generous container — should centre within 20.
	out = node.Render(20, 1)
	if w := lipgloss.Width(out); w != 20 {
		t.Fatalf("with spare: expected width 20, got %d", w)
	}
}

// Debug and an explicit border set together: explicit border wins,
// debug wrapper is suppressed. Output must fit inside (w, h) and not
// leak a stray row from the debug label reservation.
func TestDebugPlusExplicitBorderPrefersBorder(t *testing.T) {
	root := &Node{
		Debug:      true,
		ShowBorder: true,
		Border:     lipgloss.NormalBorder(),
		Name:       "panel",
		View:       func(w, h int) string { return lipgloss.NewStyle().Width(w).Height(h).Render("x") },
	}
	out := root.Render(10, 5)
	if h := lipgloss.Height(out); h != 5 {
		t.Fatalf("expected height 5, got %d", h)
	}
	if strings.Contains(out, "panel") {
		t.Fatalf("debug label leaked when explicit border is set: %q", out)
	}
}

// containsAnsi does a loose substring match over the full output; both
// lipgloss styles happen to embed colour escapes with varying separators.
func containsAnsi(s, frag string) bool {
	return strings.Contains(s, frag)
}

func TestCustomBorderPaintedWidth(t *testing.T) {
	wide := lipgloss.NormalBorder()
	wide.Left, wide.Right = "界", "界"
	asymmetric := wide
	asymmetric.Right = "|"
	pattern := wide
	pattern.Left, pattern.Right = "a界", "|界"
	zeroWidth := lipgloss.NormalBorder()
	zeroWidth.Left, zeroWidth.Right = "\u0301", "\u0301"
	emptySide := zeroWidth
	emptySide.Left = ""
	wideHorizontal := lipgloss.NormalBorder()
	wideHorizontal.Top, wideHorizontal.Bottom = "界", "界"
	wideCorners := lipgloss.NormalBorder()
	wideCorners.TopLeft, wideCorners.TopRight = "界", "界"
	wideCorners.BottomLeft, wideCorners.BottomRight = "界", "界"
	cases := []struct {
		name                                       string
		node                                       Node
		width, height, contentWidth, contentHeight int
	}{
		{"wide sides", Node{ShowBorder: true, Border: wide}, 10, 5, 6, 3},
		{"asymmetric sides", Node{ShowBorder: true, Border: asymmetric}, 10, 5, 7, 3},
		{"alternating side runes", Node{ShowBorder: true, Border: pattern}, 10, 5, 6, 3},
		{"nonempty zero width sides", Node{ShowBorder: true, Border: zeroWidth}, 10, 5, 10, 3},
		{"empty side fallback", Node{ShowBorder: true, Border: emptySide}, 10, 5, 9, 3},
		{"top only fallback", Node{ShowBorder: true, Border: lipgloss.Border{Top: "-"}}, 10, 5, 8, 3},
		{"normal", Node{ShowBorder: true, Border: lipgloss.NormalBorder()}, 10, 5, 8, 3},
		{"hidden", Node{ShowBorder: true, Border: lipgloss.HiddenBorder()}, 10, 5, 8, 3},
		{"wide horizontal runes", Node{ShowBorder: true, Border: wideHorizontal}, 10, 5, 8, 3},
		{"wide corners", Node{ShowBorder: true, Border: wideCorners}, 10, 5, 8, 3},
		{"disabled", Node{Border: wide}, 10, 5, 10, 5},
		{"zero explicit suppresses debug", Node{ShowBorder: true, Debug: true}, 10, 5, 10, 5},
		{"explicit suppresses debug", Node{ShowBorder: true, Border: wide, Debug: true}, 10, 5, 6, 3},
		{"debug only", Node{Debug: true}, 10, 5, 8, 2},
		{"padding and margin", Node{ShowBorder: true, Border: wide, Padding: 1, Margin: 1}, 14, 9, 6, 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.node
			node.Name = "b"
			node.Foreground = lipgloss.Color("7")
			node.Background = lipgloss.Color("0")
			node.BorderForeground = lipgloss.Color("4")
			node.BorderBackground = lipgloss.Color("0")
			before := node
			calls := 0
			node.View = func(w, h int) string {
				calls++
				if w != tc.contentWidth || h != tc.contentHeight {
					t.Errorf("callback %dx%d, want %dx%d", w, h, tc.contentWidth, tc.contentHeight)
				}
				return customBorderBlock(w, h)
			}
			if err := Validate(&node); err != nil {
				t.Fatal(err)
			}
			first := node.Render(tc.width, tc.height)
			second := node.Render(tc.width, tc.height)
			if first != second || calls != 2 {
				t.Fatalf("repeat changed output or callback count: %d", calls)
			}
			if w, h := lipgloss.Width(first), lipgloss.Height(first); w != tc.width || h != tc.height {
				t.Fatalf("output %dx%d, want %dx%d: %q", w, h, tc.width, tc.height, first)
			}
			node.View = nil
			if !reflect.DeepEqual(node, before) {
				t.Fatal("render mutated node styles")
			}
		})
	}
}

func TestCustomBorderNestedAllocation(t *testing.T) {
	for _, dir := range []Direction{Row, Col} {
		t.Run(map[Direction]string{Row: "row", Col: "column"}[dir], func(t *testing.T) {
			border := lipgloss.NormalBorder()
			border.Left, border.Right = "界", "界"
			calls := 0
			child := &Node{ShowBorder: true, Border: border, View: func(w, h int) string {
				calls++
				if w != 6 || h != 3 {
					t.Errorf("border child callback %dx%d, want 6x3", w, h)
				}
				return customBorderBlock(w, h)
			}}
			sibling := Leaf(customBorderBlock)
			width, height := 20, 5
			if dir == Row {
				child.Width = 10
			} else {
				child.Height = 5
				width, height = 10, 10
			}
			root := &Node{Dir: dir, Children: []*Node{child, sibling}}
			if err := Validate(root); err != nil {
				t.Fatal(err)
			}
			out := root.Render(width, height)
			if calls != 1 || lipgloss.Width(out) != width || lipgloss.Height(out) != height {
				t.Fatalf("nested output %dx%d, want %dx%d; calls=%d", lipgloss.Width(out), lipgloss.Height(out), width, height, calls)
			}
			if child.Border != border {
				t.Fatal("render mutated child border")
			}
		})
	}
}

func TestCustomBorderTinyAllocation(t *testing.T) {
	border := lipgloss.NormalBorder()
	border.Left, border.Right = "界", "界"
	for _, size := range [][2]int{{0, 0}, {1, 1}, {3, 2}, {4, 2}} {
		calls := 0
		node := &Node{ShowBorder: true, Border: border, View: func(w, h int) string {
			calls++
			if w != max(0, size[0]-4) || h != max(0, size[1]-2) {
				t.Errorf("size %v: callback %dx%d", size, w, h)
			}
			return customBorderBlock(w, h)
		}}
		node.Render(size[0], size[1])
		if calls != 1 {
			t.Fatalf("size %v: callback count %d", size, calls)
		}
	}
}

func customBorderBlock(w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	lines := make([]string, h)
	for i := range lines {
		lines[i] = strings.Repeat("x", w)
	}
	return strings.Join(lines, "\n")
}
