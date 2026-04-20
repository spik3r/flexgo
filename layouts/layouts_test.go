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

func TestDashboardShape(t *testing.T) {
	view := func(w, h int) string { return "" }
	root := Dashboard(20, 3, 1, view, view, view, view)

	if root.Dir != flexgo.Row {
		t.Fatalf("dashboard root should be Row")
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 dashboard columns, got %d", len(root.Children))
	}
	if root.Children[0].Width != 20 {
		t.Fatalf("sidebar width: expected 20, got %d", root.Children[0].Width)
	}
	right := root.Children[1]
	if right.Dir != flexgo.Col {
		t.Fatalf("dashboard right column should be Col")
	}
	if len(right.Children) != 3 {
		t.Fatalf("expected header/main/status on right, got %d", len(right.Children))
	}
}

func TestSplitPaneRatioClamping(t *testing.T) {
	root := SplitPane(flexgo.Row, 200, nil, nil)
	if root.Children[0].Flex != 99 || root.Children[1].Flex != 1 {
		t.Fatalf("ratio should clamp to 99/1, got %d/%d", root.Children[0].Flex, root.Children[1].Flex)
	}

	root = SplitPane(flexgo.Row, -1, nil, nil)
	if root.Children[0].Flex != 1 || root.Children[1].Flex != 99 {
		t.Fatalf("ratio should clamp to 1/99, got %d/%d", root.Children[0].Flex, root.Children[1].Flex)
	}
}

func TestSplitPaneFlex(t *testing.T) {
	root := SplitPaneFlex(flexgo.Row, 1, 3, nil, nil)
	if root.Children[0].Flex != 1 || root.Children[1].Flex != 3 {
		t.Fatalf("expected flex 1/3, got %d/%d", root.Children[0].Flex, root.Children[1].Flex)
	}

	root = SplitPaneFlex(flexgo.Col, 0, -5, nil, nil)
	if root.Children[0].Flex != 1 || root.Children[1].Flex != 1 {
		t.Fatalf("non-positive weights should clamp to 1/1, got %d/%d", root.Children[0].Flex, root.Children[1].Flex)
	}
}

func TestGridShape(t *testing.T) {
	root := Grid(2, 3, 1, nil)
	if root.Dir != flexgo.Col {
		t.Fatalf("grid root should be Col")
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(root.Children))
	}
	for _, row := range root.Children {
		if row.Dir != flexgo.Row {
			t.Fatalf("row should be Row")
		}
		if len(row.Children) != 3 {
			t.Fatalf("expected 3 columns, got %d", len(row.Children))
		}
	}
}
