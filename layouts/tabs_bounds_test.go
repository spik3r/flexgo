package layouts

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestTabsKeepHeadersInsideTheirRows(t *testing.T) {
	for _, title := range []string{"a title far wider than its slot", "first\nsecond\r\nthird", "界界界界界界", "\x1b[31mstyled title wider than its slot\x1b[0m"} {
		t.Run(title, func(t *testing.T) {
			panelHeight := -1
			root := Tabs(0, []Tab{{Title: title, Panel: func(w, h int) string { panelHeight = h; return strings.Repeat("P", w) }}, {Title: "Other"}})
			out := root.Render(12, 6)
			if lipgloss.Width(out) != 12 || lipgloss.Height(out) != 6 {
				t.Fatalf("header escaped 12x6: %dx%d\n%q", lipgloss.Width(out), lipgloss.Height(out), out)
			}
			if panelHeight != 4 {
				t.Fatalf("panel allocation=%d, want 4", panelHeight)
			}
			lines := strings.Split(out, "\n")
			if !strings.HasPrefix(lines[2], "PPPP") {
				t.Fatalf("panel displaced by title: %q", out)
			}
		})
	}
}

func TestTabHeaderLeavesBoundTheirOutput(t *testing.T) {
	root := Tabs(0, []Tab{{Title: "界界\nsecond"}, {Title: "B"}})
	for _, tab := range root.Children[0].Children {
		for _, leaf := range tab.Children {
			for _, size := range [][2]int{{0, 1}, {1, 0}, {-1, 1}, {1, -1}, {1, 1}, {3, 1}, {8, 1}} {
				out := leaf.View(size[0], size[1])
				if size[0] <= 0 || size[1] <= 0 {
					if out != "" {
						t.Fatalf("empty allocation %v rendered %q", size, out)
					}
					continue
				}
				if lipgloss.Width(out) > size[0] || lipgloss.Height(out) > size[1] {
					t.Fatalf("allocation %v overflowed: %q", size, out)
				}
			}
		}
	}
}

func TestTabsKeepNarrowAllocations(t *testing.T) {
	for _, width := range []int{1, 2, 3, 8, 12, 40} {
		root := Tabs(1, []Tab{{Title: "long first title"}, {Title: "界界界 second title"}, {Title: "third\ncontinued"}})
		out := root.Render(width, 6)
		if lipgloss.Width(out) > width || lipgloss.Height(out) != 6 {
			t.Fatalf("width %d escaped: %dx%d %q", width, lipgloss.Width(out), lipgloss.Height(out), out)
		}
	}
}

func TestTabsShowTruncatedTitlesAndPlainUnderlines(t *testing.T) {
	for _, tc := range []struct {
		title string
		width int
		want  string
	}{
		{"abcdefghi", 5, "abcd…"},
		{"界界界", 5, "界界…"},
		{"first\nsecond", 12, "first…"},
		{"first\rsecond", 12, "first…"},
		{"\x1b[31mabcdefghi\x1b[0m", 5, "abcd…"},
	} {
		root := Tabs(0, []Tab{{Title: tc.title}})
		leaves := root.Children[0].Children[0].Children
		got := strings.TrimSpace(ansi.Strip(leaves[0].View(tc.width, 1)))
		if got != tc.want {
			t.Fatalf("title %q: got %q want %q", tc.title, got, tc.want)
		}
		underline := strings.TrimSpace(ansi.Strip(leaves[1].View(tc.width, 1)))
		if strings.Trim(underline, "-") != "" || underline == "" {
			t.Fatalf("invalid underline: %q", underline)
		}
	}
}

func TestMultilineTabTitlesCloseTerminalFormatting(t *testing.T) {
	closeFormatting := "\x1b[m\x1b]8;;\x07"
	for _, tc := range []struct{ name, open, close string }{
		{"SGR", "\x1b[31m", "\x1b[0m"},
		{"OSC8 BEL", "\x1b]8;;https://example.invalid\x07", "\x1b]8;;\x07"},
		{"OSC8 ST", "\x1b]8;;https://example.invalid\x1b\\", "\x1b]8;;\x1b\\"},
		{"SGR and OSC8", "\x1b[1;31m\x1b]8;;https://example.invalid\x07", "\x1b[m\x1b]8;;\x07"},
	} {
		for _, separator := range []string{"\n", "\r\n", "\r"} {
			title := tc.open + "first" + separator + "second" + tc.close
			t.Run(tc.name+separator, func(t *testing.T) {
				got := tabTitle(title)
				want := tc.open + "first…" + closeFormatting
				if got != want {
					t.Fatalf("raw title = %q, want %q", got, want)
				}
				for _, width := range []int{4, 12} {
					root := Tabs(0, []Tab{{Title: title, Panel: func(w, h int) string { return "PANEL" }}, {Title: "OTHER"}})
					header := root.Children[0].Children[0].Children[0].View(width, 1)
					closing := strings.Index(header, closeFormatting)
					if closing < 0 {
						t.Fatalf("width %d lost raw closers: %q", width, header)
					}
					out := root.Render(30, 4)
					other, panel := strings.Index(out, "OTHER"), strings.Index(out, "PANEL")
					closing = strings.Index(out, closeFormatting)
					if closing < 0 || other < closing || panel < closing {
						t.Fatalf("formatting not closed before other content: %q", out)
					}
				}
			})
		}
	}
	for _, title := range []string{"plain\nsecond", "plain\rsecond"} {
		if got := tabTitle(title); got != "plain…" {
			t.Fatalf("plain title gained terminal escapes: %q", got)
		}
	}
	single := "\x1b[31mone\x1b[m"
	if got := tabTitle(single); got != single {
		t.Fatalf("single line changed: %q", got)
	}
}
