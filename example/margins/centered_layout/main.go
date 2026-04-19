// Demonstrates both horizontal and vertical auto-margin centering.
//
// MarginLeftAuto/RightAuto and MarginTopAuto/BottomAuto only activate
// when the node's natural size is smaller than its allocation. Both axes
// at once requires that the node not be exact-sized by any parent — so
// this example renders the centered node as the root directly, at the
// full terminal size. The terminal itself provides the backdrop.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

type model struct {
	width, height int
	ready         bool
	page          *flexgo.Node
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading...")
	}
	v := tea.NewView(m.page.Render(m.width, m.height))
	v.AltScreen = true
	return v
}

func box(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

func initialModel() model {
	// Root node: all four auto-margins active.
	// Rendered directly at the terminal size, so availW > Width and
	// availH > Height — spare space on both axes centres the box.
	page := &flexgo.Node{
		Width:            40,
		Height:           8,
		MarginLeftAuto:   true,
		MarginRightAuto:  true,
		MarginTopAuto:    true,
		MarginBottomAuto: true,
		Background:       lipgloss.Color("61"),
		View:             box("CENTERED"),
	}

	return model{page: page}
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(initialModel().page.Render(80, 24))
		return
	}

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
