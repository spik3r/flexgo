package flexgo

import (
	"image/color"
	"strconv"

	"charm.land/lipgloss/v2"
)

func resolveSpacing(shorthand int, sides Spacing) Spacing {
	if sides != (Spacing{}) {
		return sides
	}
	if shorthand == 0 {
		return Spacing{}
	}
	return Spacing{Top: shorthand, Right: shorthand, Bottom: shorthand, Left: shorthand}
}

func hasExplicitBorder(n *Node) bool {
	return n.Border != (lipgloss.Border{})
}

func debugLabel(name string, w, h int) string {
	label := ""
	if name != "" {
		label = name + " "
	}
	return label + "(" + strconv.Itoa(w) + "," + strconv.Itoa(h) + ")"
}

func clampMainAxis(child *Node, size int, isRow bool) int {
	if isRow {
		if child.MinWidth > 0 && size < child.MinWidth {
			size = child.MinWidth
		}
		if child.MaxWidth > 0 && size > child.MaxWidth {
			size = child.MaxWidth
		}
		return size
	}

	if child.MinHeight > 0 && size < child.MinHeight {
		size = child.MinHeight
	}
	if child.MaxHeight > 0 && size > child.MaxHeight {
		size = child.MaxHeight
	}
	return size
}

func resolveMainAxisSizes(total int, children []*Node, isRow bool, gap int) []int {
	sizes := distribute(total, children, isRow, gap)
	if len(children) == 0 {
		return sizes
	}

	totalContent := total - gap*(len(children)-1)
	if totalContent < 0 {
		totalContent = 0
	}

	frozen := make([]bool, len(children))

	for {
		changed := false

		for i, child := range children {
			if frozen[i] {
				continue
			}
			clamped := clampMainAxis(child, sizes[i], isRow)
			if clamped != sizes[i] {
				sizes[i] = clamped
				frozen[i] = true
				changed = true
			}
		}

		if !changed {
			break
		}

		used := 0
		remainingChildren := make([]*Node, 0, len(children))
		remainingIndices := make([]int, 0, len(children))
		for i, child := range children {
			if frozen[i] {
				used += sizes[i]
				continue
			}
			remainingChildren = append(remainingChildren, child)
			remainingIndices = append(remainingIndices, i)
		}

		if len(remainingChildren) == 0 {
			break
		}

		remainingContent := totalContent - used
		if remainingContent < 0 {
			remainingContent = 0
		}

		redistributed := distribute(remainingContent, remainingChildren, isRow, 0)
		for i, idx := range remainingIndices {
			sizes[idx] = redistributed[i]
		}
	}

	return sizes
}

func join(parts []string, dir Direction, justify Justify, totalW, totalH int, gap int, bg color.Color) string {
	if len(parts) == 0 {
		return ""
	}

	cross := totalH
	if dir == Col {
		cross = totalW
	}

	if justify == JustifyStart {
		if gap > 0 {
			return interleave(parts, dir, gap, bg, cross)
		}
		return concat(parts, dir)
	}

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

	switch justify {
	case JustifyCenter:
		pad := remaining / 2
		centered := interleave(parts, dir, gap, bg, cross)
		return concat([]string{spacer(dir, pad, bg, cross), centered, spacer(dir, pad, bg, cross)}, dir)
	case JustifyEnd:
		withGap := interleave(parts, dir, gap, bg, cross)
		return concat([]string{spacer(dir, remaining, bg, cross), withGap}, dir)
	case JustifySpaceBetween:
		if len(parts) == 1 {
			return parts[0]
		}
		space := remaining / (len(parts) - 1)
		return interleave(parts, dir, space, bg, cross)
	default:
		return concat(parts, dir)
	}
}

func spacer(dir Direction, size int, bg color.Color, cross int) string {
	if size <= 0 {
		return ""
	}

	style := lipgloss.NewStyle()
	if bg != nil {
		style = style.Background(bg)
	}
	if dir == Row {
		style = style.Width(size)
		if cross > 0 {
			style = style.Height(cross)
		}
	} else {
		style = style.Height(size)
		if cross > 0 {
			style = style.Width(cross)
		}
	}
	return style.Render("")
}

func asymmetricMargin(content string, top, right, bottom, left int, bg color.Color) string {
	if top == 0 && right == 0 && bottom == 0 && left == 0 {
		return content
	}

	contentW := lipgloss.Width(content)
	contentH := lipgloss.Height(content)

	bgStyle := lipgloss.NewStyle()
	if bg != nil {
		bgStyle = bgStyle.Background(bg)
	}

	middleParts := []string{}
	if left > 0 {
		middleParts = append(middleParts, bgStyle.Width(left).Height(contentH).Render(""))
	}
	middleParts = append(middleParts, content)
	if right > 0 {
		middleParts = append(middleParts, bgStyle.Width(right).Height(contentH).Render(""))
	}
	middle := lipgloss.JoinHorizontal(lipgloss.Top, middleParts...)

	totalW := contentW + left + right

	parts := []string{}
	if top > 0 {
		parts = append(parts, bgStyle.Width(totalW).Height(top).Render(""))
	}
	parts = append(parts, middle)
	if bottom > 0 {
		parts = append(parts, bgStyle.Width(totalW).Height(bottom).Render(""))
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func concat(parts []string, dir Direction) string {
	if dir == Row {
		return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func interleave(parts []string, dir Direction, gap int, bg color.Color, cross int) string {
	var out []string

	for i, p := range parts {
		out = append(out, p)
		if i < len(parts)-1 {
			out = append(out, spacer(dir, gap, bg, cross))
		}
	}

	return concat(out, dir)
}
