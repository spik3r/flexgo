package layouts

import (
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// centeredView returns a View callback that renders text centred in
// its allocated box. Used by several recipes (Tabs, Modal, Form).
func centeredView(text string) func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(text)
	}
}

func alignPtr(a flexgo.Align) *flexgo.Align {
	return &a
}
