package algorithms

// 上一次学习：2022.4.7，完成

import (
	"strings"
)

/***** 全排列 *****/
// 给定一个不含重复数字的数组 nums，返回其所有可能的全排列
// 算法思路：回溯算法（深度优先搜索）
// 使用 used 数组标记已使用的元素，path 记录当前路径
func permute(nums []int) [][]int {
	length := len(nums)
	var res [][]int
	if length == 0 {
		return res
	}

	var path []int               // 当前排列路径
	used := make([]bool, length) // 标记元素是否已使用
	dfs1(nums, length, 0, path, used, &res)
	return res
}

// dfs1 深度优先搜索生成全排列
// nums: 原始数组
// length: 数组长度
// depth: 当前递归深度（已选择的元素个数）
// path: 当前排列路径
// used: 标记数组，used[i] 表示 nums[i] 是否已使用
// res: 结果数组
func dfs1(nums []int, length int, depth int, path []int, used []bool, res *[][]int) {
	// 如果已选择所有元素，保存当前排列
	if depth == length {
		p := make([]int, length)
		copy(p, path) // 注意使用copy，避免引用问题
		*res = append(*res, p)
		return
	}

	// 遍历所有未使用的元素
	for i := 0; i < length; i++ {
		if used[i] {
			continue // 跳过已使用的元素
		}
		// 选择当前元素
		path = append(path, nums[i])
		used[i] = true
		// 递归处理下一层
		dfs1(nums, length, depth+1, path, used, res)
		// 回溯：撤销选择
		path = path[:len(path)-1]
		used[i] = false
	}
}

/***** 八皇后 *****/
// 设计一种算法，打印 N 皇后在 N × N 棋盘上的各种摆法，
// 其中每个皇后都不同行、不同列，也不在对角线上。
// 这里的"对角线"指的是所有的对角线，不只是平分整个棋盘的那两条对角线。
// PS: 注意，需要按行放置，一行放一个
// 算法思路：回溯算法，按行放置皇后，使用三个标记数组快速判断冲突
func solveNQueens(n int) [][]string {
	var res [][]string
	matrix := make([][]bool, n) // 棋盘，true 表示放置皇后
	for s := range matrix {
		matrix[s] = make([]bool, n)
	}

	// 收集结果：将棋盘转换为字符串数组
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
	cols := make([]bool, n)       // 记录访问过的列，cols[i] 表示第 i 列已有皇后
	corner1 := make(map[int]bool) // 记录左对角线（左上到右下），key = col - row
	corner2 := make(map[int]bool) // 记录右对角线（右上到左下），key = col + row
	var dfs func(row int)
	dfs = func(row int) {
		// 如果已放置完所有行，保存结果
		if row == n {
			resAdd()
			return
		}
		// 尝试在当前行的每一列放置皇后
		for i := range matrix[row] {
			// 检查列、左对角线、右对角线是否冲突
			if cols[i] == false && corner1[i-row] == false && corner2[i+row] == false {
				if row > 0 && matrix[row][i] {
					continue
				}
				// 放置皇后
				matrix[row][i] = true
				cols[i] = true
				corner1[i-row] = true
				corner2[i+row] = true
				// 递归处理下一行
				dfs(row + 1) // 去下一行
				// 回溯：撤销选择
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

/***** 复原 IP 地址 *****/
// 给定一个只包含数字的字符串，用以表示一个 IP 地址，
// 返回所有可能从 s 获得的 有效 IP 地址。
// 你可以按任何顺序返回答案。
// 算法思路：回溯算法，将字符串分割成 4 段，每段需要满足 IP 地址的规则
// IP 地址规则：每段长度为 1-3，值在 0-255 之间，不能有前导 0
func restoreIpAddresses2(s string) []string {
	var res []string
	if len(s) < 4 {
		return res
	}

	var group [4][2]int // 记录的分段，以及每段的开始和结束下标，如"123"就是[0, 3]
	groupId := 0        // 当前处于第几个段（0-3）

	// isValid 判断当前分割是否有效
	isValid := func(left, right int) bool {
		// 判断当前情况是否无效

		if groupId > 3 { // 段分组数目不得大于4，注意第4个groupId为3
			return false
		}

		if right > len(s) || len(s)-right > (4-groupId)*3 {
			// 剩余长度无法容纳（每段最多3个字符）
			return false
		}

		sem := s[left:right] // 字符串引用，降低开销
		if len(sem) == 0 || len(sem) > 3 || sem[0] == '0' && len(sem) > 1 {
			// ip段长度应该在1-3之间，且不能出现前导0（如"01"、"001"）
			return false
		}
		if len(sem) == 3 && sem > "255" {
			// ip段的值不能大于255，这里可以用字符串直接比较，用strings.Compare性能会好一丢丢
			return false
		}
		return true
	}

	// storeRes 将分割结果转换为 IP 地址字符串
	storeRes := func() {
		// 将数组翻译为字符串，可以通过buffer拼接来优化
		r := s[group[0][0]:group[0][1]]
		for i := 1; i < 4; i++ {
			r += "." + s[group[i][0]:group[i][1]]
		}
		res = append(res, r)
	}

	var dfs func(left, right int)
	dfs = func(left, right int) {
		// 如果当前分割无效，直接返回
		if !isValid(left, right) {
			return
		}
		old := group[groupId]                              // 保存group环境
		group[groupId][0], group[groupId][1] = left, right // 更新当前group情况

		// 如果已处理完所有字符且使用了4个段，保存结果
		if right == len(s) && groupId == 3 {
			// 走到尽头，并且使用了4个group -> 保存结果
			storeRes()
		}

		// 选择1：将下一个字符加到当前group（扩展当前段）
		dfs(left, right+1)

		// 选择2：下一个字符开始为新的group（开始新段）
		groupId++
		dfs(right, right+1)
		groupId--

		group[groupId] = old // 恢复group环境（回溯）
	}

	dfs(0, 1) // 输入为第一个字符
	return res
}

/***** 岛屿数量 *****/
// 给你一个由 '1'（陆地）和 '0'（水）组成的的二维网格，请你计算网格中岛屿的数量。
// 岛屿总是被水包围，并且每座岛屿只能由水平方向和/或竖直方向上相邻的陆地连接形成。
// 此外，你可以假设该网格的四条边均被水包围。
// 算法思路：深度优先搜索（DFS）
// 遍历网格，遇到 '1' 就将其标记为 '0'，并递归标记所有相邻的 '1'
// 每次遇到新的 '1' 就说明发现了一个新岛屿
func numIslands(grid [][]byte) int {
	row := len(grid)
	if row == 0 {
		return 0
	}
	col := len(grid[0])
	num := 0

	// 遍历整个网格
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			if grid[r][c] == '1' {
				// 找到一个陆地，发现新岛屿
				num++
				dfs2(grid, r, c) // 通过 DFS 将该岛屿的所有陆地标记为水
			}
		}
	}
	return num
}

// dfs2 深度优先搜索，将当前陆地及其所有相邻陆地标记为水
// 算法思路：从当前位置开始，向上下左右四个方向递归搜索
// 将访问过的 '1' 标记为 '0'，避免重复计算
func dfs2(grid [][]byte, r int, c int) {
	row := len(grid)
	col := len(grid[0])

	// 将当前陆地标记为水
	grid[r][c] = '0'

	// 向上搜索
	if r-1 >= 0 && grid[r-1][c] == '1' {
		dfs2(grid, r-1, c)
	}
	// 向下搜索
	if r+1 < row && grid[r+1][c] == '1' {
		dfs2(grid, r+1, c)
	}
	// 向左搜索
	if c-1 >= 0 && grid[r][c-1] == '1' {
		dfs2(grid, r, c-1)
	}
	// 向右搜索
	if c+1 < col && grid[r][c+1] == '1' {
		dfs2(grid, r, c+1)
	}
}
