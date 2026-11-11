package flexgo

import "math/big"

func distribute(total int, children []*Node, isRow bool, gap int) []int {
	n := len(children)
	if n == 0 {
		return nil
	}

	available := total
	if gap > 0 && n > 1 && (total <= 0 || gap > total/(n-1)) {
		available = 0
	} else {
		available -= gap * (n - 1)
	}

	sizes := make([]int, n)

	flexIndices := make([]int, 0, n)
	flexValues := make([]int, n)

	remaining := max(0, available)
	totalFlex := 0
	maxInt := int(^uint(0) >> 1)
	exactFlex := false

	for i, c := range children {
		var fixed int
		if isRow {
			fixed = c.Width
		} else {
			fixed = c.Height
		}

		if fixed > 0 {
			sizes[i] = fixed
			if fixed >= remaining {
				remaining = 0
			} else {
				remaining -= fixed
			}
		} else {
			flex := c.Flex
			if flex == 0 {
				flex = 1
			}
			flexValues[i] = flex
			flexIndices = append(flexIndices, i)
			if !exactFlex {
				if flex > 0 && totalFlex > maxInt-flex {
					exactFlex = true
				} else {
					totalFlex += flex
				}
			}
		}
	}

	// Fixed sizes do not shrink when they exceed the available space.
	if remaining <= 0 {
		return sizes
	}
	if !exactFlex {
		for _, i := range flexIndices {
			if flexValues[i] > maxInt/remaining {
				exactFlex = true
				break
			}
		}
	}

	if exactFlex && len(flexIndices) > 0 {
		total := new(big.Int)
		for _, i := range flexIndices {
			total.Add(total, big.NewInt(int64(flexValues[i])))
		}
		for _, i := range flexIndices {
			product := new(big.Int).Mul(big.NewInt(int64(remaining)), big.NewInt(int64(flexValues[i])))
			size := int(product.Quo(product, total).Int64())
			if size > 0 {
				sizes[i] = size
			}
		}
	} else {
		for _, i := range flexIndices {
			size := (remaining * flexValues[i]) / totalFlex
			if size > 0 {
				sizes[i] = size
			}
		}
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
