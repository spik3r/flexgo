package flexgo

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestValidateHappy(t *testing.T) {
	root := VBox(
		&Node{Height: 3, View: func(w, h int) string { return "" }},
		&Node{Flex: 1, View: func(w, h int) string { return "" }},
	)
	if err := Validate(root); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateDetectsConflicts(t *testing.T) {
	cases := []struct {
		name  string
		node  *Node
		want  string
	}{
		{"view+children", &Node{View: func(int, int) string { return "" }, Children: []*Node{{}}}, "View and Children"},
		{"padding+paddings", &Node{Padding: 1, Paddings: Spacing{Top: 1}}, "Padding and Paddings"},
		{"margin+margins", &Node{Margin: 1, Margins: Spacing{Top: 1}}, "Margin and Margins"},
		{"negative flex", &Node{Flex: -1}, "negative Flex"},
		{"negative width", &Node{Width: -1}, "negative Width"},
		{"min>max width", &Node{MinWidth: 10, MaxWidth: 5}, "MinWidth"},
		{"min>max height", &Node{MinHeight: 10, MaxHeight: 5}, "MinHeight"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.node)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q in error, got %q", tc.want, err.Error())
			}
		})
	}
}

func TestValidateRecursesIntoChildren(t *testing.T) {
	root := VBox(
		&Node{Name: "bad-child", Padding: 1, Paddings: Spacing{Top: 1}},
	)
	err := Validate(root)
	if err == nil || !strings.Contains(err.Error(), "bad-child") {
		t.Fatalf("expected error naming child, got %v", err)
	}
}

func TestInspectShowsStructure(t *testing.T) {
	root := &Node{
		Dir:  Col,
		Name: "root",
		Gap:  1,
		Children: []*Node{
			{Height: 3, Name: "header", View: func(int, int) string { return "" }},
			{Flex: 1, Name: "body", ShowBorder: true, Border: lipgloss.NormalBorder()},
		},
	}
	out := Inspect(root)
	for _, want := range []string{"VBox", `"root"`, "gap=1", "Leaf", `"header"`, "h=3", `"body"`, "flex=1", "border"} {
		if !strings.Contains(out, want) {
			t.Errorf("Inspect output missing %q\n--- output ---\n%s", want, out)
		}
	}
}

func TestInspectNilRoot(t *testing.T) {
	if got := Inspect(nil); got != "<nil>\n" {
		t.Errorf("Inspect(nil) = %q, want <nil>\\n", got)
	}
}
