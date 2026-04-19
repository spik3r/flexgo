package layouts

import (
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

func TestHeaderBodyFooterShape(t *testing.T) {
	leaf := func(label string) func(int, int) string {
		return func(w, h int) string {
			return lipgloss.NewStyle().Width(w).Height(h).Render(label)
		}
	}

	root := HeaderBodyFooter(
		3, leaf("H"),
		leaf("B"),
		2, leaf("F"),
	)

	out := root.Render(40, 20)
	if h := lipgloss.Height(out); h != 20 {
		t.Errorf("expected height 20, got %d", h)
	}
	if w := lipgloss.Width(out); w != 40 {
		t.Errorf("expected width 40, got %d", w)
	}
	if len(root.Children) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(root.Children))
	}
	if root.Children[0].Height != 3 {
		t.Errorf("header height: expected 3, got %d", root.Children[0].Height)
	}
	if root.Children[1].Flex != 1 {
		t.Errorf("body flex: expected 1, got %d", root.Children[1].Flex)
	}
	if root.Children[2].Height != 2 {
		t.Errorf("footer height: expected 2, got %d", root.Children[2].Height)
	}
}

func TestHeaderBodyFooterOmitsNilSections(t *testing.T) {
	body := func(w, h int) string { return "" }
	root := HeaderBodyFooter(0, nil, body, 0, nil)

	if len(root.Children) != 1 {
		t.Fatalf("expected only body when header/footer nil, got %d children", len(root.Children))
	}
	if root.Children[0].Flex != 1 {
		t.Errorf("body should be flex, got Flex=%d", root.Children[0].Flex)
	}
}

func TestWithBackground(t *testing.T) {
	n := &flexgo.Node{}
	bg := lipgloss.Color("236")
	got := WithBackground(n, bg)
	if got != n {
		t.Errorf("WithBackground should return the same node")
	}
	if n.Background != bg {
		t.Errorf("Background not applied")
	}
}
