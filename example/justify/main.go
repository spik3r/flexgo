// Demonstrates all four Justify modes.
// Justify controls where children are placed along the main axis when there
// is remaining space (i.e. children have fixed widths, not Flex).
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

func cell(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

func sectionLabel(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Left, lipgloss.Center).
			Render(text)
	}
}

// section builds one demo row: a label + a Row of three fixed-width boxes
// laid out with the given Justify mode.
func section(title string, justify flexgo.Justify) *flexgo.Node {
	return &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 3, View: sectionLabel(title)},
			{
				Dir:     flexgo.Row,
				Flex:    1,
				Justify: justify,
				Children: []*flexgo.Node{
					// Fixed-width children — the remaining space is what Justify distributes.
					{Width: 20, View: cell("A")},
					{Width: 20, View: cell("B")},
					{Width: 20, View: cell("C")},
				},
			},
		},
	}
}

func initialModel() model {
	page := &flexgo.Node{
		Dir: flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 3, View: sectionLabel("Justify Modes  (q to quit)")},
			section("JustifyStart  — children packed to the left", flexgo.JustifyStart),
			section("JustifyCenter — children centered", flexgo.JustifyCenter),
			section("JustifyEnd    — children packed to the right", flexgo.JustifyEnd),
			section("JustifySpaceBetween — space distributed between children", flexgo.JustifySpaceBetween),
		},
	}
	return model{page: page}
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
