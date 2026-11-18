package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/spik3r/flexgo"
)

func helpAppMessage(a App, msg tea.Msg) App {
	next, _ := a.Update(msg)
	return next.(App)
}

func helpAppKey(a App, code rune) App {
	return helpAppMessage(a, tea.KeyPressMsg{Code: code})
}

func helpAppSize(a App, w, h int) App {
	return helpAppMessage(a, tea.WindowSizeMsg{Width: w, Height: h})
}

func helpApp(w, h int, screen Screen) App {
	a := helpAppSize(NewApp(), w, h)
	if screen == ScreenScan {
		a = helpAppKey(a, '2')
	} else if screen == ScreenHistory {
		a = helpAppKey(a, '3')
	}
	return helpAppKey(a, '?')
}

func assertHelpFrame(t *testing.T, a App) {
	t.Helper()
	out := a.View().Content
	if lipgloss.Height(out) != a.height {
		t.Errorf("help%dx%d paints%drows", a.width, a.height, lipgloss.Height(out))
	}
	for i, row := range strings.Split(out, "\n") {
		if width := ansi.StringWidth(row); width != a.width {
			t.Errorf("help%dx%d row%d paints%dcells", a.width, a.height, i+1, width)
		}
	}
}

func helpTextRows(a App) []string {
	rows := strings.Split(ansi.Strip(a.View().Content), "\n")
	out := []string{}
	for i := 7; i < min(a.height-4, len(rows)); i++ {
		row := strings.TrimSpace(rows[i])
		row = strings.TrimSpace(strings.TrimRight(row, "│█"))
		out = append(out, row)
	}
	return out
}

func helpWithoutSpaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func assertHelpReachable(t *testing.T, a App, down rune) {
	t.Helper()
	assertHelpFrame(t, a)
	rows := helpTextRows(a)
	all := append([]string(nil), rows...)
	before := strings.Join(rows, "\n")
	for range 300 {
		a = helpAppKey(a, down)
		rows = helpTextRows(a)
		after := strings.Join(rows, "\n")
		if after == before {
			break
		}
		all = append(all, rows[len(rows)-1])
		before = after
	}
	corpus := helpWithoutSpaces(strings.Join(all, ""))
	count := 0
	for _, entry := range a.keys.HelpEntries() {
		if entry.Keys == "" && entry.Desc == "" {
			continue
		}
		count++
		for _, text := range []string{entry.Keys, entry.Desc} {
			if !strings.Contains(corpus, helpWithoutSpaces(ansi.Strip(text))) {
				t.Errorf("help%dx%d cannot reach%q", a.width, a.height, text)
			}
		}
	}
	if count != 12 {
		t.Fatalf("expected12actions got%d", count)
	}
	assertHelpFrame(t, a)
}

func TestHelpCompactTargetsAndAllActions(t *testing.T) {
	for _, screen := range []Screen{ScreenLauncher, ScreenScan, ScreenHistory} {
		for _, size := range [][2]int{{80, 24}, {40, 16}} {
			a := helpApp(size[0], size[1], screen)
			beforeLauncher, beforeScan, beforeHistory := a.launcher, a.scan, a.history
			assertHelpReachable(t, a, tea.KeyDown)
			for _, code := range []rune{'q', '1', '2', '3', tea.KeyTab, tea.KeyEnter, 'x', tea.KeySpace} {
				a = helpAppKey(a, code)
			}
			if !a.modalOpen || a.screen != screen || !reflect.DeepEqual(a.launcher, beforeLauncher) || !reflect.DeepEqual(a.scan, beforeScan) || !reflect.DeepEqual(a.history, beforeHistory) {
				t.Fatal("modal input leaked to underlying screen")
			}
		}
	}
}

func TestHelpCompactLinePageAndEnds(t *testing.T) {
	a := helpApp(40, 16, ScreenLauncher)
	initial := a.View().Content
	down := helpAppKey(a, tea.KeyDown)
	if down.View().Content == initial {
		t.Error("Down did not move help")
	}
	if helpAppKey(down, tea.KeyUp).View().Content != initial {
		t.Error("Up did not reverse one row")
	}
	ten := a
	for range 10 {
		ten = helpAppKey(ten, tea.KeyDown)
	}
	page := helpAppKey(a, tea.KeyPgDown)
	if page.View().Content != ten.View().Content {
		t.Error("PageDown differs from ten Down keys")
	}
	if page.View().Content == initial {
		t.Error("PageDown did not move help")
	}
	back := page
	for range 10 {
		back = helpAppKey(back, tea.KeyUp)
	}
	if helpAppKey(page, tea.KeyPgUp).View().Content != back.View().Content {
		t.Error("PageUp differs from ten Up keys")
	}
	bottom := a
	for range 100 {
		bottom = helpAppKey(bottom, tea.KeyDown)
	}
	if helpAppKey(a, tea.KeyEnd).View().Content != bottom.View().Content {
		t.Error("Bottom did not clamp to last page")
	}
	if !strings.Contains(ansi.Strip(bottom.View().Content), "close modal / back") {
		t.Error("last Close action not visible at Bottom")
	}
	if helpAppKey(bottom, tea.KeyUp).View().Content == bottom.View().Content {
		t.Error("Up at Bottom did not move")
	}
	if helpAppKey(bottom, tea.KeyHome).View().Content != initial {
		t.Error("Top did not restore first page")
	}
}

func TestHelpCompactResizeResetAndPurePaint(t *testing.T) {
	a := helpApp(40, 16, ScreenLauncher)
	a.launcher.cursor = 4
	initial := a.View().Content
	a = helpAppKey(a, tea.KeyEnd)
	bottom := a.View().Content
	if bottom == initial {
		t.Error("Bottom did not move before resize")
	}
	before := a
	_ = a.View()
	_ = a.View()
	if !reflect.DeepEqual(a, before) {
		t.Fatal("paint changed model")
	}
	for _, size := range [][2]int{{0, 0}, {math.MinInt, 16}, {40, math.MinInt}, {6, 16}, {40, 10}} {
		a = helpAppSize(a, size[0], size[1])
		for _, code := range []rune{tea.KeyDown, tea.KeyUp, tea.KeyPgDown, tea.KeyPgUp, tea.KeyHome, tea.KeyEnd} {
			a = helpAppKey(a, code)
		}
	}
	a = helpAppSize(a, 40, 16)
	if a.View().Content != bottom {
		t.Error("unallocated resize or input lost scroll intent")
	}
	a = helpAppSize(a, 80, 24)
	assertHelpFrame(t, a)
	a = helpAppSize(a, 40, 16)
	if a.View().Content != initial {
		t.Error("larger viewport did not clamp to its legal top row")
	}
	a = helpAppKey(a, tea.KeyEnd)
	a = helpAppSize(a, 100, 40)
	assertHelpFrame(t, a)
	a = helpAppSize(a, 40, 16)
	if a.View().Content != initial {
		t.Error("full-card resize did not reconcile help")
	}
	a = helpAppKey(a, tea.KeyEnd)
	a = helpAppKey(a, tea.KeyEscape)
	a = helpAppKey(a, '?')
	if a.View().Content != initial || a.launcher.cursor != 4 {
		t.Error("reopen failed to reset help or preserved cursor changed")
	}
}

func helpFooter(a App) string {
	rows := strings.Split(ansi.Strip(a.View().Content), "\n")
	if a.height-3 >= len(rows) {
		return ""
	}
	return strings.TrimSpace(rows[a.height-3])
}

func TestHelpCompactRemapsCollisionsAndEmptyHints(t *testing.T) {
	a := helpAppSize(NewApp(), 40, 16)
	a.keys.Help = []string{"!"}
	a.keys.Close = []string{"x"}
	a.keys.Up = []string{"x", "u"}
	a.keys.Down = []string{"!", "d"}
	a.keys.PageUp = []string{"p"}
	a.keys.PageDown = []string{"x", "n"}
	a.keys.Top = []string{"!", "t"}
	a.keys.Bottom = []string{"x", "b"}
	a = helpAppKey(a, '!')
	initial := a.View().Content
	hint := helpFooter(a)
	if !strings.Contains(hint, "x close") || !strings.Contains(hint, "u/d") || strings.Contains(hint, "x/!") {
		t.Errorf("unsafe or stale compact hint%q", hint)
	}
	assertHelpReachable(t, a, 'd')
	if helpAppKey(a, 'd').View().Content == initial {
		t.Error("safe remapped Down did not move")
	}
	if helpAppKey(a, 'b').View().Content == initial {
		t.Error("safe remapped Bottom did not move")
	}
	for _, code := range []rune{'x', '!'} {
		closed := helpAppKey(a, code)
		if closed.modalOpen {
			t.Errorf("Close/Help%q lost priority to navigation", code)
		}
	}
	a.keys = KeyMap{}
	if hint := helpFooter(a); hint != "" {
		t.Errorf("empty bindings invented hint%q", hint)
	}
	for _, code := range []rune{'q', '?', 'x', 'j', 'k', tea.KeyEnd} {
		next := helpAppKey(a, code)
		if !next.modalOpen || !reflect.DeepEqual(next, a) {
			t.Error("empty bindings invented an action")
		}
	}
}

func TestHelpCompactLongWideLabelsUseWrappedViewport(t *testing.T) {
	for _, size := range [][2]int{{40, 16}, {80, 24}, {80, 40}, {100, 40}, {120, 40}} {
		a := helpAppSize(NewApp(), size[0], size[1])
		a.keys.Next = []string{strings.Repeat("界👩‍💻", 30) + "-next-label"}
		a.keys.Prev = []string{strings.Repeat("prevlabel", 20)}
		a = helpAppKey(a, '?')
		assertHelpReachable(t, a, tea.KeyDown)
	}
}

func TestHelpCompactCardBoundaries(t *testing.T) {
	for _, w := range []int{54, 55, 56, 57, 58, 74} {
		for _, h := range []int{22, 23, 24, 25, 26, 27, 28} {
			a := helpApp(w+6, h+10, ScreenLauncher)
			assertHelpFrame(t, a)
		}
	}
}

func TestHelpFittingPaintBytes(t *testing.T) {
	want := map[string]string{"100x40-screen0": "29a018991b67d419d282deb47b46126314e3e91fa1074592329480551e8b4642",
		"100x40-screen1": "519976321e6dd6fde1077bc37ed0ca812e4d9b26c1ef8cebe725c1c596a6b70f",
		"100x40-screen2": "2877b23ecbac849a098303531be26eeac9fe7ccb39881288296ad97e4fd0d25c",
		"120x40-screen0": "e84b55c978cd70f5f7e7f794ffd2d8694502b539d94f63e6e7e235bd5eac7099",
		"120x40-screen1": "3b861c7af64cc6468baa81e284c3cb41df2523ede13898faafded7ee9f17ef2e",
		"120x40-screen2": "887cabd820226e177f87e734797f40cb2973c643d58217f1727acec70a2c6766",
	}
	for _, screen := range []Screen{ScreenLauncher, ScreenScan, ScreenHistory} {
		for _, size := range [][2]int{{100, 40}, {120, 40}} {
			a := helpApp(size[0], size[1], screen)
			out := a.View().Content
			key := fmt.Sprintf("%dx%d-screen%d", size[0], size[1], screen)
			got := fmt.Sprintf("%x", sha256.Sum256([]byte(out)))
			if got != want[key] {
				t.Errorf("%s: fitting ANSI bytes changed: sha256%s", key, got)
			}
		}
	}
}

func TestHelpCompactModelMatchesActualCallback(t *testing.T) {
	for _, size := range [][2]int{{40, 16}, {80, 24}, {60, 32}, {62, 34}, {62, 35}, {62, 36}, {80, 36}, {0, 0}, {math.MinInt, 16}, {40, math.MinInt}} {
		a := helpApp(size[0], size[1], ScreenLauncher)
		layout := a.helpLayout()
		if !layout.compact {
			continue
		}
		body := a.helpBody()
		panel := body.Children[1]
		got := viewportSize{}
		tree := pageShell("", a.screen, body, "")
		var capture func(*flexgo.Node)
		capture = func(n *flexgo.Node) {
			if n.View != nil {
				n.View = func(w, h int) string {
					if n == panel {
						got = viewportSize{w, h}
					}
					return ""
				}
			}
			for _, child := range n.Children {
				capture(child)
			}
		}
		capture(tree)
		_ = tree.Render(a.width, a.height)
		if layout.size != got {
			t.Errorf("window%v update allocation%v actual callback%v", size, layout.size, got)
		}
	}
}

func TestHelpCompactWrappedCountAndOwnedOffset(t *testing.T) {
	for _, wide := range []bool{false, true} {
		a := helpApp(40, 16, ScreenLauncher)
		if wide {
			a.keys.Next = []string{strings.Repeat("界👩‍💻", 30) + "-next-label"}
			a.keys.Prev = []string{strings.Repeat("prevlabel", 20)}
		}
		layout := a.helpLayout()
		if !layout.compact || layout.size != (viewportSize{34, 5}) {
			t.Fatalf("unexpected compact layout%+v", layout)
		}
		if wide && len(layout.rows) <= 12 {
			t.Fatal("long labels did not add wrapped rows")
		}
		a.help.Offset = math.MaxInt
		a = helpAppKey(a, tea.KeyDown)
		maximum := max(0, len(layout.rows)-layout.size.height)
		if a.help.Offset != maximum {
			t.Errorf("Bottom projected offset%d want%d", a.help.Offset, maximum)
		}
		a = helpAppKey(a, tea.KeyUp)
		if a.help.Offset != maximum-1 {
			t.Error("reverse movement did not start at legal top row")
		}
		a = helpAppKey(a, tea.KeyHome)
		a = helpAppKey(a, tea.KeyPgDown)
		if a.help.Offset != min(10, maximum) {
			t.Error("PageDown changed its ten-row step")
		}
		before := a.help.Offset
		view := a.helpBody().Children[1].View
		first := view(34, 5)
		_ = view(34, 10)
		_ = view(0, 5)
		last := view(34, 5)
		if first != last || a.help.Offset != before {
			t.Fatal("callback reuse changed captured or owned offset")
		}
		a = helpAppKey(a, tea.KeyEnd)
		a = helpAppSize(a, 40, 17)
		resized := a.helpLayout()
		if a.help.Offset != max(0, len(resized.rows)-resized.size.height) {
			t.Fatal("resize did not reconcile against painted wrapped rows")
		}
		before = a.help.Offset
		a = helpAppSize(a, 40, 16)
		if a.help.Offset != before {
			t.Fatal("shrink lost legal top-row intent")
		}
	}
}

func TestHelpResizeKeepsUnderlyingResizeContract(t *testing.T) {
	scan := helpApp(80, 22, ScreenScan)
	scan.scan.activeTab = 2
	scan.scan.logs = viewportLines(39)
	scan.scan.viewports[2].Offset = 29
	scan.scan.filesCursor = 3
	scan.launcher.cursor = 4
	scan = helpAppSize(scan, 80, 32)
	if scan.scan.viewports[2].Offset != 19 || scan.scan.filesCursor != 3 || scan.launcher.cursor != 4 || scan.screen != ScreenScan {
		t.Fatal("help resize changed underlying resize/selection contract")
	}
	history := helpApp(80, 20, ScreenHistory)
	history.history.results = []ScanResult{{Findings: 31}}
	history.history.detail.Offset = 29
	history.launcher.cursor = 4
	history = helpAppSize(history, 80, 30)
	if history.history.detail.Offset != 19 || history.history.cursor != 0 || history.launcher.cursor != 4 || history.screen != ScreenHistory {
		t.Fatal("help resize changed underlying History contract")
	}
}
