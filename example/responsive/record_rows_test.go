package main

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestMultilineTitleKeepsSelectionVisible(t *testing.T) {
	for _, title := range []string{"first\nextra row", "first\r\nextra row", "first\rextra row", "\nextra row"} {
		m := resize(initialModel(), 24, 6)
		m.records[0].title = title
		m.selected = 1
		got := m.listView(24, 2)
		if !strings.Contains(got, "> Record 02") || strings.Contains(got, "extra row") {
			t.Fatalf("selected record hidden by previous title %q: %q", title, got)
		}
	}
}

func TestRecordRowsCloseStylesAndKeepCellBounds(t *testing.T) {
	link := ansi.SetHyperlink("https://example.test")
	for _, title := range []string{
		"plain\nextra", "A\tB\r\nextra", "界界界界界\nextra", "e\u0301e\u0301e\u0301\nextra",
		"\x1b[31m" + link + "styled\nextra" + ansi.ResetStyle + ansi.ResetHyperlink(),
	} {
		for _, width := range []int{2, 3, 7, 24} {
			m := initialModel()
			m.records = []record{{title: title}, {title: "next"}}
			m.selected = 1
			before := m.records[0].title
			rows := strings.Split(m.listView(width, 2), "\n")
			if len(rows) != 2 || !strings.HasPrefix(rows[1], "> ") {
				t.Fatalf("missing selection at width %d: %q", width, rows)
			}
			for _, row := range rows {
				if lipgloss.Width(row) > width || strings.ContainsAny(row, "\r\t") {
					t.Fatalf("unbounded record at width %d: %q", width, row)
				}
			}
			if strings.Contains(title, link) {
				if !strings.HasSuffix(rows[0], ansi.ResetStyle+ansi.ResetHyperlink()) {
					t.Fatalf("title style or hyperlink escaped its row: %q", rows[0])
				}
				if strings.Contains(rows[1], "\x1b") {
					t.Fatalf("previous title style reached next record: %q", rows[1])
				}
			}
			if m.records[0].title != before {
				t.Fatal("render changed stored title")
			}
		}
	}
	m := initialModel()
	for _, size := range [][2]int{{0, 2}, {2, 0}, {-1, 2}, {2, -1}} {
		if out := m.listView(size[0], size[1]); out != "" {
			t.Fatalf("nonpositive allocation %v produced %q", size, out)
		}
	}
}

func TestRecordRowsStayVisibleThroughNavigationAndResize(t *testing.T) {
	m := initialModel()
	for i := range m.records {
		m.records[i].title += "\nextra title row"
	}
	for _, size := range [][2]int{{24, 6}, {72, 6}, {40, 9}, {80, 12}} {
		m = resize(m, size[0], size[1])
		m.detail, m.selected, m.offset = false, 0, 0
		for i := range m.records {
			before := m
			out := m.View().Content
			want := "> " + strings.Split(m.records[i].title, "\n")[0]
			if !strings.Contains(out, want) || strings.Contains(out, "extra title row") {
				t.Fatalf("selection %d hidden at %v: %q", i, size, out)
			}
			if lipgloss.Width(out) > m.width || lipgloss.Height(out) > m.height || !strings.Contains(out, "q: quit") {
				t.Fatalf("output bounds/help changed at %v: %q", size, out)
			}
			if !reflect.DeepEqual(before, m) {
				t.Fatal("View changed navigation state")
			}
			m = press(m, tea.KeyDown)
		}
		m = press(m, tea.KeyEnter)
		m = press(m, tea.KeyDown)
		before := m
		m = resize(m, 10, 3)
		m = press(m, tea.KeyDown)
		m = resize(m, size[0], size[1])
		if m.selected != before.selected || m.offset != before.offset || !m.detail {
			t.Fatal("tiny recovery changed selection, detail or scroll")
		}
		if !strings.Contains(m.View().Content, "Synthetic local data") {
			t.Fatal("detail content changed")
		}
		m = press(m, tea.KeyEscape)
		if !strings.Contains(m.View().Content, "> Record 12") {
			t.Fatal("back lost selected record")
		}
	}
}
