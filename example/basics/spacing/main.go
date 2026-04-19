// Demonstrates Gap, Padding, and Margin spacing controls.
//
//	Gap     — space inserted between siblings in the main direction.
//	Padding — inner space between a container's edge and its children.
//	Margin  — outer space that shrinks a node within its allocated slot.
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

func sectionHeader(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Bold(true).Padding(0, 1).
			Align(lipgloss.Left, lipgloss.Center).
			Render(text)
	}
}

func initialModel() model {
	// ── Gap section ──────────────────────────────────────────────────────────
	// Gap reserves space between siblings. The children are sized smaller so
	// that gap + children = total. Use JustifyStart so gaps are interleaved.
	gapSection := &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 3, View: sectionHeader("Gap: 3  — space between siblings")},
			{
				Dir:        flexgo.Row,
				Flex:       1,
				Gap:        3,
				Background: lipgloss.Color("237"),
				Children: []*flexgo.Node{
					{Flex: 1, Background: lipgloss.Color("61"), View: cell("A")},
					{Flex: 1, Background: lipgloss.Color("61"), View: cell("B")},
					{Flex: 1, Background: lipgloss.Color("61"), View: cell("C")},
				},
			},
		},
	}

	// ── Padding section ───────────────────────────────────────────────────────
	// Padding shrinks the content area inside a container.
	paddingSection := &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 3, View: sectionHeader("Padding: 2  — inner space")},
			{
				Dir:        flexgo.Col,
				Flex:       1,
				Padding:    2,
				Background: lipgloss.Color("237"),
				Children: []*flexgo.Node{
					{
						Dir:  flexgo.Row,
						Flex: 1,
						Children: []*flexgo.Node{
							{Flex: 1, Background: lipgloss.Color("61"), View: cell("A")},
							{Flex: 1, Background: lipgloss.Color("61"), View: cell("B")},
							{Flex: 1, Background: lipgloss.Color("61"), View: cell("C")},
						},
					},
				},
			},
		},
	}

	// ── Margin section ────────────────────────────────────────────────────────
	// Margin shrinks a node within its allocated slot, creating outer space.
	// Each child renders smaller than its slot with whitespace around it.
	marginSection := &flexgo.Node{
		Dir:  flexgo.Col,
		Flex: 1,
		Children: []*flexgo.Node{
			{Height: 3, View: sectionHeader("Margin: 1  — outer space around each node")},
			{
				Dir:        flexgo.Row,
				Flex:       1,
				Background: lipgloss.Color("237"),
				Children: []*flexgo.Node{
					{Flex: 1, Margin: 1, Background: lipgloss.Color("61"), View: cell("A")},
					{Flex: 1, Margin: 1, Background: lipgloss.Color("61"), View: cell("B")},
					{Flex: 1, Margin: 1, Background: lipgloss.Color("61"), View: cell("C")},
				},
			},
		},
	}

	page := &flexgo.Node{
		Dir: flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 3, View: sectionHeader("Spacing Controls  (q to quit)")},
			gapSection,
			paddingSection,
			marginSection,
		},
	}
	return model{page: page}
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
