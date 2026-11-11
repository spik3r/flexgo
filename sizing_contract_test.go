package flexgo_test

import (
	"reflect"
	"testing"

	"github.com/spik3r/flexgo"
)

func TestRenderSizingContract(t *testing.T) {
	cases := []struct {
		name     string
		total    int
		children []*flexgo.Node
		want     []int
	}{
		{"zero flex uses unit weight", 20, []*flexgo.Node{{}, {Flex: 2}}, []int{7, 13}},
		{"fixed size wins over flex", 20, []*flexgo.Node{{Width: 6, Height: 6, Flex: 100}, {Flex: 1}}, []int{6, 14}},
		{"maximum adjusts fixed size", 10, []*flexgo.Node{{Width: 8, Height: 8, MaxWidth: 3, MaxHeight: 3, Flex: 100}, {Flex: 1}}, []int{3, 7}},
		{"overconstrained sizes do not shrink", 10, []*flexgo.Node{{Width: 8, Height: 8}, {Flex: 1, MinWidth: 6, MinHeight: 6}, {Flex: 1}}, []int{8, 6, 0}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, isRow := range []bool{true, false} {
				got := make([]int, len(tc.children))
				for i, child := range tc.children {
					index := i
					child.View = func(w, h int) string {
						if isRow {
							got[index] = w
						} else {
							got[index] = h
						}
						return ""
					}
				}
				root := &flexgo.Node{Dir: flexgo.Row, Children: tc.children}
				if isRow {
					root.Render(tc.total, 1)
				} else {
					root.Dir = flexgo.Col
					root.Render(1, tc.total)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("isRow=%t: View dimensions = %v, want %v", isRow, got, tc.want)
				}
			}
		})
	}
}

func TestRootFlexDoesNotOverrideRenderSize(t *testing.T) {
	var width, height int
	root := &flexgo.Node{Flex: 100, View: func(w, h int) string {
		width, height = w, h
		return ""
	}}
	root.Render(9, 2)
	if width != 9 || height != 2 {
		t.Fatalf("root View received (%d, %d), want (9, 2)", width, height)
	}
}
