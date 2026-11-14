package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/spik3r/flexgo"
)

const (
	minWidth     = 24
	minHeight    = 6
	wideWidth    = 72
	sidebarWidth = 24
)

type record struct {
	title string
	lines []string
}

type model struct {
	width, height int
	records       []record
	selected      int
	detail        bool
	offset        int
}

func initialModel() model {
	records := make([]record, 12)
	for i := range records {
		records[i].title = fmt.Sprintf("Record %02d", i+1)
		records[i].lines = []string{records[i].title, "Synthetic local data", ""}
		for j := 1; j <= 30; j++ {
			records[i].lines = append(records[i].lines, fmt.Sprintf("Detail line %02d: record %02d", j, i+1))
		}
	}
	return model{records: records}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) usable() bool { return m.width >= minWidth && m.height >= minHeight }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = max(0, msg.Width), max(0, msg.Height)
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if !m.usable() {
			return m, nil
		}
		switch key {
		case "enter":
			if len(m.records) > 0 {
				m.detail = true
			}
		case "esc":
			m.detail = false
		case "tab":
			if len(m.records) > 0 {
				m.detail = !m.detail
			}
		case "up", "k":
			if m.detail {
				m.offset = max(0, m.visibleOffset()-1)
			} else if m.selected > 0 {
				m.selected--
				m.offset = 0
			}
		case "down", "j":
			if m.detail && len(m.records) > 0 {
				m.offset = min(m.visibleOffset()+1, max(0, len(m.records[m.selected].lines)-(m.height-4)))
			} else if m.selected+1 < len(m.records) {
				m.selected++
				m.offset = 0
			}
		}
	}
	return m, nil
}

func (m model) visibleOffset() int {
	if len(m.records) == 0 {
		return 0
	}
	return min(m.offset, max(0, len(m.records[m.selected].lines)-max(0, m.height-4)))
}

// fit joins logical records before clipping so embedded newlines also count.
func fit(lines []string, w, h int) string {
	return flexgo.ClipView(func(int, int) string {
		return strings.Join(lines, "\n")
	})(w, h)
}

func text(s string) func(int, int) string {
	return func(w, h int) string { return fit([]string{s}, w, h) }
}

func (m model) listView(w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	if len(m.records) == 0 {
		return fit([]string{"No records"}, w, h)
	}
	// Keep selection visible without storing a viewport tied to an old height.
	start := max(0, m.selected-max(0, h-1))
	lines := make([]string, 0, min(h, len(m.records)))
	for i := start; i < len(m.records) && len(lines) < h; i++ {
		prefix := "  "
		if i == m.selected {
			prefix = "> "
		}
		lines = append(lines, fit([]string{prefix + m.records[i].title}, w, 1))
	}
	return strings.Join(lines, "\n")
}

func (m model) detailView(w, h int) string {
	if len(m.records) == 0 {
		return fit([]string{"No selection"}, w, h)
	}
	lines := m.records[m.selected].lines
	// Resizing changes the visible projection, not the user's stored scroll intent.
	start := min(m.offset, max(0, len(lines)-h))
	return fit(lines[start:], w, h)
}

func (m model) pane(name string, active bool, view func(int, int) string) *flexgo.Node {
	if active {
		name = "* " + name
	}
	return &flexgo.Node{Dir: flexgo.Col, Flex: 1, Children: []*flexgo.Node{
		{Height: 1, View: text(name)}, {Flex: 1, View: view},
	}}
}

func (m model) tree() *flexgo.Node {
	list := m.pane("Records", !m.detail, m.listView)
	detail := m.pane("Details", m.detail, m.detailView)
	body := list
	mode := "compact"
	help, keys := "Enter/Tab: details", "Up/Down: move | q: quit"
	if m.detail {
		help, keys = "Tab/Esc: list", "Up/Down scroll | q: quit"
	}
	if m.width >= wideWidth {
		mode = "wide"
		help, keys = "Tab/Enter: detail | Esc: list", "Up/Down: move/scroll | q: quit"
		list.Width = sidebarWidth
		body = &flexgo.Node{Dir: flexgo.Row, Flex: 1, Gap: 1, Children: []*flexgo.Node{list, detail}}
	} else if m.detail {
		body = detail
	}
	return &flexgo.Node{Dir: flexgo.Col, Children: []*flexgo.Node{
		{Height: 1, View: text("Responsive records | " + mode)},
		body,
		{Height: 1, View: text(help)},
		{Height: 1, View: text(keys)},
	}}
}

func (m model) View() tea.View {
	var content string
	if !m.usable() {
		content = fit([]string{"Resize to at least 24x6", "q: quit"}, m.width, m.height)
	} else {
		content = m.tree().Render(m.width, m.height)
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func main() {
	m := initialModel()
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		m.width, m.height = 80, 12
		fmt.Print(m.View().Content)
		return
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
