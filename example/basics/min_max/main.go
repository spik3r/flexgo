// Demonstrates MinWidth / MaxWidth constraints.
//
// Left row: flex children with MinWidth — when the container shrinks,
// children refuse to go below the minimum and the remainder is
// redistributed to the unconstrained child.
//
// Right row: flex children with MaxWidth — when the container grows,
// children stop at the cap and extra space piles onto the
// unconstrained child.
package main

import (
	"fmt"
	"image/color"
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

func cell(label string, bg color.Color) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(fmt.Sprintf("%s\n%dx%d", label, w, h))
	}
}

func header(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).
			Padding(0, 1).
			Render(text)
	}
}

func buildRoot() *flexgo.Node {
	minSection := &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 2, View: header("MinWidth  (A/B min 20; C unbounded)")},
			{
				Dir:  flexgo.Row,
				Flex: 1,
				Gap:  1,
				Children: []*flexgo.Node{
					{Flex: 1, MinWidth: 20, View: cell("A min=20", lipgloss.Color("61"))},
					{Flex: 1, MinWidth: 20, View: cell("B min=20", lipgloss.Color("62"))},
					{Flex: 1, View: cell("C free", lipgloss.Color("66"))},
				},
			},
		},
	}

	maxSection := &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 2, View: header("MaxWidth  (A/B max 15; C soaks up the rest)")},
			{
				Dir:  flexgo.Row,
				Flex: 1,
				Gap:  1,
				Children: []*flexgo.Node{
					{Flex: 1, MaxWidth: 15, View: cell("A max=15", lipgloss.Color("61"))},
					{Flex: 1, MaxWidth: 15, View: cell("B max=15", lipgloss.Color("62"))},
					{Flex: 1, View: cell("C free", lipgloss.Color("66"))},
				},
			},
		},
	}

	return &flexgo.Node{
		Dir:        flexgo.Col,
		Background: lipgloss.Color("236"),
		Paddings:   flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2},
		Gap:        1,
		Children: []*flexgo.Node{
			{Height: 2, View: header("MinWidth / MaxWidth  (q to quit)")},
			minSection,
			maxSection,
		},
	}
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(buildRoot().Render(80, 24))
		return
	}
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
