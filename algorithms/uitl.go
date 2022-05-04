package algorithms

import "math"

func max(vals ...int) int {
	ans := math.MinInt64
	for _, v := range vals {
		if v > ans {
			ans = v
		}
	}
	return ans
}

func min(vals ...int) int {
	ans := math.MinInt64
	for _, v := range vals {
		if v < ans {
			ans = v
		}
	}
	return ans
}
