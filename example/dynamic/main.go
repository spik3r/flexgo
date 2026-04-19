package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

type model struct {
	width  int
	height int
	page   *flexgo.Node
	ready  bool
}

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
		Dir:     flexgo.Row,
		Height:  3,
		Justify: flexgo.JustifyCenter,
		View:    box("HEADER"),
	}

	content := &flexgo.Node{
		Dir:     flexgo.Row,
		Flex:    1,
		Justify: flexgo.JustifyStart,
		Children: []*flexgo.Node{
			{Flex: 3, View: box("LEFT")},
			{Flex: 7, View: box("RIGHT")},
		},
	}

	statusline := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  3,
		Justify: flexgo.JustifySpaceBetween,
		Children: []*flexgo.Node{
			{Width: 56, View: box("STATUS INFO")},
			{Width: 16, View: box("1.0.0-alpha")},
		},
	}

	contentContainer := &flexgo.Node{
		Dir:     flexgo.Col,
		Flex:    1,
		Justify: flexgo.JustifyStart,
		Children: []*flexgo.Node{
			statusline,
			content,
		},
	}

	footer := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  3,
		Justify: flexgo.JustifyCenter,
		View:    box("Footer"),
	}

	page := &flexgo.Node{
		Dir:     flexgo.Col,
		Justify: flexgo.JustifyStart,
		// Debug:   true,
		Children: []*flexgo.Node{
			header,
			contentContainer,
			footer,
		},
	}

	return model{page: page, ready: true}
}

func (m model) Init() tea.Cmd {
	return nil
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
