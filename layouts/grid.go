package layouts

import "github.com/spik3r/flexgo"

// Grid builds a uniform rows×cols grid. Each cell gets Flex:1 on both
// axes so all cells size evenly.
//
// gap is applied both between rows and between columns.
//
// Customize (top 3 overrides):
//   - root.Background — shows through the gap gutters as grid lines.
//   - root.Children[r].Children[c].ShowBorder + .Border — frame individual cells.
//   - root.Paddings — inset the grid from its container.
func Grid(rows, cols, gap int, cell func(row, col, w, h int) string) *flexgo.Node {
	if rows < 1 {
		rows = 1
	}
	if cols < 1 {
		cols = 1
	}

	rowNodes := make([]*flexgo.Node, 0, rows)
	for r := 0; r < rows; r++ {
		cells := make([]*flexgo.Node, 0, cols)
		for c := 0; c < cols; c++ {
			rr, cc := r, c
			cellView := func(w, h int) string {
				if cell == nil {
					return ""
				}
				return cell(rr, cc, w, h)
			}
			cells = append(cells, &flexgo.Node{Flex: 1, View: cellView})
		}
		rowNodes = append(rowNodes, &flexgo.Node{Flex: 1, Dir: flexgo.Row, Gap: gap, Children: cells})
	}

	return &flexgo.Node{Dir: flexgo.Col, Gap: gap, Children: rowNodes}
}
