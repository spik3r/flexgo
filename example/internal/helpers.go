package internal

import (
	"charm.land/lipgloss/v2"
)

func Box(label string) func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

func Cell(label string) func(w, h int) string {
	return Box(label)
}

func SectionHeader(text string) func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).Padding(0, 1).
			Align(lipgloss.Left, lipgloss.Center).
			Render(text)
	}
}
