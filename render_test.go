package flexgo

import (
	"reflect"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestRenderLeafNode(t *testing.T) {
	node := &Node{
		Width:  10,
		Height: 3,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	result := node.Render(20, 10)
	if result == "" {
		t.Errorf("Expected non-empty result")
	}
}

func TestRenderContainer(t *testing.T) {
	node := &Node{
		Dir: Row,
		Children: []*Node{
			{Width: 10, Height: 3},
			{Width: 10, Height: 3},
		},
	}

	result := node.Render(25, 5)
	if result == "" {
		t.Errorf("Expected non-empty result")
	}
}

func TestRenderWithMargin(t *testing.T) {
	node := &Node{
		Margin: 1,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	result := node.Render(10, 5)
	if result == "" {
		t.Errorf("Expected non-empty result")
	}
}

func TestRenderWithPadding(t *testing.T) {
	node := &Node{
		Padding: 1,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	result := node.Render(10, 5)
	if result == "" {
		t.Errorf("Expected non-empty result")
	}
}

func TestRenderAutoMarginLeftRight(t *testing.T) {
	node := &Node{
		MarginLeftAuto:  true,
		MarginRightAuto: true,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	result := node.Render(20, 5)
	result2 := node.Render(20, 5)
	if result != result2 {
		t.Fatalf("expected identical render output across repeated calls")
	}
	result = result2
	width := lipgloss.Width(result)

	if width != 20 {
		t.Errorf("Expected width 20, got %d", width)
	}
}

func TestRenderAutoMarginTopBottom(t *testing.T) {
	node := &Node{
		MarginTopAuto:    true,
		MarginBottomAuto: true,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	result := node.Render(20, 5)
	height := lipgloss.Height(result)

	if height != 5 {
		t.Errorf("Expected height 5, got %d", height)
	}
}

func TestRenderWithDebug(t *testing.T) {
	node := &Node{
		Debug:  true,
		Width:  10,
		Height: 3,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	result := node.Render(20, 10)

	if lipgloss.Width(result) == 10 {
		t.Errorf("Debug should add border, got width 10")
	}
}

func TestRenderColDirection(t *testing.T) {
	node := &Node{
		Dir: Col,
		Children: []*Node{
			{Height: 3},
			{Height: 3},
		},
	}

	result := node.Render(10, 10)
	height := lipgloss.Height(result)

	if height != 10 {
		t.Errorf("Expected height 10, got %d", height)
	}
}

func TestRenderImmutable(t *testing.T) {
	children := []*Node{
		{Flex: 1},
		{Flex: 1},
	}
	node := &Node{
		Dir:      Row,
		Children: children,
	}

	_ = node.Render(40, 5)
	_ = node.Render(40, 5)

	if children[0].Flex != 1 {
		t.Errorf("Flex should not mutate, expected 1, got %d", children[0].Flex)
	}
}

func TestRenderAutoMarginVerticalContainer(t *testing.T) {
	centered := &Node{
		Width:            10,
		Height:           5,
		MarginTopAuto:    true,
		MarginBottomAuto: true,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	node := &Node{
		Dir:   Col,
		Width: 20,
		Children: []*Node{
			centered,
		},
	}

	result := node.Render(20, 20)
	height := lipgloss.Height(result)

	if height != 20 {
		t.Errorf("Expected height 20, got %d", height)
	}
}

func TestRenderAutoMarginHorizontalContainer(t *testing.T) {
	centered := &Node{
		Width:           10,
		Height:          5,
		MarginLeftAuto:  true,
		MarginRightAuto: true,
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	node := &Node{
		Dir:    Row,
		Height: 10,
		Children: []*Node{
			centered,
		},
	}

	result := node.Render(20, 10)
	width := lipgloss.Width(result)

	if width != 20 {
		t.Errorf("Expected width 20, got %d", width)
	}
}

// TestRenderPaddingDimensions guards the "padding applied twice" bug:
// a node with Padding and Background must render at the requested size,
// not at (w × h + 2*padding).
func TestRenderPaddingDimensions(t *testing.T) {
	n := &Node{
		Padding:    3,
		Background: lipgloss.Color("160"),
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	out := n.Render(30, 8)
	if lipgloss.Width(out) != 30 {
		t.Errorf("expected width 30, got %d", lipgloss.Width(out))
	}
	if lipgloss.Height(out) != 8 {
		t.Errorf("expected height 8, got %d", lipgloss.Height(out))
	}
}

// TestRenderLeafNaturalHeightWithBackground guards recommendation R1:
// a leaf with Background must not be forced to fill its allocated height,
// otherwise cross-axis Align on the parent has nothing to position.
func TestRenderLeafNaturalHeightWithBackground(t *testing.T) {
	// Fixed-height leaf: always renders 3 rows regardless of allocation.
	fixedHeight3 := func(w, h int) string {
		return lipgloss.NewStyle().Width(w).Height(3).Render("content")
	}
	leaf := &Node{
		Background: lipgloss.Color("93"),
		View:       fixedHeight3,
	}

	out := leaf.Render(20, 10)
	if lipgloss.Height(out) != 3 {
		t.Errorf("leaf with Background should keep natural height 3, got %d", lipgloss.Height(out))
	}
}

// TestRenderMarginTransparent guards the "margin must not be painted"
// rule: a child with Margin inside a parent with a different Background
// should show the parent's bg in the margin area, not the child's.
func TestRenderMarginTransparent(t *testing.T) {
	child := &Node{
		Margin:     2,
		Background: lipgloss.Color("93"), // element colour
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}
	parent := &Node{
		Background: lipgloss.Color("240"), // parent colour (should show through margin)
		Children:   []*Node{child},
	}

	out := parent.Render(20, 10)

	// The top two rows are entirely the margin area — they must contain
	// only the parent's bg colour (240), never the child's (93).
	lines := splitLines(out)
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	for i := range 2 {
		if containsANSIBackground(lines[i], "93") {
			t.Errorf("margin row %d contains element bg '93' — margin should be transparent", i)
		}
	}
}

// splitLines / containsANSIBackground are tiny helpers used only by the
// margin-transparency test above. Kept here (not in align.go) because
// they're test-only.
func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func containsANSIBackground(line, colorCode string) bool {
	needle := "[48;5;" + colorCode + "m"
	for i := 0; i+len(needle) <= len(line); i++ {
		if line[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// T1 — Debug + Background regression. Both modes together must render at
// the Debug-expanded size (content + border + meta row), never double-pad
// or wrap the debug glyphs.
func TestRenderDebugWithBackground(t *testing.T) {
	// Debug adds 2 cols (borders) and 3 rows (borders + meta) around the
	// content box, so at allocation (20, 10) content should be 18×7.
	n := &Node{
		Debug:      true,
		Background: lipgloss.Color("93"),
		View: func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render("X")
		},
	}

	out := n.Render(20, 10)
	if w := lipgloss.Width(out); w != 20 {
		t.Errorf("Debug+Background: expected width 20, got %d", w)
	}
	if h := lipgloss.Height(out); h != 10 {
		t.Errorf("Debug+Background: expected height 10, got %d", h)
	}
}

// T3 — structural immutability. Render must not mutate any exported Node
// field on the tree (or its children), across repeated calls. The older
// test covered only `Flex`; reflect.DeepEqual catches anything new.
func TestRenderTreeImmutable(t *testing.T) {
	alignCenter := AlignCenter
	tree := &Node{
		Dir:      Row,
		Gap:      2,
		Padding:  1,
		Paddings: Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2},
		Margin:   1,
		Margins:  Spacing{Top: 1, Bottom: 1},
		Justify:  JustifyCenter,
		Align:    AlignCenter,
		Children: []*Node{
			{Flex: 1, MinWidth: 10, MaxWidth: 40, Background: lipgloss.Color("61")},
			{Flex: 2, AlignSelf: &alignCenter, Border: lipgloss.NormalBorder()},
			{Width: 15, Height: 4, Name: "fixed"},
		},
	}

	before := cloneNode(tree)

	_ = tree.Render(80, 20)
	_ = tree.Render(40, 10)
	_ = tree.Render(80, 20)

	if !reflect.DeepEqual(before, tree) {
		t.Fatalf("Render mutated the node tree\nbefore: %+v\nafter:  %+v", before, tree)
	}
}

// cloneNode deep-copies a Node tree for the immutability test. It relies
// on reflect so newly added fields are automatically included.
func cloneNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	out := *n
	out.Children = nil
	for _, c := range n.Children {
		out.Children = append(out.Children, cloneNode(c))
	}
	if n.AlignSelf != nil {
		v := *n.AlignSelf
		out.AlignSelf = &v
	}
	return &out
}
