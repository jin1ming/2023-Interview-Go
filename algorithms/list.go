package algorithms

type ListNode struct {
	Val  int
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

/***** 重排链表 *****/
//给定一个单链表 L：L0→L1→…→Ln-1→Ln ，
//将其重新排列后变为： L0→Ln→L1→Ln-1→L2→Ln-2→…
//你不能只是单纯的改变节点内部的值，而是需要实际的进行节点交换。
func reorderList(head *ListNode) {
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
	for i := 1; i < (length+1)/2; i++ {
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

func reverse(head *ListNode) *ListNode {
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

/***** 反转链表 II *****/
// 给你单链表的头指针 head 和两个整数 left 和 right ，其中 left <= right。
// 请你反转从位置 left 到位置 right 的链表节点，返回 反转后的链表 。
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	if head == nil || left <= 0 || left >= right {
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

/***** 链表中环的入口节点 *****/
func detectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}
	// 注意这里，slow和fast初始要走相应的步数
	fast, slow := head.Next.Next, head.Next
	for fast != slow {
		slow = slow.Next
		if fast == nil || fast.Next == nil {
			return nil
		}
		fast = fast.Next.Next
	}
	slow = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}
	return slow
}

type Node struct {
	Val   int
	Next  *Node
	Child *Node
	Prev  *Node
}

/***** 排序的循环链表 *****/
// 给定循环升序列表中的一个点，写一个函数向这个列表中插入一个新元素 insertVal ，使这个列表仍然是循环升序的。
// 给定的可以是这个列表中任意一个顶点的指针，并不一定是这个列表中最小元素的指针。
// 如果有多个满足条件的插入位置，可以选择任意一个位置插入新的值，插入后整个列表仍然保持有序。
// 如果列表为空（给定的节点是 null），需要创建一个循环有序列表并返回这个节点。否则。请返回原先给定的节点。
func insert(aNode *Node, x int) *Node {
	xNode := &Node{Val: x}
	xNode.Next = xNode
	if aNode == nil {
		return xNode
	}
	head := aNode
	once := true
	var maxNode, minNode *Node
	for aNode != head || once {
		if aNode == head {
			once = false
		}
		if aNode.Next.Val > x && (minNode == nil || aNode.Next.Val <= minNode.Next.Val) {
			minNode = aNode
		}
		if aNode.Val == x {
			minNode = aNode
			break
		}
		aNode = aNode.Next
		if aNode.Val < x && (maxNode == nil || aNode.Val >= maxNode.Val) {
			maxNode = aNode
		}
	}
	if minNode != nil {
		minNode.Next, xNode.Next = xNode, minNode.Next
	} else if maxNode != nil {
		maxNode.Next, xNode.Next = xNode, maxNode.Next
	}
	return head
}

/***** 展平多级双向链表 *****/
func flatten2(root *Node) *Node {
	if root == nil {
		return nil
	}
	dummyHead := &Node{}
	last := dummyHead
	var dfs func(node *Node)
	dfs = func(node *Node) {
		if node == nil {
			return
		}
		next := node.Next
		child := node.Child
		last.Next = node
		node.Prev = last
		last = last.Next

		dfs(child)
		dfs(next)
		node.Child = nil
	}
	dfs(root)
	dummyHead.Next.Prev = nil
	return dummyHead.Next
}
