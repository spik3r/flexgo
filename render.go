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

	if hasExplicitBorder(n) || n.Debug {
		w -= 2
		h -= 2
		if n.Debug {
			h--
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
	} else if n.Debug {
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

	if n.MarginLeftAuto || n.MarginRightAuto {
		innerW := lipgloss.Width(inner)
		spare := availW - innerW
		if spare > 0 {
			var left, right int
			switch {
			case n.MarginLeftAuto && n.MarginRightAuto:
				left = spare / 2
				right = spare - left
			case n.MarginLeftAuto:
				left = spare
			case n.MarginRightAuto:
				right = spare
			}
			inner = asymmetricMargin(inner, 0, right, 0, left, parentBg)
		}
	}

	if n.MarginTopAuto || n.MarginBottomAuto {
		innerH := lipgloss.Height(inner)
		spare := availH - innerH
		if spare > 0 {
			var top, bottom int
			switch {
			case n.MarginTopAuto && n.MarginBottomAuto:
				top = spare / 2
				bottom = spare - top
			case n.MarginTopAuto:
				top = spare
			case n.MarginBottomAuto:
				bottom = spare
			}
			inner = asymmetricMargin(inner, top, 0, bottom, 0, parentBg)
		}
	}

	return inner
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

		if child.MinWidth > 0 && cw < child.MinWidth {
			cw = child.MinWidth
		}
		if child.MinHeight > 0 && ch < child.MinHeight {
			ch = child.MinHeight
		}
		if child.MaxWidth > 0 && cw > child.MaxWidth {
			cw = child.MaxWidth
		}
		if child.MaxHeight > 0 && ch > child.MaxHeight {
			ch = child.MaxHeight
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
