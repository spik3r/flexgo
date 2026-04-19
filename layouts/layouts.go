// Package layouts offers ready-made flexgo node trees for common TUI
// shapes. Each recipe returns a *flexgo.Node that the caller renders
// directly. All styling (colours, borders, spacing) is set on the
// returned Node — override fields after construction to customise.
package layouts

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

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
// Customise after construction, e.g. set Background on the root for a
// page-wide backdrop, or replace Children[1].Padding to pad the body.
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
