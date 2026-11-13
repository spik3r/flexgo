package layouts

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/spik3r/flexgo"
)

func TestModalTitleKeepsBodyInItsRow(t *testing.T) {
	for _, title := range []string{"first\nsecond\nthird", "a title much longer than the available dialog width", "界界界界界界界界界界"} {
		root := Modal(title, func(w, h int) string { return "BODY" }, 12, 6, nil)
		out := root.Render(14, 8)
		if lipgloss.Width(out) != 14 || lipgloss.Height(out) != 8 {
			t.Errorf("title %q escaped 14x8: %dx%d %q", title, lipgloss.Width(out), lipgloss.Height(out), out)
		}
		if !strings.Contains(strings.Split(out, "\n")[2], "BODY") {
			t.Errorf("body displaced: %q", out)
		}
	}
}
func TestFormLabelsKeepRowsAndFieldsAligned(t *testing.T) {
	root := Form(5, []FormRow{{Label: "first\nsecond\nthird", Field: func(w, h int) string { return "FIELD" }}, {Label: "Next", Field: func(w, h int) string { return "NEXT" }}})
	out := root.Render(14, 2)
	if lipgloss.Width(out) != 14 || lipgloss.Height(out) != 2 {
		t.Fatalf("form escaped 14x2: %dx%d %q", lipgloss.Width(out), lipgloss.Height(out), out)
	}
	lines := strings.Split(out, "\n")
	if !strings.Contains(lines[0], "FIELD") || !strings.Contains(lines[1], "NEXT") {
		t.Fatalf("fields displaced: %q", out)
	}
}

func TestRecipeLabelsRespectCallbackAllocations(t *testing.T) {
	for _, text := range []string{"long label that must not wrap", "界界界界", "e\u0301e\u0301e\u0301e\u0301", "first\r\nsecond\rthird", "\x1b[31mfirst\nsecond\x1b[m"} {
		modal := Modal(text, nil, 12, 6, nil).Children[0]
		label := Form(8, []FormRow{{Label: text, Height: 2}}).Children[0].Children[0]
		for _, leaf := range []*flexgo.Node{modal, label} {
			for _, size := range [][2]int{{-1, 1}, {0, 2}, {5, 0}, {5, -1}, {1, 1}, {3, 1}, {8, 2}, {12, 4}} {
				got := leaf.View(size[0], size[1])
				if size[0] <= 0 || size[1] <= 0 {
					if got != "" {
						t.Fatalf("empty allocation %v rendered %q", size, got)
					}
					continue
				}
				if lipgloss.Width(got) > size[0] || lipgloss.Height(got) > size[1] {
					t.Fatalf("allocation %v overflowed: %q", size, got)
				}
			}
		}
	}
}

func TestLabelTruncationPreservesVisibleLinesAndGraphemes(t *testing.T) {
	for _, tc := range []struct {
		text string
		w, h int
		want string
	}{
		{"abcdef", 4, 1, "abc…"}, {"界界界", 5, 1, "界界…"}, {"e\u0301e\u0301e\u0301", 2, 1, "e\u0301…"},
		{"first\nsecond\nthird", 8, 2, "first\nsecond…"}, {"one\r\ntwo", 8, 2, "one\ntwo"},
		{"first\rsecond", 8, 1, "first…"}, {"unchanged", 12, 1, "unchanged"},
	} {
		if got := boundedLabel(tc.text, tc.w, tc.h); got != tc.want {
			t.Errorf("%q: got %q, want %q", tc.text, got, tc.want)
		}
	}
	root := Form(8, []FormRow{{Label: "one\ntwo", Height: 2, Field: func(w, h int) string { return "A\nB" }}, {Label: "Next", Field: func(w, h int) string { return "NEXT" }}})
	out := root.Render(14, 3)
	lines := strings.Split(out, "\n")
	if len(lines) != 3 || !strings.Contains(lines[0], "one") || !strings.Contains(lines[1], "two:") || !strings.Contains(lines[2], "NEXT") {
		t.Fatalf("explicit two-row label changed: %q", out)
	}
}

func TestClippedLabelsCloseFormattingBeforeOtherContent(t *testing.T) {
	closeFormatting := ansi.ResetStyle + ansi.ResetHyperlink()
	for _, escapes := range [][2]string{
		{"\x1b[31m", "\x1b[m"},
		{"\x1b]8;;https://example.invalid\x07", "\x1b]8;;\x07"},
		{"\x1b]8;;https://example.invalid\x1b\\", "\x1b]8;;\x1b\\"},
		{"\x1b[31m\x1b]8;;https://example.invalid\x07", "\x1b[m\x1b]8;;\x07"},
	} {
		for _, sep := range []string{"\n", "\r\n", "\r"} {
			text := escapes[0] + "first" + sep + "discarded" + escapes[1]
			for _, width := range []int{3, 12} {
				got := boundedLabel(text, width, 1)
				if !strings.HasSuffix(got, closeFormatting) || !strings.HasPrefix(got, escapes[0]) {
					t.Fatalf("lost raw formatting boundaries: %q", got)
				}
			}
			views := []string{
				Form(8, []FormRow{{Label: text, Field: func(w, h int) string { return "FIELD" }}}).Render(16, 1),
				Modal(text, func(w, h int) string { return "BODY" }, 12, 6, nil).Render(14, 8),
			}
			for i, out := range views {
				marker := []string{"FIELD", "BODY"}[i]
				closing, next := strings.Index(out, closeFormatting), strings.Index(out, marker)
				if closing < 0 || next < closing {
					t.Fatalf("formatting not closed before %s: %q", marker, out)
				}
			}
		}
	}
}
