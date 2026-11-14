package flexgo

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestSpaceBetweenKeepsLastChildAtFarEdge(t *testing.T) {
	for _, dir := range []Direction{Row, Col} {
		root := &Node{Dir: dir, Justify: JustifySpaceBetween, Gap: 1, Children: []*Node{
			{Width: 1, Height: 1, View: func(w, h int) string { return "A" }},
			{Width: 1, Height: 1, View: func(w, h int) string { return "B" }},
			{Width: 1, Height: 1, View: func(w, h int) string { return "C" }},
		}}
		w, h := 8, 1
		if dir == Col {
			w, h = 1, 8
		}
		got := root.Render(w, h)
		flat := strings.ReplaceAll(got, "\n", "")
		if flat != "A   B  C" {
			t.Errorf("direction %v: got %q, want %q", dir, flat, "A   B  C")
		}
		if lipgloss.Width(got) != w || lipgloss.Height(got) != h {
			t.Errorf("render escaped %dx%d: %q", w, h, got)
		}
	}
}

func TestSpaceBetweenUsesEverySpareCellBetweenChildren(t *testing.T) {
	for _, dir := range []Direction{Row, Col} {
		for _, tc := range []struct {
			parts     []string
			size, gap int
			want      string
		}{
			{[]string{"A", "B", "C"}, 8, 0, "A   B  C"},
			{[]string{"A", "B", "C"}, 8, 2, "A   B  C"},
			{[]string{"A", "B", "C"}, 9, 1, "A   B   C"},
			{[]string{"A", "B", "C", "D"}, 9, 0, "A  B  C D"},
			{[]string{"A", "B", "C", "D"}, 5, 0, "A BCD"},
			{[]string{"A", "B", "C"}, 7, 2, "A  B  C"},
			{[]string{"A", "B", "C"}, 3, 0, "ABC"},
			{[]string{"A", "B", "C"}, 2, 1, "A B C"},
			{[]string{"A", "B"}, 5, 1, "A   B"},
			{[]string{"A"}, 5, 2, "A"},
			{nil, 5, 2, ""},
		} {
			w, h := tc.size, 1
			if dir == Col {
				w, h = 1, tc.size
			}
			out := join(tc.parts, dir, JustifySpaceBetween, w, h, tc.gap, nil)
			if got := strings.ReplaceAll(out, "\n", ""); got != tc.want {
				t.Errorf("dir=%v size=%d gap=%d: got %q want %q", dir, tc.size, tc.gap, got, tc.want)
			}
			if out != "" {
				size := lipgloss.Width(out)
				if dir == Col {
					size = lipgloss.Height(out)
				}
				if size != len(tc.want) {
					t.Errorf("wrong main-axis size=%d want=%d: %q", size, len(tc.want), out)
				}
			}
		}
	}
}

func TestSpaceBetweenRemainderSpacersFillCrossAxis(t *testing.T) {
	bg := lipgloss.Color("57")
	for _, dir := range []Direction{Row, Col} {
		root := &Node{Dir: dir, Justify: JustifySpaceBetween, Background: bg, Gap: 1}
		for _, label := range []string{"A", "B", "C"} {
			root.Children = append(root.Children, &Node{Width: 1, Height: 1, View: func(w, h int) string { return label }})
		}
		w, h := 8, 3
		if dir == Col {
			w, h = 3, 8
		}
		out := root.Render(w, h)
		if lipgloss.Width(out) != w || lipgloss.Height(out) != h {
			t.Fatalf("cross-axis shape escaped %dx%d: %q", w, h, out)
		}
		plain := ansi.Strip(out)
		if dir == Row {
			if strings.Split(plain, "\n")[0] != "A   B  C" {
				t.Fatalf("wrong row positions: %q", plain)
			}
		} else {
			lines := strings.Split(plain, "\n")
			if lines[0] != "A  " || lines[4] != "B  " || lines[7] != "C  " {
				t.Fatalf("wrong column positions: %q", plain)
			}
		}
		painted := lipgloss.NewStyle().Background(bg).Render("   ")
		rows := strings.Split(out, "\n")
		if dir == Col {
			rows = rows[1:4]
		}
		for _, row := range rows {
			if !strings.Contains(row, painted) {
				t.Fatalf("expanded gap lost painted cross-axis span: %q", row)
			}
		}
	}
}
