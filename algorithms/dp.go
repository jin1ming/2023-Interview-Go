package algorithms

import "math"

//单调栈实现
func maximalRectangle2(matrix [][]byte) int {
	if matrix == nil || len(matrix) == 0 {
		return 0
	}
	//保存最终结果
	res := 0
	//行数，列数
	m, n := len(matrix), len(matrix[0])
	//高度数组（统计每一行中1的高度）
	height := make([]int, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			//每一行去找1的高度
			//如果不是‘1’，则将当前高度置为0
			if matrix[i][j] == '0' {
				height[j] = 0
			} else {
				//是‘1’，则将高度加1
				height[j] = height[j] + 1
			}
		}
		//更新最大矩形的面积
		res = int(math.Max(float64(res), float64(largestRectangleArea2(height))))
	}
	return res
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
	dp := make([]bool, len(s)+1)
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
func maximalRectangle0(matrix [][]byte) int {
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

/***** 最长递增路径 *****/
func longestIncreasingPath(matrix [][]int) int {
	res := 1 // res最小为1
	row := len(matrix)
	col := len(matrix[0])

	// dp数组，每个元素表示当前点开始所能构成的最长路径
	bitmap := make([][]int, row)
	for i := range matrix {
		bitmap[i] = make([]int, col)
	}

	directions := [][]int{[]int{1, 0}, []int{0, 1}, []int{-1, 0}, []int{0, -1}}

	isValid := func(r, c int) bool {
		if r < 0 || r >= row || c < 0 || c >= col {
			return false
		}
		return true
	}

	var dfs func(r, c, step int) int

	dfs = func(r, c, step int) int {
		// 判断该处是否走过
		if bitmap[r][c] > 0 {
			// 构成最大路径大于res进行更新
			if bitmap[r][c] > res-step {
				res = step + bitmap[r][c]
			}
			// 走过无需再走
			return bitmap[r][c]
		}

		step += 1 // 该点默认能到达

		if step > res {
			res = step
		}
		newStep := 0 // 新走的最长路径
		for _, d := range directions {
			if isValid(r+d[0], c+d[1]) && matrix[r+d[0]][c+d[1]] > matrix[r][c] {
				newStep = max(newStep, dfs(r+d[0], c+d[1], step))
			}
		}
		newStep++ // 考虑当前点
		bitmap[r][c] = newStep
		return newStep
	}

	// 对所有点扫描一遍
	for r := range matrix {
		for c := range matrix[0] {
			dfs(r, c, 0)
		}
	}

	return res
}
