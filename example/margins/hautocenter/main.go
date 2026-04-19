// Demonstrates horizontal auto-margin centering.
//
// MarginLeftAuto + MarginRightAuto center a node horizontally by filling
// the spare horizontal space on each side. They only activate when the
// node's natural width is smaller than its allocated width — so this demo
// wraps the centered node in a Col parent, which does not exact-size the
// child on the horizontal axis.
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
	centered := &flexgo.Node{
		Width:           30,
		Height:          5,
		MarginLeftAuto:  true,
		MarginRightAuto: true,
		Background:      lipgloss.Color("61"),
		View:            box("HORIZONTALLY CENTERED"),
	}

	// Col parent distributes Height, so the child's allocated Width is the
	// parent's full width — that's what gives MarginLeftAuto/RightAuto room
	// to centre the 30-wide box. A Row parent would exact-size the child
	// to 30 columns and leave no spare space.
	page := &flexgo.Node{
		Dir:        flexgo.Col,
		Background: lipgloss.Color("237"),
		Children:   []*flexgo.Node{centered},
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
