package main

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/spik3r/flexgo"
)

func viewportAppMessage(a App, msg tea.Msg) App {
	next, _ := a.Update(msg)
	return next.(App)
}

func viewportAppKey(a App, code rune) App {
	return viewportAppMessage(a, tea.KeyPressMsg{Code: code})
}

func viewportLines(count int) []string {
	lines := make([]string, count)
	for i := range lines {
		lines[i] = fmt.Sprintf("row-%02d", i)
	}
	return lines
}

func viewportScanApp(w, h int) App {
	a := NewApp()
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: w, Height: h})
	a = viewportAppKey(a, '2')
	a = viewportAppKey(a, tea.KeyTab)
	a = viewportAppKey(a, tea.KeyTab)
	a.scan.logs = viewportLines(39)
	return a
}

// Intercept the real leaf allocation while suppressing unrelated paint output.
func viewportAllocated(a App) (int, int) {
	_, _, body := a.currentScreen()
	var target *flexgo.Node
	if a.screen == ScreenHistory {
		target = body.Children[1]
	} else if a.scan.openedPath != "" {
		target = body
	} else {
		target = body.Children[1].Children[1]
	}
	tree := pageShell("", a.screen, body, "")
	var width, height int
	var replace func(*flexgo.Node)
	replace = func(n *flexgo.Node) {
		if n.View != nil {
			n.View = func(w, h int) string {
				if n == target {
					width, height = w, h
				}
				return ""
			}
		}
		for _, child := range n.Children {
			replace(child)
		}
	}
	replace(tree)
	_ = tree.Render(a.width, a.height)
	return width, height
}

func TestViewportAppLogsReverseAtBottom(t *testing.T) {
	a := viewportScanApp(80, 22)
	w, h := viewportAllocated(a)
	if w != 48 || h != 10 {
		t.Fatalf("actual log allocation %dx%d, want48x10", w, h)
	}
	for range 8 {
		a = viewportAppKey(a, tea.KeyPgDown)
	}
	if a.scan.viewports[2].Offset != 29 {
		t.Errorf("PageDown offset %d want29", a.scan.viewports[2].Offset)
	}
	a = viewportAppKey(a, tea.KeyEnd)
	if a.scan.viewports[2].Offset != 29 {
		t.Errorf("Bottom offset %d want29", a.scan.viewports[2].Offset)
	}
	before := a.scan.activePanelView()(w, h)
	a = viewportAppKey(a, tea.KeyUp)
	if a.scan.viewports[2].Offset != 28 {
		t.Errorf("reverse offset %d want28", a.scan.viewports[2].Offset)
	}
	after := a.scan.activePanelView()(w, h)
	if before == after || !strings.HasPrefix(ansi.Strip(after), "row-28") {
		t.Fatal("one Up did not move visible top row")
	}
}

func TestViewportAppResizeIntent(t *testing.T) {
	a := viewportScanApp(80, 22)
	a = viewportAppKey(a, tea.KeyEnd)
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 32})
	if a.scan.viewports[2].Offset != 19 {
		t.Errorf("grown viewport offset%d want19", a.scan.viewports[2].Offset)
	}
	a = viewportAppKey(a, tea.KeyUp)
	if a.scan.viewports[2].Offset != 18 {
		t.Errorf("grown reverse offset%d want18", a.scan.viewports[2].Offset)
	}
	a = viewportScanApp(80, 22)
	a = viewportAppKey(a, tea.KeyEnd)
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 17})
	if a.scan.viewports[2].Offset != 29 {
		t.Errorf("shrink changed legal top row to%d", a.scan.viewports[2].Offset)
	}
	a = viewportAppKey(a, tea.KeyDown)
	if a.scan.viewports[2].Offset != 30 {
		t.Errorf("shrunken Down offset%d want30", a.scan.viewports[2].Offset)
	}
	for _, size := range [][2]int{{0, 0}, {32, 22}, {80, 12}, {-1, 22}} {
		a = viewportAppMessage(a, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		for _, code := range []rune{tea.KeyDown, tea.KeyUp, tea.KeyPgDown, tea.KeyPgUp, tea.KeyHome, tea.KeyEnd} {
			a = viewportAppKey(a, code)
		}
		if a.scan.viewports[2].Offset != 30 {
			t.Errorf("unallocated %v changed logical position to%d", size, a.scan.viewports[2].Offset)
		}
	}
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 22})
	if a.scan.viewports[2].Offset != 29 {
		t.Errorf("restore offset%d want29", a.scan.viewports[2].Offset)
	}
}

func TestViewportAppGeometryAndAllScanTabs(t *testing.T) {
	for _, tab := range []int{0, 1, 2} {
		for _, rows := range []int{1, 2, 10, 50} {
			for _, width := range []int{-1, 0, 32, 33, 80} {
				a := viewportScanApp(width, rows+12)
				a.scan.activeTab = tab
				a.scan.viewports[tab].Offset = 99
				w, h := viewportAllocated(a)
				a = viewportAppKey(a, tea.KeyEnd)
				want := 99
				if w > 0 && h > 0 {
					want = max(0, a.scan.contentLen()-h)
				}
				if a.scan.viewports[tab].Offset != want {
					t.Errorf("tab%d alloc%dx%d Bottom=%d want%d", tab, w, h, a.scan.viewports[tab].Offset, want)
				}
				a = viewportAppKey(a, tea.KeyUp)
				if w > 0 && h > 0 {
					want = max(0, want-1)
				}
				if a.scan.viewports[tab].Offset != want {
					t.Errorf("tab%d alloc%dx%d Up=%d want%d", tab, w, h, a.scan.viewports[tab].Offset, want)
				}
			}
		}
	}
}

func TestViewportAppInactiveAndFileTransitions(t *testing.T) {
	a := viewportScanApp(80, 22)
	a.scan.viewports[0].Offset = 5
	a.scan.viewports[1].Offset = 6
	a = viewportAppKey(a, tea.KeyEnd)
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 32})
	if a.scan.viewports[0].Offset != 5 || a.scan.viewports[1].Offset != 6 {
		t.Fatal("resize changed inactive tab positions")
	}
	a = viewportAppKey(a, tea.KeyLeft)
	if a.scan.activeTab != 1 || a.scan.viewports[1].Offset != 0 || a.scan.viewports[2].Offset != 19 {
		t.Error("tab activation failed to reconcile only the active viewport")
	}
	a = viewportAppKey(a, tea.KeyRight)
	a = viewportAppKey(a, tea.KeyRight)
	a = viewportAppKey(a, tea.KeyEnd)
	if a.scan.filesCursor != len(a.scan.files)-1 {
		t.Fatal("Files Bottom cursor changed")
	}
	beforeCursor := a.scan.filesCursor
	a = viewportAppKey(a, tea.KeyPgUp)
	if a.scan.filesCursor != beforeCursor {
		t.Fatal("PageUp added list paging")
	}
	a = viewportAppKey(a, tea.KeyEnter)
	if a.scan.openedPath == "" {
		t.Fatal("Open failed")
	}
	a.scan.openedBody = viewportLines(39)
	a.scan.viewports[3].Offset = 99
	w, h := viewportAllocated(a)
	if w != 72 || h != 20 {
		t.Fatalf("actual opened allocation%dx%d want72x20", w, h)
	}
	a = viewportAppKey(a, tea.KeyUp)
	if a.scan.viewports[3].Offset != 18 {
		t.Errorf("opened reverse=%d want18", a.scan.viewports[3].Offset)
	}
	a = viewportAppKey(a, tea.KeyPgUp)
	if a.scan.viewports[3].Offset != 8 {
		t.Errorf("opened PageUp=%d want8", a.scan.viewports[3].Offset)
	}
	a = viewportAppKey(a, tea.KeyPgDown)
	if a.scan.viewports[3].Offset != 18 {
		t.Errorf("opened PageDown=%d want18", a.scan.viewports[3].Offset)
	}
	for _, size := range [][2]int{{0, 32}, {80, 0}} {
		a = viewportAppMessage(a, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		a = viewportAppKey(a, tea.KeyDown)
		if a.scan.viewports[3].Offset != 18 {
			t.Fatal("unallocated opened viewer moved")
		}
	}
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 22})
	a = viewportAppKey(a, tea.KeyDown)
	if a.scan.viewports[3].Offset != 19 {
		t.Fatal("restored opened viewer did not move")
	}
	a = viewportAppKey(a, tea.KeyEscape)
	if a.scan.openedPath != "" || a.scan.filesCursor != beforeCursor {
		t.Fatal("Close changed file selection")
	}
	if a.scan.viewports[3].Offset != max(0, len(a.scan.files)-10) {
		t.Fatal("Close failed to reconcile Files state")
	}
}

func TestViewportAppHistoryActivationAndResize(t *testing.T) {
	a := viewportScanApp(80, 20)
	a.history.results = []ScanResult{{ID: "history", Findings: 31}}
	a.history.detail.Offset = 99
	a = viewportAppKey(a, '3')
	w, h := viewportAllocated(a)
	if w != 44 || h != 10 {
		t.Fatalf("actual history allocation%dx%d want44x10", w, h)
	}
	if a.history.detail.Offset != 29 {
		t.Errorf("History activation offset%d want29", a.history.detail.Offset)
	}
	a = viewportAppKey(a, tea.KeyPgUp)
	if a.history.detail.Offset != 19 {
		t.Errorf("History PageUp offset%d want19", a.history.detail.Offset)
	}
	for range 5 {
		a = viewportAppKey(a, tea.KeyPgDown)
	}
	if a.history.detail.Offset != 29 {
		t.Errorf("History PageDown offset%d want29", a.history.detail.Offset)
	}
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 30})
	if a.history.detail.Offset != 19 {
		t.Errorf("History resize offset%d want19", a.history.detail.Offset)
	}
	a = viewportAppKey(a, tea.KeyPgUp)
	if a.history.detail.Offset != 9 {
		t.Errorf("History resized PageUp offset%d want9", a.history.detail.Offset)
	}
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 36, Height: 30})
	a = viewportAppKey(a, tea.KeyPgDown)
	if a.history.detail.Offset != 9 {
		t.Fatal("zero-width History moved")
	}
	a = viewportAppMessage(a, tea.WindowSizeMsg{Width: 80, Height: 60})
	if a.history.detail.Offset != 0 {
		t.Fatal("full-content History did not reconcile")
	}
	a = viewportAppKey(a, tea.KeyPgDown)
	if a.history.detail.Offset != 0 {
		t.Fatal("full-content History moved")
	}
}

func TestViewportAppHistoryAndOpenedGeometry(t *testing.T) {
	for _, history := range []bool{false, true} {
		for _, rows := range []int{1, 2, 10, 50} {
			for _, width := range []int{-1, 0, 8, 9, 36, 37, 80} {
				var a App
				if history {
					a = NewApp()
					a.history.results = []ScanResult{{ID: "history", Findings: 31}}
					a.history.detail.Offset = 99
					a = viewportAppMessage(a, tea.WindowSizeMsg{Width: width, Height: rows + 10})
					a = viewportAppKey(a, '3')
				} else {
					a = viewportScanApp(width, rows+12)
					a.scan.activeTab = 3
					a = viewportAppKey(a, tea.KeyEnter)
					a.scan.openedBody = viewportLines(39)
					a.scan.viewports[3].Offset = 99
					a = viewportAppMessage(a, tea.WindowSizeMsg{Width: width, Height: rows + 12})
				}
				w, h := viewportAllocated(a)
				want := 99
				if w > 0 && h > 0 {
					want = max(0, 39-h)
				}
				offset := a.scan.viewports[3].Offset
				if history {
					offset = a.history.detail.Offset
				}
				if offset != want {
					t.Errorf("history%t allocation%dx%d activation offset%d want%d", history, w, h, offset, want)
				}
				a = viewportAppKey(a, tea.KeyPgUp)
				if w > 0 && h > 0 {
					want = max(0, want-10)
				}
				offset = a.scan.viewports[3].Offset
				if history {
					offset = a.history.detail.Offset
				}
				if offset != want {
					t.Errorf("history%t allocation%dx%d PageUp offset%d want%d", history, w, h, offset, want)
				}
				if !history {
					for _, code := range []rune{tea.KeyHome, tea.KeyEnd} {
						a = viewportAppKey(a, code)
						if a.scan.viewports[3].Offset != want {
							t.Fatal("added a Top/Bottom binding to the opened viewer")
						}
					}
				}
			}
		}
	}
}

func TestViewportAppRemappedKeysAndCaptiveHelp(t *testing.T) {
	a := viewportScanApp(80, 22)
	a.keys.Bottom = []string{"b"}
	a.keys.Up = []string{"u"}
	a.keys.PageDown = []string{"d"}
	a = viewportAppKey(a, 'b')
	if a.scan.viewports[2].Offset != 29 {
		t.Fatal("remapped Bottom failed")
	}
	a = viewportAppKey(a, 'u')
	if a.scan.viewports[2].Offset != 28 {
		t.Fatal("remapped Up failed")
	}
	a = viewportAppKey(a, '?')
	before := a.scan
	for _, code := range []rune{'u', 'd', 'b', '3', tea.KeyTab} {
		a = viewportAppKey(a, code)
	}
	if !a.modalOpen || a.screen != ScreenScan || !reflect.DeepEqual(a.scan, before) {
		t.Fatal("help did not capture navigation")
	}
	a = viewportAppKey(a, tea.KeyEscape)
	a = viewportAppKey(a, 'd')
	if a.scan.viewports[2].Offset != 29 {
		t.Fatal("remapped PageDown failed after closing help")
	}
}

func TestViewportPaintReuseIsPure(t *testing.T) {
	lines := viewportLines(39)
	for _, offset := range []int{-10, 0, 29, 99} {
		state := ViewportState{Offset: offset}
		view := viewportView(&state, lines)
		first := view(12, 10)
		_ = view(12, 20)
		_ = view(0, 10)
		last := view(12, 10)
		if first != last || state.Offset != offset {
			t.Errorf("offset%d: repeated paint changed output or state%d", offset, state.Offset)
		}
		want := max(0, min(offset, 29))
		if !strings.HasPrefix(ansi.Strip(first), fmt.Sprintf("row-%02d", want)) {
			t.Errorf("offset%d: bad local projection", offset)
		}
	}
}

func TestViewportLegalPaintBytes(t *testing.T) {
	state := ViewportState{Offset: 1}
	const want = "\u001b[38;2;169;177;214;48;2;36;40;59mrow-01\u001b[m\u001b[38;2;169;177;214;48;2;36;40;59m\u001b[m\u001b[48;2;36;40;59m  \u001b[m\u001b[38;2;31;61;95;48;2;36;40;59m\u2588\u001b[m\n\u001b[38;2;169;177;214;48;2;36;40;59mrow-02\u001b[m\u001b[38;2;169;177;214;48;2;36;40;59m\u001b[m\u001b[48;2;36;40;59m  \u001b[m\u001b[38;2;59;66;97;48;2;36;40;59m\u2502\u001b[m"
	if got := viewportView(&state, viewportLines(4))(9, 2); got != want {
		t.Fatal("ordinary legal-position paint bytes changed")
	}
}

func TestViewportSizeMatchesActualLayout(t *testing.T) {
	for _, screen := range []Screen{ScreenScan, ScreenHistory} {
		for _, opened := range []bool{false, true} {
			if screen == ScreenHistory && opened {
				continue
			}
			for _, size := range [][2]int{{80, 22}, {80, 12}, {80, 10}, {80, 0}, {37, 20}, {36, 20}, {33, 20}, {32, 20}, {9, 20}, {8, 20}, {0, 0}, {-1, -1}} {
				a := NewApp()
				a.screen = screen
				if opened {
					a.scan.activeTab = 3
					a.scan.openedPath = "synthetic.go"
					a.scan.openedBody = viewportLines(39)
				}
				a = viewportAppMessage(a, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
				w, h := viewportAllocated(a)
				got := a.scan.viewportSize()
				if screen == ScreenHistory {
					got = a.history.viewportSize()
				}
				if got.width != w || got.height != h {
					t.Errorf("screen%d opened%t window%v: update helper%v actual leaf%dx%d", screen, opened, size, got, w, h)
				}
			}
		}
	}
}

func TestViewportNonpositiveSizesDoNotWrap(t *testing.T) {
	for _, mode := range []string{"scan", "opened", "history"} {
		for _, size := range [][2]int{{math.MinInt, 22}, {80, math.MinInt}, {math.MinInt, math.MinInt}} {
			a := viewportScanApp(80, 22)
			if mode == "opened" {
				a.scan.activeTab = 3
				a = viewportAppKey(a, tea.KeyEnter)
				a.scan.openedBody = viewportLines(39)
				a.scan.viewports[3].Offset = 17
			} else if mode == "history" {
				a.history.results = []ScanResult{{Findings: 31}}
				a = viewportAppKey(a, '3')
				a.history.detail.Offset = 17
			} else {
				a.scan.viewports[2].Offset = 17
			}
			a = viewportAppMessage(a, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
			a = viewportAppKey(a, tea.KeyPgUp)
			offset := a.scan.viewports[a.scan.activeTab].Offset
			if mode == "history" {
				offset = a.history.detail.Offset
			}
			if offset != 17 {
				t.Errorf("%s WindowSize%v changed unallocated offset to%d", mode, size, offset)
			}
		}
	}
	for _, size := range []viewportSize{{math.MinInt, 12}, {74, math.MinInt}, {math.MinInt, math.MinInt}} {
		scan := NewScanScreen(DefaultProfile())
		scan.activeTab = 2
		scan.logs = viewportLines(39)
		scan.viewports[2].Offset = 17
		scan = scan.withBodySize(size).Update(tea.KeyPressMsg{Code: tea.KeyPgUp}, DefaultKeys())
		if scan.viewports[2].Offset != 17 {
			t.Errorf("direct Scan size%v changed unallocated offset", size)
		}
		scan.openedPath = "synthetic.go"
		scan.openedBody = viewportLines(39)
		scan.activeTab = 3
		scan.viewports[3].Offset = 17
		scan = scan.withBodySize(size).Update(tea.KeyPressMsg{Code: tea.KeyPgUp}, DefaultKeys())
		if scan.viewports[3].Offset != 17 {
			t.Errorf("direct opened size%v changed unallocated offset", size)
		}
		history := HistoryScreen{results: []ScanResult{{Findings: 31}}, detail: ViewportState{Offset: 17}}
		history = history.withBodySize(size).Update(tea.KeyPressMsg{Code: tea.KeyPgUp}, DefaultKeys())
		if history.detail.Offset != 17 {
			t.Errorf("direct History size%v changed unallocated offset", size)
		}
	}
}
