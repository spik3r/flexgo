package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
	"github.com/spik3r/flexgo/layouts"
)

const sample = `package main

import (
	"fmt"
	"os"
)

// The body below intentionally runs longer than the viewport so
// j / k scrolling has something to do. Press SPACE to open the
// example modal, ESC to dismiss it.

func main() {
	greet("world")
	greet("flexgo")
	for i := range 20 {
		fmt.Printf("line %d: the quick brown fox\n", i)
	}
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func greet(name string) {
	fmt.Println("hello,", name)
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no args")
	}
	for _, a := range args {
		fmt.Println(a)
	}
	return nil
}`

// Palette — one place to tweak the whole demo.
var (
	pageBg      = lipgloss.Color("236")
	panelBg     = lipgloss.Color("234")
	panelFg     = lipgloss.Color("252")
	accentBg    = lipgloss.Color("61")
	accentFg    = lipgloss.Color("230")
	mutedFg     = lipgloss.Color("244")
	buttonBg    = lipgloss.Color("205")
	borderFg    = lipgloss.Color("240")
	thumbFg     = lipgloss.Color("205")
	trackFg     = lipgloss.Color("240")
	modalBg     = lipgloss.Color("237")
	modalBorder = lipgloss.Color("205")
	modalTextFg = lipgloss.Color("252")
)

type model struct {
	width, height int
	ready         bool
	scroll        int
	lines         []string
	modalOpen     bool
}

func initialModel() model {
	return model{lines: strings.Split(sample, "\n")}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	case tea.KeyPressMsg:
		key := msg.String()
		if m.modalOpen {
			if key == "esc" || key == "q" {
				m.modalOpen = false
			}
			return m, nil
		}
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.scroll = clamp(m.scroll+1, 0, m.maxScroll())
		case "k", "up":
			m.scroll = clamp(m.scroll-1, 0, m.maxScroll())
		case "g", "home":
			m.scroll = 0
		case "G", "end":
			m.scroll = m.maxScroll()
		case "space", " ", "enter":
			m.modalOpen = true
		}
	}
	return m, nil
}

// panelSize returns the (width, height) of the bordered code panel
// within the body allocation. Half the body height and ~two-thirds the
// width, so there's always spare space around it and enough content to
// actually exercise the scrollbar.
func (m model) panelSize() (int, int) {
	bodyH := m.height - 1 - 3 - 2 // header + footer + outer padding
	bodyW := m.width - 4          // outer left/right padding
	panelH := max(6, bodyH/2)
	panelW := max(30, bodyW*2/3)
	return panelW, panelH
}

// viewportHeight is the code-panel allocation minus its border (2 rows).
func (m model) viewportHeight() int {
	_, panelH := m.panelSize()
	return max(1, panelH-2)
}

func (m model) maxScroll() int {
	maxScroll := len(m.lines) - m.viewportHeight()
	if maxScroll < 0 {
		return 0
	}
	return maxScroll
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading...")
	}
	v := tea.NewView(renderPage(m))
	v.AltScreen = true
	return v
}

func renderPage(m model) string {
	page := layouts.HeaderBodyFooter(
		1, m.headerView,
		m.bodyView,
		3, m.footerView,
	)
	page.Background = pageBg
	page.Paddings = flexgo.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}
	page.Gap = 1

	out := page.Render(m.width, m.height)

	if m.modalOpen {
		modal := layouts.Modal("ExampleModal", modalBodyView, 40, 7, modalBg)
		modal.BorderForeground = modalBorder
		out = modal.Render(m.width, m.height)
	}
	return out
}

func (m model) headerView(w, h int) string {
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Background(accentBg).
		Foreground(accentFg).
		Bold(true).
		Align(lipgloss.Center, lipgloss.Center).
		Render("flexgo · modal demo")
}

func (m model) bodyView(w, h int) string {
	panelW, panelH := m.panelSize()
	if panelW > w {
		panelW = w
	}
	if panelH > h {
		panelH = h
	}

	panel := m.renderCodePanel(panelW, panelH)

	return lipgloss.Place(
		w, h,
		lipgloss.Center, lipgloss.Center,
		panel,
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(pageBg)),
	)
}

func (m model) renderCodePanel(w, h int) string {
	innerW := w - 2
	innerH := h - 2
	if innerW < 4 || innerH < 1 {
		return lipgloss.NewStyle().Width(w).Height(h).Background(panelBg).Render("")
	}

	scrollbarW := 1
	codeW := innerW - scrollbarW - 1 // -1 for a 1-col spacer before the bar

	total := len(m.lines)
	start := clamp(m.scroll, 0, max(0, total))
	end := min(start+innerH, total)
	visible := m.lines[start:end]
	for len(visible) < innerH {
		visible = append(visible, "")
	}

	lineStyle := lipgloss.NewStyle().
		Background(panelBg).
		Foreground(panelFg).
		Width(codeW).
		PaddingLeft(1)

	codeRows := make([]string, innerH)
	for i, line := range visible {
		codeRows[i] = lineStyle.Render(line)
	}
	code := lipgloss.JoinVertical(lipgloss.Left, codeRows...)

	spacerCol := lipgloss.NewStyle().Background(panelBg).Width(1).Render("")
	spacer := lipgloss.JoinVertical(lipgloss.Left, repeat(spacerCol, innerH)...)

	thumbStart, thumbLen := thumbRange(innerH, total, m.scroll)
	sbRows := make([]string, innerH)
	thumb := lipgloss.NewStyle().Background(panelBg).Foreground(thumbFg).Render("█")
	track := lipgloss.NewStyle().Background(panelBg).Foreground(trackFg).Render("░")
	for i := range innerH {
		if i >= thumbStart && i < thumbStart+thumbLen {
			sbRows[i] = thumb
		} else {
			sbRows[i] = track
		}
	}
	scrollbar := lipgloss.JoinVertical(lipgloss.Left, sbRows...)

	inner := lipgloss.JoinHorizontal(lipgloss.Top, code, spacer, scrollbar)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderFg).
		Background(panelBg).
		Render(inner)
}

func (m model) footerView(w, h int) string {
	button := lipgloss.NewStyle().
		Background(buttonBg).
		Foreground(accentFg).
		Padding(0, 3).
		Bold(true).
		Render("[ Open Dialog ]")

	hint := lipgloss.NewStyle().
		Background(pageBg).
		Foreground(mutedFg).
		Render("space · j/k scroll · g/G top/end · q quit")

	pos := lipgloss.NewStyle().
		Background(pageBg).
		Foreground(mutedFg).
		Render(fmt.Sprintf("line %d/%d", m.scroll+1, len(m.lines)))

	gap := lipgloss.NewStyle().Background(pageBg).Render("  ")
	row := lipgloss.JoinHorizontal(lipgloss.Center, button, gap, hint, gap, pos)

	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Background(pageBg).
		Align(lipgloss.Center, lipgloss.Center).
		Render(row)
}

func modalBodyView(w, h int) string {
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Background(modalBg).
		Foreground(modalTextFg).
		Align(lipgloss.Center, lipgloss.Center).
		Render("Press ESC to close")
}

// thumbRange returns the scrollbar thumb's start offset and length in
// rows, sized proportionally to the fraction of content visible.
func thumbRange(viewH, total, scroll int) (start, length int) {
	if total <= viewH {
		return 0, viewH
	}
	length = max(1, viewH*viewH/total)
	maxScroll := total - viewH
	if maxScroll <= 0 {
		return 0, length
	}
	start = scroll * (viewH - length) / maxScroll
	return start, length
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func repeat(s string, n int) []string {
	out := make([]string, n)
	for i := range n {
		out[i] = s
	}
	return out
}

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		m := initialModel()
		m.width = 80
		m.height = 24
		m.ready = true
		fmt.Print(renderPage(m))
		return
	}
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
