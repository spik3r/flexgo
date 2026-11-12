package layouts

import "github.com/spik3r/flexgo"

// AutoGrid reflows cells in row-major order using the width allocated to its
// content at render time. It fits as many columns as possible at minCellWidth,
// including gap gutters, and shares remaining width evenly. A narrow allocation
// uses one column even when minCellWidth cannot be met. The final row keeps the
// same column slots as earlier rows; unused slots are blank. Nil cells are blank.
// Nonpositive minCellWidth becomes 1; negative gap becomes 0.
//
// Like Grid, rows share available height, callbacks must bound their own output,
// and vertical gutters can overflow an allocation too short for all rows. This
// is width-driven reflow, not scrolling, clipping, or the core flex-wrap algorithm.
// Keep selection/focus/component state outside callbacks and route by cell index.
//
// The returned node is a leaf: its View builds a fresh Grid on each render.
// Inspect/Validate see the leaf, not the generated rows. Do not add Children.
// The callback slice is copied; captured state remains owned by the caller.
//
// Customize (top 3 overrides):
//   - root.Paddings — reduce the content area before column selection.
//   - root.Background — paint the grid's blank slots and gutters.
//   - root.ShowBorder + root.Border — frame the entire grid.
func AutoGrid(minCellWidth, gap int, cells ...func(w, h int) string) *flexgo.Node {
	minCellWidth = max(1, minCellWidth)
	gap = max(0, gap)
	cells = append([]func(int, int) string(nil), cells...)
	return &flexgo.Node{View: func(w, h int) string {
		if w <= 0 || h <= 0 || len(cells) == 0 {
			return ""
		}
		cols := autoGridColumns(w, minCellWidth, gap, len(cells))
		rows := 1 + (len(cells)-1)/cols
		return Grid(rows, cols, gap, func(row, col, w, h int) string {
			index := row*cols + col
			if index >= len(cells) || cells[index] == nil || w <= 0 || h <= 0 {
				return ""
			}
			return cells[index](w, h)
		}).Render(w, h)
	}}
}

func autoGridColumns(width, minWidth, gap, count int) int {
	if count == 0 {
		return 0
	}
	// 1 + (width-minWidth)/(minWidth+gap) is equivalent to the usual
	// (width+gap)/(minWidth+gap), without overflowing width+gap. If the
	// denominator would overflow, even two columns cannot fit in an int width.
	maxInt := int(^uint(0) >> 1)
	if width < minWidth || gap > maxInt-minWidth {
		return 1
	}
	return min(count, 1+(width-minWidth)/(minWidth+gap))
}
