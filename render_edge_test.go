package flexgo

import (
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
