package layouts

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// FormRow describes one row in Form.
type FormRow struct {
	Label  string
	Field  func(w, h int) string
	Height int
}

// Form builds vertically stacked Label:Field rows with aligned labels.
//
// Each row is rendered as:
//
//	[label (fixed width, right-aligned)] [field (flex)]
//
// Row height defaults to 1 when FormRow.Height is not set.
//
// Customize (top 3 overrides):
//   - root.Gap — vertical space between rows.
//   - root.Paddings / root.Background — frame and inset the whole form.
//   - row.Children[0].View — wrap the label to restyle (bold, coloured, etc.).
func Form(labelWidth int, rows []FormRow) *flexgo.Node {
	children := make([]*flexgo.Node, 0, len(rows))
	for _, row := range rows {
		h := row.Height
		if h <= 0 {
			h = 1
		}

		labelText := row.Label
		if labelText != "" && !strings.HasSuffix(labelText, ":") {
			labelText += ":"
		}

		labelView := func(text string) func(int, int) string {
			return func(w, h int) string {
				return lipgloss.NewStyle().
					Width(w).
					Height(h).
					Align(lipgloss.Right, lipgloss.Center).
					Render(text)
			}
		}(labelText)

		children = append(children, &flexgo.Node{
			Height: h,
			Dir:    flexgo.Row,
			Gap:    1,
			Children: []*flexgo.Node{
				{Width: labelWidth, View: labelView},
				{Flex: 1, View: row.Field},
			},
		})
	}

	return &flexgo.Node{Dir: flexgo.Col, Children: children}
}
