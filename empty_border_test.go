package flexgo

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestEmptyExplicitBorderKeepsLeafAllocation(t *testing.T) {
	for _, debug := range []bool{false, true} {
		var gotW, gotH int
		root := &Node{ShowBorder: true, Debug: debug, Name: "hidden-debug", View: func(w, h int) string {
			gotW, gotH = w, h
			rows := make([]string, h)
			for i := range rows {
				rows[i] = strings.Repeat("x", w)
			}
			return strings.Join(rows, "\n")
		}}
		out := root.Render(8, 4)
		if gotW != 8 || gotH != 4 {
			t.Errorf("debug=%v callback got %dx%d, want 8x4", debug, gotW, gotH)
		}
		if out != "xxxxxxxx\nxxxxxxxx\nxxxxxxxx\nxxxxxxxx" {
			t.Errorf("empty border changed output: %q", out)
		}
	}
}

func TestEmptyBorderMatchesUnframedSpacingAndFixedSizes(t *testing.T) {
	for _, tc := range []struct {
		name      string
		configure func(*Node)
	}{
		{"plain", func(n *Node) {}},
		{"padding and margin", func(n *Node) { n.Padding = 1; n.Margin = 1 }},
		{"asymmetric spacing", func(n *Node) {
			n.Paddings = Spacing{Top: 1, Right: 2, Bottom: 1, Left: 1}
			n.Margins = Spacing{Top: 1, Right: 1, Bottom: 2, Left: 1}
		}},
		{"fixed dimensions", func(n *Node) { n.Width = 6; n.Height = 3 }},
		{"automatic margins", func(n *Node) {
			n.Width = 6
			n.Height = 3
			n.MarginLeftAuto = true
			n.MarginRightAuto = true
			n.MarginTopAuto = true
			n.MarginBottomAuto = true
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var sizes [][2]int
			leaf := &Node{View: func(w, h int) string {
				sizes = append(sizes, [2]int{w, h})
				return lipgloss.NewStyle().Width(w).Height(h).Render("x")
			}}
			tc.configure(leaf)
			plain := leaf.Render(12, 9)
			leaf.ShowBorder = true
			leaf.Debug = true
			leaf.Name = "hidden-debug"
			framed := leaf.Render(12, 9)
			if sizes[0] != sizes[1] || plain != framed {
				t.Fatalf("empty border changed allocation/output: sizes=%v plain=%q empty=%q", sizes, plain, framed)
			}
		})
	}
}

func TestEmptyBordersKeepNestedChildAllocations(t *testing.T) {
	for _, dir := range []Direction{Row, Col} {
		var sizes [][2]int
		root := &Node{Dir: dir, Gap: 1, Children: []*Node{
			{Flex: 1, View: func(w, h int) string {
				sizes = append(sizes, [2]int{w, h})
				return lipgloss.NewStyle().Width(w).Height(h).Render("A")
			}},
			{Flex: 2, View: func(w, h int) string {
				sizes = append(sizes, [2]int{w, h})
				return lipgloss.NewStyle().Width(w).Height(h).Render("B")
			}},
		}}
		plain := root.Render(12, 7)
		root.ShowBorder = true
		root.Debug = true
		for _, child := range root.Children {
			child.ShowBorder = true
			child.Debug = true
			child.Name = "hidden-debug"
		}
		empty := root.Render(12, 7)
		if sizes[0] != sizes[2] || sizes[1] != sizes[3] || plain != empty {
			t.Fatalf("direction=%v sizes=%v plain=%q empty=%q", dir, sizes, plain, empty)
		}
	}
}

func TestNonemptyBordersAndDebugKeepFrameReservations(t *testing.T) {
	for _, tc := range []struct {
		name         string
		border       lipgloss.Border
		show, debug  bool
		wantW, wantH int
	}{
		{"normal", lipgloss.NormalBorder(), true, false, 6, 3},
		{"hidden", lipgloss.HiddenBorder(), true, false, 6, 3},
		{"normal overrides debug", lipgloss.NormalBorder(), true, true, 6, 3},
		{"hidden overrides debug", lipgloss.HiddenBorder(), true, true, 6, 3},
		{"border disabled", lipgloss.NormalBorder(), false, false, 8, 5},
		{"debug only", lipgloss.Border{}, false, true, 6, 2},
	} {
		var gotW, gotH int
		n := &Node{ShowBorder: tc.show, Border: tc.border, Debug: tc.debug, View: func(w, h int) string { gotW, gotH = w, h; return lipgloss.NewStyle().Width(w).Height(h).Render("x") }}
		out := n.Render(8, 5)
		if gotW != tc.wantW || gotH != tc.wantH {
			t.Errorf("%s callback=%dx%d want=%dx%d", tc.name, gotW, gotH, tc.wantW, tc.wantH)
		}
		if lipgloss.Width(out) != 8 || lipgloss.Height(out) != 5 {
			t.Errorf("%s frame changed size: %q", tc.name, out)
		}
	}
}
