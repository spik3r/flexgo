// Builder version of layouts.Dashboard — sidebar + header + main +
// status, constructed via NodeBuilder.
package main

import (
	"fmt"
	"image/color"
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
	right := flexgo.NewNode().
		Flex(1).
		Dir(flexgo.Col).
		Gap(1).
		Children(
			flexgo.NewNode().Height(3).View(section("HEADER", lipgloss.Color("61"))).Build(),
			flexgo.NewNode().Flex(1).View(section("MAIN", lipgloss.Color("238"))).Build(),
			flexgo.NewNode().Height(1).View(section("STATUS: all systems nominal", lipgloss.Color("66"))).Build(),
		).
		Build()

	root := flexgo.NewNode().
		Dir(flexgo.Row).
		Background(lipgloss.Color("236")).
		Paddings(flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}).
		Gap(1).
		Children(
			flexgo.NewNode().
				Width(22).
				View(section("SIDEBAR", lipgloss.Color("59"))).
				Build(),
			right,
		).
		Build()

	return model{page: root}
}

func section(title string, bg color.Color) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(title)
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
