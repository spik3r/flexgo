// Builder version of layouts.Grid — a uniform NxM grid where every
// cell gets Flex:1 on both axes, built with nested NodeBuilder loops.
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

const (
	rows = 3
	cols = 4
	gap  = 1
)

func initialModel() model {
	rowNodes := make([]*flexgo.Node, 0, rows)
	for r := 0; r < rows; r++ {
		cells := make([]*flexgo.Node, 0, cols)
		for c := 0; c < cols; c++ {
			rr, cc := r, c
			cells = append(cells, flexgo.NewNode().
				Flex(1).
				View(cellView(rr, cc)).
				Build())
		}
		rowNodes = append(rowNodes, flexgo.NewNode().
			Flex(1).
			Dir(flexgo.Row).
			Gap(gap).
			Children(cells...).
			Build())
	}

	root := flexgo.NewNode().
		Dir(flexgo.Col).
		Gap(gap).
		Background(lipgloss.Color("236")).
		Paddings(flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}).
		Children(rowNodes...).
		Build()

	return model{page: root}
}

func cellView(r, c int) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(lipgloss.Color("61")).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(fmt.Sprintf("%d,%d", r, c))
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
