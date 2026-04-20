// Visual reference for Debug mode and the Inspect() tree dump. The
// page renders a small nested layout with Debug on every container,
// then prints Inspect(root) into the footer leaf so you can compare
// the rendered frames to the structural tree.
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
	v := tea.NewView(buildRoot().Render(m.width, m.height))
	v.AltScreen = true
	return v
}

func leaf(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

func buildRoot() *flexgo.Node {
	inner := &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Name: "inner",
		Gap:  1,
		Children: []*flexgo.Node{
			{Height: 3, Name: "header", View: leaf("HEADER")},
			{
				Dir:  flexgo.Row,
				Flex: 1,
				Name: "body",
				Gap:  1,
				Children: []*flexgo.Node{
					{Flex: 1, Name: "left", View: leaf("left")},
					{Flex: 2, Name: "right", View: leaf("right")},
				},
			},
		},
	}

	tree := &flexgo.Node{
		Dir:      flexgo.Col,
		Flex:     1,
		Name:     "root",
		Paddings: flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2},
		Children: []*flexgo.Node{inner},
	}
	flexgo.DebugAll(tree)

	footer := &flexgo.Node{
		Height: 10,
		View: func(w, h int) string {
			dump := flexgo.Inspect(tree)
			return lipgloss.NewStyle().
				Width(w).Height(h).
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("252")).
				Padding(0, 2).
				Render("Inspect(root):\n\n" + dump)
		},
	}

	return &flexgo.Node{
		Dir: flexgo.Col,
		Children: []*flexgo.Node{
			tree,
			footer,
		},
	}
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(buildRoot().Render(80, 30))
		return
	}
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
