package algorithms

import (
	"container/heap"
	"math"
	"strings"
)

/***** 字符串转换整数 (atoi) *****/

func myAtoi(str string) int {
	return convert(clean(str))
}

func clean(s string) (sign int, abs string) {
	// 先去除首尾空格
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	// 判断第一个字符
	switch s[0] {
	// 有效的
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		sign, abs = 1, s
	// 有效的，正号
	case '+':
		sign, abs = 1, s[1:]
	// 有效的，负号
	case '-':
		sign, abs = -1, s[1:]
	// 无效的，当空字符处理，并且直接返回
	default:
		abs = ""
		return
	}
	for i, b := range abs {
		// 遍历第一波处理过的字符，如果直到第i个位置有效，那就取s[:i]，
		// 从头到这个有效的字符，剩下的就不管了，也就是break掉
		// 比如 s=123abc，那么就取123，也就是s[:3]
		if b < '0' || '9' < b {
			abs = abs[:i]
			// 一定要break，因为后面的就没用了
			break
		}
	}
	return
}

// 接收的输入是已经处理过的纯数字
func convert(sign int, absStr string) int {
	absNum := 0
	for _, b := range absStr {
		// b - '0' ==> 得到这个字符类型的数字的真实数值的绝对值
		absNum = absNum*10 + int(b-'0')
		// 检查溢出
		switch {
		case sign == 1 && absNum > math.MaxInt32:
			return math.MaxInt32
		// 这里和正数不一样的是，必须和负号相乘，也就是变成负数，否则永远走不到里面
		case sign == -1 && absNum*sign < math.MinInt32:
			return math.MinInt32
		}
	}
	return sign * absNum
}

/***** Pow(x, n) *****/

func myPow(x float64, n int) float64 {
	if n >= 0 {
		return quickMul(x, n)
	}
	return 1.0 / quickMul(x, -n)
}

func quickMul(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	y := quickMul(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return y * y * x
}

/***** 柱状图中的最大矩形 *****/

func largestRectangleArea(heights []int) int {
	// 首尾添加负数高度，这样原本的第一个高度能形成升序，原本的最后一个高度也能得到处理
	heights = append([]int{-2}, heights...)
	heights = append(heights, -1)
	size:=len(heights)
	// 递增栈
	s:=make([]int,1,size)

	res:=0
	i:=1
	for i < len(heights) {
		// 递增则入栈
		if heights[s[len(s)-1]]<heights[i]{
			s=append(s,i)
			i++
			continue
		}
		// s[len(s)-2]是矩形的左边界
		res=max(res, heights[s[len(s)-1]]*(i-s[len(s)-2]-1))
		s=s[:len(s)-1]
	}
	return res
}
func max(a,b int)int{
	if a>b{return a}
	return b
}

/***** 单词拆分
 * 给定一个非空字符串 s 和一个包含非空单词的列表 wordDict，
 * 判定 s 是否可以被空格拆分为一个或多个在字典中出现的单词。
 * 说明：
 * 拆分时可以重复使用字典中的单词。
 * 你可以假设字典中没有重复的单词。
 */

func wordBreak(s string, wordDict []string) bool {
	wordDictSet := make(map[string]bool)
	for _, w := range wordDict {
		wordDictSet[w] = true
	}
	dp := make([]bool, len(s) + 1)
	dp[0] = true
	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			if dp[j] && wordDictSet[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	return dp[len(s)]
}


/***** 前K个高频元素 *****/

func topKFrequent(nums []int, k int) []int {
	occurrences := map[int]int{}
	for _, num := range nums {
		occurrences[num]++
	}
	h := &IHeap{}
	heap.Init(h)
	for key, value := range occurrences {
		heap.Push(h, [2]int{key, value})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	ret := make([]int, k)
	for i := 0; i < k; i++ {
		ret[k - i - 1] = heap.Pop(h).([2]int)[0]
	}
	return ret
}

type IHeap [][2]int

func (h IHeap) Len() int           { return len(h) }
func (h IHeap) Less(i, j int) bool { return h[i][1] < h[j][1] }
func (h IHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IHeap) Push(x interface{}) {
	*h = append(*h, x.([2]int))
}

func (h *IHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
