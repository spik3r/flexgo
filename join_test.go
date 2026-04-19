package flexgo

import (
	"reflect"
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
