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
	input := func(value string) func(int, int) string {
		return func(w, h int) string {
			return lipgloss.NewStyle().
				Width(w).
				Height(h).
				Padding(0, 1).
				Foreground(lipgloss.Color("252")).
				Background(lipgloss.Color("238")).
				Render(value)
		}
	}

	root := layouts.Form(12, []layouts.FormRow{
		{Label: "Name", Field: input("Ada Lovelace")},
		{Label: "Email", Field: input("ada@example.com")},
		{Label: "Region", Field: input("eu-west-1")},
		{Label: "Team", Field: input("Platform")},
	})

	root.Gap = 1
	root.Paddings = flexgo.Spacing{Top: 2, Right: 4, Bottom: 2, Left: 4}
	root.Background = lipgloss.Color("236")
	for _, row := range root.Children {
		row.Children[0].View = labelView(row.Children[0].View)
	}

	return model{root: root}
}

func labelView(base func(int, int) string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			Bold(true).
			Render(base(w, h))
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
