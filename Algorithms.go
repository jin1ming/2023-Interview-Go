package algorithms

import (
	"container/heap"
	"math"
	"sort"
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

/***** 两数之和 *****/
func twoSum(nums []int, target int) []int {
	left, right := 0, len(nums) - 1
	indexs := make([]int, len(nums))
	for i, _ := range indexs {
		indexs[i] = i
	}
	sort.Slice(indexs, func(i, j int) bool {
		return nums[indexs[i]] < nums[indexs[j]]
	})
	for left < right {
		sum := nums[indexs[left]] + nums[indexs[right]]
		if sum == target {
			return []int{indexs[left], indexs[right]}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	return nil
}

func twoSum2(nums []int, target int) []int {
	hashTable := map[int]int{}
	// map存储，挨着找
	for i, x := range nums {
		if p, ok := hashTable[target-x]; ok {
			return []int{p, i}
		}
		hashTable[x] = i
	}
	return nil
}

/***** 无重复字符的最长子串 *****/
// 类型: 滑动窗口
func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}
	bitmap := make(map[uint8]int)
	// bitmap存储字符最后出现的位置，用于判断是否重复
	maxLen := 1
	dp := make([]int, len(s))
	// dp记录的是当前窗口的大小
	dp[0] = 1
	bitmap[s[0]] = 0
	for i := 1; i < len(s); i++ {
		if v, ok := bitmap[s[i]]; ok && v >= i - dp[i-1] {
			// 该字符在当前窗口曾经出现过
			dp[i] = i - v
		} else {
			dp[i] = dp[i-1] + 1
			if dp[i] > maxLen {
				maxLen = dp[i]
			}
		}
		bitmap[s[i]] = i
	}
	return maxLen
}

/***** 最长回文子串 *****/
// 中心扩展算法
func longestPalindrome(s string) string {
	if s == "" {
		return ""
	}
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		left1, right1 := expandAroundCenter(s, i, i)
		left2, right2 := expandAroundCenter(s, i, i + 1)
		if right1 - left1 > end - start {
			start, end = left1, right1
		}
		if right2 - left2 > end - start {
			start, end = left2, right2
		}
	}
	return s[start:end+1]
}

func expandAroundCenter(s string, left, right int) (int, int) {
	for ; left >= 0 && right < len(s) && s[left] == s[right]; left, right = left-1 , right+1 { }
	return left + 1, right - 1
}

type ListNode struct {
    Val int
    Next *ListNode
}
/***** 反转链表 *****/
func reverseList(head *ListNode) *ListNode {
	var root *ListNode
	var tmp *ListNode

	for head != nil {
		tmp = head.Next
		head.Next = root
		root = head
		head = tmp
	}

	return root
}

/***** 全排列 *****/
func permute(nums []int) [][]int {
	length := len(nums)
	var res [][]int
	if length == 0 {
		return res
	}

	var path []int
	used := make([]bool, length)
	dfs(nums, length, 0, path, used, &res)
	return res
}

func dfs(nums []int, length int, depth int, path []int, used []bool, res *[][]int){
	if  depth == length {
		p := make([]int, length)
		copy(p, path)  // 注意使用copy
		*res = append(*res, p)
		return
	}

	for i := 0; i < length; i++ {
		if used[i]{
			continue
		}
		path = append(path, nums[i])
		used[i] = true
		dfs(nums, length, depth + 1, path, used, res)
		path = path[:len(path) - 1]
		used[i] = false
	}
}

/***** 跳跃游戏 II *****/
// 给定一个非负整数数组，你最初位于数组的第一个位置。
// 数组中的每个元素代表你在该位置可以跳跃的最大长度。
// 你的目标是使用最少的跳跃次数到达数组的最后一个位置。
// 假设你总是可以到达数组的最后一个位置。
func jump(nums []int) int {
	length := len(nums)
	end := 0
	maxPosition := 0
	steps := 0
	for i := 0; i < length - 1; i++ {
		maxPosition = max(maxPosition, i + nums[i])
		// 当前可到达最远位置
		if i == end {
			// 已经到达可走的最远位置
			end = maxPosition
			steps++
		}
	}
	return steps
}

/***** 三数之和 *****/
// 给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，
// 使得 a + b + c = 0 ？请你找出所有和为 0 且不重复的三元组。
func threeSum(nums []int) [][]int {
	if len(nums) < 3 {
		return [][]int{}
	}

	sort.Ints(nums)
	var res [][]int

	var ptrLeft, ptrRight int
	for k, _ := range nums {
		switch {
		case nums[k] > 0:
			return res
		case k > 0 && nums[k-1] == nums[k]:
			continue
		default:
			ptrLeft = k + 1
			ptrRight = len(nums) - 1
			for ptrLeft < ptrRight {
				sum := nums[k] + nums[ptrLeft] + nums[ptrRight]
				if sum == 0 {
					r := []int {nums[k], nums[ptrLeft], nums[ptrRight]}
					res = append(res, r)
					for ptrLeft < ptrRight && nums[ptrLeft] == nums[ptrLeft + 1]{
						ptrLeft += 1
					}
					for ptrLeft < ptrRight && nums[ptrRight] == nums[ptrRight - 1]{
						ptrRight -= 1
					}
				}

				if sum > 0 {
					ptrRight -= 1
				} else {
					ptrLeft += 1
				}
			}
		}
	}
	return res
}

/***** 爬楼梯 *****/
//每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？
func climbStairs(n int) int {
	switch n {
	case 0,1:
		return 1
	default:
		tmp := []int{1, 1}
		res := 0
		for i := 2; i <= n; i++ {
			res = tmp[0] + tmp[1]
			tmp[0], tmp[1] = tmp[1], res
		}
		return res
	}
}

/***** 重排链表 *****/
//给定一个单链表 L：L0→L1→…→Ln-1→Ln ，
//将其重新排列后变为： L0→Ln→L1→Ln-1→L2→Ln-2→…
//你不能只是单纯的改变节点内部的值，而是需要实际的进行节点交换。
func reorderList(head *ListNode)  {
	if head == nil {
		return
	}
	length := 0

	p := head

	for p != nil {
		length += 1
		p = p.Next
	}
	p = head
	for i := 1; i < (length+1) / 2; i++ {
		p = p.Next
	}

	head2 := p.Next
	p.Next = nil
	// 将后半段逆转
	head2 = reverse(head2)

	p = head
	// 依次插入
	var tmp *ListNode
	for p != nil && head2 != nil {
		tmp = p.Next
		p.Next = head2
		head2 = head2.Next
		p = p.Next
		p.Next = tmp
		p = p.Next
	}
}

func reverse(head *ListNode) *ListNode{
	var tmp *ListNode
	var root *ListNode
	for head != nil {
		tmp = head.Next
		head.Next = root
		root = head
		head = tmp
	}
	return root
}

/***** 反转每对括号间的子串 *****/
func reverseParentheses(s string) string {
	// 字符串无法直接修改，转换为byte slice
	brr := []byte(s)
	var stack []int
	for i := 0; i < len(brr); i ++ {
		if brr[i] == '(' {
			// 遇到左括号，加入栈中
			stack = append(stack, i)
		} else if brr[i] == ')'{
			// 题目保证括号左右匹配，所以不用检验stack中是否有左括号
			lastIdx := stack[len(stack)-1]
			// 反转左括号位置+1到右括号位置-1之间的字符
			for lj, rj := lastIdx + 1, i - 1; lj < rj; lj, rj = lj +1, rj -1 {
				brr[lj], brr[rj] = brr[rj], brr[lj]
			}
			// 已匹配的左括号退栈
			stack = stack[:len(stack)-1]
		}
	}

	// 去掉所有括号字符
	sb := strings.Builder{}
	for i := 0; i < len(brr); i ++ {
		if brr[i] != '(' && brr[i] !=')' {
			sb.WriteByte(brr[i])
		}
	}

	return sb.String()
}

/***** 接雨水 *****/
func trap(height []int) int {
	maxHeight := 0
	for _, k := range height {
		if k > maxHeight {
			maxHeight = k
		}
	}

	capacity := 0

	for k := 1; k <= maxHeight; k++ {
		pre := -1
		for i := 0; i < len(height); i++ {
			if height[i] < k {
				continue
			}
			if pre == -1 {
				pre = i
			} else {
				capacity += i - pre -1
				pre = i
			}
		}
	}
	return capacity
}

/***** K 个一组翻转链表 *****/
//给你一个链表，每 k 个节点一组进行翻转，请你返回翻转后的链表。
//k 是一个正整数，它的值小于或等于链表的长度。
//如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
func reverseKGroup(head *ListNode, k int) *ListNode {
	hair := &ListNode{Next: head}
	pre := hair

	for head != nil {
		tail := pre
		for i := 0; i < k; i++ {
			tail = tail.Next
			if tail == nil {
				return hair.Next
			}
		}
		nex := tail.Next
		head, tail = myReverse(head, tail)
		pre.Next = head
		tail.Next = nex
		pre = tail
		head = tail.Next
	}
	return hair.Next
}

func myReverse(head, tail *ListNode) (*ListNode, *ListNode) {
	prev := tail.Next
	p := head
	for prev != tail {
		nex := p.Next
		p.Next = prev
		prev = p
		p = nex
	}
	return tail, head
}