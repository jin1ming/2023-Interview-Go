package algorithms

import (
	"sort"
	"strings"
)

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

	for i := 1; i < len(intervals); i++ {
		cur := intervals[i]
		// 上一个区间的结束位置 在 当前区间的开始位置的左边
		// 说明没有一点重合
		if prev[1] < cur[0] {
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
	l, r := 0, len(height)-1 // 双指针移动
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

/***** 括号生成 *****/
// 数字 n 代表生成括号的对数，请你设计一个函数，
// 用于能够生成所有可能的并且 有效的 括号组合。
func generateParenthesis(n int) []string {
	var res = &[]string{}
	// 注意切片在函数中进行append操作后，是无法返回append生成的切片的
	add(res, "", n, n)
	return *res
}

// left 和 right 代表还需要添加几个左括号和几个右括号
func add(res *[]string, str string, left int, right int) {
	if left == 0 && right == 0 {
		*res = append(*res, str)
	}
	if left == right {
		str += "("
		add(res, str, left-1, right)
		return
	}
	if left > 0 {
		add(res, str+"(", left-1, right)
		add(res, str+")", left, right-1)
		return
	}
	if right > 0 {
		add(res, str+")", left, right-1)
	}
}

/***** 旋转数组的最小数字 *****/
func minArray(numbers []int) int {
	left := 0
	right := len(numbers) - 1
	// 类似二分查找的方法去寻找
	for left < right {
		pivot := left + (right-left)/2 // 中点
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
	n := 2*numRows - 2 // 循环周期
	for i, char := range s {
		x := i % n
		// min(x, n - x) 是行号
		// 将每行的字符拼接到一块
		rows[min(x, n-x)] += string(char)
	}
	return strings.Join(rows, "")
}

/***** 买卖股票的最佳时机 *****/
// 寻找历史最低点，每天判断卖出赚多少
// TODO: 进阶题
func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	res := 0
	big := prices[0]
	small := big
	for _, k := range prices {
		if k-small > res {
			res = k - small
			continue
		}
		if k < small {
			small = k
		}
	}
	return res
}

/***** 有效的数独 *****/
// 请你判断一个 9x9 的数独是否有效。
// 数字 1-9 在每一行只能出现一次。
// 数字 1-9 在每一列只能出现一次。
// 数字 1-9 在每一个以粗实线分隔的 3x3 宫内只能出现一次。
// 数独部分空格内已填入了数字，空白格用 '.' 表示。
func isValidSudoku(board [][]byte) bool {
	// 行列检测
	for i := 0; i < 9; i++ {
		mp1 := map[byte]bool{}
		mp2 := map[byte]bool{}
		mp3 := map[byte]bool{}
		for j := 0; j < 9; j++ {
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

/***** 会议室 II *****/
// TODO: 遗忘
func minMeetingRooms(intervals [][]int) int {
	nums := make([]int, 0, 2*len(intervals))
	for _, v := range intervals {
		nums = append(nums, v[0]*10+2)
		nums = append(nums, v[1]*10+1)
	}
	sort.Ints(nums)
	maxRoom := 0
	curNeedRoom := 0
	for _, v := range nums {
		if v%10 == 1 {
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
			if bitmap[i]&bitmap[j] == 0 {
				mul := len(words[i]) * len(words[j])
				if mul > res {
					res = mul
				}
			}
		}
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
		if count, ok := bitmap[sum-k]; ok {
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
// TODO: 遗忘
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
func (nm *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	// 画图表示好理解些
	return nm.preSums[row2+1][col2+1] - nm.preSums[row1][col2+1] -
		nm.preSums[row2+1][col1] + nm.preSums[row1][col1]
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
		if ok && i-k > res {
			res = i - k
		}
		if !ok {
			offsetMap[offset] = i
		}
	}
	return res
}

/***** 优势洗牌 *****/
// TODO: 遗忘
// 给定两个大小相等的数组 A 和 B，A 相对于 B 的优势可以用满足 A[i] > B[i] 的索引 i 的数目来描述。
// 返回 A 的任意排列，使其相对于 B 的优势最大化。
func advantageCount(nums1 []int, nums2 []int) []int {
	indexs2 := make([]int, len(nums1))
	for i := range nums1 {
		indexs2[i] = i
	}
	sort.Ints(nums1)
	sort.Slice(indexs2, func(i, j int) bool {
		return nums2[indexs2[i]] < nums2[indexs2[j]]
	})

	left1, left2 := 0, 0
	res := make([]int, len(nums1))
	for left2 < len(nums2) {
		v2 := nums2[indexs2[left2]]
		ok := false
		for left1 < len(nums1) {
			if nums1[left1] > v2 {
				res[indexs2[left2]] = nums1[left1]
				nums1[left1] = -1
				ok = true
				break
			}
			left1++
		}
		if !ok {
			res[indexs2[left2]] = -1
		}
		left2++
	}

	j := 0
	for i, v := range res {
		if v != -1 {
			continue
		}
		for nums1[j] == -1 {
			j++
		}
		res[i] = nums1[j]
		j++
	}
	return res
}

/***** 任务调度器 *****/
func leastInterval(tasks []byte, n int) int {
	table := make([]int, 26)
	for _, c := range tasks {
		table[c-'A']++
	}
	sort.Slice(table, func(i, j int) bool {
		return table[i] > table[j]
	})
	cnt := 1
	for cnt < len(table) && table[cnt] == table[0] {
		cnt++
	}
	return max(len(tasks), cnt+(n+1)*(table[0]-1))
}

func candy(ratings []int) int {

}
