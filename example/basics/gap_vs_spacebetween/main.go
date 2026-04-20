// Side-by-side comparison of Gap and JustifySpaceBetween — the most
// common flexbox confusion.
//
//	Gap                 — fixed amount of space between siblings, at
//	                      both ends packed against the start.
//	JustifySpaceBetween — first child at the start, last child at the
//	                      end, remaining siblings evenly spaced in
//	                      between; the "gap" scales with container size.
//	Gap + SpaceBetween  — SpaceBetween treats Gap as a minimum; extra
//	                      slack is still distributed between children.
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

func box(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(lipgloss.Color("61")).
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
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

func row(title string, gap int, justify flexgo.Justify) *flexgo.Node {
	cells := []*flexgo.Node{
		{Width: 10, View: box("A")},
		{Width: 10, View: box("B")},
		{Width: 10, View: box("C")},
	}
	return &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 2, View: header(title)},
			{
				Dir:        flexgo.Row,
				Flex:       1,
				Gap:        gap,
				Justify:    justify,
				Background: lipgloss.Color("237"),
				Children:   cells,
			},
		},
	}
}

func buildRoot() *flexgo.Node {
	return &flexgo.Node{
		Dir:        flexgo.Col,
		Background: lipgloss.Color("236"),
		Paddings:   flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2},
		Gap:        1,
		Children: []*flexgo.Node{
			{Height: 2, View: header("Gap vs JustifySpaceBetween  (q to quit)")},
			row("Gap: 3 (JustifyStart)", 3, flexgo.JustifyStart),
			row("JustifySpaceBetween (no gap)", 0, flexgo.JustifySpaceBetween),
			row("Gap: 3 + JustifySpaceBetween (gap is a min)", 3, flexgo.JustifySpaceBetween),
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
