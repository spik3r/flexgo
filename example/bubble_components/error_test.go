package main

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestInputValidationErrorIsVisible(t *testing.T) {
	m := initialModel()
	m.input.Validate = func(value string) error {
		if strings.Contains(value, "!") {
			return errors.New("Use plain text")
		}
		return nil
	}
	updated, _ := m.Update(tea.PasteMsg{Content: "!"})
	m = updated.(model)
	if m.input.Err == nil {
		t.Fatal("test did not produce a component error")
	}
	rows := strings.Split(ansi.Strip(m.render()), "\n")
	if !strings.HasPrefix(rows[2], "Error: Use plain text") || !strings.HasPrefix(rows[3], "01  Bubbles viewport") {
		t.Fatal("component error is missing or displaced the viewport")
	}
	if m.input.Value() != "!" || m.input.Position() != 1 || !m.input.Focused() {
		t.Fatal("feedback changed input state")
	}
	updated, _ = m.Update(key(tea.KeyBackspace, ""))
	m = updated.(model)
	if m.input.Err != nil || strings.Contains(m.render(), "Error:") || m.input.Value() != "" {
		t.Fatal("feedback did not clear after the corrected input")
	}
}

func TestInputErrorClearsOnSyntheticPasteRecovery(t *testing.T) {
	m := initialModel()
	m.input.SetValue("saved")
	m.input.Err = errors.New("Clipboard unavailable")
	if !strings.Contains(m.render(), "Error: Clipboard unavailable") {
		t.Fatal("stored component failure is not visible")
	}
	updated, _ := m.Update(tea.PasteMsg{Content: " retry"})
	m = updated.(model)
	if m.input.Err != nil || strings.Contains(m.render(), "Error:") || m.input.Value() != "saved retry" {
		t.Fatal("successful synthetic paste did not replace the failure state")
	}
}

func TestInputErrorIsBoundedWithoutChangingState(t *testing.T) {
	for _, detail := range []string{
		"Synthetic input failure",
		strings.Repeat("long error ", 30),
		"first error\nsecond line\nthird line",
		"\x1b[31m" + ansi.SetHyperlink("https://example.test") + "styled error\nclosure" + ansi.ResetStyle + ansi.ResetHyperlink(),
		"界\twide error\r\nmore",
	} {
		for _, size := range []tea.WindowSizeMsg{{Width: 20, Height: 5}, {Width: 42, Height: 12}, {Width: 80, Height: 24}} {
			m := initialModel()
			m.input.SetValue("saved")
			m.input.SetCursor(2)
			m.input.Err = errors.New(detail)
			m.view.ScrollDown(3)
			m.resize(size.Width, size.Height)
			beforeErr := m.input.Err
			out := m.render()
			assertRenderBounds(t, m)
			rows := strings.Split(ansi.Strip(out), "\n")
			if len(rows) != size.Height || !strings.HasPrefix(rows[2], "Error:") || !strings.HasPrefix(rows[3], "04") ||
				!strings.Contains(rows[len(rows)-1], "Esc:quit Tab:focus") {
				t.Fatalf("error displaced content at %v: %q", size, out)
			}
			if m.input.Err != beforeErr || m.input.Value() != "saved" || m.input.Position() != 2 || !m.input.Focused() || m.view.YOffset() != 3 {
				t.Fatal("rendering feedback changed component state")
			}
			if out != m.render() {
				t.Fatal("feedback rendering is not stable")
			}
		}
	}
}

func TestInputErrorPersistsAcrossFocusAndTinyRecovery(t *testing.T) {
	m := initialModel()
	m.input.SetValue("saved")
	m.input.Err = errors.New("Synthetic input failure")
	failure := m.input.Err
	for range 3 {
		updated, _ := m.Update(key(tea.KeyTab, ""))
		m = updated.(model)
		if m.input.Focused() || !strings.Contains(m.render(), "Error: Synthetic input failure") {
			t.Fatal("focus switch hid or cleared the input failure")
		}
		updated, _ = m.Update(tea.WindowSizeMsg{Width: 19, Height: 4})
		m = updated.(model)
		updated, _ = m.Update(tea.PasteMsg{Content: "ignored"})
		m = updated.(model)
		if m.input.Err != failure || m.input.Value() != "saved" || strings.Contains(m.render(), "Error:") {
			t.Fatal("hidden error state changed or replaced the resize hint")
		}
		updated, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
		m = updated.(model)
		if m.input.Focused() || m.input.Err != failure || !strings.Contains(m.render(), "Error: Synthetic input failure") {
			t.Fatal("recovery lost the failure or selected focus")
		}
		updated, _ = m.Update(key(tea.KeyTab, ""))
		m = updated.(model)
	}
}
