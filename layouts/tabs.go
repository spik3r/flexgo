package layouts

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

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
// underline is centered under its title using AlignSelf. Long titles are
// truncated to their allocated column with an ellipsis. Multiline titles use
// their first line with an ellipsis; the strip always keeps its two rows.
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
		title := tabTitle(tab.Title)
		underlineText := " "
		underlineWidth := 1
		if i == active {
			titleWidth := max(1, lipgloss.Width(title))
			underlineText = strings.Repeat("-", titleWidth)
			underlineWidth = titleWidth
		}

		strip = append(strip, &flexgo.Node{
			Flex: 1,
			Dir:  flexgo.Col,
			Children: []*flexgo.Node{
				{Height: 1, View: tabHeaderView(title, "…")},
				{Height: 1, Width: underlineWidth, AlignSelf: alignPtr(flexgo.AlignCenter), View: tabHeaderView(underlineText, "")},
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

// Header callbacks cannot wrap: the strip reserves one row for each leaf.
func tabHeaderView(text, tail string) func(int, int) string {
	return func(w, h int) string {
		if w <= 0 || h <= 0 {
			return ""
		}
		return lipgloss.NewStyle().Width(w).MaxWidth(w).MaxHeight(1).
			AlignHorizontal(lipgloss.Center).Render(ansi.Truncate(text, w, tail))
	}
}

func tabTitle(title string) string {
	first, _, multiline := strings.Cut(strings.ReplaceAll(title, "\r\n", "\n"), "\n")
	first, _, carriageReturn := strings.Cut(first, "\r")
	if multiline || carriageReturn {
		first += "…"
		// The discarded lines may contain the original style/link closers.
		// Close both before padding, neighboring tabs or panel output.
		if ansi.Strip(first) != first {
			first += ansi.ResetStyle + ansi.ResetHyperlink()
		}
	}
	return first
}
