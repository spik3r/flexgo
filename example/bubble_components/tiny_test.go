package main

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
)

func TestTinyWindowIgnoresHiddenInteractions(t *testing.T) {
	for _, size := range []tea.WindowSizeMsg{{Width: 19, Height: 8}, {Width: 20, Height: 4}, {Width: 0, Height: 8}, {Width: 20, Height: 0}, {Width: -1, Height: -1}} {
		for _, inputFocus := range []bool{true, false} {
			for _, msg := range []tea.Msg{
				key('X', "X"), key(tea.KeyBackspace, ""), key(tea.KeyTab, ""),
				tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift},
				tea.KeyPressMsg{Code: 'v', Mod: tea.ModCtrl},
				key(tea.KeyDown, ""), tea.PasteMsg{Content: "hidden"},
				tea.MouseWheelMsg{Button: tea.MouseWheelDown}, cursor.Blink(),
			} {
				m := initialModel()
				m.input.SetValue("saved")
				m.input.SetCursor(2)
				m.view.ScrollDown(3)
				if !inputFocus {
					m.input.Blur()
				}
				updated, _ := m.Update(size)
				m = updated.(model)
				updated, cmd := m.Update(msg)
				after := updated.(model)
				if after.input.Value() != "saved" || after.input.Position() != 2 ||
					after.input.Focused() != inputFocus || after.view.YOffset() != 3 || cmd != nil {
					t.Errorf("%v focus=%v hidden %T changed state or issued command: value=%q cursor=%d focus=%v offset=%d cmd=%v",
						size, inputFocus, msg, after.input.Value(), after.input.Position(), after.input.Focused(), after.view.YOffset(), cmd != nil)
				}
				assertRenderBounds(t, after)
			}
		}
	}
}

func TestTinyWindowQuitAndSpinnerRemainActive(t *testing.T) {
	for _, inputFocus := range []bool{true, false} {
		m := initialModel()
		if !inputFocus {
			m.input.Blur()
		}
		for _, size := range []tea.WindowSizeMsg{{Width: 19, Height: 8}, {Width: 80, Height: 4}, {Width: 0, Height: 0}} {
			updated, _ := m.Update(size)
			m = updated.(model)
			for _, quit := range []tea.KeyPressMsg{{Code: tea.KeyEsc}, {Code: 'c', Mod: tea.ModCtrl}} {
				_, cmd := m.Update(quit)
				if cmd == nil {
					t.Fatalf("%v: quit unavailable", size)
				}
				if _, ok := cmd().(tea.QuitMsg); !ok {
					t.Fatal("not a quit command")
				}
			}
			updated, cmd := m.Update(m.spinner.Tick())
			if cmd == nil || updated.(model).spinner.View() == m.spinner.View() {
				t.Fatal("hidden spinner stopped its update chain")
			}
			if updated.(model).input.Focused() != inputFocus {
				t.Fatal("spinner changed focus")
			}
		}
	}
}

func TestTinyWindowRecoveryResumesFocusedComponent(t *testing.T) {
	for _, inputFocus := range []bool{true, false} {
		m := initialModel()
		m.input.SetValue("saved")
		m.input.SetCursor(2)
		m.view.ScrollDown(3)
		if !inputFocus {
			m.input.Blur()
		}
		for range 3 {
			updated, cmd := m.Update(tea.WindowSizeMsg{Width: 1, Height: 1})
			m = updated.(model)
			if cmd != nil {
				t.Fatal("hiding components restarted the cursor")
			}
			updated, cmd = m.Update(tea.WindowSizeMsg{Width: 20, Height: 5})
			m = updated.(model)
			if (cmd != nil) != inputFocus {
				t.Fatal("recovery did not restart only the focused input cursor")
			}
			if m.input.Focused() != inputFocus || m.input.Value() != "saved" || m.input.Position() != 2 || m.view.YOffset() != 3 {
				t.Fatal("hide/restore changed component state")
			}
			if inputFocus && cursorCells(strings.Split(m.render(), "\n")[1]) != 1 {
				t.Fatal("recovered input did not show its cursor cell")
			}
			updated, cmd = m.Update(tea.WindowSizeMsg{Width: 42, Height: 12})
			m = updated.(model)
			if cmd != nil {
				t.Fatal("ordinary resize restarted the cursor")
			}
			assertRenderBounds(t, m)
		}
		if inputFocus {
			updated, _ := m.Update(key('X', "X"))
			m = updated.(model)
			if m.input.Value() != "saXved" || m.input.Position() != 3 {
				t.Fatal("editing did not resume at the retained cursor")
			}
		} else {
			updated, _ := m.Update(key(tea.KeyDown, ""))
			m = updated.(model)
			if m.view.YOffset() != 4 || m.input.Value() != "saved" {
				t.Fatal("scrolling did not resume in the retained viewport")
			}
		}
	}
}

func TestTinyWindowHintAndBounds(t *testing.T) {
	m := initialModel()
	for _, size := range []tea.WindowSizeMsg{{Width: 19, Height: 4}, {Width: 80, Height: 4}, {Width: 14, Height: 2}, {Width: 13, Height: 2}, {Width: 8, Height: 2}, {Width: 1, Height: 1}, {Width: 0, Height: 0}} {
		updated, _ := m.Update(size)
		m = updated.(model)
		assertRenderBounds(t, m)
		if size.Width >= len("Resize to 20x5") && !strings.Contains(m.render(), "Resize to 20x5") {
			t.Fatal("missing actionable resize hint")
		}
		if size.Width >= 8 && size.Height >= 2 && !strings.Contains(m.render(), "Esc:quit") {
			t.Fatal("missing available quit hint")
		}
	}
}
