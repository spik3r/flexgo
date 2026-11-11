package flexgo

import (
	"math/big"
	"reflect"
	"testing"
)

// exactFlexSizes is a test oracle for positive flex-only children. It uses
// unbounded integers, then applies distribute's leading-child remainder order.
func exactFlexSizes(total int, weights ...int) []int {
	sum := new(big.Int)
	for _, weight := range weights {
		sum.Add(sum, big.NewInt(int64(weight)))
	}
	sizes := make([]int, len(weights))
	used := 0
	for i, weight := range weights {
		product := new(big.Int).Mul(big.NewInt(int64(total)), big.NewInt(int64(weight)))
		sizes[i] = int(product.Div(product, sum).Int64())
		used += sizes[i]
	}
	for rem, i := total-used, 0; rem > 0; rem, i = rem-1, (i+1)%len(sizes) {
		sizes[i]++
	}
	return sizes
}

func TestDistributeLargePositiveFlexWeights(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	cases := []struct {
		name    string
		total   int
		weights []int
		want    []int
	}{
		{"sum overflow", 20, []int{maxInt, 1}, []int{20, 0}},
		{"product overflow", 21, []int{maxInt / 4, maxInt / 8}, []int{15, 6}},
		{"scaled ratio", 19, []int{maxInt / 3, maxInt / 6, maxInt / 6}, []int{10, 5, 4}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			children := make([]*Node, len(tc.weights))
			for i, weight := range tc.weights {
				children[i] = &Node{Flex: weight}
			}
			want := exactFlexSizes(tc.total, tc.weights...)
			if !reflect.DeepEqual(want, tc.want) {
				t.Fatalf("exact oracle = %v, hand-checked %v", want, tc.want)
			}
			for _, isRow := range []bool{true, false} {
				if got := distribute(tc.total, children, isRow, 0); !reflect.DeepEqual(got, want) {
					t.Errorf("isRow=%t: got %v, want %v", isRow, got, want)
				}
			}
		})
	}
}

func TestRenderLargePositiveFlexWeights(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, isRow := range []bool{true, false} {
		weights := []int{maxInt / 4, maxInt / 8}
		got := make([]int, len(weights))
		children := make([]*Node, len(weights))
		for i, weight := range weights {
			index := i
			children[i] = &Node{Flex: weight, View: func(w, h int) string {
				if isRow {
					got[index] = w
				} else {
					got[index] = h
				}
				return "x"
			}}
		}
		root := &Node{Dir: Col, Children: children}
		if isRow {
			root.Dir = Row
		}
		if isRow {
			root.Render(21, 1)
		} else {
			root.Render(1, 21)
		}
		if want := exactFlexSizes(21, weights...); !reflect.DeepEqual(got, want) {
			t.Errorf("isRow=%t: View allocations = %v, want %v", isRow, got, want)
		}
	}
}

func TestDistributeLargeFlexWithFixedAndGap(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, isRow := range []bool{true, false} {
		fixed := &Node{Width: 6, Height: 6}
		children := []*Node{fixed, {Flex: maxInt}, {Flex: 1}}
		// 30 - two 2-cell gaps - 6 fixed cells leaves 20 for MaxInt:1.
		want := []int{6, 20, 0}
		if got := distribute(30, children, isRow, 2); !reflect.DeepEqual(got, want) {
			t.Errorf("isRow=%t: got %v, want %v", isRow, got, want)
		}
	}
}

func TestDistributeFlexOverflowKeepsExistingRules(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	cases := []struct {
		name     string
		children []*Node
		want     []int
	}{
		{"ordinary remainder", []*Node{{Flex: 1}, {Flex: 2}}, []int{7, 13}},
		{"zero uses unit weight", []*Node{{Flex: 0}, {Flex: 2}}, []int{7, 13}},
		{"fixed width wins", []*Node{{Width: 6, Flex: maxInt}, {Flex: 1}}, []int{6, 14}},
		{"fixed overflow does not shrink", []*Node{{Width: 15}, {Width: 15}, {Flex: maxInt}}, []int{15, 15, 0}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := distribute(20, tc.children, true, 0); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDistributeBudgetArithmeticOverflow(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	cases := []struct {
		name     string
		total    int
		gap      int
		children []*Node
		want     []int
	}{
		{"fixed sum overflow", 20, 0, []*Node{{Width: maxInt, Height: maxInt}, {Width: maxInt, Height: maxInt}, {Flex: 1}}, []int{maxInt, maxInt, 0}},
		{"fixed sum overflow with gap", 20, 2, []*Node{{Width: maxInt, Height: maxInt}, {Width: maxInt, Height: maxInt}, {Flex: 1}}, []int{maxInt, maxInt, 0}},
		{"gap product overflow", 20, maxInt, []*Node{{Flex: 1}, {Flex: 1}, {Flex: 1}}, []int{0, 0, 0}},
		{"gap exceeds budget", 20, 11, []*Node{{Flex: 1}, {Flex: 1}, {Flex: 1}}, []int{0, 0, 0}},
		{"zero budget with gap", 0, maxInt, []*Node{{Flex: 1}, {Flex: 1}}, []int{0, 0}},
		{"negative budget keeps fixed", -1, 0, []*Node{{Width: maxInt, Height: maxInt}, {Width: maxInt, Height: maxInt}, {Flex: 1}}, []int{maxInt, maxInt, 0}},
		{"fixed and gap leave one cell", maxInt, 1, []*Node{{Width: maxInt - 2, Height: maxInt - 2}, {Flex: 1}}, []int{maxInt - 2, 1}},
		{"large gap leaves one cell for exact weights", maxInt, maxInt / 2, []*Node{{Flex: maxInt}, {Flex: 1}, {Flex: 1}}, []int{1, 0, 0}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, isRow := range []bool{true, false} {
				if got := distribute(tc.total, tc.children, isRow, tc.gap); !reflect.DeepEqual(got, tc.want) {
					t.Errorf("isRow=%t: got %v, want %v", isRow, got, tc.want)
				}
			}
		})
	}
}

func TestResolveMainAxisBudgetArithmeticOverflow(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	cases := []struct {
		name     string
		total    int
		gap      int
		children []*Node
		want     []int
	}{
		{"gap product overflow with minimum", 20, maxInt, []*Node{{Flex: 1, MinWidth: 1, MinHeight: 1}, {Flex: 1}, {Flex: 1}}, []int{1, 0, 0}},
		{"minimum sum overflow", 20, 0, []*Node{{Flex: 1, MinWidth: maxInt, MinHeight: maxInt}, {Flex: 1, MinWidth: maxInt, MinHeight: maxInt}, {Flex: 1}}, []int{maxInt, maxInt, 0}},
		{"exact budget with minimum", 20, 2, []*Node{{Flex: 1, MinWidth: 8, MinHeight: 8}, {Flex: 1}, {Flex: 1}}, []int{8, 4, 4}},
		{"over constrained minima", 20, 2, []*Node{{Flex: 1, MinWidth: 12, MinHeight: 12}, {Flex: 1, MinWidth: 12, MinHeight: 12}, {Flex: 1}}, []int{12, 12, 0}},
		{"mixed fixed and minimum", 20, 2, []*Node{{Width: 6, Height: 6}, {Flex: 1, MinWidth: 8, MinHeight: 8}, {Flex: 1}}, []int{6, 8, 2}},
		{"mixed fixed minimum and maximum", 30, 2, []*Node{{Width: 6, Height: 6}, {Flex: 1, MinWidth: 8, MinHeight: 8}, {Flex: 1, MaxWidth: 3, MaxHeight: 3}, {Flex: 1}}, []int{6, 8, 3, 7}},
		{"maximum releases space", 20, 2, []*Node{{Flex: 1, MaxWidth: 2, MaxHeight: 2}, {Flex: 1}, {Flex: 1}}, []int{2, 7, 7}},
		{"large exact boundary", maxInt, 1, []*Node{{Width: maxInt - 2, Height: maxInt - 2}, {Flex: 1, MinWidth: 1, MinHeight: 1}}, []int{maxInt - 2, 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, isRow := range []bool{true, false} {
				if got := resolveMainAxisSizes(tc.total, tc.children, isRow, tc.gap); !reflect.DeepEqual(got, tc.want) {
					t.Errorf("isRow=%t: got %v, want %v", isRow, got, tc.want)
				}
			}
		})
	}
}

func TestDistributeFixedChildren(t *testing.T) {
	children := []*Node{
		{Width: 20},
		{Width: 20},
		{Width: 20},
	}

	sizes := distribute(80, children, true, 0)

	if sizes[0] != 20 || sizes[1] != 20 || sizes[2] != 20 {
		t.Errorf("Expected [20, 20, 20], got %v", sizes)
	}
}

func TestDistributeFlexChildren(t *testing.T) {
	children := []*Node{
		{Flex: 1},
		{Flex: 1},
		{Flex: 1},
	}

	sizes := distribute(90, children, true, 0)

	if sizes[0] != 30 || sizes[1] != 30 || sizes[2] != 30 {
		t.Errorf("Expected [30, 30, 30], got %v", sizes)
	}
}

func TestDistributeMixedChildren(t *testing.T) {
	children := []*Node{
		{Width: 20},
		{Flex: 1},
		{Width: 20},
	}

	sizes := distribute(80, children, true, 0)

	if sizes[0] != 20 {
		t.Errorf("Child 0: expected 20, got %d", sizes[0])
	}
	if sizes[2] != 20 {
		t.Errorf("Child 2: expected 20, got %d", sizes[2])
	}
}

func TestDistributeWithGap(t *testing.T) {
	children := []*Node{
		{Flex: 1},
		{Flex: 1},
		{Flex: 1},
	}

	sizes := distribute(86, children, true, 3)

	if sizes[0] != 27 || sizes[1] != 27 || sizes[2] != 26 {
		t.Errorf("Expected [27, 27, 26], got %v", sizes)
	}
}

func TestDistributeFixedChildrenWithGap(t *testing.T) {
	children := []*Node{
		{Width: 20},
		{Width: 20},
		{Width: 20},
	}

	sizes := distribute(80, children, true, 3)

	if sizes[0] != 20 || sizes[1] != 20 || sizes[2] != 20 {
		t.Errorf("Fixed children should not get extra space, expected [20, 20, 20], got %v", sizes)
	}
}

func TestDistributeRemainderOnlyToFlexChildren(t *testing.T) {
	children := []*Node{
		{Width: 20},
		{Flex: 1},
		{Flex: 1},
	}

	sizes := distribute(86, children, true, 0)

	if sizes[0] != 20 {
		t.Errorf("Child 0 (fixed): expected 20, got %d", sizes[0])
	}
	if sizes[1] != 33 {
		t.Errorf("Child 1 (flex): expected 33, got %d", sizes[1])
	}
	if sizes[2] != 33 {
		t.Errorf("Child 2 (flex): expected 33, got %d", sizes[2])
	}
}

func TestDistributeEmptyChildren(t *testing.T) {
	children := []*Node{}
	sizes := distribute(80, children, true, 0)
	if sizes != nil {
		t.Errorf("Expected nil for empty children, got %v", sizes)
	}
}

func TestDistributeSingleChild(t *testing.T) {
	children := []*Node{{Flex: 1}}
	sizes := distribute(80, children, true, 0)
	if sizes[0] != 80 {
		t.Errorf("Expected [80], got %v", sizes)
	}
}

func TestDistributeSingleChildWithFlex(t *testing.T) {
	children := []*Node{
		{Width: 20},
		{Flex: 1},
	}
	sizes := distribute(80, children, true, 0)
	if sizes[0] != 20 {
		t.Errorf("Child 0 (fixed): expected 20, got %d", sizes[0])
	}
	if sizes[1] != 60 {
		t.Errorf("Child 1 (flex): expected 60, got %d", sizes[1])
	}
}

func TestDistributeZeroTotal(t *testing.T) {
	children := []*Node{
		{Flex: 1},
		{Flex: 1},
	}
	sizes := distribute(0, children, true, 0)
	if sizes[0] != 0 || sizes[1] != 0 {
		t.Errorf("Expected [0, 0], got %v", sizes)
	}
}

func TestDistributeWithMinWidth(t *testing.T) {
	children := []*Node{
		{Flex: 1, MinWidth: 30},
		{Flex: 1},
	}
	sizes := distribute(80, children, true, 0)
	if sizes[0] != 40 {
		t.Errorf("Child 0: expected 40 from distribute, got %d (min/max applied in renderChildren)", sizes[0])
	}
}

func TestDistributeWithMaxWidth(t *testing.T) {
	children := []*Node{
		{Flex: 1},
		{Flex: 1, MaxWidth: 30},
	}
	sizes := distribute(100, children, true, 0)
	if sizes[1] != 50 {
		t.Errorf("Child 1: expected 50 from distribute, got %d (min/max applied in renderChildren)", sizes[1])
	}
}

func TestDistributeVertical(t *testing.T) {
	children := []*Node{
		{Height: 5},
		{Height: 10},
	}
	sizes := distribute(20, children, false, 0)
	if sizes[0] != 5 || sizes[1] != 10 {
		t.Errorf("Expected [5, 10], got %v", sizes)
	}
}

func TestDistributeFlexWithZeroFlex(t *testing.T) {
	children := []*Node{
		{Flex: 1},
		{Width: 20},
	}
	sizes := distribute(60, children, true, 0)
	if sizes[0] != 40 {
		t.Errorf("Child 0 (flex): expected 40, got %d", sizes[0])
	}
	if sizes[1] != 20 {
		t.Errorf("Child 1 (fixed): expected 20, got %d", sizes[1])
	}
}

// B1 — min/max redistribution. resolveMainAxisSizes must detect clamped
// children and re-distribute the remaining space among the unclamped ones
// so the sum still equals the available total.

func TestResolveMainAxisMinBumpsRedistribute(t *testing.T) {
	children := []*Node{
		{Flex: 1, MinWidth: 50},
		{Flex: 1},
	}
	sizes := resolveMainAxisSizes(60, children, true, 0)
	if sum(sizes) != 60 {
		t.Errorf("sizes must sum to 60, got %v (sum=%d)", sizes, sum(sizes))
	}
	if sizes[0] != 50 {
		t.Errorf("clamped child: expected 50, got %d", sizes[0])
	}
	if sizes[1] != 10 {
		t.Errorf("unclamped child should absorb the delta: expected 10, got %d", sizes[1])
	}
}

func TestResolveMainAxisMaxClipRedistribute(t *testing.T) {
	children := []*Node{
		{Flex: 1, MaxWidth: 20},
		{Flex: 1},
	}
	sizes := resolveMainAxisSizes(100, children, true, 0)
	if sum(sizes) != 100 {
		t.Errorf("sizes must sum to 100, got %v (sum=%d)", sizes, sum(sizes))
	}
	if sizes[0] != 20 {
		t.Errorf("clipped child: expected 20, got %d", sizes[0])
	}
	if sizes[1] != 80 {
		t.Errorf("unclipped child should absorb spare: expected 80, got %d", sizes[1])
	}
}

func TestResolveMainAxisMinAndMaxTogether(t *testing.T) {
	// min on child 0 bumps up, max on child 2 clips down, child 1 absorbs.
	children := []*Node{
		{Flex: 1, MinWidth: 40},
		{Flex: 1},
		{Flex: 1, MaxWidth: 10},
	}
	sizes := resolveMainAxisSizes(60, children, true, 0)
	if sum(sizes) != 60 {
		t.Errorf("sizes must sum to 60, got %v (sum=%d)", sizes, sum(sizes))
	}
	if sizes[0] != 40 {
		t.Errorf("min child: expected 40, got %d", sizes[0])
	}
	if sizes[2] != 10 {
		t.Errorf("max child: expected 10, got %d", sizes[2])
	}
	if sizes[1] != 10 {
		t.Errorf("unconstrained child: expected 10, got %d", sizes[1])
	}
}

func TestResolveMainAxisChainedClamp(t *testing.T) {
	// Child 0 has MinWidth 50 (bumped up); redistributing the rest of
	// 60 - 50 = 10 across two Flex:1 children would give each 5, but
	// child 1's MinWidth 8 now fires — triggering another redistribution
	// pass. Final: [50, 8, 2].
	children := []*Node{
		{Flex: 1, MinWidth: 50},
		{Flex: 1, MinWidth: 8},
		{Flex: 1},
	}
	sizes := resolveMainAxisSizes(60, children, true, 0)
	if sum(sizes) != 60 {
		t.Errorf("sizes must sum to 60, got %v (sum=%d)", sizes, sum(sizes))
	}
	if sizes[0] != 50 || sizes[1] != 8 || sizes[2] != 2 {
		t.Errorf("expected [50, 8, 2], got %v", sizes)
	}
}

func TestResolveMainAxisColDirectionMinHeight(t *testing.T) {
	children := []*Node{
		{Flex: 1, MinHeight: 15},
		{Flex: 1},
	}
	sizes := resolveMainAxisSizes(20, children, false, 0)
	if sum(sizes) != 20 {
		t.Errorf("sizes must sum to 20, got %v (sum=%d)", sizes, sum(sizes))
	}
	if sizes[0] != 15 || sizes[1] != 5 {
		t.Errorf("expected [15, 5], got %v", sizes)
	}
}

func TestResolveMainAxisWithGap(t *testing.T) {
	// Total 60, gap 2 between 3 children → 56 for content.
	// Child 0 clamped to 30; remaining 26 split between two flex children → [13, 13].
	children := []*Node{
		{Flex: 1, MinWidth: 30},
		{Flex: 1},
		{Flex: 1},
	}
	sizes := resolveMainAxisSizes(60, children, true, 2)
	// Gap totals 4 (2 * 2 interior gaps). Content should sum to 56.
	if sum(sizes) != 56 {
		t.Errorf("sizes must sum to 56 (content only), got %v (sum=%d)", sizes, sum(sizes))
	}
	if sizes[0] != 30 {
		t.Errorf("min child: expected 30, got %d", sizes[0])
	}
}

func sum(xs []int) int {
	n := 0
	for _, x := range xs {
		n += x
	}
	return n
}
