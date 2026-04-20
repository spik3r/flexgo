// Builder version of layouts.Modal — a centred, bordered dialog built
// from auto-margins + fixed size via NodeBuilder.
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

func initialModel() model {
	bg := lipgloss.Color("237")

	titleView := func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Confirm action")
	}
	bodyView := func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(bg).
			Foreground(lipgloss.Color("252")).
			Padding(1, 2).
			Render("Discard 3 unsaved changes?\n\n[y] Yes   [n] No")
	}

	root := flexgo.NewNode().
		Dir(flexgo.Col).
		Width(40).
		Height(10).
		MarginTopAuto(true).
		MarginBottomAuto(true).
		MarginLeftAuto(true).
		MarginRightAuto(true).
		Border(lipgloss.RoundedBorder()).
		Background(bg).
		Children(
			flexgo.NewNode().Height(1).View(titleView).Build(),
			flexgo.NewNode().Flex(1).View(bodyView).Build(),
		).
		Build()

	return model{page: root}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
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
