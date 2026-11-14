package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/spik3r/flexgo"
)

type model struct {
	input    textinput.Model
	view     viewport.Model
	spinner  spinner.Model
	focusCmd tea.Cmd
	width    int
	height   int
	ready    bool
}

func initialModel() model {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Type here (q is text)"
	focusCmd := input.Focus()

	view := viewport.New()
	var lines []string
	for i := 1; i <= 40; i++ {
		lines = append(lines, fmt.Sprintf("%02d  Bubbles viewport inside a Flexgo leaf", i))
	}
	view.SetContent(strings.Join(lines, "\n"))

	m := model{input: input, view: view, spinner: spinner.New(), focusCmd: focusCmd, width: 80, height: 24, ready: true}
	m.resize(80, 24)
	return m
}

func (m *model) resize(w, h int) {
	m.width, m.height = max(0, w), max(0, h)
	setInputWidth(&m.input, max(0, m.width-3)) // Two prompt cells plus the virtual cursor cell.
	m.view.SetWidth(m.width)
	m.view.SetHeight(max(0, m.height-4)) // Three input rows and one status row.
	m.view.SetYOffset(m.view.YOffset())
}

func setInputWidth(input *textinput.Model, width int) {
	position := input.Position()
	input.SetWidth(width)
	// SetWidth leaves the old visible range in place. Rebuild it at the new width.
	input.CursorEnd()
	input.SetCursor(position)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.focusCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		m.ready = true
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab":
			if m.input.Focused() {
				m.input.Blur()
				return m, nil
			}
			return m, m.input.Focus()
		}
		if m.input.Focused() {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		} else {
			m.view, _ = m.view.Update(msg)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		// Cursor and mouse messages belong to their respective components.
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		if !m.input.Focused() {
			m.view, _ = m.view.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

func (m model) render() string {
	if !m.ready {
		return "Loading..."
	}
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.width < 20 || m.height < 5 {
		line := "Terminal too small"
		if len(line) > m.width {
			line = line[:m.width]
		}
		rows := make([]string, m.height)
		rows[0] = line + strings.Repeat(" ", m.width-len(line))
		for i := 1; i < len(rows); i++ {
			rows[i] = strings.Repeat(" ", m.width)
		}
		return strings.Join(rows, "\n")
	}
	root := &flexgo.Node{Dir: flexgo.Col, Children: []*flexgo.Node{
		{Height: 3, View: flexgo.ClipView(func(w, h int) string {
			input := m.input
			setInputWidth(&input, max(0, w-3))
			return "Input (Tab to switch)\n" + input.View()
		})},
		{Flex: 1, View: flexgo.ClipView(func(w, h int) string {
			view := m.view
			view.SetWidth(w)
			view.SetHeight(h)
			return view.View()
		})},
		{Height: 1, View: flexgo.ClipView(func(w, h int) string {
			focus := "viewport (↑/↓ to scroll)"
			if m.input.Focused() {
				focus = "input"
			}
			// Keep both controls before optional detail at the minimum width.
			return m.spinner.View() + " Esc:quit Tab:focus | " + focus
		})},
	}}
	return root.Render(m.width, m.height)
}

func (m model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(initialModel().render())
		return
	}
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
