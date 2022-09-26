package algorithms

// 上一次学习：2022.4.7，看到了133

import "math"

/***** 单词拆分 *****/
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
	dp := make([]bool, len(s)+1)
	dp[0] = true
	for right := 1; right <= len(s); right++ {
		for left := 0; left < right; left++ {
			if dp[left] && wordDictSet[s[left:right]] {
				dp[right] = true
				break
			}
		}
	}
	return dp[len(s)]
}

/***** 爬楼梯 *****/
//每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？
func climbStairs(n int) int {
	switch n {
	case 0, 1:
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

/***** 最长递增子序列 *****/
func lengthOfLIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	dp := make([]int, len(nums))
	result := 0
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[j] >= nums[i] {
				continue
			}
			dp[i] = max(dp[i], dp[j]+1)
		}
	}
	for _, v := range dp {
		if v > result {
			result = v
		}
	}
	return result + 1
}

/***** 编辑距离 *****/
// 给你两个单词 word1 和 word2，请你计算出将 word1 转换成 word2 所使用的最少操作数。
// 你可以对一个单词进行如下三种操作：
// 插入一个字符
// 删除一个字符
// 替换一个字符
func minDistance(word1 string, word2 string) int {
	if len(word1)*len(word2) == 0 {
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

/***** 最小路径之和 *****/
// 给定一个包含非负整数的 m x n 网格 grid ，
// 请找出一条从左上角到右下角的路径，使得路径上的数字总和为最小。
// 一个机器人每次只能向下或者向右移动一步
func minPathSum(grid [][]int) int {
	row := len(grid)
	if row < 2 {
		return 0
	}
	col := len(grid)
	if col < 2 {
		return 0
	}

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
func maximalRectangle(matrix [][]byte) int {
	row := len(matrix)
	if row == 0 {
		return 0
	}
	col := len(matrix[0])
	if col == 0 {
		return 0
	}

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
			// 从下往上扫描
			for k := i; k >= 0; k-- {
				w = min(w, dp[k][j])
				if w == 0 {
					break
				}
				res = max(res, w*(i-k+1))
			}
		}
	}
	return res
}

/***** 乘积最大子数组 *****/
// 给你一个整数数组 nums ，请你找出数组中乘积最大的非空连续子数组
// （该子数组中至少包含一个数字），并返回该子数组所对应的乘积。
// 测试用例的答案是一个 32-位 整数。
// 子数组 是数组的连续子序列。
func maxProduct2(nums []int) int {
	preMax, preMin, ans := 1, 1, math.MinInt32
	for _, num := range nums {
		preMax, preMin = max(preMax*num, preMin*num, num), min(preMax*num, preMin*num, num)
		ans = max(preMax, ans)
	}
	return ans
}

/***** 最长有效括号 *****/
// 给你一个只包含 '(' 和 ')' 的字符串
// 找出最长有效（格式正确且连续）括号子串的长度。
func longestValidParentheses(s string) int {
	maxAns := 0
	dp := make([]int, len(s))
	for i := 1; i < len(s); i++ {
		if s[i] == ')' {
			if s[i-1] == '(' {
				// 找同级关系的一串字串，然后合并
				if i >= 2 {
					dp[i] = dp[i-2] + 2
				} else {
					dp[i] = 2
				}
			} else if i-dp[i-1] > 0 && s[i-dp[i-1]-1] == '(' {
				// 找下一级的子串，注意这里的子串已经合并过
				// i-dp[i-1]是子串left，子串的left必须是'('
				if i-dp[i-1] >= 2 {
					// 子串旁边可能有别的子串
					// dp[i-dp[i-1]-2]代表着子串左侧的字串
					dp[i] = dp[i-1] + dp[i-dp[i-1]-2] + 2
				} else {
					dp[i] = dp[i-1] + 2
				}
			}
			maxAns = max(maxAns, dp[i])
		}
	}
	return maxAns
}

/***** 可被3整除的最大和 *****/
// 给你一个整数数组 nums，请你找出并返回能被三整除的元素最大和。
func maxSumDivThree(nums []int) int {
	rest := [3]int{}
	for _, num := range nums {
		a := rest[0] + num
		b := rest[1] + num
		c := rest[2] + num
		rest[a%3] = max(rest[a%3], a)
		rest[b%3] = max(rest[b%3], b)
		rest[c%3] = max(rest[c%3], c)
	}
	return rest[0]
}

/***** 交错字符串 *****/
// 帮忙验证 s3 是否是由 s1 和 s2 交错 组成的。
func isInterleave(s1 string, s2 string, s3 string) bool {
	n, m, t := len(s1), len(s2), len(s3)
	if (n + m) != t {
		return false
	}
	f := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]bool, m+1)
	}
	f[0][0] = true
	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			p := i + j - 1
			if i > 0 {
				f[i][j] = f[i][j] || (f[i-1][j] && s1[i-1] == s3[p])
			}
			if j > 0 {
				f[i][j] = f[i][j] || (f[i][j-1] && s2[j-1] == s3[p])
			}
		}
	}
	return f[n][m]
}

/***** 交错字符串 *****/
func rob(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0]
	}
	if n == 2 {
		return max(nums[0], nums[1])
	}
	return max(_rob(nums[:n-1]), _rob(nums[1:]))
}

func _rob(nums []int) int {
	first, second := nums[0], max(nums[0], nums[1])
	for _, v := range nums[2:] {
		first, second = second, max(first+v, second)
	}
	return second
}
