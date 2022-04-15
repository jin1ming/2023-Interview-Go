package algorithms

// 上一次学习：2022.4.8，完成

import "math"

type TreeNode struct {
	Val   int
	Left  *TreeNode
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

/***** 二叉树的层序遍历 *****/
// 算法思想：
// 借助于一个队列，先将根结点入队，然后出队，访问该结点，
// 若它有左子树，则将左子树根结点入队，若有右子树，则将右子树根节点入队。
// 然后出队，对出队结点访问，如此往复，直到队列为空。
func levelOrder(root *TreeNode) [][]int {
	var res [][]int

	if root == nil {
		return res
	}

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
		if level%2 == 0 {
			for i := 0; i < len(values)/2; i++ {
				values[i], values[len(values)-1-i] = values[len(values)-1-i], values[i]
			}
		}
		res = append(res, values)
		level++
	}
	return res
}

/***** 二叉树中的最大路径和 *****/
func maxPathSum(root *TreeNode) int {
	if root == nil {
		return 0
	}

	maxSum := math.MinInt64

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
		if node == nil {
			return
		}
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

/***** 二叉树剪枝 *****/
// 给定一个二叉树 根节点 root ，树的每个节点的值要么是 0，要么是 1。
// 请剪除该二叉树中所有节点的值为 0 的子树。
//节点 node 的子树为 node 本身，以及所有 node 的后代。
func pruneTree(root *TreeNode) *TreeNode {
	var dfs func(*TreeNode) bool
	dfs = func(nd *TreeNode) bool {
		if nd == nil {
			return true
		}
		l := dfs(nd.Left)
		r := dfs(nd.Right)
		if l {
			nd.Left = nil
		}
		if r {
			nd.Right = nil
		}
		return l && r && nd.Val == 0
	}
	dfs(root)
	if root != nil && root.Val == 0 && root.Left == nil && root.Right == nil {
		return nil
	}
	return root
}
