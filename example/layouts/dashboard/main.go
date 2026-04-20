package main

import (
	"fmt"
	"image/color"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
	"github.com/spik3r/flexgo/layouts"
)

type model struct {
	width, height int
	ready         bool
	root          *flexgo.Node
}

func initialModel() model {
	root := layouts.WithBackground(
		layouts.Dashboard(
			22,
			3,
			1,
			section("SIDEBAR", lipgloss.Color("59")),
			section("HEADER", lipgloss.Color("61")),
			section("MAIN", lipgloss.Color("238")),
			section("STATUS: all systems nominal", lipgloss.Color("66")),
		),
		lipgloss.Color("236"),
	)
	root.Gap = 1
	root.Paddings = flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}
	root.Children[1].Gap = 1
	return model{root: root}
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
	v := tea.NewView(m.root.Render(m.width, m.height))
	v.AltScreen = true
	return v
}

func section(title string, bg color.Color) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(title)
	}
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(initialModel().root.Render(80, 24))
		return
	}
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
