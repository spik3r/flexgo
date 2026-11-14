package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
	"github.com/spik3r/flexgo/layouts"
)

const tileCount = 7

type model struct {
	width, height int
	root          *flexgo.Node
}

func label(text string) func(int, int) string {
	return flexgo.ClipView(func(w, h int) string { return text })
}

func tile(index int) func(int, int) string {
	return func(w, h int) string {
		if w <= 0 || h <= 0 {
			return ""
		}
		text := fmt.Sprintf("Tile %02d\nSynthetic local data", index+1)
		return lipgloss.NewStyle().Width(w).Height(h).
			Background(lipgloss.Color("60")).Foreground(lipgloss.Color("230")).
			Render(label(text)(w, h))
	}
}

func initialModel() model {
	cells := make([]func(int, int) string, tileCount)
	for i := range cells {
		cells[i] = tile(i)
	}
	grid := layouts.AutoGrid(14, 1, cells...)
	root := &flexgo.Node{Dir: flexgo.Col, Padding: 1, Background: lipgloss.Color("236"), Children: []*flexgo.Node{
		{Height: 1, View: label("Responsive cards")},
		grid,
		{Height: 1, View: label("Resize: reflow | q: quit")},
	}}
	return model{root: root}
}

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = max(0, msg.Width), max(0, msg.Height)
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m model) View() tea.View {
	var content string
	// Seven one-column tiles need 7 content rows + 6 gutters + 4 frame rows.
	if m.width < 18 || m.height < 17 {
		content = label("Resize to 18x17\nq: quit")(m.width, m.height)
	} else {
		content = m.root.Render(m.width, m.height)
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
func main() {
	m := initialModel()
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		m.width, m.height = 80, 20
		fmt.Print(m.View().Content)
		return
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
