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
		Children(
			flexgo.NewNode().
				Height(3).
				Justify(flexgo.JustifyCenter).
				Background(lipgloss.Color("237")).
				View(box("BUILDER HEADER")).
				Build(),
			flexgo.NewNode().
				Dir(flexgo.Row).
				Flex(1).
				Gap(2).
				Paddings(flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}).
				Background(lipgloss.Color("237")).
				Children(
					flexgo.NewNode().
						Flex(2).
						Background(lipgloss.Color("61")).
						View(box("NAV")).
						Build(),
					flexgo.NewNode().
						Flex(5).
						Background(lipgloss.Color("69")).
						View(box("CONTENT")).
						Build(),
				).
				Build(),
			flexgo.NewNode().
				Height(3).
				Justify(flexgo.JustifyEnd).
				Background(lipgloss.Color("237")).
				View(box("BUILDER FOOTER")).
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

func box(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
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
