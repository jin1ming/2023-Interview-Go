package algorithms

// 构建无向图
func buildGraph(n int, dislikes [][]int) [][]int {
	noPeace := make([][]int, n+1)
	for _, pair := range dislikes {
		x, y := pair[0], pair[1]
		noPeace[x] = append(noPeace[x], y)
		noPeace[y] = append(noPeace[y], x)
	}

	return noPeace
}

/***** 可能的二分法 *****/
func possibleBipartition(n int, dislikes [][]int) bool {
	ok := true
	visited := make([]bool, n+1)
	color := make([]bool, n+1)

	graph := buildGraph(n, dislikes)
	var dfs func(x int)
	dfs = func(x int) {
		// 如果不是二分图，就不再遍历了
		if !ok {
			return
		}
		visited[x] = true
		for _, y := range graph[x] {
			// 相邻节点w未被访问过，给w涂v不同的色
			// 相邻节点w被访问过，若w与v同色，则不是二分图
			if !visited[y] {
				color[y] = !color[x]
				// 继续遍历
				dfs(y)
			} else if color[y] == color[x] {
				ok = false
			}
		}
	}
	// 遍历每一个节点，验证所有子图都不是二分图
	for i := 0; i <= n; i++ {
		if !visited[i] {
			dfs(i)
		}
	}

	return ok
}
