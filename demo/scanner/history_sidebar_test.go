package main

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func assertHistoryBounds(t *testing.T, out string, w, h int) {
	t.Helper()
	if w <= 0 || h <= 0 {
		if out != "" {
			t.Errorf("%dx%d: expected empty output", w, h)
		}
		return
	}
	if lipgloss.Width(out) != w || lipgloss.Height(out) != h {
		t.Errorf("%dx%d: painted %dx%d", w, h, lipgloss.Width(out), lipgloss.Height(out))
	}
}

func TestHistorySidebarSelectedWindow(t *testing.T) {
	s := NewHistoryScreen()
	for _, cursor := range []int{0, 2, 4} {
		s.cursor = cursor
		for height := 1; height <= 12; height++ {
			out := s.sidebarView()(28, height)
			assertHistoryBounds(t, out, 28, height)
			if !strings.Contains(ansi.Strip(out), s.results[cursor].ID) {
				t.Errorf("cursor %d height %d: selected run missing", cursor, height)
			}
			selected := lipgloss.NewStyle().Background(colSelected).Foreground(colHighlight).Render("X")
			prefix, _, _ := strings.Cut(selected, "X")
			if prefix == "" || !strings.Contains(out, prefix) {
				t.Errorf("cursor %d height %d: selection colors missing", cursor, height)
			}
		}
	}
}

func TestHistorySidebarBoundsAndText(t *testing.T) {
	for _, id := range []string{"short", strings.Repeat("long", 30), strings.Repeat("界👩‍💻", 20), "\x1b[31m" + strings.Repeat("styled", 20) + "\x1b[0m", "first\nsecond", "tabs\twide界"} {
		s := HistoryScreen{results: []ScanResult{{ID: id, Findings: 9}}}
		for _, w := range []int{-1, 0, 1, 2, 3, 4, 5, 6, 10, 28} {
			for _, h := range []int{-1, 0, 1, 2, 3, 4, 5, 9, 12} {
				assertHistoryBounds(t, s.sidebarView()(w, h), w, h)
			}
		}
	}
}

func TestHistorySidebarNarrowWidePrefix(t *testing.T) {
	for _, tc := range []struct{ id, prefix string }{{"界界", "界"}, {"👩‍💻name", "👩‍💻"}} {
		s := HistoryScreen{results: []ScanResult{{ID: tc.id}}}
		out := s.sidebarView()(7, 1)
		assertHistoryBounds(t, out, 7, 1)
		if !strings.HasPrefix(ansi.Strip(out), tc.prefix) {
			t.Errorf("%q: leading grapheme missing: %q", tc.id, ansi.Strip(out))
		}
		selected := lipgloss.NewStyle().Background(colSelected).Foreground(colHighlight).Render("X")
		prefix, _, _ := strings.Cut(selected, "X")
		if prefix == "" || !strings.Contains(out, prefix) {
			t.Errorf("%q: selection colors missing", tc.id)
		}
	}
}

func TestHistorySidebarNavigationAndProjection(t *testing.T) {
	s := NewHistoryScreen()
	keys := DefaultKeys()
	for _, step := range []struct {
		code   rune
		cursor int
	}{{tea.KeyEnd, 4}, {tea.KeyUp, 3}, {tea.KeyDown, 4}, {tea.KeyHome, 0}} {
		s.detail.Offset = 7
		s = s.Update(tea.KeyPressMsg{Code: step.code}, keys)
		if s.cursor != step.cursor || s.detail.Offset != 0 {
			t.Fatal("navigation changed its selection/reset contract")
		}
		for _, height := range []int{1, 2, 5, 9} {
			out := s.sidebarView()(28, height)
			if !strings.Contains(ansi.Strip(out), s.results[s.cursor].ID) {
				t.Errorf("key %d height %d hid selection", step.code, height)
			}
		}
	}
	s.cursor = 4
	s.detail.Offset = 7
	before := s
	before.results = append([]ScanResult(nil), s.results...)
	for _, height := range []int{9, 1, 5, 12, 1, 9} {
		for _, width := range []int{28, 1, 28} {
			assertHistoryBounds(t, s.sidebarView()(width, height), width, height)
		}
		_ = s.Body(80, height).Render(80, height)
	}
	if !reflect.DeepEqual(s, before) {
		t.Fatal("render or resize mutated model")
	}
}

func TestHistorySidebarEmptyAndShortened(t *testing.T) {
	for _, cursor := range []int{-5, 0, 99} {
		for _, results := range [][]ScanResult{nil, {{ID: "only"}}} {
			s := HistoryScreen{results: results, cursor: cursor, detail: ViewportState{Offset: 7}}
			for _, height := range []int{1, 3, 5, 12} {
				out := s.sidebarView()(28, height)
				assertHistoryBounds(t, out, 28, height)
				want := "(no history)"
				if len(results) > 0 {
					want = "only"
				}
				if !strings.Contains(ansi.Strip(out), want) {
					t.Errorf("cursor %d height %d: missing %q", cursor, height, want)
				}
			}
			if s.cursor != cursor || s.detail.Offset != 7 {
				t.Fatal("stale projection mutated state")
			}
			next := s.Update(tea.KeyPressMsg{Code: 'x'}, DefaultKeys())
			if next.cursor != 0 {
				t.Fatal("Update failed to reconcile cursor")
			}
			if (cursor != 0 || len(results) == 0) && next.detail.Offset != 0 {
				t.Fatal("Update failed to reset stale details")
			}
			if cursor == 0 && len(results) > 0 && next.detail.Offset != 7 {
				t.Fatal("Update reset valid details")
			}
		}
	}
}

func TestHistorySidebarFittingBytes(t *testing.T) {
	s := NewHistoryScreen()
	s.cursor = 4
	const want = "\u001b[48;2;31;35;53m                            \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[1;38;2;115;122;162;48;2;31;35;53mRUNS\u001b[m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m                    \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m                        \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[38;2;169;177;214;48;2;31;35;53m  2026-04-20-a  (12)\u001b[m\u001b[48;2;31;35;53m    \u001b[m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[38;2;169;177;214;48;2;31;35;53m  2026-04-19-b  (9)\u001b[m\u001b[48;2;31;35;53m     \u001b[m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[38;2;169;177;214;48;2;31;35;53m  2026-04-19-a  (21)\u001b[m\u001b[48;2;31;35;53m    \u001b[m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[38;2;169;177;214;48;2;31;35;53m  2026-04-18-a  (4)\u001b[m\u001b[48;2;31;35;53m     \u001b[m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\n\u001b[48;2;31;35;53m  \u001b[m\u001b[48;2;31;35;53m\u001b[38;2;201;164;90;48;2;57;75;112m  2026-04-17-a  (11)\u001b[m\u001b[48;2;57;75;112m    \u001b[m\u001b[m\u001b[48;2;31;35;53m  \u001b[m\n\u001b[48;2;31;35;53m                            \u001b[m\n\u001b[48;2;31;35;53m                            \u001b[m\n\u001b[48;2;31;35;53m                            \u001b[m\n\u001b[48;2;31;35;53m                            \u001b[m"
	if got := s.sidebarView()(28, 12); got != want {
		t.Fatal("normal fitting sidebar bytes changed")
	}
}

func TestHistorySidebarPageFits(t *testing.T) {
	s := NewHistoryScreen()
	s.cursor = 4
	out := pageShell(s.Subtitle(), ScreenHistory, s.Body(80, 15), s.Footer()).Render(80, 15)
	assertHistoryBounds(t, out, 80, 15)
	if !strings.Contains(ansi.Strip(out), s.results[s.cursor].ID) {
		t.Fatal("page hid selected run")
	}
}
