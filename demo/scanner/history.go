package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// HistoryScreen is the "previous scans" view. Same split shape as the
// Scan screen (sidebar + main), different data. Reusing the shape is
// the point — switch screens, keep the user oriented.
type HistoryScreen struct {
	results []ScanResult
	cursor  int
	detail  ViewportState
}

func NewHistoryScreen() HistoryScreen {
	return HistoryScreen{results: SampleHistory()}
}

func (h HistoryScreen) Title() string    { return "History" }
func (h HistoryScreen) Subtitle() string { return "previous scans" }
func (h HistoryScreen) Footer() string {
	return "j/k browse runs  ·  pgup/pgdn scroll details"
}

func (h HistoryScreen) Update(msg tea.KeyPressMsg, keys KeyMap) HistoryScreen {
	key := msg.String()
	switch {
	case Matches(key, keys.Up), Matches(key, keys.Prev):
		if h.cursor > 0 {
			h.cursor--
			h.detail = ViewportState{}
		}
	case Matches(key, keys.Down), Matches(key, keys.Next):
		if h.cursor < len(h.results)-1 {
			h.cursor++
			h.detail = ViewportState{}
		}
	case Matches(key, keys.Top):
		h.cursor = 0
		h.detail = ViewportState{}
	case Matches(key, keys.Bottom):
		h.cursor = len(h.results) - 1
		h.detail = ViewportState{}
	case Matches(key, keys.PageUp):
		h.detail.ScrollBy(-10, len(h.detailLines()), 1)
	case Matches(key, keys.PageDown):
		h.detail.ScrollBy(10, len(h.detailLines()), 1)
	}
	return h
}

func (h HistoryScreen) Body(w, height int) *flexgo.Node {
	sidebar := &flexgo.Node{
		Width:      28,
		Background: colPanel,
		View:       h.sidebarView(),
	}

	lines := h.detailLines()
	detailState := h.detail
	main := &flexgo.Node{
		Flex:       1,
		Background: colBg,
		View:       viewportView(&detailState, lines),
	}

	return &flexgo.Node{
		Flex:       1,
		Dir:        flexgo.Row,
		Gap:        2,
		Background: colBg,
		Children:   []*flexgo.Node{sidebar, main},
	}
}

func (h HistoryScreen) sidebarView() func(int, int) string {
	return func(w, hgt int) string {
		rows := []string{
			headingStyle().Render("RUNS"),
			"",
		}
		for i, r := range h.results {
			line := fmt.Sprintf("  %s  (%d)", r.ID, r.Findings)
			bg := colPanel
			fg := colText
			if i == h.cursor {
				bg = colSelected
				fg = colHighlight
			}
			rows = append(rows,
				lipgloss.NewStyle().
					Background(bg).Foreground(fg).
					Width(w-4).
					Render(line))
		}
		return lipgloss.NewStyle().
			Width(w).Height(hgt).
			Background(colPanel).
			Padding(1, 2).
			Render(strings.Join(rows, "\n"))
	}
}

func (h HistoryScreen) detailLines() []string {
	if len(h.results) == 0 {
		return []string{"(no history)"}
	}
	idx := h.cursor
	if idx >= len(h.results) {
		idx = len(h.results) - 1
	}
	r := h.results[idx]
	out := []string{
		"Run:        " + r.ID,
		"Started:    " + r.Started.Format("2006-01-02 15:04"),
		"Duration:   " + r.Duration.String(),
		"Profile:    " + r.Profile,
		fmt.Sprintf("Findings:   %d", r.Findings),
		"",
		lipgloss.NewStyle().Background(colBg).Foreground(colSubtle).Bold(true).Render("summary"),
		"",
	}
	for i := 0; i < r.Findings; i++ {
		out = append(out, fmt.Sprintf("  [%s] finding %02d", sev(i).Label(), i+1))
	}
	return out
}

func sev(i int) Severity {
	switch i % 3 {
	case 0:
		return SevHigh
	case 1:
		return SevMed
	}
	return SevLow
}
