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
	cursor := max(0, min(h.cursor, len(h.results)-1))
	if cursor != h.cursor || len(h.results) == 0 {
		h.cursor, h.detail = cursor, ViewportState{}
	}
	if len(h.results) == 0 {
		return h
	}
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
		if w <= 0 || hgt <= 0 {
			return ""
		}

		padY, padX := 0, 0
		rows := []string{}
		if hgt >= 5 {
			padY = 1
			rows = append(rows, headingStyle().Render("RUNS"), "")
		}
		if w >= 8 {
			padX = 2
		}
		contentW := w - 2*padX
		visible := hgt - 2*padY - len(rows)
		cursor := max(0, min(h.cursor, len(h.results)-1))
		offset := max(0, cursor-visible+1)
		for i := offset; i < min(offset+visible, len(h.results)); i++ {
			r := h.results[i]
			line := fmt.Sprintf("%s  (%d)", r.ID, r.Findings)
			if padX > 0 {
				line = "  " + line
			}
			bg, fg := colPanel, colText
			if i == cursor {
				bg, fg = colSelected, colHighlight
			}
			clipLine := flexgo.ClipView(func(int, int) string { return line })
			rows = append(rows, lipgloss.NewStyle().
				Background(bg).Foreground(fg).
				Width(contentW).
				Render(clipLine(contentW, 1)))
		}
		if len(h.results) == 0 {
			rows = append(rows, padLineBg("(no history)", contentW, colPanel, colText))
		}
		// Heading text can exceed the content area when the sidebar is narrow.
		clipped := flexgo.ClipView(func(int, int) string { return strings.Join(rows, "\n") })(contentW, hgt-2*padY)
		return lipgloss.NewStyle().
			Width(w).Height(hgt).
			Background(colPanel).
			Padding(padY, padX).
			Render(clipped)
	}
}

func (h HistoryScreen) detailLines() []string {
	if len(h.results) == 0 {
		return []string{"(no history)"}
	}
	idx := max(0, min(h.cursor, len(h.results)-1))
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
