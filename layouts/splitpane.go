package layouts

import "github.com/spik3r/flexgo"

// SplitPane builds a two-pane split with a fixed ratio.
//
// firstRatio is interpreted as a percentage in [1,99]. Values outside
// this range are clamped.
//
// Customize (top 3 overrides):
//   - root.Gap — visible gutter between the two panes.
//   - root.Children[i].ShowBorder + .Border — frame a pane.
//   - root.Children[i].Paddings — inset pane content.
func SplitPane(
	dir flexgo.Direction,
	firstRatio int,
	first func(w, h int) string,
	second func(w, h int) string,
) *flexgo.Node {
	if firstRatio < 1 {
		firstRatio = 1
	}
	if firstRatio > 99 {
		firstRatio = 99
	}
	return SplitPaneFlex(dir, firstRatio, 100-firstRatio, first, second)
}

// SplitPaneFlex is like SplitPane but expresses the ratio as two flex
// weights, matching the rest of the library. SplitPaneFlex(Row, 1, 2,
// …) gives the second pane twice the width of the first.
//
// Non-positive weights are clamped to 1 to avoid zero-size panes.
func SplitPaneFlex(
	dir flexgo.Direction,
	firstFlex, secondFlex int,
	first func(w, h int) string,
	second func(w, h int) string,
) *flexgo.Node {
	if firstFlex < 1 {
		firstFlex = 1
	}
	if secondFlex < 1 {
		secondFlex = 1
	}
	return &flexgo.Node{
		Dir: dir,
		Children: []*flexgo.Node{
			{Flex: firstFlex, View: first},
			{Flex: secondFlex, View: second},
		},
	}
}
