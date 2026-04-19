package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

// maxPageWidth is the fixed width of the centered page.
const maxPageWidth = 120

type model struct {
	width  int
	height int
	page   *flexgo.Node
	ready  bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading...")
	}
	// Pass the full terminal width so the page can center itself.
	v := tea.NewView(m.page.Render(m.width, m.height))
	v.AltScreen = true
	return v
}

func box(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

func initialModel() model {
	header := &flexgo.Node{
		Dir:    flexgo.Row,
		Height: 3,
		View:   box("HEADER"),
	}

	sidebar := &flexgo.Node{Flex: 3, View: box("SIDEBAR")}
	main := &flexgo.Node{Flex: 7, View: box("MAIN")}

	content := &flexgo.Node{
		Dir:      flexgo.Row,
		Flex:     1,
		Children: []*flexgo.Node{sidebar, main},
	}

	footer := &flexgo.Node{
		Dir:    flexgo.Row,
		Height: 3,
		View:   box("FOOTER"),
	}

	// Width caps the page at maxPageWidth chars.
	// MarginLeftAuto + MarginRightAuto centers it within whatever
	// terminal width is passed to Render().
	page := &flexgo.Node{
		Dir:             flexgo.Col,
		Width:           maxPageWidth,
		MarginLeftAuto:  true,
		MarginRightAuto: true,
		Children:        []*flexgo.Node{header, content, footer},
	}

	return model{page: page}
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(initialModel().page.Render(80, 24))
		return
	}

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
