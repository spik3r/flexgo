package main

import (
	"fmt"
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
	root := layouts.Grid(3, 4, 1, func(r, c, w, h int) string {
		bg := lipgloss.Color("60")
		if (r+c)%2 == 1 {
			bg = lipgloss.Color("59")
		}
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Render(fmt.Sprintf("%d,%d", r+1, c+1))
	})
	root.Paddings = flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}
	root.Background = lipgloss.Color("236")
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
