package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func key(code rune, text string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Text: text}
}

func TestInputFocusAndViewportScroll(t *testing.T) {
	m := initialModel()
	if !m.input.Focused() {
		t.Fatal("input should start focused")
	}
	updated, _ := m.Update(key('q', "q"))
	m = updated.(model)
	if got := m.input.Value(); got != "q" {
		t.Fatalf("q should be input text, got %q", got)
	}
	updated, _ = m.Update(key(tea.KeyTab, ""))
	m = updated.(model)
	if m.input.Focused() {
		t.Fatal("Tab should focus the viewport")
	}
	updated, _ = m.Update(key(tea.KeyDown, ""))
	m = updated.(model)
	if got := m.view.YOffset(); got != 1 {
		t.Fatalf("Down should scroll the viewport, got offset %d", got)
	}
	if got := m.input.Value(); got != "q" {
		t.Fatalf("viewport key changed input to %q", got)
	}
	updated, _ = m.Update(key(tea.KeyTab, ""))
	m = updated.(model)
	if !m.input.Focused() {
		t.Fatal("second Tab should return focus to the input")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = updated.(model)
	if m.input.Focused() {
		t.Fatal("Shift+Tab should switch focus")
	}
}

func TestResizePreservesStateAndRenderIsStable(t *testing.T) {
	m := initialModel()
	m.input.SetValue("saved")
	m.view.ScrollDown(3)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 42, Height: 12})
	m = updated.(model)
	if m.input.Value() != "saved" || m.view.YOffset() != 3 {
		t.Fatal("resize lost component state")
	}
	if m.input.Width() != 40 || m.view.Width() != 42 || m.view.Height() != 8 {
		t.Fatal("resize did not update component dimensions")
	}
	first := m.render()
	if first != m.render() {
		t.Fatal("render changed without an update")
	}
	if m.input.Value() != "saved" || !m.input.Focused() || m.view.YOffset() != 3 {
		t.Fatal("render changed component state")
	}
	if !strings.Contains(first, "saved") || !strings.Contains(first, "04  Bubbles viewport") {
		t.Fatal("resized render omitted component state")
	}
	if m.view.YOffset() != 3 {
		t.Fatal("render changed viewport position")
	}
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 1, Height: 2})
	m = updated.(model)
	if m.view.Height() != 0 {
		t.Fatalf("tiny layout should give viewport zero rows, got %d", m.view.Height())
	}
	assertRenderBounds(t, m)
}

func assertRenderBounds(t *testing.T, m model) {
	t.Helper()
	output := m.render()
	if m.width == 0 || m.height == 0 {
		if output != "" {
			t.Fatalf("zero-size render should be empty, got %q", output)
		}
		return
	}
	rows := strings.Split(output, "\n")
	if len(rows) > m.height {
		t.Fatalf("rendered %d rows into height %d: %q", len(rows), m.height, output)
	}
	for i, row := range rows {
		if width := ansi.StringWidth(row); width > m.width {
			t.Fatalf("row %d has width %d, window width %d", i, width, m.width)
		}
	}
}

func TestRenderBounds(t *testing.T) {
	for _, size := range []tea.WindowSizeMsg{
		{Width: 80, Height: 24},
		{Width: 42, Height: 12},
		{Width: 20, Height: 5},
		{Width: 19, Height: 4},
		{Width: 1, Height: 2},
		{Width: 0, Height: 2},
		{Width: 2, Height: 0},
		{Width: -1, Height: -2},
	} {
		m := initialModel()
		updated, _ := m.Update(size)
		assertRenderBounds(t, updated.(model))
	}
}

func TestQuitKeys(t *testing.T) {
	for _, msg := range []tea.KeyPressMsg{
		{Code: tea.KeyEsc},
		{Code: 'c', Mod: tea.ModCtrl},
	} {
		_, cmd := initialModel().Update(msg)
		if cmd == nil {
			t.Fatalf("%s did not return a quit command", msg.String())
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf("%s returned a command other than Quit", msg.String())
		}
	}
}

func TestSpinnerTickContinuesWithoutFocus(t *testing.T) {
	m := initialModel()
	m.input.Blur()
	updated, cmd := m.Update(m.spinner.Tick())
	if cmd == nil {
		t.Fatal("spinner tick should schedule the next tick")
	}
	if updated.(model).spinner.View() == m.spinner.View() {
		t.Fatal("spinner frame did not advance")
	}
}

func TestGrowingViewportClampsScroll(t *testing.T) {
	m := initialModel()
	m.resize(40, 8)
	m.view.ScrollDown(100)
	if m.view.YOffset() == 0 {
		t.Fatal("test needs a scrolled viewport")
	}
	m.resize(40, 40)
	if got := m.view.YOffset(); got != 4 {
		t.Fatalf("growing viewport should clamp offset to 4, got %d", got)
	}
}

func TestUnicodeAndPasteInput(t *testing.T) {
	long := strings.Repeat("界e\u0301🙂", 30)
	for _, tc := range []struct {
		name string
		msg  tea.Msg
		want string
	}{
		{name: "typed wide and combining", msg: key('界', "界e\u0301🙂"), want: "界e\u0301🙂"},
		{name: "pasted single line sanitization", msg: tea.PasteMsg{Content: "界\te\u0301\n🙂"}, want: "界 e\u0301 🙂"},
		{name: "long pasted text scrolls", msg: tea.PasteMsg{Content: long}, want: long},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := initialModel()
			updated, _ := m.Update(tc.msg)
			m = updated.(model)
			if got := m.input.Value(); got != tc.want {
				t.Fatalf("input value = %q, want %q", got, tc.want)
			}
			if !m.input.Focused() || m.view.YOffset() != 0 {
				t.Fatal("input changed focus or viewport scroll")
			}
			for _, size := range []tea.WindowSizeMsg{{Width: 80, Height: 24}, {Width: 20, Height: 5}, {Width: 1, Height: 2}, {Width: 42, Height: 12}} {
				updated, _ = m.Update(size)
				m = updated.(model)
				assertRenderBounds(t, m)
				if got := m.input.Value(); got != tc.want {
					t.Fatalf("resize changed input value to %q", got)
				}
			}
			if !strings.Contains(m.render(), "🙂") {
				t.Fatal("resized input no longer shows the visible end of its value")
			}
		})
	}
}

func TestPasteFollowsFocusAfterResize(t *testing.T) {
	m := initialModel()
	updated, _ := m.Update(tea.PasteMsg{Content: "start"})
	m = updated.(model)
	updated, _ = m.Update(key(tea.KeyTab, ""))
	m = updated.(model)
	updated, _ = m.Update(key(tea.KeyDown, ""))
	m = updated.(model)
	if m.input.Focused() || m.view.YOffset() != 1 {
		t.Fatal("viewport should be focused and scrolled")
	}
	updated, _ = m.Update(tea.PasteMsg{Content: "ignored\n界"})
	m = updated.(model)
	if m.input.Value() != "start" || m.view.YOffset() != 1 {
		t.Fatal("paste while viewport-focused changed component state")
	}
	for _, size := range []tea.WindowSizeMsg{{Width: 1, Height: 2}, {Width: 42, Height: 12}} {
		updated, _ = m.Update(size)
		m = updated.(model)
		assertRenderBounds(t, m)
	}
	if m.input.Focused() || m.input.Value() != "start" || m.view.YOffset() != 1 {
		t.Fatal("resize changed focus, input, or scroll")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = updated.(model)
	updated, _ = m.Update(tea.PasteMsg{Content: "界\te\u0301"})
	m = updated.(model)
	if !m.input.Focused() || m.input.Value() != "start界 e\u0301" || m.view.YOffset() != 1 {
		t.Fatal("paste after focus recovery did not reach the input")
	}
	assertRenderBounds(t, m)
}

func TestNarrowInputWithInteriorCursor(t *testing.T) {
	m := initialModel()
	long := strings.Repeat("界e\u0301🙂", 2) + "CURSOR" + strings.Repeat("界e\u0301🙂", 28)
	updated, _ := m.Update(tea.PasteMsg{Content: long})
	m = updated.(model)
	updated, _ = m.Update(key(tea.KeyHome, ""))
	m = updated.(model)
	for range 10 {
		updated, _ = m.Update(key(tea.KeyRight, ""))
		m = updated.(model)
	}
	before := m.input.Position()
	if before != 10 {
		t.Fatalf("cursor position = %d, want 10", before)
	}
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 20, Height: 5})
	m = updated.(model)
	if m.input.Value() != long || m.input.Position() != before || !m.input.Focused() {
		t.Fatal("narrowing changed input value, cursor, or focus")
	}
	assertRenderBounds(t, m)
	if !strings.Contains(ansi.Strip(m.input.View()), "RSOR") {
		t.Fatalf("resized input did not keep the cursor vicinity visible: %q", ansi.Strip(m.input.View()))
	}
}

func TestMinimumLayoutKeepsControlsVisible(t *testing.T) {
	for _, inputFocus := range []bool{true, false} {
		m := initialModel()
		m.resize(20, 5)
		if !inputFocus {
			m.input.Blur()
		}
		rows := strings.Split(ansi.Strip(m.render()), "\n")
		footer := rows[len(rows)-1]
		if !strings.Contains(footer, "Esc:quit") || !strings.Contains(footer, "Tab:focus") {
			t.Fatalf("essential controls hidden with input focus %v: %q", inputFocus, footer)
		}
	}
}

func TestInitialCommandsStartBothComponents(t *testing.T) {
	cmd := initialModel().Init()
	if cmd == nil {
		t.Fatal("Init did not start spinner and cursor commands")
	}
	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok || len(batch) != 2 || batch[0] == nil || batch[1] == nil {
		t.Fatalf("expected spinner and cursor commands, got %T", msg)
	}
}

func TestRepeatedResizePreservesCursorAndEditing(t *testing.T) {
	m := initialModel()
	value := strings.Repeat("界e\u0301🙂", 20)
	updated, _ := m.Update(tea.PasteMsg{Content: value})
	m = updated.(model)
	m.input.SetCursor(8)
	m.view.ScrollDown(3)
	for _, size := range []tea.WindowSizeMsg{
		{Width: 20, Height: 5}, {Width: 1, Height: 1}, {Width: 0, Height: 0},
		{Width: 42, Height: 12}, {Width: 80, Height: 24}, {Width: 20, Height: 5},
	} {
		updated, _ = m.Update(size)
		m = updated.(model)
		before := m.render()
		assertRenderBounds(t, m)
		if before != m.render() || m.input.Value() != value || m.input.Position() != 8 ||
			!m.input.Focused() || m.view.YOffset() != 3 {
			t.Fatalf("render/resize changed state at %dx%d", size.Width, size.Height)
		}
	}
	updated, _ = m.Update(key('X', "X"))
	m = updated.(model)
	runes := []rune(value)
	want := string(runes[:8]) + "X" + string(runes[8:])
	if m.input.Value() != want || m.input.Position() != 9 || m.view.YOffset() != 3 {
		t.Fatal("input did not resume editing at the retained cursor")
	}
	assertRenderBounds(t, m)
}

func TestViewportContentCannotDisplaceControls(t *testing.T) {
	m := initialModel()
	m.input.Blur()
	m.view.SetContent("\x1b[31m" + strings.Repeat("界\tlong line\r\n", 80) + ansi.ResetStyle)
	for _, size := range []tea.WindowSizeMsg{{Width: 20, Height: 5}, {Width: 42, Height: 12}, {Width: 80, Height: 24}} {
		updated, _ := m.Update(size)
		m = updated.(model)
		m.view.ScrollDown(2)
		offset := m.view.YOffset()
		assertRenderBounds(t, m)
		rows := strings.Split(ansi.Strip(m.render()), "\n")
		if len(rows) != size.Height || !strings.Contains(rows[len(rows)-1], "Esc:quit Tab:focus") {
			t.Fatalf("viewport displaced footer at %dx%d: %q", size.Width, size.Height, rows)
		}
		if m.view.YOffset() != offset || m.input.Focused() {
			t.Fatal("render changed viewport state or focus")
		}
	}
}

func TestQuitFromViewportFocus(t *testing.T) {
	for _, msg := range []tea.KeyPressMsg{{Code: tea.KeyEsc}, {Code: 'c', Mod: tea.ModCtrl}} {
		m := initialModel()
		m.input.Blur()
		_, cmd := m.Update(msg)
		if cmd == nil {
			t.Fatalf("%s did not quit from viewport", msg.String())
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf("%s returned a non-quit command", msg.String())
		}
	}
}
