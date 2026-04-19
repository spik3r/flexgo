package flexgo

import "testing"

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
