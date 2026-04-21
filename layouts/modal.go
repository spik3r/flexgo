package layouts

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// Modal builds a centred, bordered dialog of the requested size. It
// uses auto-margins on all four sides, so when rendered at the screen's
// (w, h) the dialog floats in the middle with spare space around it.
//
// Layout: a one-row bold title, then the body fills the remaining
// height. Both the title's fill area and the modal's border area are
// painted with bg so the dialog has a single consistent backdrop —
// pass nil for terminal-default.
//
// Typical usage — swap the modal in for the root tree when open:
//
//	if m.modalOpen {
//	    return layouts.Modal("Title", body, 40, 10,
//	        lipgloss.Color("237")).Render(w, h)
//	}
//
// Customize (top 3 overrides):
//   - root.Border / root.BorderForeground — swap the frame style.
//   - root.Children[0] — replace the title leaf (e.g. for tabs or an icon row).
//   - root.Children[1].Paddings — inset the body from the border.
func Modal(
	title string,
	body func(w, h int) string,
	width, height int,
	bg color.Color,
) *flexgo.Node {
	titleView := func(w, h int) string {
		style := lipgloss.NewStyle().
			Width(w).
			Height(h).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center)
		if bg != nil {
			style = style.Background(bg)
		}
		return style.Render(title)
	}
	return &flexgo.Node{
		MarginTopAuto:    true,
		MarginBottomAuto: true,
		MarginLeftAuto:   true,
		MarginRightAuto:  true,
		Width:            width,
		Height:           height,
		ShowBorder:       true,
		Border:           lipgloss.RoundedBorder(),
		Background:       bg,
		Dir:              flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 1, View: titleView},
			{Flex: 1, View: body},
		},
	}
}

// WithBackground sets n.Background and returns n, for fluent
// post-construction tweaks to a recipe:
//
//	root := layouts.WithBackground(
//	    layouts.HeaderBodyFooter(...),
//	    lipgloss.Color("236"),
//	)
func WithBackground(n *flexgo.Node, bg color.Color) *flexgo.Node {
	n.Background = bg
	return n
}

// WithForeground sets n.Foreground and returns n, for fluent
// post-construction tweaks to a recipe:
//
//	root := layouts.WithForeground(
//	    layouts.HeaderBodyFooter(...),
//	    lipgloss.Color("250"),
//	)
func WithForeground(n *flexgo.Node, fg color.Color) *flexgo.Node {
	n.Foreground = fg
	return n
}
