// Demonstrates cross-axis Align modes and a 3-column layout.
//
// Align controls where children sit on the cross axis (vertical for Row
// containers, horizontal for Col containers). Children must render shorter
// than the container's cross-axis size for the positioning to be visible —
// fixedCell() achieves this by using a fixed height regardless of the
// allocated height.
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

func header(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).Padding(0, 1).
			Align(lipgloss.Left, lipgloss.Center).
			Render(text)
	}
}

// fixedCell renders a bordered box at a fixed height regardless of the
// allocated height. The parent's Align then positions this shorter box
// at the top, center, or bottom of the container.
func fixedCell(label string, fixedH int) func(int, int) string {
	return func(w, h int) string {
		height := fixedH
		if h > 0 && h < fixedH {
			height = h // don't overflow the container
		}
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

// column builds one of the three demo columns.
// The inner Row has the given Align, and its children are fixed-height boxes
// so the alignment positioning is clearly visible.
func column(title string, align flexgo.Align) *flexgo.Node {
	return &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 3, View: header(title)},
			{
				Dir:   flexgo.Row,
				Flex:  1,
				Align: align,
				Children: []*flexgo.Node{
					// Each child renders at its own fixed height (10 or 5 lines).
					// The parent Align places them at top / center / bottom.
					{Flex: 1, View: fixedCell("TALL", 10)},
					{Flex: 1, View: fixedCell("SHORT", 5)},
					{Flex: 1, View: fixedCell("TALL", 10)},
				},
			},
		},
	}
}

func initialModel() model {
	page := &flexgo.Node{
		Dir: flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 3, View: header("Align Modes — 3-column layout  (q to quit)")},
			// Three columns side by side — each a Col with Flex:1.
			{
				Dir:  flexgo.Row,
				Flex: 1,
				Children: []*flexgo.Node{
					column("AlignStart  (top)", flexgo.AlignStart),
					column("AlignCenter (middle)", flexgo.AlignCenter),
					column("AlignEnd    (bottom)", flexgo.AlignEnd),
				},
			},
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
