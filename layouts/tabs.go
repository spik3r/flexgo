package layouts

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// Tab describes one tab title and panel content in Tabs.
type Tab struct {
	Title string
	Panel func(w, h int) string
}

// Tabs builds a tab-strip plus active panel layout.
//
// The strip is two rows high (titles + underline). The active tab's
// underline is centered under its title using AlignSelf.
//
// Customize (top 3 overrides):
//   - root.Children[0].Gap — space between tab titles on the strip.
//   - root.Background / root.Paddings — frame the whole widget.
//   - root.Children[1] — swap the active panel for a bordered wrapper.
func Tabs(active int, tabs []Tab) *flexgo.Node {
	if len(tabs) == 0 {
		return &flexgo.Node{Dir: flexgo.Col, Children: []*flexgo.Node{{Flex: 1}}}
	}
	if active < 0 {
		active = 0
	}
	if active >= len(tabs) {
		active = len(tabs) - 1
	}

	strip := make([]*flexgo.Node, 0, len(tabs))
	for i, tab := range tabs {
		underlineText := " "
		underlineWidth := 1
		if i == active {
			titleWidth := max(1, lipgloss.Width(tab.Title))
			underlineText = strings.Repeat("-", titleWidth)
			underlineWidth = titleWidth
		}

		strip = append(strip, &flexgo.Node{
			Flex: 1,
			Dir:  flexgo.Col,
			Children: []*flexgo.Node{
				{Height: 1, View: centeredView(tab.Title)},
				{Height: 1, Width: underlineWidth, AlignSelf: alignPtr(flexgo.AlignCenter), View: centeredView(underlineText)},
			},
		})
	}

	panel := tabs[active].Panel
	return &flexgo.Node{
		Dir: flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 2, Dir: flexgo.Row, Children: strip},
			{Flex: 1, View: panel},
		},
	}
}
