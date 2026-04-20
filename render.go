package flexgo

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Render is the public entry point. It renders the node tree into a string
// fitting (w, h). Use this at the root of your layout.
func (n *Node) Render(w, h int) string {
	return n.render(w, h, nil)
}

// render is the internal entry point used when a parent renders its child.
// parentBg is the parent container's Background, used to paint this node's
// margin area so it visually belongs to the parent.
func (n *Node) render(w, h int, parentBg color.Color) string {
	availW := w
	availH := h
	allocW := w
	allocH := h

	// Border and Debug both wrap the content. ShowBorder wins when both
	// are set so we reserve exactly the space the outer wrapper will use.
	showDebugFrame := n.Debug && !hasExplicitBorder(n)
	if hasExplicitBorder(n) || showDebugFrame {
		w -= 2
		h -= 2
		if showDebugFrame {
			h-- // one extra row for the debug label
		}
		if w < 0 {
			w = 0
		}
		if h < 0 {
			h = 0
		}
	}

	margin := resolveSpacing(n.Margin, n.Margins)
	w -= margin.TotalWidth()
	h -= margin.TotalHeight()
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	if n.Width > 0 && n.Width < w {
		w = n.Width
	}
	if n.Height > 0 && n.Height < h {
		h = n.Height
	}

	padding := resolveSpacing(n.Padding, n.Paddings)
	contentW := w - padding.TotalWidth()
	contentH := h - padding.TotalHeight()
	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	isLeaf := n.View != nil && len(n.Children) == 0

	ambientBg := n.Background
	if ambientBg == nil {
		ambientBg = parentBg
	}

	var inner string
	if isLeaf {
		inner = n.View(contentW, contentH)
	} else {
		inner = n.renderChildren(contentW, contentH, ambientBg)
	}

	boxStyle := lipgloss.NewStyle().Padding(padding.Top, padding.Right, padding.Bottom, padding.Left).Width(w)
	if !isLeaf {
		boxStyle = boxStyle.Height(h)
	}
	if n.Background != nil {
		boxStyle = boxStyle.Background(n.Background)
	}
	inner = boxStyle.Render(inner)

	if hasExplicitBorder(n) {
		borderStyle := lipgloss.NewStyle().Border(n.Border)
		if n.BorderForeground != nil {
			borderStyle = borderStyle.BorderForeground(n.BorderForeground)
		}
		if n.BorderBackground != nil {
			borderStyle = borderStyle.BorderBackground(n.BorderBackground)
		}
		inner = borderStyle.Render(inner)
	} else if showDebugFrame {
		meta := lipgloss.NewStyle().Bold(true).Width(w).MaxWidth(w).Render(debugLabel(n.Name, allocW, allocH))
		inner = lipgloss.JoinVertical(lipgloss.Left, meta, inner)
		inner = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(inner)
	}

	if margin != (Spacing{}) {
		inner = asymmetricMargin(inner, margin.Top, margin.Right, margin.Bottom, margin.Left, parentBg)
	}

	return applyAutoMargins(inner, n, availW, availH, parentBg)
}

// applyAutoMargins pushes a rendered node within its allocated slot
// when MarginLeftAuto/RightAuto/TopAuto/BottomAuto are set. Auto-margins
// only engage when there is spare space on the relevant axis.
func applyAutoMargins(inner string, n *Node, availW, availH int, parentBg color.Color) string {
	if n.MarginLeftAuto || n.MarginRightAuto {
		spare := availW - lipgloss.Width(inner)
		if spare > 0 {
			left, right := splitAuto(spare, n.MarginLeftAuto, n.MarginRightAuto)
			inner = asymmetricMargin(inner, 0, right, 0, left, parentBg)
		}
	}

	if n.MarginTopAuto || n.MarginBottomAuto {
		spare := availH - lipgloss.Height(inner)
		if spare > 0 {
			top, bottom := splitAuto(spare, n.MarginTopAuto, n.MarginBottomAuto)
			inner = asymmetricMargin(inner, top, 0, bottom, 0, parentBg)
		}
	}

	return inner
}

// splitAuto divides spare space between two auto-margin ends. Centering
// (both ends) gives the extra pixel to the trailing side.
func splitAuto(spare int, leading, trailing bool) (int, int) {
	switch {
	case leading && trailing:
		l := spare / 2
		return l, spare - l
	case leading:
		return spare, 0
	case trailing:
		return 0, spare
	}
	return 0, 0
}

func (n *Node) renderChildren(w, h int, ambientBg color.Color) string {
	isRow := n.Dir == Row

	total := w
	if !isRow {
		total = h
	}

	sizes := resolveMainAxisSizes(total, n.Children, isRow, n.Gap)

	var parts []string

	for i, child := range n.Children {
		var cw, ch int

		if isRow {
			cw = sizes[i]
			ch = h
		} else {
			cw = w
			ch = sizes[i]
		}

		// Main-axis min/max is already applied by resolveMainAxisSizes.
		// Clamp the cross axis here so a tall row or wide column child
		// can opt into a tighter/looser container size.
		if isRow {
			if child.MinHeight > 0 && ch < child.MinHeight {
				ch = child.MinHeight
			}
			if child.MaxHeight > 0 && ch > child.MaxHeight {
				ch = child.MaxHeight
			}
		} else {
			if child.MinWidth > 0 && cw < child.MinWidth {
				cw = child.MinWidth
			}
			if child.MaxWidth > 0 && cw > child.MaxWidth {
				cw = child.MaxWidth
			}
		}

		childView := child.render(cw, ch, ambientBg)

		align := n.Align
		if child.AlignSelf != nil {
			align = *child.AlignSelf
		}
		childView = applyAlign(childView, cw, ch, align, ambientBg)

		parts = append(parts, childView)
	}

	return join(parts, n.Dir, n.Justify, w, h, n.Gap, ambientBg)
}
