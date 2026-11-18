package main

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func launcherAppMessage(a App, msg tea.Msg) App {
	next, _ := a.Update(msg)
	return next.(App)
}

func launcherAppKey(a App, code rune) App {
	return launcherAppMessage(a, tea.KeyPressMsg{Code: code})
}

func launcherAppSize(a App, w, h int) App {
	return launcherAppMessage(a, tea.WindowSizeMsg{Width: w, Height: h})
}

func launcherSnapshot(l LauncherScreen) LauncherScreen {
	l.fields = append([]string(nil), l.fields...)
	l.profile.Include = append([]string(nil), l.profile.Include...)
	l.profile.Exclude = append([]string(nil), l.profile.Exclude...)
	return l
}

func assertLauncherFrame(t *testing.T, out string, w, h int) {
	t.Helper()
	if lipgloss.Height(out) != h {
		t.Errorf("window%dx%d painted %d rows", w, h, lipgloss.Height(out))
	}
	for i, line := range strings.Split(out, "\n") {
		if width := ansi.StringWidth(line); width != w {
			t.Errorf("window%dx%d row%d painted %d cells", w, h, i+1, width)
		}
	}
}

func assertLauncherSelected(t *testing.T, a App) {
	t.Helper()
	out := a.View().Content
	assertLauncherFrame(t, out, a.width, a.height)
	label := a.launcher.fields[a.launcher.cursor] + ":"
	selected := lipgloss.NewStyle().Background(colSelected).Foreground(colHighlight).Render("X")
	prefix, _, _ := strings.Cut(selected, "X")
	found := false
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(ansi.Strip(line), label) && prefix != "" && strings.Contains(line, prefix) {
			found = true
		}
	}
	if !found {
		t.Errorf("selected %q missing its visible highlight", label)
	}
	for _, field := range []string{"Name:", "Root:", "Include:", "Exclude:", "Deep:"} {
		if !strings.Contains(ansi.Strip(out), field) {
			t.Errorf("field%q not visible", field)
		}
	}
	plain := ansi.Strip(out)
	if !strings.Contains(plain, "2 scan") && !strings.Contains(plain, "2 start scan") && !strings.Contains(plain, "press 2 to start scanning") {
		t.Error("missing short scan start hint")
	}
	if !strings.Contains(plain, "? help") || !strings.Contains(plain, "1/2/3 screens") || !strings.Contains(plain, "q quit") {
		t.Error("global navigation hint missing")
	}
}

func TestLauncherCompactAppBoundsAndCursor(t *testing.T) {
	for _, size := range [][2]int{{80, 24}, {40, 16}} {
		for _, keys := range [][2]rune{{tea.KeyDown, tea.KeyUp}, {tea.KeyTab, tea.KeyLeft}, {'j', 'k'}} {
			a := launcherAppSize(NewApp(), size[0], size[1])
			original := launcherSnapshot(a.launcher)
			for cursor := 0; cursor < 5; cursor++ {
				if a.launcher.cursor != cursor {
					t.Fatalf("forward cursor%d want%d", a.launcher.cursor, cursor)
				}
				assertLauncherSelected(t, a)
				a = launcherAppKey(a, keys[0])
			}
			if a.launcher.cursor != 4 {
				t.Fatal("cursor went past final field")
			}
			for cursor := 4; cursor >= 0; cursor-- {
				if a.launcher.cursor != cursor {
					t.Fatalf("backward cursor%d want%d", a.launcher.cursor, cursor)
				}
				assertLauncherSelected(t, a)
				a = launcherAppKey(a, keys[1])
			}
			if !reflect.DeepEqual(a.launcher, original) {
				t.Fatal("navigation changed profile or field data")
			}
		}
	}
}

func TestLauncherCompactResizeAndPaintPurity(t *testing.T) {
	for _, cursor := range []int{0, 1, 2, 3, 4} {
		a := launcherAppSize(NewApp(), 100, 40)
		for range cursor {
			a = launcherAppKey(a, tea.KeyTab)
		}
		before := launcherSnapshot(a.launcher)
		normal := a.View().Content
		for _, size := range [][2]int{{80, 24}, {40, 16}, {100, 40}, {40, 16}, {120, 40}, {100, 40}} {
			a = launcherAppSize(a, size[0], size[1])
			first := a.View().Content
			if first != a.View().Content {
				t.Fatal("repeated paint changed output")
			}
			if !reflect.DeepEqual(a.launcher, before) || a.screen != ScreenLauncher {
				t.Fatal("resize or paint changed Launcher state")
			}
			assertLauncherFrame(t, first, size[0], size[1])
		}
		if a.View().Content != normal {
			t.Fatal("normal bytes changed after resize round trip")
		}
	}
}

func TestLauncherCompactCustomProfileAndInertInput(t *testing.T) {
	profile := Profile{Name: "[red]" + strings.Repeat("界👩‍💻", 30), Root: strings.Repeat("root/", 40), Include: []string{strings.Repeat("long", 30), "界.go"}, Exclude: []string{strings.Repeat("wide界", 30)}, Deep: false}
	for _, size := range [][2]int{{80, 24}, {40, 16}} {
		a := NewApp()
		a.launcher = NewLauncherScreen(profile)
		a = launcherAppSize(a, size[0], size[1])
		for cursor := 0; cursor < 5; cursor++ {
			before := launcherSnapshot(a.launcher)
			assertLauncherSelected(t, a)
			if !strings.Contains(ansi.Strip(a.View().Content), "[red]界👩‍💻") {
				t.Error("custom value markup or leading wide graphemes changed")
			}
			for _, msg := range []tea.KeyPressMsg{{Code: 'x', Text: "x"}, {Code: tea.KeyEnter}, {Code: tea.KeySpace, Text: " "}} {
				a = launcherAppMessage(a, msg)
				if !reflect.DeepEqual(a.launcher, before) {
					t.Fatal("text, Enter or Space edited profile/cursor")
				}
			}
			if !reflect.DeepEqual(a.launcher.Profile(), profile) {
				t.Fatal("custom profile changed")
			}
			a = launcherAppKey(a, tea.KeyTab)
		}
	}
}

func TestLauncherCompactRoutingAndRemappedKeys(t *testing.T) {
	a := launcherAppSize(NewApp(), 40, 16)
	a.keys.Down = []string{"n"}
	a.keys.Up = []string{"p"}
	a.keys.Next = []string{"f"}
	a.keys.Prev = []string{"b"}
	a.keys.Help = []string{"?"}
	a.keys.Scan = []string{"s"}
	a.keys.Launch = []string{"v"}
	a.keys.History = []string{"r"}
	a.keys.Quit = []string{"z"}
	a.keys.Close = []string{"c"}
	for _, code := range []rune{'n', 'f', 'n', 'f'} {
		a = launcherAppKey(a, code)
	}
	if a.launcher.cursor != 4 {
		t.Fatal("remapped forward keys failed")
	}
	for _, code := range []rune{'p', 'b'} {
		a = launcherAppKey(a, code)
	}
	if a.launcher.cursor != 2 {
		t.Fatal("remapped backward keys failed")
	}
	before := launcherSnapshot(a.launcher)
	a = launcherAppKey(a, '?')
	for _, code := range []rune{'n', 'p', 's', 'r', 'z', 'v'} {
		a = launcherAppKey(a, code)
	}
	if !a.modalOpen || a.screen != ScreenLauncher || !reflect.DeepEqual(a.launcher, before) {
		t.Fatal("help did not capture navigation")
	}
	a = launcherAppKey(a, 'c')
	assertLauncherSelected(t, a)
	a = launcherAppKey(a, 's')
	if a.screen != ScreenScan || !reflect.DeepEqual(a.scan.profile, before.Profile()) {
		t.Fatal("synthetic Scan handoff changed profile")
	}
	a = launcherAppKey(a, 'v')
	if a.screen != ScreenLauncher || !reflect.DeepEqual(a.launcher, before) {
		t.Fatal("return to Launcher changed state")
	}
	a = launcherAppKey(a, 'r')
	if a.screen != ScreenHistory {
		t.Fatal("remapped History navigation failed")
	}
	a = launcherAppKey(a, 'v')
	_, cmd := a.Update(tea.KeyPressMsg{Code: 'z'})
	if cmd == nil {
		t.Fatal("remapped quit lost its command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("quit command changed")
	}
}

func TestLauncherFittingPaintBytes(t *testing.T) {
	want := map[string]string{"100x40-cursor0": "d9b5eeb081b1bbf57519ef5d2416ad2fccc2f5c763ea49cf947146842c1d0dd8",
		"100x40-cursor1": "e7b81f8bcb73b4fd091ff9bb017d35cd5e66c88b2a06f0a93188e5c5b26298e0",
		"100x40-cursor2": "5193a80eff580a3ebbdab8d7dd03e47ad83d765dd6d1fdefb280f32f63c5aa71",
		"100x40-cursor3": "62f318e5fe001d6bed3c6f9c75c1bb17894efa18169b0b53b0fd3b7201ad43ff",
		"100x40-cursor4": "e64f909b5cf3723b90c66bfd386f0fcb9156a34e95312e19d1c18a1ed655982e",
		"120x40-cursor0": "a3805ad34bcd5bdc590ecd1191ee95d726c25b208929bf4178e9ce5a864db025",
		"120x40-cursor1": "cf81a8dccb8419c11453ccd04156eb5fb207685cf75757ef3bcbe1ff034b324d",
		"120x40-cursor2": "b6451fc55ee479d1cc1635d095f95d0a7eaab1a6374c7fe19b93d57a179596e6",
		"120x40-cursor3": "cb0627394e1660d5498b4e0548f0ea1fad34a05a247e9f13bfdf07fe1de729a8",
		"120x40-cursor4": "cbc0cf3830a8de9e7bb6acacb5def1017ad1b8732491ba5bbd724c9cb19a5195",
	}
	for _, size := range [][2]int{{100, 40}, {120, 40}} {
		a := launcherAppSize(NewApp(), size[0], size[1])
		for cursor := 0; cursor < 5; cursor++ {
			out := a.View().Content
			key := fmt.Sprintf("%dx%d-cursor%d", size[0], size[1], cursor)
			got := fmt.Sprintf("%x", sha256.Sum256([]byte(out)))
			if got != want[key] {
				t.Errorf("%s: fitting ANSI bytes changed: sha256%s", key, got)
			}
			a = launcherAppKey(a, tea.KeyTab)
		}
	}
}

func TestLauncherCenteredCardBoundaryFitsBodyAndApp(t *testing.T) {
	for _, w := range []int{58, 59, 60, 61, 62, 63, 64, 74, 94} {
		for _, h := range []int{19, 20, 21, 22, 23} {
			l := NewLauncherScreen(DefaultProfile())
			l.cursor = 4
			out := l.Body(w, h).Render(w, h)
			assertLauncherFrame(t, out, w, h)
			a := NewApp()
			a.launcher = l
			a = launcherAppSize(a, w+6, h+10)
			assertLauncherSelected(t, a)
		}
	}
}
