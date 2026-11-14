package main

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestEmptyFilesEndThenOpen(t *testing.T) {
	s := NewScanScreen(DefaultProfile())
	s.activeTab, s.files = 3, nil
	s = s.Update(tea.KeyPressMsg{Code: tea.KeyEnd}, DefaultKeys())
	s = s.Update(tea.KeyPressMsg{Code: tea.KeyEnter}, DefaultKeys())
	if s.filesCursor != 0 || s.openedPath != "" {
		t.Fatal("empty navigation changed selection")
	}
}

func TestEmptyHistoryEndThenRecovery(t *testing.T) {
	h := HistoryScreen{}
	h = h.Update(tea.KeyPressMsg{Code: tea.KeyEnd}, DefaultKeys())
	h.results = []ScanResult{{ID: "synthetic-run"}}
	_ = h.detailLines()
	if h.cursor != 0 {
		t.Fatal("empty history left a negative selection")
	}
}

func TestFilesRepeatedEmptyAndRecovery(t *testing.T) {
	keys := DefaultKeys()
	s := NewScanScreen(DefaultProfile())
	s.activeTab = 3
	for range 4 {
		s.files = []string{"  alpha.go", "  beta.go", "  gamma.go"}
		s = s.Update(tea.KeyPressMsg{Code: tea.KeyEnd}, keys)
		s.files = s.files[:1]
		s = s.Update(tea.KeyPressMsg{Code: tea.KeyEnter}, keys)
		if s.filesCursor != 0 || s.openedPath != "alpha.go" || len(s.openedBody) == 0 {
			t.Fatal("shortened list did not open valid item")
		}
		s.files = nil
		before := append([]string(nil), s.openedBody...)
		s = s.Update(tea.KeyPressMsg{Code: tea.KeyDown}, keys)
		if s.openedPath != "alpha.go" || !reflect.DeepEqual(before, s.openedBody) {
			t.Fatal("list removal changed open file")
		}
		s = s.Update(tea.KeyPressMsg{Code: tea.KeyEscape}, keys)
		for _, code := range []rune{tea.KeyEnd, tea.KeyEnter, tea.KeyUp, tea.KeyDown, tea.KeyHome, tea.KeyPgDown, tea.KeyPgUp} {
			s = s.Update(tea.KeyPressMsg{Code: code}, keys)
			if s.filesCursor != 0 || s.openedPath != "" || s.viewports[3].Offset != 0 {
				t.Fatal("empty file navigation changed state")
			}
		}
		if !strings.Contains(ansi.Strip(s.activePanelView()(20, 3)), "(no files)") {
			t.Fatal("missing empty feedback")
		}
		s.files = []string{"recovered.go"}
		s = s.Update(tea.KeyPressMsg{Code: tea.KeyEnter}, keys)
		if s.openedPath != "recovered.go" {
			t.Fatal("repopulated list did not recover")
		}
		s = s.Update(tea.KeyPressMsg{Code: tea.KeyEscape}, keys)
	}
}

func TestHistoryRepeatedEmptyAndRecovery(t *testing.T) {
	keys := DefaultKeys()
	h := HistoryScreen{}
	for range 4 {
		h.results = []ScanResult{{ID: "first", Findings: 15}, {ID: "second", Findings: 20}, {ID: "third", Findings: 25}}
		h = h.Update(tea.KeyPressMsg{Code: tea.KeyEnd}, keys)
		h = h.Update(tea.KeyPressMsg{Code: tea.KeyPgDown}, keys)
		h.results = h.results[:1]
		h = h.Update(tea.KeyPressMsg{Code: 'x'}, keys)
		if h.cursor != 0 || h.detail.Offset != 0 || h.detailLines()[0] != "Run:        first" {
			t.Fatal("shortened history retained stale selection")
		}
		h.results = nil
		for _, code := range []rune{tea.KeyEnd, tea.KeyHome, tea.KeyDown, tea.KeyUp, tea.KeyPgDown, tea.KeyPgUp, tea.KeyTab} {
			h = h.Update(tea.KeyPressMsg{Code: code}, keys)
			if h.cursor != 0 || h.detail.Offset != 0 || h.detailLines()[0] != "(no history)" {
				t.Fatal("empty history state changed")
			}
		}
		h.results = []ScanResult{{ID: "recovered"}}
		if h.detailLines()[0] != "Run:        recovered" {
			t.Fatal("history failed to recover")
		}
	}
}

func TestEmptyRecoveryRenderProjectionAndBindings(t *testing.T) {
	for _, cursor := range []int{-5, -1, 0, 4, 99} {
		h := HistoryScreen{results: []ScanResult{{ID: "only"}}, cursor: cursor, detail: ViewportState{Offset: 4}}
		_ = h.sidebarView()(28, 8)
		if h.detailLines()[0] != "Run:        only" || h.cursor != cursor || h.detail.Offset != 4 {
			t.Fatal("history projection changed stored state")
		}
		out := ansi.Strip(filetreeView([]string{"only.go"}, cursor)(12, 2))
		if !strings.Contains(out, "only.go") {
			t.Fatal("file projection hid recovered item")
		}
	}
	for _, size := range [][2]int{{1, 1}, {5, 2}, {20, 3}, {0, 3}, {20, 0}, {-1, 2}, {5, -1}} {
		out := filetreeView(nil, -1)(size[0], size[1])
		if size[0] <= 0 || size[1] <= 0 {
			if out != "" {
				t.Fatal("nonpositive empty view rendered")
			}
		} else if lipgloss.Width(out) > size[0] || lipgloss.Height(out) > size[1] {
			t.Fatal("empty hint exceeded bounds")
		}
	}
	keys := DefaultKeys()
	keys.Bottom, keys.Open = []string{"b"}, []string{"o"}
	s := NewScanScreen(DefaultProfile())
	s.activeTab, s.files = 3, nil
	s = s.Update(tea.KeyPressMsg{Code: 'b', Text: "b"}, keys)
	s = s.Update(tea.KeyPressMsg{Code: 'o', Text: "o"}, keys)
	s.files = []string{"mapped.go"}
	s = s.Update(tea.KeyPressMsg{Code: 'o', Text: "o"}, keys)
	if s.openedPath != "mapped.go" {
		t.Fatal("recovery bypassed keymap")
	}
}
