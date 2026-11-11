package flexgo

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

var renderBenchmarkOutput string

type renderFixture struct {
	name string
	root *Node
	w, h int
	want []string
}

func benchmarkView(text string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().Width(w).Height(h).Render(text)
	}
}

func renderFixtures() []renderFixture {
	return []renderFixture{
		{
			name: "leaf_80x24", w: 80, h: 24, want: []string{"status"},
			root: Leaf(benchmarkView("status")),
		},
		{
			name: "fixed_flex_row_80x24", w: 80, h: 24, want: []string{"menu", "content"},
			root: HBox(
				&Node{Width: 18, View: benchmarkView("menu")},
				&Node{Flex: 1, View: benchmarkView("content")},
			),
		},
		{
			name: "nested_100x30", w: 100, h: 30,
			want: []string{"header", "navigation", "summary", "detail", "footer"},
			root: VBox(
				&Node{Height: 2, View: benchmarkView("header")},
				&Node{Flex: 1, Children: []*Node{
					{Width: 20, View: benchmarkView("navigation")},
					{Flex: 1, Dir: Col, Children: []*Node{
						{Height: 3, View: benchmarkView("summary")},
						{Flex: 1, View: benchmarkView("detail")},
					}},
				}},
				&Node{Height: 1, View: benchmarkView("footer")},
			),
		},
		{
			name: "styled_unicode_80x24", w: 80, h: 24, want: []string{"東京"},
			root: &Node{
				Padding: 1, Foreground: lipgloss.Color("33"),
				ShowBorder: true, Border: lipgloss.RoundedBorder(),
				View: benchmarkView("東京 — status"),
			},
		},
	}
}

func checkRenderFixture(t testing.TB, fixture renderFixture) {
	t.Helper()
	if err := Validate(fixture.root); err != nil {
		t.Fatalf("invalid benchmark tree: %v", err)
	}
	out := fixture.root.Render(fixture.w, fixture.h)
	for _, want := range fixture.want {
		if !strings.Contains(out, want) {
			t.Fatalf("render missing %q", want)
		}
	}
	if width, height := lipgloss.Width(out), lipgloss.Height(out); width != fixture.w || height != fixture.h {
		t.Fatalf("render size = %dx%d, want %dx%d", width, height, fixture.w, fixture.h)
	}
}

func TestRenderBenchmarkFixtures(t *testing.T) {
	for _, fixture := range renderFixtures() {
		t.Run(fixture.name, func(t *testing.T) { checkRenderFixture(t, fixture) })
	}
}

func BenchmarkRender(b *testing.B) {
	for _, fixture := range renderFixtures() {
		b.Run(fixture.name, func(b *testing.B) {
			checkRenderFixture(b, fixture)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderBenchmarkOutput = fixture.root.Render(fixture.w, fixture.h)
			}
		})
	}
}
