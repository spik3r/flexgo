package main

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestPadLineBgMeasuresExpandedTabs(t *testing.T) {
	for _, tc := range []struct {
		line  string
		width int
		want  string
	}{
		{"A\tB", 3, "A  "}, {"A\tB", 6, "A    B"}, {"\tB", 3, "   "},
		{"A\t\tB", 7, "A      "}, {"界\tB", 4, "界  "}, {"e\u0301\tB", 3, "e\u0301  "},
		{"plain", 7, "plain  "},
	} {
		got := padLineBg(tc.line, tc.width, colBg, colText)
		if lipgloss.Width(got) != tc.width || lipgloss.Height(got) != 1 {
			t.Errorf("%q escaped %dx1: %dx%d %q", tc.line, tc.width, lipgloss.Width(got), lipgloss.Height(got), got)
		}
		if visible := ansi.Strip(got); visible != tc.want {
			t.Errorf("%q at width %d: got %q want %q", tc.line, tc.width, visible, tc.want)
		}
	}
}

func TestScannerRowsKeepTabContentInsideAllocation(t *testing.T) {
	lines := []string{"A\tB", "\t\tname", "界\tB", "e\u0301\tfile"}
	for _, width := range []int{1, 2, 3, 4, 6, 12} {
		for _, height := range []int{1, 2, 4} {
			for _, offset := range []int{0, 1, 3} {
				state := ViewportState{Offset: offset}
				viewport := viewportView(&state, lines)(width, height)
				tree := filetreeView(lines, offset)(width, height)
				for _, out := range []string{viewport, tree} {
					if lipgloss.Width(out) != width || lipgloss.Height(out) != height {
						t.Errorf("allocation %dx%d escaped: %dx%d %q", width, height, lipgloss.Width(out), lipgloss.Height(out), out)
					}
					if strings.Contains(out, "\t") {
						t.Errorf("unexpanded tab: %q", out)
					}
				}
				for _, row := range strings.Split(ansi.Strip(viewport), "\n") {
					if !strings.HasSuffix(row, "│") && !strings.HasSuffix(row, "█") {
						t.Errorf("scroll indicator displaced: %q", row)
					}
				}
			}
		}
	}
	for _, width := range []int{-1, 0} {
		if got := padLineBg("\tB", width, colBg, colText); got != "" {
			t.Errorf("empty allocation rendered %q", got)
		}
	}
}
