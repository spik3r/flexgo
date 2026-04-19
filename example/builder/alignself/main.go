package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

type model struct {
	width  int
	height int
	ready  bool
	page   *flexgo.Node
}

func initialModel() model {
	root := flexgo.NewNode().
		Dir(flexgo.Col).
		Paddings(flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}).
		Background(lipgloss.Color("236")).
		Children(
			flexgo.NewNode().
				Height(3).
				View(panel("Builder AlignSelf + Border")).
				Build(),
			flexgo.NewNode().
				Dir(flexgo.Row).
				Flex(1).
				Align(flexgo.AlignStart).
				Gap(1).
				Children(
					flexgo.NewNode().
						Flex(1).
						Height(8).
						Border(lipgloss.NormalBorder()).
						BorderForeground(lipgloss.Color("240")).
						Background(lipgloss.Color("61")).
						View(panel("Top (inherit AlignStart) ")).
						Build(),
					flexgo.NewNode().
						Flex(1).
						Height(8).
						Border(lipgloss.NormalBorder()).
						BorderForeground(lipgloss.Color("240")).
						Background(lipgloss.Color("69")).
						AlignSelf(flexgo.AlignCenter).
						View(panel("Center (AlignSelf) ")).
						Build(),
					flexgo.NewNode().
						Flex(1).
						Height(8).
						Border(lipgloss.NormalBorder()).
						BorderForeground(lipgloss.Color("240")).
						Background(lipgloss.Color("75")).
						AlignSelf(flexgo.AlignEnd).
						View(panel("Bottom (AlignSelf) ")).
						Build(),
				).
				Build(),
		).
		Build()

	return model{page: root}
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

func panel(title string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(title)
	}
}

func render(w, h int) string {
	return initialModel().page.Render(w, h)
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(render(80, 24))
		return
	}

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
