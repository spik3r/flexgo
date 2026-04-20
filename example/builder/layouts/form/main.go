// Builder version of layouts.Form — aligned Label:Field rows built
// from a loop of NodeBuilder calls.
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

type field struct {
	label, value string
}

func initialModel() model {
	fields := []field{
		{"Name", "Ada Lovelace"},
		{"Email", "ada@example.com"},
		{"Region", "eu-west-1"},
		{"Team", "Platform"},
	}

	const labelWidth = 12

	rows := make([]*flexgo.Node, 0, len(fields))
	for _, f := range fields {
		rows = append(rows, flexgo.NewNode().
			Dir(flexgo.Row).
			Height(1).
			Gap(1).
			Children(
				flexgo.NewNode().
					Width(labelWidth).
					View(labelView(f.label+":")).
					Build(),
				flexgo.NewNode().
					Flex(1).
					View(inputView(f.value)).
					Build(),
			).
			Build())
	}

	root := flexgo.NewNode().
		Dir(flexgo.Col).
		Gap(1).
		Paddings(flexgo.Spacing{Top: 2, Right: 4, Bottom: 2, Left: 4}).
		Background(lipgloss.Color("236")).
		Children(rows...).
		Build()

	return model{page: root}
}

func labelView(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).
			Foreground(lipgloss.Color("250")).
			Align(lipgloss.Right, lipgloss.Center).
			Render(text)
	}
}

func inputView(value string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Padding(0, 1).
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("238")).
			Render(value)
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
