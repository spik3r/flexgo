package layouts

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

func TestAutoGridReflowsAtAllocatedWidth(t *testing.T) {
	for _, tc := range []struct{ width, columns int }{
		{1, 1}, {9, 1}, {10, 1}, {20, 1}, {21, 2}, {31, 2}, {32, 3}, {42, 3}, {43, 4}, {100, 5},
	} {
		t.Run(fmt.Sprint(tc.width), func(t *testing.T) {
			type size struct{ w, h int }
			got := make([]size, 5)
			cells := make([]func(int, int) string, len(got))
			for i := range cells {
				cells[i] = func(w, h int) string { got[i] = size{w, h}; return "" }
			}
			root := AutoGrid(10, 1, cells...)
			root.Render(tc.width, 24)
			rows := 1 + (len(cells)-1)/tc.columns
			widths := make([]int, tc.columns)
			for c := range widths {
				widths[c] = (tc.width - (tc.columns - 1)) / tc.columns
				if c < (tc.width-(tc.columns-1))%tc.columns {
					widths[c]++
				}
			}
			for i, s := range got {
				if s.w != widths[i%tc.columns] {
					t.Fatalf("cell %d width %d, want %d", i, s.w, widths[i%tc.columns])
				}
				base := (24 - (rows - 1)) / rows
				extra := (24 - (rows - 1)) % rows
				wantH := base
				if i/tc.columns < extra {
					wantH++
				}
				if s.h != wantH {
					t.Fatalf("cell %d height %d, want %d", i, s.h, wantH)
				}
			}
		})
	}
}

func TestAutoGridUsesParentAndDecorationBudget(t *testing.T) {
	widths := make([]int, 3)
	cells := make([]func(int, int) string, 3)
	for i := range cells {
		cells[i] = func(w, h int) string { widths[i] = w; return "" }
	}
	grid := AutoGrid(10, 1, cells...)
	grid.Paddings = flexgo.Spacing{Left: 2, Right: 2}
	grid.ShowBorder = true
	grid.Border = lipgloss.NormalBorder()
	root := flexgo.HBox(&flexgo.Node{Width: 20}, grid)
	root.Gap = 1
	// 60 - sidebar 20 - gap 1 - padding 4 - border 2 = 33 content columns.
	root.Render(60, 20)
	if fmt.Sprint(widths) != "[11 10 10]" {
		t.Fatalf("wrong content budget: %v", widths)
	}
	root.Render(47, 20) // 20 content columns: only one 10-column slot fits with gap.
	if fmt.Sprint(widths) != "[20 20 20]" {
		t.Fatalf("stale parent width: %v", widths)
	}
}

func TestAutoGridEmptyNilAndNormalization(t *testing.T) {
	calls := 0
	cell := func(w, h int) string { calls++; return "x" }
	for _, size := range [][2]int{{0, 10}, {10, 0}, {-1, 2}, {2, -1}} {
		if got := AutoGrid(10, 1, cell).View(size[0], size[1]); got != "" {
			t.Fatalf("zero content area: %q", got)
		}
	}
	if calls != 0 {
		t.Fatal("zero area reached cell")
	}
	if out := AutoGrid(10, 1).View(40, 10); out != "" {
		t.Fatal("empty grid is not empty")
	}
	grid := AutoGrid(0, -4, nil, cell)
	if err := flexgo.Validate(grid); err != nil {
		t.Fatal(err)
	}
	out := grid.Render(2, 1)
	if out != " x" || calls != 1 {
		t.Fatalf("nil/normalized cells: %q, calls %d", out, calls)
	}
}

func TestAutoGridSnapshotsCallbackSliceAndPreservesSlots(t *testing.T) {
	cells := []func(int, int) string{
		func(w, h int) string { return strings.Repeat("A", w) },
		func(w, h int) string { return strings.Repeat("B", w) },
		func(w, h int) string { return strings.Repeat("C", w) },
	}
	grid := AutoGrid(3, 1, cells...)
	cells[0] = func(w, h int) string { return "changed" }
	for _, width := range []int{7, 3, 11, 7} {
		out := grid.Render(width, 5)
		if strings.Contains(out, "changed") {
			t.Fatal("callback slice not copied")
		}
		if lipgloss.Width(out) > width || lipgloss.Height(out) > 5 {
			t.Fatalf("bounded callbacks overflow: %q", out)
		}
	}
	out := grid.Render(7, 3)
	if out != "AAA BBB\n       \nCCC    " {
		t.Fatalf("final row stretched or reordered: %q", out)
	}
}

func TestAutoGridColumnArithmeticLimits(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, tc := range []struct{ width, min, gap, count, want int }{
		{maxInt, 1, 0, 5, 5}, {maxInt, maxInt, 0, 5, 1},
		{maxInt, 1, maxInt, 5, 1}, {maxInt, maxInt, maxInt, 5, 1},
		{maxInt, maxInt / 2, 1, 5, 2}, {10, 3, 1, 0, 0},
	} {
		if got := autoGridColumns(tc.width, tc.min, tc.gap, tc.count); got != tc.want {
			t.Fatalf("%+v: got %d", tc, got)
		}
	}
}

func TestAutoGridConcurrentRender(t *testing.T) {
	grid := AutoGrid(4, 1, func(w, h int) string { return strings.Repeat("a", w) }, func(w, h int) string { return "b" })
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 10 {
				if got := grid.Render(20, 5); !strings.Contains(got, "b") {
					t.Error("lost callback")
				}
			}
		})
	}
	wg.Wait()
}

func TestAutoGridKeepsGridOverflowContract(t *testing.T) {
	cell := func(w, h int) string { return "x" }
	// Three rows need two two-row gutters even if callbacks have no content area.
	if out := AutoGrid(3, 2, cell, cell, cell).Render(3, 1); lipgloss.Height(out) <= 1 {
		t.Fatal("expected overconstrained vertical gutters to remain visible")
	}
	// A callback that ignores its height is not silently clipped by the recipe.
	tall := func(w, h int) string { return "one\ntwo\nthree" }
	if out := AutoGrid(5, 0, tall).Render(5, 1); lipgloss.Height(out) <= 1 {
		t.Fatal("callback output was silently clipped")
	}
}
