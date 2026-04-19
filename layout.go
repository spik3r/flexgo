package flexgo

func distribute(total int, children []*Node, isRow bool, gap int) []int {
	n := len(children)
	if n == 0 {
		return nil
	}

	totalGap := gap * (n - 1)
	available := total - totalGap

	sizes := make([]int, n)

	flexIndices := make([]int, 0, n)
	flexValues := make([]int, n)

	remaining := available
	totalFlex := 0

	for i, c := range children {
		var fixed int
		if isRow {
			fixed = c.Width
		} else {
			fixed = c.Height
		}

		if fixed > 0 {
			sizes[i] = fixed
			remaining -= fixed
		} else {
			flex := c.Flex
			if flex == 0 {
				flex = 1
			}
			flexValues[i] = flex
			flexIndices = append(flexIndices, i)
			totalFlex += flex
		}
	}

	for i := range children {
		if sizes[i] > 0 {
			continue
		}

		flex := flexValues[i]
		size := (remaining * flex) / totalFlex
		if size < 0 {
			size = 0
		}
		sizes[i] = size
	}

	used := 0
	for _, s := range sizes {
		used += s
	}

	rem := available - used

	for rem > 0 && len(flexIndices) > 0 {
		for _, idx := range flexIndices {
			if rem <= 0 {
				break
			}
			sizes[idx]++
			rem--
		}
	}

	return sizes
}
