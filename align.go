package flexgo

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

func applyAlign(view string, w, h int, align Align, bg color.Color) string {
	style := lipgloss.NewStyle().Width(w).Height(h)
	if bg != nil {
		style = style.Background(bg)
	}

	switch align {
	case AlignCenter:
		style = style.Align(lipgloss.Center, lipgloss.Center)
	case AlignEnd:
		style = style.Align(lipgloss.Right, lipgloss.Bottom)
	default:
		style = style.Align(lipgloss.Left, lipgloss.Top)
	}

	return style.Render(view)
}
