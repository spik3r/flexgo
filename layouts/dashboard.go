package layouts

import "github.com/spik3r/flexgo"

// Dashboard builds a common operator-console shape:
//
//	┌ sidebar ┬──────── header ────────┐
//	│         │         main           │
//	│         ├──────── status ────────┤
//	└─────────┴────────────────────────┘
//
// The sidebar/header/status sections are fixed-size, while main fills
// remaining space. Pass nil views for optional sections.
//
// Customize (top 3 overrides):
//   - root.Background + root.Gap — backdrop that shows through as gutters.
//   - root.Children[0].ShowBorder + .Border — frame the sidebar.
//   - root.Children[1].Gap — space between header/main/status.
func Dashboard(
	sidebarWidth int,
	headerHeight int,
	statusHeight int,
	sidebar func(w, h int) string,
	header func(w, h int) string,
	main func(w, h int) string,
	status func(w, h int) string,
) *flexgo.Node {
	rightChildren := []*flexgo.Node{}
	if header != nil {
		rightChildren = append(rightChildren, &flexgo.Node{Height: headerHeight, View: header})
	}
	rightChildren = append(rightChildren, &flexgo.Node{Flex: 1, View: main})
	if status != nil {
		rightChildren = append(rightChildren, &flexgo.Node{Height: statusHeight, View: status})
	}

	return &flexgo.Node{
		Dir: flexgo.Row,
		Children: []*flexgo.Node{
			{Width: sidebarWidth, View: sidebar},
			{Flex: 1, Dir: flexgo.Col, Children: rightChildren},
		},
	}
}
