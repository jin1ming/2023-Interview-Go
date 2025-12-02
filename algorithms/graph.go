package algorithms

// buildGraph 构建无向图
// 使用邻接表表示图，每个节点的邻居列表存储在数组中
// n: 节点数量（节点编号从 1 到 n）
// dislikes: 边列表，每个边 [x, y] 表示节点 x 和 y 之间有边
func buildGraph(n int, dislikes [][]int) [][]int {
	noPeace := make([][]int, n+1) // 邻接表，索引从 1 开始
	for _, pair := range dislikes {
		x, y := pair[0], pair[1]
		// 无向图：x 和 y 互为邻居
		noPeace[x] = append(noPeace[x], y)
		noPeace[y] = append(noPeace[y], x)
	}

	return noPeace
}

/***** 可能的二分法 *****/
// 给定一组 n 个人（编号为 1, 2, ..., n），以及一些不喜欢关系
// 判断是否可以将这些人分成两组，使得每组内的人都不互相不喜欢
// 算法思路：判断图是否为二分图
// 使用深度优先搜索（DFS）进行图的着色，相邻节点必须不同色
// 如果能成功着色，说明是二分图，可以分成两组
func possibleBipartition(n int, dislikes [][]int) bool {
	ok := true                   // 标记是否为二分图
	visited := make([]bool, n+1) // 标记节点是否已访问
	color := make([]bool, n+1)   // 标记节点的颜色（true/false 表示两种颜色）

	graph := buildGraph(n, dislikes)
	var dfs func(x int)
	dfs = func(x int) {
		// 如果不是二分图，就不再遍历了（提前剪枝）
		if !ok {
			return
		}
		visited[x] = true
		// 遍历当前节点的所有邻居
		for _, y := range graph[x] {
			// 相邻节点 y 未被访问过，给 y 涂与 x 不同的颜色
			// 相邻节点 y 被访问过，若 y 与 x 同色，则不是二分图
			if !visited[y] {
				color[y] = !color[x] // 给邻居涂上不同的颜色
				// 继续遍历邻居节点
				dfs(y)
			} else if color[y] == color[x] {
				// 如果相邻节点颜色相同，说明不是二分图
				ok = false
			}
		}
	}
	// 遍历每一个节点，验证所有连通分量都是二分图
	// 注意：图可能不连通，需要遍历所有连通分量
	for i := 1; i <= n; i++ {
		if !visited[i] {
			dfs(i)
		}
	}

	return ok
}
