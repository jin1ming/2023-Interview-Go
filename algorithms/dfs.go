package algorithms

// 上一次学习：2022.4.7，完成

import (
	"strings"
)

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

func dfs1(nums []int, length int, depth int, path []int, used []bool, res *[][]int) {
	if depth == length {
		p := make([]int, length)
		copy(p, path) // 注意使用copy
		*res = append(*res, p)
		return
	}

	for i := 0; i < length; i++ {
		if used[i] {
			continue
		}
		path = append(path, nums[i])
		used[i] = true
		dfs1(nums, length, depth+1, path, used, res)
		path = path[:len(path)-1]
		used[i] = false
	}
}

/***** 八皇后 *****/
// 设计一种算法，打印 N 皇后在 N × N 棋盘上的各种摆法，
// 其中每个皇后都不同行、不同列，也不在对角线上。
// 这里的“对角线”指的是所有的对角线，不只是平分整个棋盘的那两条对角线。
// PS: 注意，需要按行放置，一行放一个
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
	cols := make([]bool, n)       // 记录访问过的列
	corner1 := make(map[int]bool) // 记录该左对角线
	corner2 := make(map[int]bool) // 记录该右对角线
	var dfs func(row int)
	dfs = func(row int) {
		if row == n {
			resAdd()
			return
		}
		for i := range matrix[row] {
			if cols[i] == false && corner1[i-row] == false && corner2[i+row] == false {
				if row > 0 && matrix[row][i] {
					continue
				}
				matrix[row][i] = true
				cols[i] = true
				corner1[i-row] = true
				corner2[i+row] = true
				dfs(row + 1) // 去下一行
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
func restoreIpAddresses2(s string) []string {
	var res []string
	if len(s) < 4 {
		return res
	}

	var group [4][2]int //记录的分段，以及每段的开始和结束下表，如"123"就是[0, 3]
	groupId := 0        //当前处于第几个段

	isValid := func(left, right int) bool {
		// 判断当前情况是否无效

		if groupId > 3 { // 段分组数目不得大于4，注意第4个groupId为3
			return false
		}

		if right > len(s) || len(s)-right > (4-groupId)*3 {
			// 剩余长度无法容纳
			return false
		}

		sem := s[left:right] // 字符串引用，降低开销
		if len(sem) == 0 || len(sem) > 3 || sem[0] == '0' && len(sem) > 1 {
			// ip段长度应该在1-3之间，且不能出现连续的0
			return false
		}
		if len(sem) == 3 && sem > "255" {
			// ip段的值不能大于255，这里可以用字符串直接比较，用strings.Compare性能会好一丢丢
			return false
		}
		return true
	}

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
		if !isValid(left, right) {
			return
		}
		old := group[groupId]                              // 保存group环境
		group[groupId][0], group[groupId][1] = left, right // 更新当前group情况

		if right == len(s) && groupId == 3 {
			// 走到尽头，并且使用了4个group -> 保存结果
			storeRes()
		}

		// 选择1：将下一个字符加到当前group
		dfs(left, right+1)

		// 选择2：下一个字符开始为新的group
		groupId++
		dfs(right, right+1)
		groupId--

		group[groupId] = old // 恢复group环境
	}

	dfs(0, 1) // 输入为第一个字符
	return res
}

/***** 岛屿数量 *****/
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

func dfs2(grid [][]byte, r int, c int) {
	row := len(grid)
	col := len(grid[0])

	grid[r][c] = '0'

	if r-1 >= 0 && grid[r-1][c] == '1' {
		dfs2(grid, r-1, c)
	}
	if r+1 < row && grid[r+1][c] == '1' {
		dfs2(grid, r+1, c)
	}
	if c-1 >= 0 && grid[r][c-1] == '1' {
		dfs2(grid, r, c-1)
	}
	if c+1 < col && grid[r][c+1] == '1' {
		dfs2(grid, r, c+1)
	}
}
