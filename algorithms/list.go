package algorithms

// 上一次学习：2022.4.8，完成

type ListNode struct {
	Val  int
	Next *ListNode
}

/***** 反转链表 *****/
func reverseList(head *ListNode) *ListNode {
	var pre *ListNode
	for head != nil {
		next := head.Next
		head.Next = pre
		pre = head
		head = next
	}
	return pre
}

/***** 重排链表 *****/
//给定一个单链表 L：L0→L1→…→Ln-1→Ln ，
//将其重新排列后变为： L0→Ln→L1→Ln-1→L2→Ln-2→…
//你不能只是单纯地改变节点内部的值，而是需要实际的进行节点交换。
func reorderList(head *ListNode) {
	if head == nil {
		return
	}
	mid := middleNode(head)
	l1 := head
	l2 := mid.Next
	mid.Next = nil
	l2 = reverseList(l2)
	mergeList(l1, l2)
}

func middleNode(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

func mergeList(l1, l2 *ListNode) {
	var l1Tmp, l2Tmp *ListNode
	for l1 != nil && l2 != nil {
		l1Tmp = l1.Next
		l2Tmp = l2.Next

		l1.Next = l2
		l1 = l1Tmp

		l2.Next = l1
		l2 = l2Tmp
	}
}

/***** K 个一组翻转链表 *****/
// 给你一个链表，每 k 个节点一组进行翻转，请你返回翻转后的链表。
// k 是一个正整数，它的值小于或等于链表的长度。
// 如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
func reverseKGroup(head *ListNode, k int) *ListNode {
	dummyHead := &ListNode{Next: head}
	lastEnd := dummyHead

out:
	for head != nil {
		tail := head
		for i := 0; i < k-1; i++ {
			tail = tail.Next
			if tail == nil {
				break out
			}
		}
		lastEnd.Next, lastEnd = myReverse(head, tail)
		head = lastEnd.Next
	}
	return dummyHead.Next
}

func myReverse(head, tail *ListNode) (*ListNode, *ListNode) {
	end := tail.Next
	prev := end

	p := head
	for p != end {
		next := p.Next
		p.Next = prev
		prev = p
		p = next
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

	// 找需要翻转的开始节点和结束节点
	right = right - left + 1
	dummyHead := &ListNode{Next: head}
	preNode := dummyHead
	for left > 1 {
		left--
		preNode = preNode.Next
	}
	endNode := preNode
	for right > 0 {
		right--
		endNode = endNode.Next
	}

	h, _ := myReverse(preNode.Next, endNode)
	preNode.Next = h

	return dummyHead.Next
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

/***** 链表中环的入口节点 *****/
// https://leetcode-cn.com/problems/c32eOV/solution/tu-jie-kuai-man-zhi-zhen-ji-qiao-yuan-li-rdih/
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
	// 快指针走了2k步，慢指针走了k步
	// 多出来的k步，就是n倍的环周长
	// 假定环入口到相遇点距离为：m，那么：
	// head到环入口的距离是k-m，相遇位置再走k-m步，便可以到达环入口
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

	maxNode := aNode

	var p *Node
	for p != aNode {
		if p == nil {
			p = aNode
		}
		if p.Val <= x && (p.Next.Val >= x || p.Next.Val < p.Val) {
			p.Next, xNode.Next = xNode, p.Next
			return aNode
		}
		if maxNode.Val < p.Val {
			maxNode = p
		}
		p = p.Next
	}
	// 找不到合适位置时，加到最大值节点后面
	maxNode.Next, xNode.Next = xNode, maxNode.Next
	return aNode
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

/***** 旋转链表 *****/
func rotateRight(head *ListNode, k int) *ListNode {
	if k == 0 || head == nil || head.Next == nil {
		return head
	}
	length := 1
	iter := head
	for iter.Next != nil {
		iter = iter.Next
		length++
	}
	add := length - k%length
	if add == length {
		return head
	}
	iter.Next = head
	for add > 0 {
		iter = iter.Next
		add--
	}
	ret := iter.Next
	iter.Next = nil
	return ret
}

func mergeList2(head1, head2 *ListNode) *ListNode {
	dummyHead := &ListNode{}
	temp, temp1, temp2 := dummyHead, head1, head2
	for temp1 != nil && temp2 != nil {
		if temp1.Val <= temp2.Val {
			temp.Next = temp1
			temp1 = temp1.Next
		} else {
			temp.Next = temp2
			temp2 = temp2.Next
		}
		temp = temp.Next
	}
	if temp1 != nil {
		temp.Next = temp1
	} else if temp2 != nil {
		temp.Next = temp2
	}
	return dummyHead.Next
}

func sortList(head *ListNode) *ListNode {
	if head == nil {
		return head
	}

	length := 0
	for node := head; node != nil; node = node.Next {
		length++
	}

	dummyHead := &ListNode{Next: head}
	for subLength := 1; subLength < length; subLength <<= 1 {
		prev, cur := dummyHead, dummyHead.Next
		for cur != nil {
			head1 := cur
			for i := 1; i < subLength && cur.Next != nil; i++ {
				cur = cur.Next
			}

			head2 := cur.Next
			cur.Next = nil
			cur = head2
			for i := 1; i < subLength && cur != nil && cur.Next != nil; i++ {
				cur = cur.Next
			}

			var next *ListNode
			if cur != nil {
				next = cur.Next
				cur.Next = nil
			}

			prev.Next = mergeList2(head1, head2)

			for prev.Next != nil {
				prev = prev.Next
			}
			cur = next
		}
	}
	return dummyHead.Next
}

func cutList(head *ListNode) *ListNode {
	slow := head
	fast := head
	preSlow := slow
	for fast != nil && fast.Next != nil {
		preSlow = slow
		slow = slow.Next
		fast = fast.Next.Next
	}
	preSlow.Next = nil
	return slow
}

func mergeList3(l, r *ListNode) *ListNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}

	head := &ListNode{}
	cur := head
	for l != nil && r != nil {
		if l.Val < r.Val {
			cur.Next = l
			cur = l
			l = l.Next
		} else {
			cur.Next = r
			cur = r
			r = r.Next
		}
	}

	if l != nil {
		cur.Next = l
	}
	if r != nil {
		cur.Next = r
	}

	return head.Next
}

func sortList2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	mid := cutList(head)
	l := sortList(head)
	r := sortList(mid)
	return mergeList3(l, r)
}
