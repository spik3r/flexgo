// Demonstrates vertical auto-margin centering.
//
// MarginTopAuto + MarginBottomAuto center a node vertically by filling
// the spare vertical space above and below. They only activate when the
// node's natural height is smaller than its allocated height — so this
// demo wraps the centered node in a Row parent, which does not
// exact-size the child on the vertical axis (Row distributes Width).
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
		Width:            30,
		Height:           7,
		MarginTopAuto:    true,
		MarginBottomAuto: true,
		Background:       lipgloss.Color("61"),
		View:             box("VERTICALLY\nCENTERED"),
	}

	// Row parent distributes Width, so the child's allocated Height is
	// the parent's full height — that's what gives MarginTopAuto /
	// MarginBottomAuto room to centre the 7-tall box. A Col parent would
	// exact-size the child to 7 rows and leave no spare space.
	page := &flexgo.Node{
		Dir:        flexgo.Row,
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
