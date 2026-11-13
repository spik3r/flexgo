package flexgo

import (
	"strings"
	"testing"
)

func TestDiagnosticsHandleCycles(t *testing.T) {
	t.Run("self", func(t *testing.T) {
		root := &Node{Name: "loop", Dir: Row}
		root.Children = []*Node{root, {Width: -1}}
		err := Validate(root)
		if err == nil || !strings.Contains(err.Error(), "root.loop: cycle to root") || !strings.Contains(err.Error(), "root.Children[1]: negative Width") {
			t.Fatalf("missing cycle or sibling diagnostic: %v", err)
		}
		want := "HBox \"loop\"\n├─ <cycle>\n└─ HBox\n"
		if got := Inspect(root); got != want {
			t.Fatalf("inspect = %q, want %q", got, want)
		}
	})
	t.Run("ancestor", func(t *testing.T) {
		root := &Node{Dir: Col}
		child := &Node{Name: "child", Children: []*Node{root}}
		root.Children = []*Node{child}
		if err := Validate(root); err == nil || err.Error() != "root.child.Children[0]: cycle to root" {
			t.Fatalf("got %v", err)
		}
		want := "VBox\n└─ HBox \"child\"\n   └─ <cycle>\n"
		for range 2 {
			if got := Inspect(root); got != want {
				t.Fatalf("got %q", got)
			}
		}
	})
}

func TestDiagnosticsAllowSharedAcyclicSubtrees(t *testing.T) {
	shared := &Node{Name: "shared", Dir: Col, Children: []*Node{{Name: "leaf"}}}
	root := HBox(shared, shared)
	if err := Validate(root); err != nil {
		t.Fatal(err)
	}
	got := Inspect(root)
	if strings.Contains(got, "<cycle>") || strings.Count(got, `"shared"`) != 2 || strings.Count(got, `"leaf"`) != 2 {
		t.Fatalf("shared subtree omitted: %q", got)
	}
	// The same invalid shared node is diagnosed at each occurrence.
	shared.Width = -1
	err := Validate(root)
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok || len(joined.Unwrap()) != 2 {
		t.Fatalf("want two occurrence diagnostics: %v", err)
	}
}
