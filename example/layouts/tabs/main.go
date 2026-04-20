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
	active        int
}

func initialModel() model { return model{active: 1} }

func buildRoot(active int) *flexgo.Node {
	root := layouts.Tabs(active, []layouts.Tab{
		{Title: "Logs", Panel: panel("Live logs stream", lipgloss.Color("238"))},
		{Title: "Metrics", Panel: panel("CPU 14%  RAM 38%", lipgloss.Color("59"))},
		{Title: "Alerts", Panel: panel("No active alerts", lipgloss.Color("95"))},
	})
	root.Background = lipgloss.Color("236")
	root.Paddings = flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}
	root.Children[0].Gap = 1
	return root
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ready = true
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "left":
			m.active = max(0, m.active-1)
		case "l", "right":
			m.active = min(2, m.active+1)
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading...")
	}
	root := buildRoot(m.active)
	v := tea.NewView(root.Render(m.width, m.height))
	v.AltScreen = true
	return v
}

func panel(text string, bg color.Color) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Render(text)
	}
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(buildRoot(1).Render(80, 24))
		return
	}
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
