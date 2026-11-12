package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

func TestResizeBoundsAndTileRetention(t *testing.T) {
	m := initialModel()
	root := m.root
	for _, size := range [][2]int{{0, 0}, {-1, 20}, {20, -1}, {1, 1}, {17, 20}, {18, 16}, {18, 17}, {31, 20}, {32, 20}, {46, 20}, {47, 20}, {80, 20}, {120, 30}, {18, 17}} {
		updated, _ := m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		m = updated.(model)
		before := m
		out := m.View().Content
		if !reflect.DeepEqual(m, before) || m.root != root {
			t.Fatal("resize/render replaced component state")
		}
		if m.width == 0 || m.height == 0 {
			if out != "" {
				t.Fatalf("zero area: %q", out)
			}
			continue
		}
		if lipgloss.Width(out) > m.width || lipgloss.Height(out) > m.height {
			t.Fatalf("%v: rendered %dx%d", size, lipgloss.Width(out), lipgloss.Height(out))
		}
		if m.width >= 18 && m.height >= 17 {
			for i := 1; i <= tileCount; i++ {
				if !strings.Contains(out, fmt.Sprintf("Tile %02d", i)) {
					t.Fatalf("lost tile %d at %v", i, size)
				}
			}
		} else if !strings.Contains(out, "R") {
			t.Fatal("missing resize hint")
		}
	}
	if err := flexgo.Validate(root); err != nil {
		t.Fatal(err)
	}
}

func TestQuitAtTinyAndUsableSizes(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {80, 20}} {
		m := initialModel()
		m.width, m.height = size[0], size[1]
		_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
		if cmd == nil {
			t.Fatal("quit key ignored")
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatal("not a quit command")
		}
	}
}
