package flexgo

import (
	"fmt"
	"sync"

	"charm.land/lipgloss/v2"
)

func (n *Node) Render(w, h int) string {
	availW := w // preserve original width for auto-margin centering

	// Reserve space for debug border (1 char per side = 2 per axis)
	// before any other sizing so the border fits within the allocated space.
	if n.Debug {
		w -= 2
		h -= 2
		if w < 0 {
			w = 0
		}
		if h < 0 {
			h = 0
		}
	}

	// Apply margin (outer shrink)
	w -= n.Margin * 2
	h -= n.Margin * 2
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	// Clamp to fixed Width when narrower than available space.
	// This lets a root node declare its own width (e.g. for centering)
	// without relying on a parent's distribute() call.
	if n.Width > 0 && n.Width < w {
		w = n.Width
	}

	// Apply padding (inner shrink)
	contentW := w - n.Padding*2
	contentH := h - n.Padding*2
	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	var inner string

	// Leaf node
	if n.View != nil && len(n.Children) == 0 {
		inner = n.View(contentW, contentH)
	} else {
		inner = n.renderChildren(contentW, contentH)
	}

	// Apply padding.
	// Container nodes are forced to fill their allocated height so parent
	// joins stay consistent. Leaf nodes are left at their natural height so
	// the parent's applyAlign can position them on the cross axis.
	isLeaf := n.View != nil && len(n.Children) == 0
	paddingStyle := lipgloss.NewStyle().Padding(n.Padding).Width(w)
	if !isLeaf {
		paddingStyle = paddingStyle.Height(h)
	}
	inner = paddingStyle.Render(inner)

	// Debug overlay
	if n.Debug {
		inner = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(inner)
	}

	// Apply margin (outer spacing)
	if n.Margin > 0 {
		inner = lipgloss.NewStyle().
			Margin(n.Margin).
			Render(inner)
	}

	// Auto horizontal centering via MarginLeftAuto / MarginRightAuto.
	// Centers (or right-aligns) the rendered node within the originally
	// available width. Designed for fixed-width root nodes.
	if n.MarginLeftAuto || n.MarginRightAuto {
		innerW := lipgloss.Width(inner)
		spare := availW - innerW
		if spare > 0 {
			switch {
			case n.MarginLeftAuto && n.MarginRightAuto:
				left := spare / 2
				right := spare - left
				inner = lipgloss.NewStyle().MarginLeft(left).MarginRight(right).Render(inner)
			case n.MarginLeftAuto:
				inner = lipgloss.NewStyle().MarginLeft(spare).Render(inner)
			case n.MarginRightAuto:
				inner = lipgloss.NewStyle().MarginRight(spare).Render(inner)
			}
		}
	}

	return inner
}

func (n *Node) renderChildren(w, h int) string {
	isRow := n.Dir == Row

	// determine main axis size
	total := w
	if !isRow {
		total = h
	}

	// ⚠️ warning: flex + space-between conflict
	totalFlex := 0
	for _, c := range n.Children {
		if c.Flex > 0 {
			totalFlex += c.Flex
		}
	}

	if totalFlex > 0 && n.Justify == JustifySpaceBetween {
		// non-fatal but very useful for debugging layouts
		// you can replace with log.Println if you prefer
		var warnOnce sync.Once
		warnOnce.Do(func() {
			fmt.Println("[flexgo warning] Flex + SpaceBetween in same container is discouraged")
		})
	}

	// compute sizes (flex distribution)
	sizes := distribute(total, n.Children, isRow, n.Gap)

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

		// min/max constraints
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

		childView := child.Render(cw, ch)

		// cross-axis alignment
		childView = applyAlign(childView, cw, ch, n.Align)

		parts = append(parts, childView)
	}

	return join(parts, n.Dir, n.Justify, w, h, n.Gap)
}

func join(parts []string, dir Direction, justify Justify, totalW, totalH int, gap int) string {
	if len(parts) == 0 {
		return ""
	}

	// Simple case first
	if justify == JustifyStart {
		if gap > 0 {
			return interleave(parts, dir, gap)
		}
		if dir == Row {
			return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
		}
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	// Compute total size
	totalSize := 0
	for _, p := range parts {
		if dir == Row {
			totalSize += lipgloss.Width(p)
		} else {
			totalSize += lipgloss.Height(p)
		}
	}

	containerSize := totalW
	if dir == Col {
		containerSize = totalH
	}

	remaining := containerSize - totalSize

	// distribute space
	switch justify {

	case JustifyCenter:
		pad := remaining / 2
		return spacer(dir, pad) + concat(parts, dir) + spacer(dir, pad)

	case JustifyEnd:
		return spacer(dir, remaining) + concat(parts, dir)

	case JustifySpaceBetween:
		if len(parts) == 1 {
			return parts[0]
		}
		space := remaining / (len(parts) - 1)
		return interleave(parts, dir, space)
	}

	return concat(parts, dir)
}

func spacer(dir Direction, size int) string {
	if size <= 0 {
		return ""
	}

	if dir == Row {
		return lipgloss.NewStyle().Width(size).Render("")
	}
	return lipgloss.NewStyle().Height(size).Render("")
}

func concat(parts []string, dir Direction) string {
	if dir == Row {
		return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func interleave(parts []string, dir Direction, gap int) string {
	var out []string

	for i, p := range parts {
		out = append(out, p)
		if i < len(parts)-1 {
			out = append(out, spacer(dir, gap))
		}
	}

	return concat(out, dir)
}
