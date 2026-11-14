package main

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

func resize(m model, w, h int) model {
	updated, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return updated.(model)
}
func press(m model, code rune) model {
	updated, _ := m.Update(tea.KeyPressMsg{Code: code})
	return updated.(model)
}

func TestBreakpoints(t *testing.T) {
	for _, tc := range []struct {
		w, h         int
		wide, usable bool
	}{
		{23, 12, false, false}, {24, 5, false, false}, {24, 6, false, true},
		{71, 12, false, true}, {72, 12, true, true}, {120, 30, true, true},
	} {
		m := resize(initialModel(), tc.w, tc.h)
		if m.usable() != tc.usable {
			t.Fatalf("%dx%d usable", tc.w, tc.h)
		}
		out := m.View().Content
		if !tc.usable {
			if !strings.Contains(out, "Resize") {
				t.Fatal(out)
			}
			continue
		}
		if err := flexgo.Validate(m.tree()); err != nil {
			t.Fatal(err)
		}
		if got := strings.Contains(out, "Details"); got != tc.wide {
			t.Fatalf("%dx%d details=%v", tc.w, tc.h, got)
		}
	}
}

func TestRenderedBounds(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {-1, 10}, {10, -1}, {1, 1}, {23, 5}, {24, 6}, {71, 7}, {72, 6}, {80, 24}, {120, 40}} {
		for _, detail := range []bool{false, true} {
			m := resize(initialModel(), size[0], size[1])
			m.detail = detail
			m.selected = 11
			m.offset = 29
			before := m
			out := m.View().Content
			if !reflect.DeepEqual(before, m) {
				t.Fatal("View mutated model")
			}
			if m.width == 0 || m.height == 0 {
				if out != "" {
					t.Fatalf("zero area: %q", out)
				}
				continue
			}
			if lipgloss.Width(out) > m.width || lipgloss.Height(out) > m.height {
				t.Fatalf("%v detail=%v rendered %dx%d", size, detail, lipgloss.Width(out), lipgloss.Height(out))
			}
		}
	}
}

func TestResizePreservesNavigation(t *testing.T) {
	m := resize(initialModel(), 80, 12)
	for range 9 {
		m = press(m, tea.KeyDown)
	}
	m = press(m, tea.KeyEnter)
	for range 15 {
		m = press(m, tea.KeyDown)
	}
	if m.selected != 9 || m.offset != 15 || !m.detail {
		t.Fatalf("setup %+v", m)
	}
	for _, size := range [][2]int{{40, 9}, {12, 3}, {120, 50}, {72, 6}, {80, 12}} {
		m = resize(m, size[0], size[1])
		_ = m.View()
		if m.selected != 9 || m.offset != 15 || !m.detail {
			t.Fatalf("resize %v lost navigation", size)
		}
	}
	m = resize(m, 40, 9)
	m = press(m, tea.KeyEscape)
	if m.detail || !strings.Contains(m.View().Content, "> Record 10") {
		t.Fatal("back lost selection")
	}
	m = press(m, tea.KeyEnter)
	if m.offset != 15 {
		t.Fatal("reopening detail lost scroll")
	}
	m = press(m, tea.KeyEscape)
	m = press(m, tea.KeyDown)
	if m.selected != 10 || m.offset != 0 {
		t.Fatal("new selection did not reset scroll")
	}
}

func TestScrollClampsToVisibleProjection(t *testing.T) {
	m := resize(initialModel(), 40, 6)
	m.detail = true
	for range 100 {
		m = press(m, tea.KeyDown)
	}
	if m.offset != 31 {
		t.Fatalf("bottom=%d", m.offset)
	}
	m = resize(m, 40, 24)
	// A key operates from the visible position after growing the viewport.
	m = press(m, tea.KeyUp)
	if m.offset != 12 {
		t.Fatalf("up after resize=%d", m.offset)
	}
	for range 100 {
		m = press(m, tea.KeyUp)
	}
	if m.offset != 0 {
		t.Fatal("negative scroll")
	}
	m = press(m, tea.KeyEscape)
	for range 100 {
		m = press(m, tea.KeyDown)
	}
	if m.selected != 11 {
		t.Fatal("selection past end")
	}
	for range 100 {
		m = press(m, tea.KeyUp)
	}
	if m.selected != 0 {
		t.Fatal("selection before start")
	}
}

func TestEmptyAndSmallNavigation(t *testing.T) {
	m := resize(model{}, 80, 12)
	for _, k := range []rune{tea.KeyEnter, tea.KeyTab, tea.KeyDown, tea.KeyUp, tea.KeyEscape} {
		m = press(m, k)
	}
	if m.detail || m.selected != 0 || m.offset != 0 {
		t.Fatal("empty model changed")
	}
	if !strings.Contains(m.View().Content, "No records") {
		t.Fatal("missing empty state")
	}
	m = resize(initialModel(), 10, 3)
	for _, k := range []rune{tea.KeyEnter, tea.KeyTab, tea.KeyDown, tea.KeyUp} {
		m = press(m, k)
	}
	if m.detail || m.selected != 0 {
		t.Fatal("tiny navigation changed state")
	}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
	if cmd == nil {
		t.Fatal("quit unavailable in tiny terminal")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("not quit")
	}
	m = resize(initialModel(), 40, 9)
	m = press(m, tea.KeyTab)
	m = press(m, tea.KeyTab)
	if m.detail {
		t.Fatal("repeated toggle did not return")
	}
}

func TestLeafClippingUnicode(t *testing.T) {
	out := fit([]string{"界界界界界", "long line", "third"}, 5, 2)
	if lipgloss.Width(out) > 5 || lipgloss.Height(out) > 2 {
		t.Fatalf("unbounded %q", out)
	}
}

func TestMinimumWidthHelp(t *testing.T) {
	m := resize(initialModel(), minWidth, minHeight)
	for _, detail := range []bool{false, true} {
		m.detail = detail
		lines := strings.Split(m.View().Content, "\n")
		help := strings.Join(lines[len(lines)-2:], "\n")
		if !strings.Contains(help, "q: quit") || !strings.Contains(help, "Tab") {
			t.Fatalf("incomplete help: %q", help)
		}
		if detail && !strings.Contains(help, "Esc: list") {
			t.Fatalf("missing back: %q", help)
		}
	}
}

func TestLeafClippingCountsEmbeddedLines(t *testing.T) {
	got := fit([]string{"first\nsecond", "third"}, 5, 1)
	if got != "first" {
		t.Fatalf("embedded newline escaped height limit: %q", got)
	}
}
