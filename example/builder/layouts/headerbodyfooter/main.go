// Builder version of layouts.HeaderBodyFooter — the classic three-row
// TUI shape expressed via the fluent NodeBuilder API.
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
	root := flexgo.NewNode().
		Dir(flexgo.Col).
		Background(lipgloss.Color("236")).
		Children(
			flexgo.NewNode().
				Height(3).
				Background(lipgloss.Color("61")).
				View(banner("HEADER")).
				Build(),
			flexgo.NewNode().
				Flex(1).
				Background(lipgloss.Color("238")).
				View(banner("BODY")).
				Build(),
			flexgo.NewNode().
				Height(2).
				Background(lipgloss.Color("66")).
				View(banner("FOOTER")).
				Build(),
		).
		Build()

	return model{page: root}
}

func banner(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(text)
	}
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
