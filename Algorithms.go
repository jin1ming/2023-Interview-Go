package algorithms

import (
	"container/heap"
	"math"
	"sort"
	"strconv"
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

/***** 单词拆分  *****/
// 给定一个非空字符串 s 和一个包含非空单词的列表 wordDict，
// 判定 s 是否可以被空格拆分为一个或多个在字典中出现的单词。
// 说明：
// 拆分时可以重复使用字典中的单词。
// 你可以假设字典中没有重复的单词。
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
func lengthOfLongestSubstring1(s string) int {
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

func lengthOfLongestSubstring2(s string) int {
	res, left := 0, 0
	bitmap := make(map[byte]int)

	for right := 0; right < len(s); right++ {
		if c, ok := bitmap[s[right]]; ok {
			for left <= c {
				delete(bitmap, s[left])
				left++
			}
		}
		bitmap[s[right]] = right
		if right - left > res {
			res = right - left
		}
	}
	return res + 1
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
	dfs1(nums, length, 0, path, used, &res)
	return res
}

func dfs1(nums []int, length int, depth int, path []int, used []bool, res *[][]int){
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
		dfs1(nums, length, depth + 1, path, used, res)
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

/***** 数组中的逆序对 *****/
// 在数组中的两个数字，如果前面一个数字大于后面的数字，则这两个数字组成一个逆序对。
// 输入一个数组，求出这个数组中的逆序对的总数。
// 分解： 待排序的区间为 [l,r]，令 m = (l+r) / 2,
//       我们把 [l,r] 分成 [l,m] 和 [m+1,r]
// 解决： 使用归并排序递归地排序两个子序列
// 合并： 把两个已经排好序的子序列 [l,m] 和 [m+1,r] 合并起来
func reversePairs(nums []int) int {
	return mergeSort(nums, 0, len(nums)-1)
}

func mergeSort(nums []int, start, end int) int {
	if start >= end {
		return 0
	}
	mid := start + (end - start)/2 // 防止start和end相加引起的数组越界
	cnt := mergeSort(nums, start, mid) + mergeSort(nums, mid + 1, end)
	// 左右分别是排好序的数组
	// cnt 是返回的逆序对的数量
	var tmp []int
	i, j := start, mid + 1
	// i是左边数组的指针，j是右边数组的指针
	for i <= mid && j <= end { // 加判断防止越界
		if nums[i] <= nums[j] {
			tmp = append(tmp, nums[i]) // 将最小的元素放入tmp
			cnt += j - (mid + 1)
			// 当前右边数组被存入 tmp 的数量就是右边有几个元素小于左边数组的当前元素
			i++
		} else {
			tmp = append(tmp, nums[j]) // 将最小的元素放入tmp
			j++
		}
	}
	// 将左边数组剩余的加入
	for ; i <= mid; i++ {
		tmp = append(tmp, nums[i])
		cnt += end - (mid + 1) + 1
		// 右边数组全部被存入 tmp, 说明左边数组剩余元素都比右边数组中所有元素要大
	}
	// 将右边数组剩余的加入
	for ; j <= end; j++ {
		tmp = append(tmp, nums[j])
	}
	// 将排好序的 tmp 拷贝到当前数组片段中
	for i = start; i <= end; i++ {
		nums[i] = tmp[i - start]
	}
	return cnt
}

/***** 连续子数组的最大和 *****/
// 输入一个整型数组，数组中的一个或连续多个整数组成一个子数组。
// 求所有子数组的和的最大值。
func maxSubArray(nums []int) int {
	res := -101
	sum := 0
	for _, k := range nums {
		if sum < 0{
			sum = 0
		}
		sum += k
		if sum > res {
			res = sum
		}
	}
	return res
}

type TreeNode struct {
    Val int
    Left *TreeNode
    Right *TreeNode
}
/***** 从前序与中序遍历序列构造二叉树 *****/
// 给定一棵树的前序遍历 preorder 与中序遍历 inorder。
// 请构造二叉树并返回其根节点。
func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	i := 0
	for inorder[i] != preorder[0] {
		i++
	}

	left := buildTree(preorder[1:i+1], inorder[:i])
	right := buildTree(preorder[i+1:], inorder[i+1:])

	return &TreeNode{Val: preorder[0], Left: left, Right: right}
}

/***** 反转链表 II *****/
// 给你单链表的头指针 head 和两个整数 left 和 right ，其中 left <= right。
// 请你反转从位置 left 到位置 right 的链表节点，返回 反转后的链表 。
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	if head == nil || left <= 0 || left >= right{
		return head
	}
	ps := head
	m := 2
	for m < left {
		m++
		ps = ps.Next
	}

	var tmp, head2, q *ListNode
	if left == 1 {
		q = ps
		m--
	} else {
		q = ps.Next
	}
	tail := q

	for m = m - 1; m < right && q != nil; m++ {
		tmp = q.Next
		q.Next = head2
		head2 = q
		q = tmp
	}

	if tail != nil {
		tail.Next = q
	}
	if left == 1 {
		return head2
	}
	ps.Next = head2
	return head
}

/***** 编辑距离 *****/
// 给你两个单词 word1 和 word2，请你计算出将 word1 转换成 word2 所使用的最少操作数。
// 你可以对一个单词进行如下三种操作：
// 插入一个字符
// 删除一个字符
// 替换一个字符
func minDistance(word1 string, word2 string) int {
	if len(word1) * len(word2) == 0 {
		return len(word1) + len(word2)
	}
	dp := make([][]int, len(word1))
	// dp[i][j] 代表 word1 到 i 位置转换成 word2 到 j 位置需要最少步数
	var i, j int
	// 初始化边界
	// j 为 0 时，转化成到 i 的步数为 i
	for i = 0; i < len(word1)+1; i++ {
		dp[i] = make([]int, len(word2)+1)
		dp[i][0] = i
	}
	// i 为 0 时，转化成到 j 的步数为 j
	for j = 0; j < len(word2)+1; j++ {
		dp[0][j] = j
	}
	for i = 1; i < len(word1)+1; i++ {
		for j = 1; j < len(word2)+1; j++ {
			if word1[i-1] != word2[j-1] {
				// 当前字符不一致，就对 dp 值加一
				dp[i-1][j-1] += 1
			}
			dp[i][j] = tMin(dp[i-1][j-1], dp[i-1][j]+1, dp[i][j-1]+1)
			// dp[i-1][j-1] 表示替换操作，
			// dp[i-1][j] 表示删除操作，
			// dp[i][j-1] 表示插入操作。
		}
	}
	return dp[i-1][j-1]
}

func tMin(a, b, c int) int {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		return c
	} else {
		return min
	}
}

/***** 链表中倒数第k个节点 *****/
func getKthFromEnd(head *ListNode, k int) *ListNode {
	slow, fast := head, head
	for k > 0 && fast != nil {
		fast = fast.Next
		k--
	}
	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}
	return slow
}

/***** 链表中倒数第k个节点 *****/
// 给你一个由 '1'（陆地）和 '0'（水）组成的的二维网格，请你计算网格中岛屿的数量。
// 岛屿总是被水包围，并且每座岛屿只能由水平方向和/或竖直方向上相邻的陆地连接形成。
// 此外，你可以假设该网格的四条边均被水包围。
func numIslands(grid [][]byte) int {
	row := len(grid)
	if row == 0 {
		return 0
	}
	col := len(grid[0])
	num := 0

	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			if grid[r][c] == '1' {
				// 找到一个陆地
				num++
				dfs2(grid, r, c) // 让该陆地变成水
			}
		}
	}
	return num
}

func dfs2(grid [][]byte, r int, c int)  {
	row := len(grid)
	col := len(grid[0])

	grid[r][c] = '0'

	if r - 1 >= 0 && grid[r-1][c] == '1' {
		dfs2(grid, r-1, c)
	}
	if r + 1 < row && grid[r+1][c] == '1' {
		dfs2(grid, r+1, c)
	}
	if c - 1 >= 0 && grid[r][c-1] == '1' {
		dfs2(grid, r, c-1)
	}
	if c + 1 < col && grid[r][c+1] == '1' {
		dfs2(grid, r, c+1)
	}
}

/***** 螺旋矩阵 *****/
// 给你一个 m 行 n 列的矩阵 matrix ，
// 请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。
func spiralOrder(matrix [][]int) []int {
	var res []int
	top, bottom, left, right := 0, len(matrix)-1, 0, len(matrix[0])-1
	for top <= bottom && left <= right {
		// 往右走
		for i := left; i <= right; i++ {
			res = append(res, matrix[top][i])
		}
		// 上边距+1（因为再也不会走这行了）
		top++
		// 往下走
		for i := top; i <= bottom; i++ {
			res = append(res, matrix[i][right])
		}
		// 右边距-1
		right--
		// 判断是否到达终点
		// 放最后可能出现越界
		if top > bottom || right < left {
			break
		}
		// 向左走
		for i := right; i >= left; i-- {
			res = append(res, matrix[bottom][i])
		}
		// 下边距-1
		bottom--
		// 向上走
		for i := bottom; i >= top; i-- {
			res = append(res, matrix[i][left])
		}
		// 左边距+1
		left++
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


type LRUCache struct {
	size int
	capacity int
	cache map[int]*DLinkedNode
	head, tail *DLinkedNode
}

type DLinkedNode struct {
	key, value int
	prev, next *DLinkedNode
}

func initDLinkedNode(key, value int) *DLinkedNode {
	return &DLinkedNode{
		key: key,
		value: value,
	}
}

// Constructor /***** LRU 缓存 *****/
func Constructor(capacity int) LRUCache {
	l := LRUCache{
		cache: map[int]*DLinkedNode{},
		head: initDLinkedNode(0, 0),
		tail: initDLinkedNode(0, 0),
		capacity: capacity,
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

func (this *LRUCache) Get(key int) int {
	if _, ok := this.cache[key]; !ok {
		return -1
	}
	node := this.cache[key]
	this.moveToHead(node)
	return node.value
}


func (this *LRUCache) Put(key int, value int)  {
	if _, ok := this.cache[key]; !ok {
		node := initDLinkedNode(key, value)
		this.cache[key] = node
		this.addToHead(node)
		this.size++
		if this.size > this.capacity {
			removed := this.removeTail()
			delete(this.cache, removed.key)
			this.size--
		}
	} else {
		node := this.cache[key]
		node.value = value
		this.moveToHead(node)
	}
}

func (this *LRUCache) addToHead(node *DLinkedNode) {
	node.prev = this.head
	node.next = this.head.next
	this.head.next.prev = node
	this.head.next = node
}

func (this *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (this *LRUCache) moveToHead(node *DLinkedNode) {
	this.removeNode(node)
	this.addToHead(node)
}

func (this *LRUCache) removeTail() *DLinkedNode {
	node := this.tail.prev
	this.removeNode(node)
	return node
}

/***** 二叉树的层序遍历 *****/
// 算法思想：
// 借助于一个队列，先将根结点入队，然后出队，访问该结点，
// 若它有左子树，则将左子树根结点入队，若有右子树，则将右子树根节点入队。
// 然后出队，对出队结点访问，如此往复，直到队列为空。
func levelOrder(root *TreeNode) [][]int {
	var res [][]int

	if root == nil { return res	}

	queue := []*TreeNode{root}
	res = append(res, []int{root.Val})
	for len(queue) > 0 {
		var nodes []*TreeNode
		var resTemp []int
		for _, n := range queue {
			if n.Left != nil {
				nodes = append(nodes, n.Left)
				resTemp = append(resTemp, n.Left.Val)
			}
			if n.Right != nil {
				nodes = append(nodes, n.Right)
				resTemp = append(resTemp, n.Right.Val)
			}
		}
		queue = nodes
		if len(resTemp) > 0 {
			res = append(res, resTemp)
		}
	}
	return res
}

/***** 合并区间 *****/
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func merge(intervals [][]int) [][]int {
	// 按照区间开始位置进行排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var res [][]int
	prev := intervals[0]

	for i := 1; i < len(intervals); i++{
		cur := intervals[i]
		// 上一个区间的结束位置 在 当前区间的开始位置的左边
		// 说明没有一点重合
		if prev[1] < cur[0]{
			res = append(res, prev) // 直接将prev保存
			prev = cur
		} else {
			prev[1] = max(prev[1], cur[1])
			// 只需要将结束位置进行合并
		}
	}
	res = append(res, prev)
	// 别忘记将最后一个区间加入 res
	return res
}

/***** 盛最多水的容器 *****/
// 给你 n 个非负整数 a1，a2，...，an，每个数代表坐标中的一个点 (i, ai)。
// 在坐标内画 n 条垂直线，垂直线 i 的两个端点分别为 (i, ai) 和 (i, 0)。
// 找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。
func maxArea(height []int) int {
	l, r := 0, len(height) - 1 // 双指针移动
	store := 0
	for l < r {
		area := height[r]
		if height[l] < area {
			area = height[l]
		}
		area *= r - l
		if store < area {
			store = area
		}
		// 小的往前走一格
		if height[l] <= height[r] {
			l++
		} else {
			r--
		}
	}
	return store
}

/***** 盛最多水的容器 *****/
// 数字 n 代表生成括号的对数，请你设计一个函数，
// 用于能够生成所有可能的并且 有效的 括号组合。
func generateParenthesis(n int) []string {
	var res = &[]string{}
	// 注意切片在函数中进行append操作后，是无法返回append生成的切片的
	add(res, "", n, n)
	return *res
}

// left 和 right 代表还需要添加几个左括号和几个右括号
func add(res *[]string, str string, left int, right int){
	if left == 0 && right == 0 {
		*res = append(*res, str)
	}
	if left == right {
		str += "("
		add(res, str, left-1, right)
		return
	}
	if left > 0 {
		add(res, str + "(", left-1, right)
		add(res, str + ")", left, right-1)
		return
	}
	if right > 0 {
		add(res, str + ")", left, right-1)
	}
}

/***** 二叉树的右视图 *****/
// 给定一个二叉树的 根节点 root，想象自己站在它的右侧，
// 按照从顶部到底部的顺序，返回从右侧所能看到的节点值。
func rightSideView(root *TreeNode) []int {
	var res []int
	if root == nil {
		return res
	}
	var queue []*TreeNode
	p := root
	res = append(res, root.Val)
	queue = append(queue, root)
	for len(queue) > 0 {
		qLen := len(queue)
		right := math.MinInt32
		// 保存每一层最右边的值
		for i := 0; i < qLen; i++ {
			p = queue[i]
			if p.Left != nil {
				queue = append(queue, p.Left)
				right = p.Left.Val
			}
			if p.Right != nil {
				queue = append(queue, p.Right)
				right = p.Right.Val
			}
		}
		if right != math.MinInt32 {
			res = append(res, right)
		}
		newQueue := make([]*TreeNode, len(queue) - qLen)
		copy(newQueue, queue[qLen:])
		queue = newQueue
	}
	return res
}

/***** 字符串解码 *****/
// 示例：
// 输入：s = "3[a2[c]]"
// 输出："accaccacc"
// 如果当前的字符为数位，解析出一个数字（连续的多个数位）并进栈
// 如果当前的字符为字母或者左括号，直接进栈
// 如果当前的字符为右括号，开始出栈，一直到左括号出栈
func decodeString(s string) string {
	var stack []string
	ptr := 0
	for ptr < len(s) {
		cur := s[ptr]
		if cur >= '0' && cur <= '9' {
			digits := getDigits(s, &ptr)
			stack = append(stack, digits)
		} else if (cur >= 'a' && cur <= 'z' || cur >= 'A' && cur <= 'Z') || cur == '[' {
			stack = append(stack, string(cur))
			ptr++
		} else {
			ptr++
			var sub []string
			for stack[len(stack)-1] != "[" {
				sub = append(sub, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			for i := 0; i < len(sub)/2; i++ {
				sub[i], sub[len(sub)-i-1] = sub[len(sub)-i-1], sub[i]
			}
			stack = stack[:len(stack)-1]
			repTime, _ := strconv.Atoi(stack[len(stack)-1])
			stack = stack[:len(stack)-1]
			t := strings.Repeat(getString(sub), repTime)
			stack = append(stack, t)
		}
	}
	return getString(stack)
}

func getDigits(s string, ptr *int) string {
	ret := ""
	for ; s[*ptr] >= '0' && s[*ptr] <= '9'; *ptr++ {
		ret += string(s[*ptr])
	}
	return ret
}

func getString(v []string) string {
	ret := ""
	for _, s := range v {
		ret += s
	}
	return ret
}

/***** 八皇后 *****/
// 设计一种算法，打印 N 皇后在 N × N 棋盘上的各种摆法，
// 其中每个皇后都不同行、不同列，也不在对角线上。
// 这里的“对角线”指的是所有的对角线，不只是平分整个棋盘的那两条对角线。
func solveNQueens(n int) [][]string {
	var res [][]string
	matrix := make([][]bool, n)
	for s := range matrix {
		matrix[s] = make([]bool, n)
	}

	// 收集结果
	resAdd := func() {
		r := make([]string, 0, n)
		for s := range matrix {
			buf := strings.Builder{}
			for i := range matrix[s] {
				if matrix[s][i] == true {
					buf.WriteByte('Q')
				} else {
					buf.WriteByte('.')
				}
			}
			r = append(r, buf.String())
		}
		res = append(res, r)
	}
	cols := make([]bool, n)  // 记录访问过的列
	corner1 := make(map[int]bool)  // 记录该左对角线
	corner2 := make(map[int]bool)  // 记录该右对角线
	var dfs func(row int)
	dfs = func(row int) {
		if row == n {
			resAdd()
			return
		}
		for i, v := range matrix[row] {
			if v == false && cols[i] == false && corner1[i-row] == false && corner2[i+row] == false {
				if row > 0 && matrix[row][i] {
					continue
				}
				matrix[row][i] = true
				cols[i] = true
				corner1[i-row] = true
				corner2[i+row] = true
				dfs(row+1)  // 去下一行
				matrix[row][i] = false
				cols[i] = false
				delete(corner1, i-row)
				delete(corner2, i+row)
			}
		}
	}
	dfs(0)
	return res
}

/***** 二叉树的锯齿形层序遍历 *****/
func zigzagLevelOrder(root *TreeNode) [][]int {
	var queue [][]*TreeNode
	var res [][]int
	if root == nil {
		return res
	}
	queue = append(queue, []*TreeNode{root})

	level := 1
	for len(queue) != 0 {
		var values []int
		curNodes := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		var nodes []*TreeNode

		for _, n := range curNodes {
			values = append(values, n.Val)
			if n.Left != nil {
				nodes = append(nodes, n.Left)
			}
			if n.Right != nil {
				nodes = append(nodes, n.Right)
			}
		}
		if len(nodes) != 0 {
			queue = append(queue, nodes)
		}
		if level % 2 == 0{
			for i := 0; i < len(values) / 2; i++ {
				values[i], values[len(values)-1-i] = values[len(values)-1-i], values[i]
			}
		}
		res = append(res, values)
		level++
	}
	return res
}

/***** 旋转数组的最小数字 *****/
func minArray(numbers []int) int {
	left := 0
	right := len(numbers) - 1
	// 类似二分查找的方法去寻找
	for left < right {
		pivot := left + (right-left) / 2 // 中点
		if numbers[pivot] < numbers[right] {
			// 中点比 right 指向的值小
			// 说明中点往右不存在最小值
			right = pivot
		} else if numbers[pivot] > numbers[right] {
			// 中点比 right 指向的值要大
			// 说明最小值必然存在于中点和 right 的中间
			left = pivot + 1
		} else {
			// 中点和 right 指向的值相等
			right--
		}
	}
	return numbers[left]
}

/***** Z 字形变换 *****/
// 将一个给定字符串 s 根据给定的行数 numRows ，以从上往下、从左到右进行 Z 字形排列。
// 比如输入字符串为 "PAYPALISHIRING" 行数为 3 时，排列如下：
// P   A   H   N
// A P L S I I G
// Y   I   R
// 之后，你的输出需要从左往右逐行读取，产生出一个新的字符串，比如："PAHNAPLSIIGYIR"。
func convert2(s string, numRows int) string {
	if numRows == 1 {
		return s
	}
	rows := make([]string, numRows)
	n := 2 * numRows - 2 // 循环周期
	for i, char := range s {
		x := i % n
		// min(x, n - x) 是行号
		// 将每行的字符拼接到一块
		rows[min(x, n - x)] += string(char)
	}
	return strings.Join(rows, "")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/***** 买卖股票的最佳时机 *****/
func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	res := 0
	big := prices[0]
	small := big
	for _, k := range prices {
		if k - small > res {
			res = k - small
			continue
		}
		if k < small {
			small = k
		}
	}
	return res
}

/***** 二叉树中的最大路径和 *****/
func maxPathSum(root *TreeNode) int {
	if root == nil {
		return 0
	}

	maxSum := math.MinInt64
	if root.Left != nil {
		maxSum = root.Left.Val
	}
	if root.Right != nil {
		maxSum = root.Right.Val
	}

	var dfsPath func(*TreeNode) int
	dfsPath = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		leftPath := max(dfsPath(root.Left), 0)
		rightPath := max(dfsPath(root.Right), 0)

		tmpSum := root.Val + leftPath + rightPath
		if tmpSum > maxSum {
			maxSum = tmpSum
		}
		return max(leftPath, rightPath) + root.Val
	}
	dfsPath(root)
	return maxSum
}

/***** 最接近的三数之和 *****/
// 给定一个包括 n 个整数的数组 nums 和 一个目标值 target。
// 找出 nums 中的三个整数，使得它们的和与 target 最接近。
// 返回这三个数的和。假定每组输入只存在唯一答案。
func threeSumClosest(nums []int, target int) int {
	sort.Ints(nums)
	var (
		n = len(nums)
		best = math.MaxInt32
	)

	// 根据差值的绝对值来更新答案
	update := func(cur int) {
		if abs(cur - target) < abs(best - target) {
			best = cur
		}
	}

	// 枚举 a
	for i := 0; i < n; i++ {
		// 保证和上一次枚举的元素不相等
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		// 使用双指针枚举 b 和 c
		j, k := i + 1, n - 1
		for j < k {
			sum := nums[i] + nums[j] + nums[k]
			// 如果和为 target 直接返回答案
			if sum == target {
				return target
			}
			update(sum)
			if sum > target {
				// 如果和大于 target，移动 c 对应的指针
				k0 := k - 1
				// 移动到下一个不相等的元素
				for j < k0 && nums[k0] == nums[k] {
					k0--
				}
				k = k0
			} else {
				// 如果和小于 target，移动 b 对应的指针
				j0 := j + 1
				// 移动到下一个不相等的元素
				for j0 < k && nums[j0] == nums[j] {
					j0++
				}
				j = j0
			}
		}
	}
	return best
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

/***** 复原 IP 地址 *****/
// 给定一个只包含数字的字符串，用以表示一个 IP 地址，
// 返回所有可能从 s 获得的 有效 IP 地址。
// 你可以按任何顺序返回答案。
func restoreIpAddresses(s string) []string {
	const SEG_COUNT = 4
	var (
		ans []string
		segments []int
	)
	var dfs func(s string, segId, segStart int)
	dfs = func(s string, segId, segStart int) {
		// 如果找到了 4 段 IP 地址并且遍历完了字符串，那么就是一种答案
		if segId == SEG_COUNT {
			if segStart == len(s) {
				ipAddr := ""
				for i := 0; i < SEG_COUNT; i++ {
					ipAddr += strconv.Itoa(segments[i])
					if i != SEG_COUNT - 1 {
						ipAddr += "."
					}
				}
				ans = append(ans, ipAddr)
			}
			return
		}

		// 如果还没有找到 4 段 IP 地址就已经遍历完了字符串，那么提前回溯
		if segStart == len(s) {
			return
		}
		// 由于不能有前导零，如果当前数字为 0，那么这一段 IP 地址只能为 0
		if s[segStart] == '0' {
			segments[segId] = 0
			dfs(s, segId + 1, segStart + 1)
		}
		// 一般情况，枚举每一种可能性并递归
		addr := 0
		for segEnd := segStart; segEnd < len(s); segEnd++ {
			addr = addr * 10 + int(s[segEnd] - '0')
			if addr > 0 && addr <= 0xFF {
				segments[segId] = addr
				dfs(s, segId + 1, segEnd + 1)
			} else {
				break
			}
		}
	}

	segments = make([]int, SEG_COUNT)
	ans = []string{}
	dfs(s, 0, 0)
	return ans
}

/***** 只出现一次的数字 *****/
// 给定一个非空整数数组，除了某个元素只出现一次以外，
// 其余每个元素均出现两次。找出那个只出现了一次的元素。
func singleNumber(nums []int) int {
	single := 0
	for _, num := range nums {
		// 出现两次的数字会抵消为 0
		// 最后剩的就是最终的结果
		single ^= num
	}
	return single
}

/***** 快速排序 *****/
func sortArray(nums []int) []int {

	var quickSort func(left, right int)
	var findPosition func(left, right int) int

	quickSort = func(left, right int) {
		if left >= right {
			return
		}
		pos := findPosition(left, right)
		quickSort(left, pos-1)
		quickSort(pos+1, right)
	}

	findPosition = func(left, right int) int {
		temp := nums[left]

		for left < right {
			for left < right && nums[right] >= temp {
				right--
			}
			nums[left] = nums[right]
			for left < right && nums[left] <= temp {
				left++
			}
			nums[right] = nums[left]
		}

		nums[left] = temp
		return left
	}

	quickSort(0, len(nums) - 1)

	return nums
}

/***** 在排序数组中查找元素的第一个和最后一个位置 *****/
// 给定一个按照升序排列的整数数组 nums，和一个目标值 target。
// 找出给定目标值在数组中的开始位置和结束位置。
// 如果数组中不存在目标值 target，返回 [-1, -1]。
func searchRange(nums []int, target int) []int {
	left, right := 0, len(nums) - 1
	var mid int
	for left < right {
		mid = left + (right - left) / 2
		if nums[mid] < target {
			left = mid + 1 // 不加 1 可能死循环
		} else if nums[mid] > target {
			right = mid
		} else {
			break
		}
	}
	/* 注意 */
	mid = left + (right - left) / 2
	if len(nums) == 0 || mid < 0 || mid - 1 > len(nums) || nums[mid] != target {
		return []int{-1, -1}
	}
	/* --- */
	res := []int{mid, mid}
	for res[0] > 0 && nums[res[0]-1] == target {
		res[0]--
	}
	for res[1] < len(nums) - 1 && nums[res[1]+1] == target {
		res[1]++
	}
	return res
}

func searchRange2(nums []int, target int) []int {
	leftmost := sort.SearchInts(nums, target)
	if leftmost == len(nums) || nums[leftmost] != target {
		return []int{-1, -1}
	}
	rightmost := sort.SearchInts(nums, target + 1) - 1
	return []int{leftmost, rightmost}
}

/***** 有效的数独 *****/
// 请你判断一个 9x9 的数独是否有效。
// 数字 1-9 在每一行只能出现一次。
// 数字 1-9 在每一列只能出现一次。
// 数字 1-9 在每一个以粗实线分隔的 3x3 宫内只能出现一次。
// 数独部分空格内已填入了数字，空白格用 '.' 表示。
func isValidSudoku(board [][]byte) bool {
	// 行列检测
	for i:=0;i<9;i++ {
		mp1 := map[byte]bool{}
		mp2 := map[byte]bool{}
		mp3 := map[byte]bool{}
		for j:=0;j<9;j++ {
			// row
			if board[i][j] != '.' {
				if mp1[board[i][j]] {
					return false
				}
				mp1[board[i][j]] = true
			}
			// column
			if board[j][i] != '.' {
				if mp2[board[j][i]] {
					return false
				}
				mp2[board[j][i]] = true
			}
			// part
			row := (i%3)*3 + j%3
			col := (i/3)*3 + j/3
			if board[row][col] != '.' {
				if mp3[board[row][col]] {
					return false
				}
				mp3[board[row][col]] = true
			}
		}
	}
	return true
}

/***** 二叉树的最近公共祖先 *****/
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root
	}
	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	return root
}

/***** 把二叉搜索树转换为累加树 *****/
// 给出二叉 搜索 树的根节点，该树的节点值各不相同，请你将其转换为累加树，
// 使每个节点 node 的新值等于原树中大于或等于 node.val 的值之和。
// 思路:
// 访问每个节点时，时刻维护变量 sum 。
// 反向中序遍历（按值从大到小遍历）
func convertBST(root *TreeNode) *TreeNode {
	sum := 0
	var dfs func(*TreeNode)
	dfs = func(node *TreeNode) {
		if node != nil {
			dfs(node.Right)
			sum += node.Val
			node.Val = sum
			dfs(node.Left)
		}
	}
	dfs(root)
	return root
}

/***** 二叉树展开为链表 *****/
func flatten(root *TreeNode) {
	var lastNode *TreeNode
	var dfs func(*TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil { return }
		if lastNode != nil {
			lastNode.Right = node
			lastNode.Left = nil
		}
		lastNode = node
		if node.Left != nil {
			dfs(node.Left)
		}
		if node.Right != nil {
			dfs(node.Right)
		}
	}
	dfs(root)
	lastNode.Left = nil
	lastNode.Right = nil
}

/***** 汉明距离 *****/
func hammingDistance(x int, y int) int {
	xor := x ^ y
	res := 0
	for xor != 0 {
		if xor % 2 == 1 {
			res++
		}
		xor >>= 1
	}
	return res
}

/***** 找到字符串中所有字母异位词 *****/
// 滑动窗口
func findAnagrams2(s string, p string) []int {
	n, m := len(s), len(p)
	if n < m {
		return nil
	}

	var res []int
	cntS, cntP := [26]int{}, [26]int{}
	for i := 0; i < m; i++ {
		cntP[p[i]-'a']++
	}

	left, right := 0, 0
	// 右窗口开始不断向右移动
	for ; right < n; right++ {
		curRight := s[right] - 'a'
		// 将右窗口当前访问到的元素个数加1
		cntS[curRight]++
		// 当前窗口中 curRight 比 cntP 数组中对应元素的个数
		// 要多的时候就该移动左窗口指针
		for cntS[curRight] > cntP[curRight] {
			curLeft := s[left] - 'a'
			// 将左窗口当前访问到的元素个数减1
			cntS[curLeft]--
			left++
		}
		if right-left+1 == m {
			res = append(res, left)
		}
	}
	return res
}

/***** 会议室 II *****/
func minMeetingRooms(intervals [][]int) int {
	nums := make([]int, 0, 2 * len(intervals))
	for _, v := range intervals {
		nums = append(nums, v[0] * 10 + 2)
		nums = append(nums, v[1] * 10 + 1)
	}
	sort.Ints(nums)
	maxRoom := 0
	curNeedRoom := 0
	for _, v := range nums {
		if v % 10 == 1 {
			curNeedRoom--
		} else {
			curNeedRoom++
		}
		if curNeedRoom > maxRoom {
			maxRoom = curNeedRoom
		}
	}
	return maxRoom
}

/***** 最多删除一个字符得到回文 *****/
func validPalindrome(s string) bool {
	var isValid func(s string, flag bool) bool
	isValid = func(s string, flag bool) bool {
		if len(s) < 2 {
			return true
		}
		i, j := 0, len(s) - 1
		for i < j {
			if s[i] == s[j] {
				i++
				j--
			} else if !flag {
				return false
			} else {
				return isValid(s[i+1:j+1], false) ||
					isValid(s[i:j], false)
			}
		}
		return true
	}
	return isValid(s, true)
}

/***** 最小路径之和 *****/
// 一个机器人每次只能向下或者向右移动一步
func minPathSum(grid [][]int) int {
	row := len(grid)
	if row < 2 { return 0 }
	col := len(grid)
	if col < 2 { return 0 }

	for r := 1; r < row; r++ {
		grid[r][0] += grid[r-1][0]
	}
	for c := 1; c < col; c++ {
		grid[0][c] += grid[0][c-1]
	}

	getMinDis := func(r, c int) int {
		left := grid[r][c-1]
		top := grid[r-1][c]
		if left < top {
			return left
		}
		return top
	}

	for r := 1; r < row; r++ {
		for c := 1; c < col; c++ {
			grid[r][c] += getMinDis(r, c)
		}
	}
	return grid[row-1][col-1]
}

/***** 矩阵中最大的矩形 *****/
func maximalRectangle0(matrix [][]byte) int {
	row := len(matrix)
	if row == 0 {return 0}
	col := len(matrix[0])
	if col == 0 {return 0}

	dp := make([][]int, row)
	// 保存的是左边有几个连续的 1
	// 避免每次都要遍历
	for i := range dp {
		dp[i] = make([]int, col)
		dp[i][0] = int(matrix[i][0] - '0')
	}
	for r := 0; r < row; r++ {
		for c := 1; c < col; c++ {
			if matrix[r][c] == '0' {
				continue
			}
			dp[r][c] = dp[r][c-1] + 1
		}
	}

	res := 0
	// 以(i, j)为右下角，寻找左上角可能存在的最大矩形
	// 高度不断增加，随着更新宽度，判断是否需要更新最大面积
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			w := dp[i][j]
			for k := i; k >= 0; k-- {
				w = min(w, dp[k][j])
				if w == 0 {
					break
				}
				res = max(res, w * (i-k+1))
			}
		}
	}
	return res
}

/***** 两整数之和 *****/
// 不允许使用 +、-
func getSum(a, b int) int {
	for a != 0 {
		temp := a ^ b
		a = (a&b) << 1
		b = temp
	}
	return b
}

/***** 单词长度的最大乘积 *****/
func maxProduct(words []string) int {
	// rune 用26位保存26个字母
	bitmap := make([]rune, len(words))
	for k, v := range words {
		for i := 0; i < len(v); i++ {
			bitmap[k] |= 1 << (v[i] - 'a')
		}
	}

	res := 0
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			if bitmap[i] & bitmap[j] == 0 {
				mul := len(words[i]) * len(words[j])
				if mul > res {
					res = mul
				}
			}
		}
	}
	return res
}

/***** 和大于等于 target 的最短子数组 *****/
// 滑动窗口一定要 right 考虑为先！！！！！！！
func minSubArrayLen(target int, nums []int) int {
	res := math.MaxInt64
	left, sum := 0, 0
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		if sum >= target {
			for sum - nums[left] >= target {
				sum -= nums[left]
				left++
			}
			if right - left + 1 < res {
				res = right - left
			}
		}
	}
	if res == math.MaxInt64 {
		return 0
	}
	return res
}

/***** 乘积小于 K 的子数组 *****/
// 滑动窗口一定要 right 考虑为先！！！！！！！
func numSubarrayProductLessThanK(nums []int, k int) int {
	left := 0
	sum, res := 1, 0

	// 每次循环找出满足条件的最大的 left，再将 right 加 1
	// 因为 nums 中每个数都大于等于 1
	// 所以每次 right 右移后，left 向左移动时不会满足条件
	for right := 0; right < len(nums); right++ {
		sum *= nums[right]
		for left <= right && sum >= k {
			sum /= nums[left]
			left++
		}
		res += right - left + 1
		right++
	}
	return res
}

/***** 和为 k 的子数组 *****/
// 给定一个整数数组（可能为负数）和一个整数 k ，
// 请找到该数组中和为 k 的连续子数组的个数。
func subarraySum(nums []int, k int) int {
	res := 0
	// bitmap 存储从数组起始位置到当前位置出现了几次
	bitmap := make(map[int]int)
	sum := 0
	// 注意要将开始 sum 为 0 时标记，不然之后统计会缺少
	bitmap[0] = 1
	for _, v := range nums {
		sum += v
		if count, ok := bitmap[sum - k]; ok {
			// 从当前位置往前看，sum - k 出现了几次
			res += count
		}
		// 注意更新 map 放在循环最后，不然当 k == 0 时，
		// 会出现任何时候 sum - 0 == sum，错误统计
		bitmap[sum] += 1
	}
	return res
}

type NumMatrix struct {
	preSums [][]int
}

// Constructor1 /***** 二维子矩阵的和 *****/
func Constructor1(matrix [][]int) NumMatrix {
	row := len(matrix)
	// preSums[i+1][j+1] 保存(0,0)到(i,j)形成矩阵内所有元素的和
	preSums := make([][]int, row+1)
	col := len(matrix[0])
	for i := range preSums {
		preSums[i] = make([]int, col+1)
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			preSums[i+1][j+1] = preSums[i][j+1] + preSums[i+1][j] -
				preSums[i][j] + matrix[i][j]
		}
	}
	return NumMatrix{preSums: preSums}
}

// SumRegion 返回左上角 (row1, col1) 、右下角 (row2, col2) 的子矩阵的元素总和。
func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	// 画图表示好理解些
	return this.preSums[row2+1][col2+1] - this.preSums[row1][col2+1] -
		this.preSums[row2+1][col1] + this.preSums[row1][col1]
}

/***** 0 和 1 个数相同的子数组 *****/
func findMaxLength(nums []int) int {
	offsetMap := make(map[int]int)
	offsetMap[0] = -1
	res := 0
	offset := 0
	var k int
	var ok bool
	for i, v := range nums {
		if v == 0 {
			offset--
		} else {
			offset++
		}
		if offset == 0 {
			res = i + 1
			continue
		}
		k, ok = offsetMap[offset]
		if ok && i - k > res {
			res = i - k
		}
		if !ok {
			offsetMap[offset] = i
		}
	}
	return res
}
