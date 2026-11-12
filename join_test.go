package flexgo

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestResolveMainAxisSizesRedistribution(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		isRow    bool
		gap      int
		children []*Node
		want     []int
	}{
		{
			name:  "row_min_redistributes",
			total: 60,
			isRow: true,
			children: []*Node{
				{Flex: 1, MinWidth: 50},
				{Flex: 1},
			},
			want: []int{50, 10},
		},
		{
			name:  "row_max_redistributes",
			total: 60,
			isRow: true,
			children: []*Node{
				{Flex: 1, MaxWidth: 20},
				{Flex: 1},
			},
			want: []int{20, 40},
		},
		{
			name:  "col_min_redistributes",
			total: 20,
			isRow: false,
			children: []*Node{
				{Flex: 1, MinHeight: 12},
				{Flex: 1},
			},
			want: []int{12, 8},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveMainAxisSizes(tc.total, tc.children, tc.isRow, tc.gap)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("unexpected sizes: got=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestJoinJustifyStart(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
		box.Render("C"),
	}

	result := join(parts, Row, JustifyStart, 30, 0, 0, nil)

	if lipgloss.Width(result) != 30 {
		t.Errorf("Expected width 30, got %d", lipgloss.Width(result))
	}
}

func TestJoinJustifyCenter(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
		box.Render("C"),
	}

	result := join(parts, Row, JustifyCenter, 50, 0, 0, nil)

	width := lipgloss.Width(result)
	if width != 50 {
		t.Errorf("Expected width 50, got %d", width)
	}
}

func TestJoinJustifyEnd(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
		box.Render("C"),
	}

	result := join(parts, Row, JustifyEnd, 50, 0, 0, nil)

	width := lipgloss.Width(result)
	if width != 50 {
		t.Errorf("Expected width 50, got %d", width)
	}
}

func TestJoinJustifySpaceBetween(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
		box.Render("C"),
	}

	result := join(parts, Row, JustifySpaceBetween, 50, 0, 0, nil)

	width := lipgloss.Width(result)
	if width != 50 {
		t.Errorf("Expected width 50, got %d", width)
	}
}

func TestJoinWithGapJustifyStart(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
	}

	result := join(parts, Row, JustifyStart, 22, 0, 2, nil)

	width := lipgloss.Width(result)
	if width < 22 {
		t.Errorf("Expected width >= 22, got %d", width)
	}
}

func TestJoinWithGapJustifyCenter(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
	}

	result := join(parts, Row, JustifyCenter, 30, 0, 2, nil)

	width := lipgloss.Width(result)
	if width < 30 {
		t.Errorf("Expected width >= 30, got %d", width)
	}
}

func TestJoinWithGapJustifyEnd(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
	}

	result := join(parts, Row, JustifyEnd, 30, 0, 2, nil)

	width := lipgloss.Width(result)
	if width < 30 {
		t.Errorf("Expected width >= 30, got %d", width)
	}
}

func TestJoinEmpty(t *testing.T) {
	result := join([]string{}, Row, JustifyStart, 10, 0, 0, nil)
	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestJoinVertical(t *testing.T) {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(10).Height(3)

	parts := []string{
		box.Render("A"),
		box.Render("B"),
	}

	result := join(parts, Col, JustifyStart, 0, 6, 0, nil)

	height := lipgloss.Height(result)
	if height != 6 {
		t.Errorf("Expected height 6, got %d", height)
	}
}

func TestJoinGapCenterAndEndFitContainer(t *testing.T) {
	tests := []struct {
		name    string
		dir     Direction
		justify Justify
		want    string
	}{
		{"row_center", Row, JustifyCenter, "   A  B   "},
		{"row_end", Row, JustifyEnd, "      A  B"},
		{"column_center", Col, JustifyCenter, "   A  B   "},
		{"column_end", Col, JustifyEnd, "      A  B"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, h := 10, 1
			if tc.dir == Col {
				w, h = 1, 10
			}
			got := join([]string{"A", "B"}, tc.dir, tc.justify, w, h, 2, nil)
			if tc.dir == Col {
				if height := lipgloss.Height(got); height != len(tc.want) {
					t.Fatalf("height=%d, want %d: %q", height, len(tc.want), got)
				}
				got = strings.ReplaceAll(got, "\n", "")
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestJoinGapBoundaryCases(t *testing.T) {
	tests := []struct {
		name    string
		dir     Direction
		justify Justify
		size    int
		gap     int
		want    string
	}{
		{"row_odd_center", Row, JustifyCenter, 9, 2, "  A  B   "},
		{"column_odd_center", Col, JustifyCenter, 9, 2, "  A  B   "},
		{"row_no_gap_center", Row, JustifyCenter, 9, 0, "   AB    "},
		{"column_no_gap_center", Col, JustifyCenter, 9, 0, "   AB    "},
		{"row_center_overflow", Row, JustifyCenter, 3, 2, "A  B"},
		{"column_end_overflow", Col, JustifyEnd, 3, 2, "A  B"},
		{"row_space_between", Row, JustifySpaceBetween, 10, 2, "A        B"},
		{"column_space_between", Col, JustifySpaceBetween, 10, 2, "A        B"},
		{"row_space_between_overflow", Row, JustifySpaceBetween, 3, 2, "A  B"},
		{"column_space_between_overflow", Col, JustifySpaceBetween, 3, 2, "A  B"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, h := tc.size, 1
			if tc.dir == Col {
				w, h = 1, tc.size
			}
			got := join([]string{"A", "B"}, tc.dir, tc.justify, w, h, tc.gap, nil)
			if tc.dir == Col {
				if height := lipgloss.Height(got); height != len(tc.want) {
					t.Fatalf("height=%d, want %d: %q", height, len(tc.want), got)
				}
				got = strings.ReplaceAll(got, "\n", "")
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRenderFixedChildrenWithGapJustification(t *testing.T) {
	for _, dir := range []Direction{Row, Col} {
		for _, justify := range []Justify{JustifyCenter, JustifyEnd} {
			t.Run(strings.Join([]string{strconv.Itoa(int(dir)), strconv.Itoa(int(justify))}, "_"), func(t *testing.T) {
				root := &Node{
					Dir: dir, Justify: justify, Gap: 2,
					Children: []*Node{
						{Flex: 1, Width: 1, Height: 1, View: func(w, h int) string { return "A" }},
						{Flex: 1, Width: 1, Height: 1, View: func(w, h int) string { return "B" }},
					},
				}
				w, h := 10, 1
				if dir == Col {
					w, h = 1, 10
				}
				got := root.Render(w, h)
				want := "   A  B   "
				if justify == JustifyEnd {
					want = "      A  B"
				}
				flat := got
				if dir == Col {
					flat = strings.ReplaceAll(got, "\n", "")
				}
				if flat != want {
					t.Fatalf("rendered %q, want %q", flat, want)
				}
				if width, height := lipgloss.Width(got), lipgloss.Height(got); width != w || height != h {
					t.Fatalf("size=(%d,%d), want (%d,%d): %q", width, height, w, h, got)
				}
			})
		}
	}
}

func TestJoinGapSpacersUseBackground(t *testing.T) {
	bg := lipgloss.Color("57")
	for _, dir := range []Direction{Row, Col} {
		for _, justify := range []Justify{JustifyCenter, JustifyEnd} {
			w, h := 10, 1
			if dir == Col {
				w, h = 1, 10
			}
			got := join([]string{"A", "B"}, dir, justify, w, h, 2, bg)
			if !containsAnsi(got, "48;5;57") && !containsAnsi(got, "48:5:57") {
				t.Fatalf("missing gap or padding background for dir=%d justify=%d: %q", dir, justify, got)
			}
		}
	}
}
