package flexgo

func distribute(total int, children []*Node, isRow bool, gap int) []int {
	n := len(children)
	if n == 0 {
		return nil
	}

	totalGap := gap * (n - 1)
	available := total - totalGap

	sizes := make([]int, n)

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
			if c.Flex == 0 {
				c.Flex = 1
			}
			totalFlex += c.Flex
		}
	}

	for i, c := range children {
		if sizes[i] > 0 {
			continue
		}

		size := (remaining * c.Flex) / totalFlex
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
	for i := 0; rem > 0; i++ {
		sizes[i%len(sizes)]++
		rem--
	}

	return sizes
}
