package flexgo

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestClipViewText(t *testing.T) {
	for _, tc := range []struct {
		name, text, want string
		w, h             int
	}{
		{"short", "short", "short", 8, 3},
		{"empty", "", "", 8, 3},
		{"blank lines", "a\n\nb\n", "a\n\nb\n", 8, 4},
		{"crop both axes", "abcdef\nghijkl\nlast", "abc\nghi", 3, 2},
		{"line endings", "a\r\nb\rc\nd", "a\nb\nc", 5, 3},
		{"tabs", "A\tB", "A  ", 3, 1},
		{"wide", "界界x", "界", 3, 1},
		{"wide cannot fit", "界x", "", 1, 1},
		{"combining", "e\u0301xy", "e\u0301x", 2, 1},
		{"emoji", "👩‍💻x", "👩‍💻", 2, 1},
		{"emoji cannot fit", "👩‍💻x", "", 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			view := ClipView(func(int, int) string { return tc.text })
			if got := view(tc.w, tc.h); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestClipViewCallbackContract(t *testing.T) {
	calls := 0
	view := ClipView(func(w, h int) string {
		calls++
		if w != 3 || h != 2 {
			t.Fatalf("callback allocation %dx%d", w, h)
		}
		return "abcdef\nsecond\nthird"
	})
	for _, size := range [][2]int{{0, 3}, {3, 0}, {-1, 3}, {3, -1}} {
		if got := view(size[0], size[1]); got != "" {
			t.Fatalf("nonpositive allocation: %q", got)
		}
	}
	if calls != 0 {
		t.Fatal("called view with no usable area")
	}
	if got := view(3, 2); got != "abc\nsec" || calls != 1 {
		t.Fatalf("got %q in %d calls", got, calls)
	}
	if got := ClipView(nil)(3, 2); got != "" {
		t.Fatalf("nil view: %q", got)
	}
}

func TestClipViewANSI(t *testing.T) {
	red, blue, reset := "\x1b[31m", "\x1b[34m", ansi.ResetStyle
	link, closeLink := ansi.SetHyperlink("https://example.test", "id=synthetic"), ansi.ResetHyperlink()
	for _, tc := range []struct {
		name, text, want string
		w, h             int
	}{
		{"balanced", red + "abc" + reset, red + "ab" + reset, 2, 1},
		{"discarded closures", red + link + "first\nsecond" + closeLink + reset, red + link + "fir" + reset + closeLink, 3, 1},
		{"style across rows", red + "ab\ncd" + reset, red + "ab" + reset + "\n" + red + "cd" + reset, 2, 2},
		{"link across rows", link + "ab\ncd" + closeLink, link + "ab" + closeLink + "\n" + link + "cd" + closeLink, 2, 2},
		{"style beyond horizontal cut", red + "abcd" + blue + "\nef" + reset, red + "ab" + blue + reset + "\n" + red + blue + "ef" + reset, 2, 2},
		{"unclosed style", red + "x", red + "x" + reset, 2, 1},
		{"closed style stays closed", red + "a" + reset + "\nb", red + "a" + reset + "\nb", 2, 2},
		{"closed link stays closed", link + "a" + closeLink + "\nb", link + "a" + closeLink + "\nb", 2, 2},
		{"ST link across rows", "\x1b]8;;https://example.test\x1b\\a\nb\x1b]8;;\x1b\\", "\x1b]8;;https://example.test\x1b\\a" + closeLink + "\n\x1b]8;;https://example.test\x1b\\b\x1b]8;;\x1b\\", 2, 2},
		{"partial reset", red + "\x1b[1mab\n\x1b[22mcd", red + "\x1b[1mab" + reset + "\n" + red + "\x1b[1m\x1b[22mcd" + reset, 2, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ClipView(func(int, int) string { return tc.text })(tc.w, tc.h)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestClipViewHorizontalSiblingIsolation(t *testing.T) {
	styled := "\x1b[31m" + ansi.SetHyperlink("https://example.test") + "left\nnext\nlast" + ansi.ResetStyle + ansi.ResetHyperlink()
	root := &Node{Dir: Row, Children: []*Node{
		{Width: 3, View: ClipView(func(int, int) string { return styled })},
		{Width: 5, View: ClipView(func(int, int) string { return "RIGHT\nAGAIN" })},
	}}
	out := root.Render(8, 2)
	if got := ansi.Strip(out); got != "lefRIGHT\nnexAGAIN" {
		t.Fatalf("sibling displaced: %q", out)
	}
	for i, line := range strings.Split(out, "\n") {
		marker := []string{"RIGHT", "AGAIN"}[i]
		prefix := line[:strings.Index(line, marker)]
		if !strings.HasSuffix(prefix, ansi.ResetStyle+ansi.ResetHyperlink()) {
			t.Fatalf("sibling inherits callback state: %q", line)
		}
	}
}

func TestClipViewNestedCompositionAndResize(t *testing.T) {
	text := "A\t界e\u0301👩‍💻very long line\r\nsecond\nthird\nfourth\nfifth\nsixth\nseventh"
	view := ClipView(func(int, int) string { return text })
	root := &Node{Dir: Col, Children: []*Node{
		{Height: 1, View: ClipView(func(int, int) string { return "HEADER\nexcess" })},
		{Dir: Row, Gap: 1, Children: []*Node{
			{Flex: 1, Padding: 1, ShowBorder: true, Border: lipgloss.NormalBorder(), View: view},
			NewNode().Flex(1).View(view).Build(),
		}},
		{Height: 1, View: ClipView(func(int, int) string { return "FOOTER\nexcess" })},
	}}
	for _, size := range [][2]int{{24, 9}, {40, 14}, {13, 8}, {24, 9}} {
		out := root.Render(size[0], size[1])
		if w, h := lipgloss.Width(out), lipgloss.Height(out); w != size[0] || h != size[1] {
			t.Fatalf("%v rendered %dx%d: %q", size, w, h, out)
		}
		lines := strings.Split(ansi.Strip(out), "\n")
		if !strings.HasPrefix(lines[0], "HEADER") || !strings.HasPrefix(lines[len(lines)-1], "FOOTER") {
			t.Fatalf("header/footer displaced: %q", out)
		}
	}
	plain := &Node{View: func(int, int) string { return "long\nsecond" }}
	if lipgloss.Height(plain.Render(2, 1)) <= 1 {
		t.Fatal("unwrapped callback unexpectedly clipped")
	}
}

func TestClipViewConcurrentAllocations(t *testing.T) {
	view := ClipView(func(w, h int) string { return fmt.Sprintf("%d/%d\nmore", w, h) })
	var wg sync.WaitGroup
	for w := 1; w <= 12; w++ {
		for h := 1; h <= 4; h++ {
			wg.Go(func() {
				got := view(w, h)
				if lipgloss.Width(got) > w || lipgloss.Height(got) > h {
					t.Errorf("%dx%d: %q", w, h, got)
				}
			})
		}
	}
	wg.Wait()
}

func ExampleClipView() {
	view := ClipView(func(w, h int) string {
		return "first line\nsecond line\nthird line"
	})
	fmt.Println(view(5, 2))
	// Output:
	// first
	// secon
}
