package algorithms

// 上一次学习：2022.4.7，看到了133

import "math"

/***** 单词拆分 *****/
// 给定一个非空字符串 s 和一个包含非空单词的列表 wordDict，
// 判定 s 是否可以被空格拆分为一个或多个在字典中出现的单词。
// 说明：
// 拆分时可以重复使用字典中的单词。
// 你可以假设字典中没有重复的单词。
// 算法思路：使用动态规划，dp[i] 表示字符串前 i 个字符能否被拆分
// 对于每个位置 right，检查所有可能的 left，如果 dp[left] 为 true 且 s[left:right] 在字典中，则 dp[right] 为 true
func wordBreak(s string, wordDict []string) bool {
	// 将字典转换为 map，提高查找效率
	wordDictSet := make(map[string]bool)
	for _, w := range wordDict {
		wordDictSet[w] = true
	}
	// dp[i] 表示字符串前 i 个字符能否被拆分
	dp := make([]bool, len(s)+1)
	dp[0] = true // 空字符串可以被拆分
	// 遍历字符串的每个位置
	for right := 1; right <= len(s); right++ {
		// 检查所有可能的分割点
		for left := 0; left < right; left++ {
			// 如果前 left 个字符可以被拆分，且 s[left:right] 在字典中
			if dp[left] && wordDictSet[s[left:right]] {
				dp[right] = true
				break // 找到一种拆分方式即可
			}
		}
	}
	return dp[len(s)]
}

/***** 爬楼梯 *****/
// 每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？
// 算法思路：动态规划，dp[i] = dp[i-1] + dp[i-2]
// 到达第 i 阶的方法数 = 从第 i-1 阶爬 1 步 + 从第 i-2 阶爬 2 步
// 使用滚动数组优化空间复杂度为 O(1)
func climbStairs(n int) int {
	switch n {
	case 0, 1:
		return 1 // 0 阶或 1 阶只有一种方法
	default:
		// 使用两个变量保存前两个状态，优化空间
		tmp := []int{1, 1} // tmp[0] = dp[i-2], tmp[1] = dp[i-1]
		res := 0
		for i := 2; i <= n; i++ {
			// dp[i] = dp[i-1] + dp[i-2]
			res = tmp[0] + tmp[1]
			// 更新状态，为下一次迭代做准备
			tmp[0], tmp[1] = tmp[1], res
		}
		return res
	}
}

/***** 最长递增子序列 *****/
// 给定一个整数数组，找到其中最长严格递增子序列的长度
// 算法思路：动态规划，dp[i] 表示以 nums[i] 结尾的最长递增子序列的长度
// 对于每个位置 i，遍历所有 j < i，如果 nums[j] < nums[i]，则更新 dp[i]
func lengthOfLIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	// dp[i] 表示以 nums[i] 结尾的最长递增子序列的长度
	dp := make([]int, len(nums))
	result := 0
	// 遍历数组的每个位置
	for i := 1; i < len(nums); i++ {
		// 检查所有在 i 之前的元素
		for j := 0; j < i; j++ {
			// 如果 nums[j] < nums[i]，说明可以将 nums[i] 接在以 nums[j] 结尾的子序列后面
			if nums[j] >= nums[i] {
				continue
			}
			// 更新以 nums[i] 结尾的最长递增子序列长度
			dp[i] = max(dp[i], dp[j]+1)
		}
	}
	// 找到所有 dp[i] 中的最大值
	for _, v := range dp {
		if v > result {
			result = v
		}
	}
	// 返回结果 + 1，因为 dp[i] 表示的是相对长度（不包括自身），实际长度需要 +1
	return result + 1
}

/***** 编辑距离 *****/
// 给你两个单词 word1 和 word2，请你计算出将 word1 转换成 word2 所使用的最少操作数。
// 你可以对一个单词进行如下三种操作：
// 插入一个字符
// 删除一个字符
// 替换一个字符
// 算法思路：动态规划，dp[i][j] 表示 word1 前 i 个字符转换成 word2 前 j 个字符需要的最少操作数
func minDistance(word1 string, word2 string) int {
	// 如果其中一个字符串为空，编辑距离就是另一个字符串的长度
	if len(word1)*len(word2) == 0 {
		return len(word1) + len(word2)
	}
	dp := make([][]int, len(word1)+1)
	// dp[i][j] 代表 word1 前 i 个字符转换成 word2 前 j 个字符需要最少步数
	var i, j int
	// 初始化边界
	// j 为 0 时，word1 前 i 个字符转换成空字符串需要删除 i 个字符
	for i = 0; i < len(word1)+1; i++ {
		dp[i] = make([]int, len(word2)+1)
		dp[i][0] = i
	}
	// i 为 0 时，空字符串转换成 word2 前 j 个字符需要插入 j 个字符
	for j = 0; j < len(word2)+1; j++ {
		dp[0][j] = j
	}
	// 填充 dp 表
	for i = 1; i < len(word1)+1; i++ {
		for j = 1; j < len(word2)+1; j++ {
			if word1[i-1] != word2[j-1] {
				// 当前字符不一致，需要替换操作，替换成本为 1
				dp[i][j] = tMin(dp[i-1][j-1]+1, dp[i-1][j]+1, dp[i][j-1]+1)
			} else {
				// 当前字符一致，不需要替换，直接继承之前的状态
				dp[i][j] = tMin(dp[i-1][j-1], dp[i-1][j]+1, dp[i][j-1]+1)
			}
			// dp[i-1][j-1] 表示替换操作（如果字符不同）或保持不变（如果字符相同），
			// dp[i-1][j] 表示删除 word1[i-1] 操作，
			// dp[i][j-1] 表示在 word1 中插入 word2[j-1] 操作。
		}
	}
	return dp[len(word1)][len(word2)]
}

// tMin 返回三个整数中的最小值
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
// 算法思路：动态规划，直接在原数组上修改，grid[i][j] 表示从 (0,0) 到 (i,j) 的最小路径和
// 状态转移：grid[i][j] = grid[i][j] + min(grid[i-1][j], grid[i][j-1])
func minPathSum(grid [][]int) int {
	row := len(grid)
	if row == 0 {
		return 0
	}
	col := len(grid[0])
	if col == 0 {
		return 0
	}
	if row == 1 && col == 1 {
		return grid[0][0]
	}

	// 初始化第一列：只能从上往下走
	for r := 1; r < row; r++ {
		grid[r][0] += grid[r-1][0]
	}
	// 初始化第一行：只能从左往右走
	for c := 1; c < col; c++ {
		grid[0][c] += grid[0][c-1]
	}

	// 获取从左边和上边来的最小路径和
	getMinDis := func(r, c int) int {
		left := grid[r][c-1] // 从左边来的路径和
		top := grid[r-1][c]  // 从上边来的路径和
		if left < top {
			return left
		}
		return top
	}

	// 填充剩余位置
	for r := 1; r < row; r++ {
		for c := 1; c < col; c++ {
			// 当前位置的最小路径和 = 当前值 + min(左边路径和, 上边路径和)
			grid[r][c] += getMinDis(r, c)
		}
	}
	return grid[row-1][col-1]
}

/***** 矩阵中最大的矩形 *****/
// 给定一个仅包含 0 和 1 的二维二进制矩阵，找出只包含 1 的最大矩形，并返回其面积
// 算法思路：使用动态规划预处理每行从左到右连续 1 的个数，然后对每个位置向上扩展计算最大矩形
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
	// dp[i][j] 保存的是第 i 行第 j 列左边有几个连续的 1（包括自身）
	// 这样避免每次计算矩形宽度时都要遍历
	for i := range dp {
		dp[i] = make([]int, col)
		dp[i][0] = int(matrix[i][0] - '0')
	}
	// 计算每行从左到右连续 1 的个数
	for r := 0; r < row; r++ {
		for c := 1; c < col; c++ {
			if matrix[r][c] == '0' {
				continue // 如果当前是 0，则连续 1 的个数为 0（默认值）
			}
			// 如果当前是 1，则连续 1 的个数 = 左边连续 1 的个数 + 1
			dp[r][c] = dp[r][c-1] + 1
		}
	}

	res := 0
	// 以 (i, j) 为右下角，向上扫描寻找可能存在的最大矩形
	// 高度不断增加，随着更新宽度，判断是否需要更新最大面积
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			w := dp[i][j] // 当前行的宽度
			// 从当前行向上扫描，计算以 (i, j) 为右下角的最大矩形面积
			for k := i; k >= 0; k-- {
				// 更新宽度为当前行和扫描行的最小宽度（矩形宽度由最窄的行决定）
				w = min(w, dp[k][j])
				if w == 0 {
					break // 如果宽度为 0，无法形成矩形
				}
				// 计算当前矩形的面积：宽度 * 高度（i-k+1）
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
// 算法思路：由于存在负数，乘积可能由最大变最小或由最小变最大
// 需要同时维护以当前位置结尾的最大乘积和最小乘积
func maxProduct2(nums []int) int {
	// preMax: 以当前位置结尾的最大乘积
	// preMin: 以当前位置结尾的最小乘积（用于处理负数情况）
	// ans: 全局最大乘积
	preMax, preMin, ans := 1, 1, math.MinInt32
	for _, num := range nums {
		// 由于存在负数，最大值可能来自：preMax*num, preMin*num, num
		// 最小值可能来自：preMax*num, preMin*num, num
		// 同时更新最大值和最小值
		preMax, preMin = max(preMax*num, preMin*num, num), min(preMax*num, preMin*num, num)
		// 更新全局最大乘积
		ans = max(preMax, ans)
	}
	return ans
}

/***** 最长有效括号 *****/
// 给你一个只包含 '(' 和 ')' 的字符串
// 找出最长有效（格式正确且连续）括号子串的长度。
// 算法思路：动态规划，dp[i] 表示以 s[i] 结尾的最长有效括号子串的长度
func longestValidParentheses(s string) int {
	maxAns := 0
	dp := make([]int, len(s))
	for i := 1; i < len(s); i++ {
		if s[i] == ')' {
			if s[i-1] == '(' {
				// 情况1：s[i-1] == '(' 且 s[i] == ')'，形成一对有效括号
				// 找同级关系的一串子串，然后合并
				if i >= 2 {
					// 如果前面还有字符，加上前面的最长有效括号长度
					dp[i] = dp[i-2] + 2
				} else {
					// 如果前面没有字符，只有当前这一对括号
					dp[i] = 2
				}
			} else if i-dp[i-1] > 0 && s[i-dp[i-1]-1] == '(' {
				// 情况2：s[i-1] == ')'，需要找到与当前 ')' 匹配的 '('
				// i-dp[i-1]-1 是当前有效子串左侧的位置，必须是 '(' 才能与当前 ')' 匹配
				// 注意这里的子串已经合并过
				if i-dp[i-1] >= 2 {
					// 子串旁边可能有别的子串
					// dp[i-dp[i-1]-2] 代表着子串左侧的子串长度
					dp[i] = dp[i-1] + dp[i-dp[i-1]-2] + 2
				} else {
					// 如果左侧没有其他子串，只加上当前匹配的一对括号
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
// 算法思路：动态规划，rest[i] 表示当前余数为 i 的最大和（i = 0, 1, 2）
// 对于每个数字，更新三种余数状态的最大和
func maxSumDivThree(nums []int) int {
	// rest[0], rest[1], rest[2] 分别表示余数为 0, 1, 2 的最大和
	rest := [3]int{}
	for _, num := range nums {
		// 计算将当前数字加到三种余数状态后的新和
		a := rest[0] + num // 从余数为 0 的状态转移
		b := rest[1] + num // 从余数为 1 的状态转移
		c := rest[2] + num // 从余数为 2 的状态转移
		// 更新新余数状态的最大和（取原值和转移后的较大值）
		rest[a%3] = max(rest[a%3], a)
		rest[b%3] = max(rest[b%3], b)
		rest[c%3] = max(rest[c%3], c)
	}
	// 返回余数为 0 的最大和，即能被 3 整除的最大和
	return rest[0]
}

/***** 交错字符串 *****/
// 帮忙验证 s3 是否是由 s1 和 s2 交错 组成的。
// 算法思路：动态规划，f[i][j] 表示 s1 的前 i 个字符和 s2 的前 j 个字符能否交错组成 s3 的前 i+j 个字符
func isInterleave(s1 string, s2 string, s3 string) bool {
	n, m, t := len(s1), len(s2), len(s3)
	// 如果长度不匹配，直接返回 false
	if (n + m) != t {
		return false
	}
	// f[i][j] 表示 s1 的前 i 个字符和 s2 的前 j 个字符能否交错组成 s3 的前 i+j 个字符
	f := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]bool, m+1)
	}
	f[0][0] = true // 空字符串可以组成空字符串
	// 填充 dp 表
	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			p := i + j - 1 // s3 中对应的位置索引
			if i > 0 {
				// 如果 s1 的第 i-1 个字符等于 s3 的第 p 个字符，且 f[i-1][j] 为 true
				// 则 f[i][j] 可以为 true（从 s1 取字符）
				f[i][j] = f[i][j] || (f[i-1][j] && s1[i-1] == s3[p])
			}
			if j > 0 {
				// 如果 s2 的第 j-1 个字符等于 s3 的第 p 个字符，且 f[i][j-1] 为 true
				// 则 f[i][j] 可以为 true（从 s2 取字符）
				f[i][j] = f[i][j] || (f[i][j-1] && s2[j-1] == s3[p])
			}
		}
	}
	return f[n][m]
}

/***** 打家劫舍 II *****/
// 你是一个专业的小偷，计划偷窃沿街的房屋，每间房内都藏有一定的现金。
// 这个地方所有的房屋都围成一圈，这意味着第一个房屋和最后一个房屋是紧挨着的。
// 同时，相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。
// 算法思路：由于房屋围成圈，第一个和最后一个不能同时偷
// 分两种情况：1. 偷第一个，不偷最后一个；2. 不偷第一个，偷最后一个
// 取两种情况的最大值
func rob(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0]
	}
	if n == 2 {
		return max(nums[0], nums[1])
	}
	// 情况1：偷第一个，不偷最后一个（范围：nums[0:n-1]）
	// 情况2：不偷第一个，偷最后一个（范围：nums[1:n]）
	return max(_rob(nums[:n-1]), _rob(nums[1:]))
}

// _rob 解决线性排列的房屋打家劫舍问题
// 算法思路：动态规划，dp[i] = max(dp[i-1], dp[i-2] + nums[i])
// 使用滚动数组优化空间复杂度为 O(1)
func _rob(nums []int) int {
	// first 表示 dp[i-2]，second 表示 dp[i-1]
	first, second := nums[0], max(nums[0], nums[1])
	for _, v := range nums[2:] {
		// dp[i] = max(dp[i-1], dp[i-2] + nums[i])
		// 更新状态：first = dp[i-1], second = dp[i]
		first, second = second, max(first+v, second)
	}
	return second
}
