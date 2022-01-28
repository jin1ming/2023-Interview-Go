package algorithms

/***** 柱状图中的最大矩形 *****/
func largestRectangleArea(heights []int) int {
	// 首尾添加负数高度，这样原本的第一个高度能形成升序，原本的最后一个高度也能得到处理
	heights = append([]int{-2}, heights...)
	heights = append(heights, -1)
	size := len(heights)
	// 递增栈
	s := make([]int, 1, size)

	res := 0
	i := 1
	for i < len(heights) {
		// 递增则入栈
		if heights[s[len(s)-1]] < heights[i] {
			s = append(s, i)
			i++
			continue
		}
		// s[len(s)-2]是矩形的左边界
		res = max(res, heights[s[len(s)-1]]*(i-s[len(s)-2]-1))
		s = s[:len(s)-1]
	}
	return res
}

func largestRectangleArea2(heights []int) int {
	N := len(heights)
	if N == 0 {
		return 0
	}

	// 栈的简易实现
	st, pos := make([]int, N+2), 0
	push := func(v int) {
		st[pos] = v
		pos++
	}
	pop := func() int {
		pos--
		return st[pos]
	}
	top := func() int {
		return st[pos-1]
	}

	// 首尾各加一个哨兵
	get := func(i int) int {
		if i == 0 || i == N+1 {
			return 0
		}
		return heights[i-1]
	}

	// 这里才开始
	res := 0
	for i := 0; i < N+2; i++ {
		for pos > 0 && get(top()) > get(i) {
			res = max(get(pop())*(i-top()-1), res)
		}
		push(i)
	}
	return res
}

/***** 每日温度 *****/
// 请根据每日 气温 列表 temperatures ，请计算在每一天需要等几天才会有更高的温度。
// 如果气温在这之后都不会升高，请在该位置用 0 来代替。
func dailyTemperatures(T []int) []int {
	res := make([]int, len(T))
	var stack []int
	for i, v := range T {
		for len(stack) != 0 && v > T[stack[len(stack)-1]] {
			t := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[t] = i - t
		}
		stack = append(stack, i)
	}
	return res
}
