package layouts

import "github.com/spik3r/flexgo"

// HeaderBodyFooter builds the classic three-row TUI shape:
//
//	┌────────────── header ──────────────┐  (headerHeight rows, fixed)
//	├──────────────  body  ──────────────┤  (flex: fills remaining)
//	└────────────── footer ──────────────┘  (footerHeight rows, fixed)
//
// Each parameter is a leaf View callback. Pass nil to omit a section.
// Fixed-height sections render at exactly the height requested; the
// body expands to fill whatever remains.
//
// Customize (top 3 overrides):
//   - root.Background — page-wide backdrop.
//   - root.Children[1].Paddings — pad the body without affecting header/footer.
//   - root.Gap — blank rows between the three sections.
func HeaderBodyFooter(
	headerHeight int,
	header func(w, h int) string,
	body func(w, h int) string,
	footerHeight int,
	footer func(w, h int) string,
) *flexgo.Node {
	var children []*flexgo.Node
	if header != nil {
		children = append(children, &flexgo.Node{
			Height: headerHeight,
			View:   header,
		})
	}
	children = append(children, &flexgo.Node{
		Flex: 1,
		View: body,
	})
	if footer != nil {
		children = append(children, &flexgo.Node{
			Height: footerHeight,
			View:   footer,
		})
	}
	return &flexgo.Node{
		Dir:      flexgo.Col,
		Children: children,
	}
}
