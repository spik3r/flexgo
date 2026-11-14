package main

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestPageNavigationMovesActivePane(t *testing.T) {
	for _, width := range []int{24, 40, 72, 120} {
		m := resize(initialModel(), width, 9)
		for _, step := range []struct {
			key  rune
			want int
		}{
			{tea.KeyPgDown, 5}, {tea.KeyPgDown, 10}, {tea.KeyPgDown, 11},
			{tea.KeyPgDown, 11}, {tea.KeyPgUp, 6}, {tea.KeyHome, 0},
			{tea.KeyPgUp, 0}, {tea.KeyEnd, 11},
		} {
			m = press(m, step.key)
			if m.selected != step.want || m.detail {
				t.Fatalf("width %d key %d: selection=%d want %d", width, step.key, m.selected, step.want)
			}
			out := m.View().Content
			if !strings.Contains(out, "> "+m.records[m.selected].title) {
				t.Fatalf("selection not visible: %q", out)
			}
		}
		m = press(m, tea.KeyEnter)
		for _, step := range []struct {
			key  rune
			want int
		}{
			{tea.KeyPgDown, 5}, {tea.KeyEnd, 28}, {tea.KeyPgDown, 28},
			{tea.KeyPgUp, 23}, {tea.KeyHome, 0}, {tea.KeyPgUp, 0},
		} {
			m = press(m, step.key)
			if m.offset != step.want || m.selected != 11 || !m.detail {
				t.Fatalf("width %d key %d: offset=%d want %d", width, step.key, m.offset, step.want)
			}
		}
	}
}

func TestPageNavigationUsesVisibleProjectionAfterResize(t *testing.T) {
	m := resize(initialModel(), 40, 6)
	m = press(m, tea.KeyEnter)
	m = press(m, tea.KeyEnd)
	if m.offset != 31 {
		t.Fatalf("bottom=%d", m.offset)
	}
	m = resize(m, 80, 24)
	if m.offset != 31 || m.visibleOffset() != 13 {
		t.Fatal("resize lost stored intent")
	}
	m = press(m, tea.KeyPgUp)
	if m.offset != 0 {
		t.Fatalf("page up did not start at visible offset: %d", m.offset)
	}
	m = press(m, tea.KeyPgDown)
	if m.offset != 13 {
		t.Fatalf("page down did not clamp: %d", m.offset)
	}
	m = resize(m, 24, 6)
	m = press(m, tea.KeyPgDown)
	if m.offset != 15 {
		t.Fatalf("page size did not follow resize: %d", m.offset)
	}
	m = press(m, tea.KeyEscape)
	m = press(m, tea.KeyHome)
	if m.offset != 15 {
		t.Fatal("unchanged selection reset detail offset")
	}
	m = press(m, tea.KeyPgDown)
	if m.selected != 2 || m.offset != 0 {
		t.Fatal("changed selection failed to reset detail offset")
	}
	m = press(m, tea.KeyEnd)
	m = press(m, tea.KeyEnter)
	m = press(m, tea.KeyPgDown)
	m = press(m, tea.KeyEscape)
	m = press(m, tea.KeyEnd)
	if m.selected != 11 || m.offset != 2 {
		t.Fatal("boundary selection lost detail offset")
	}
}

func TestPageNavigationEmptyTinyAndShortContent(t *testing.T) {
	keys := []rune{tea.KeyPgUp, tea.KeyPgDown, tea.KeyHome, tea.KeyEnd}
	for _, detail := range []bool{false, true} {
		m := resize(model{detail: detail}, 80, 24)
		for _, key := range keys {
			before := m
			m = press(m, key)
			if !reflect.DeepEqual(before, m) {
				t.Fatal("empty navigation changed state")
			}
		}
		m = resize(initialModel(), 80, 24)
		m.detail, m.selected, m.offset = detail, 8, 10
		for _, size := range [][2]int{{0, 0}, {23, 12}, {80, 5}} {
			m = resize(m, size[0], size[1])
			before := m
			for _, key := range keys {
				m = press(m, key)
			}
			if !reflect.DeepEqual(before, m) {
				t.Fatal("hidden navigation changed state")
			}
		}
		m = resize(m, 80, 24)
		if m.selected != 8 || m.offset != 10 || m.detail != detail {
			t.Fatal("recovery lost navigation")
		}
	}
	m := resize(model{records: []record{{title: "Only", lines: []string{"short"}}}}, 24, 6)
	for _, detail := range []bool{false, true} {
		m.detail = detail
		for _, key := range keys {
			m = press(m, key)
			if m.selected != 0 || m.offset != 0 {
				t.Fatal("short content moved outside bounds")
			}
		}
	}
}

func TestPageNavigationHelpIsVisible(t *testing.T) {
	for _, size := range [][2]int{{24, 6}, {40, 9}, {72, 6}, {80, 12}} {
		for _, detail := range []bool{false, true} {
			m := resize(initialModel(), size[0], size[1])
			m.detail = detail
			out := m.View().Content
			rows := strings.Split(out, "\n")
			help := strings.Join(rows[len(rows)-2:], "\n")
			for _, hint := range []string{"PgUp/PgDn", "Home/End", "q: quit", "Tab"} {
				if !strings.Contains(help, hint) {
					t.Fatalf("missing %q at %v: %q", hint, size, help)
				}
			}
			if detail && !strings.Contains(help, "Esc: list") {
				t.Fatal("back control hidden")
			}
			if lipgloss.Width(out) > m.width || lipgloss.Height(out) > m.height {
				t.Fatal("help exceeded bounds")
			}
		}
	}
}
