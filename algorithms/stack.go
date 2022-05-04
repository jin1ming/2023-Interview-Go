package algorithms

import (
	"math"
	"strings"
)

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

/***** 接雨水 *****/
func trap2(height []int) int {
	if len(height) == 0 {
		return 0
	}
	res := 0
	var stack []int
	for r, rightH := range height {
		// 注意stack存储的是下标，需要比较高度的时候需要访问height数组
		for len(stack) > 0 && rightH > height[stack[len(stack)-1]] {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) > 0 { // 存在比cur处高度高的墙
				left := stack[len(stack)-1]
				leftH := height[left]
				// r - left - 1 是当前低洼区两侧墙之间的距离
				// 低洼区的高度计算为：两侧墙的最小海拔高度-低洼的海拔高度
				res += (min(leftH, rightH) - height[cur]) * (r - left - 1)
			}
		}
		stack = append(stack, r)
	}
	return res
}

/***** 移掉 K 位数字 *****/
// 给你一个以字符串表示的非负整数 num 和一个整数 k ，
// 移除这个数中的 k 位数字，使得剩下的数字最小。
// 请你以字符串形式返回这个最小的数字。
func removeKdigits(num string, k int) string {
	var stack []byte
	if len(num) == 0 {
		return ""
	}
	remain := len(num) - k
	for _, c := range append([]byte(num), '0') {
		for k > 0 && len(stack) > 0 && c < stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, c)
	}
	res := strings.TrimLeft(string(stack[:remain]), "0")
	if len(res) == 0 {
		return "0"
	}
	return res
}

/***** 去除重复字母 *****/
// 给你一个字符串 s ，请你去除字符串中重复的字母，
// 使得每个字母只出现一次。需保证 返回结果的字典序最小
// （要求不能打乱其他字符的相对位置）。
func removeDuplicateLetters(s string) string {
	if len(s) == 0 {
		return ""
	}
	left := make([]byte, 26)  // 记录每个字母剩余出现的次数
	exist := make([]bool, 26) // 记录每个字母是否在栈中出现
	for _, v := range []byte(s) {
		left[v-'a']++
	}
	var stack []byte
	for _, v := range []byte(s) {
		if !exist[v-'a'] {
			for len(stack) > 0 && v < stack[len(stack)-1] {
				chi := stack[len(stack)-1] - 'a'
				if left[chi] == 0 {
					// 该字母后续没出现，无法舍弃
					break
				}
				exist[chi] = false
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, v)
			exist[v-'a'] = true
		} // 对于已经出现的字母，直接舍弃
		left[v-'a']--
	}
	return string(stack)
}

type MinStack struct {
	stack    []int
	minStack []int
}

/***** 最小栈 *****/
func ConstructorS() MinStack {
	return MinStack{
		stack:    []int{},
		minStack: []int{math.MaxInt64},
	}
}

func (this *MinStack) Push(x int) {
	this.stack = append(this.stack, x)
	top := this.minStack[len(this.minStack)-1]
	this.minStack = append(this.minStack, min(x, top))
}

func (this *MinStack) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
	this.minStack = this.minStack[:len(this.minStack)-1]
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1]
}

func (this *MinStack) GetMin() int {
	return this.minStack[len(this.minStack)-1]
}
