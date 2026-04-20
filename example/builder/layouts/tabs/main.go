// Builder version of layouts.Tabs — tab strip + active panel with an
// underline centred via AlignSelf, all built with NodeBuilder.
package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

type model struct {
	width, height int
	ready         bool
	active        int
}

func initialModel() model { return model{active: 1} }

type tab struct {
	title string
	panel func(w, h int) string
}

func buildRoot(active int) *flexgo.Node {
	tabs := []tab{
		{"Logs", panel("Live logs stream", lipgloss.Color("238"))},
		{"Metrics", panel("CPU 14%  RAM 38%", lipgloss.Color("59"))},
		{"Alerts", panel("No active alerts", lipgloss.Color("95"))},
	}

	strip := make([]*flexgo.Node, 0, len(tabs))
	for i, t := range tabs {
		underlineText := " "
		underlineWidth := 1
		if i == active {
			w := max(1, lipgloss.Width(t.title))
			underlineText = strings.Repeat("-", w)
			underlineWidth = w
		}

		strip = append(strip, flexgo.NewNode().
			Dir(flexgo.Col).
			Flex(1).
			Children(
				flexgo.NewNode().Height(1).View(centered(t.title)).Build(),
				flexgo.NewNode().
					Height(1).
					Width(underlineWidth).
					AlignSelf(flexgo.AlignCenter).
					View(centered(underlineText)).
					Build(),
			).
			Build())
	}

	return flexgo.NewNode().
		Dir(flexgo.Col).
		Background(lipgloss.Color("236")).
		Paddings(flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}).
		Children(
			flexgo.NewNode().
				Height(2).
				Dir(flexgo.Row).
				Gap(1).
				Children(strip...).
				Build(),
			flexgo.NewNode().
				Flex(1).
				View(tabs[active].panel).
				Build(),
		).
		Build()
}

func centered(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(text)
	}
}

func panel(text string, bg color.Color) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Background(bg).
			Foreground(lipgloss.Color("230")).
			Render(text)
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ready = true
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "left":
			m.active = max(0, m.active-1)
		case "l", "right":
			m.active = min(2, m.active+1)
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading...")
	}
	v := tea.NewView(buildRoot(m.active).Render(m.width, m.height))
	v.AltScreen = true
	return v
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(buildRoot(1).Render(80, 24))
		return
	}
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
