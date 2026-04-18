package flexgo

import "charm.land/lipgloss/v2"

func applyAlign(view string, w, h int, align Align) string {
	style := lipgloss.NewStyle().Width(w).Height(h)

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
