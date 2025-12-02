package algorithms

// 上一次学习：2022.4.8，完成

type ListNode struct {
	Val  int
	Next *ListNode
}

/***** 反转链表 *****/
// 反转一个单链表
// 算法思路：迭代法，使用三个指针 pre、head、next
// 每次将当前节点的 Next 指向前一个节点，然后移动指针
func reverseList(head *ListNode) *ListNode {
	var pre *ListNode // 前一个节点
	for head != nil {
		next := head.Next // 保存下一个节点
		head.Next = pre   // 反转当前节点的指针
		pre = head        // 移动 pre 指针
		head = next       // 移动 head 指针
	}
	return pre // pre 就是新的头节点
}

/***** 重排链表 *****/
// 给定一个单链表 L：L0→L1→…→Ln-1→Ln ，
// 将其重新排列后变为： L0→Ln→L1→Ln-1→L2→Ln-2→…
// 你不能只是单纯地改变节点内部的值，而是需要实际的进行节点交换。
// 算法思路：
// 1. 找到链表中点，将链表分成两部分
// 2. 反转后半部分链表
// 3. 交替合并两个链表
func reorderList(head *ListNode) {
	if head == nil {
		return
	}
	mid := middleNode(head) // 找到链表中点
	l1 := head              // 前半部分
	l2 := mid.Next          // 后半部分
	mid.Next = nil          // 断开两部分
	l2 = reverseList(l2)    // 反转后半部分
	mergeList(l1, l2)       // 交替合并
}

// middleNode 找到链表的中间节点（如果长度为偶数，返回第一个中间节点）
// 算法思路：快慢指针，快指针每次走两步，慢指针每次走一步
func middleNode(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

// mergeList 交替合并两个链表
// 将 l2 的节点交替插入到 l1 中
func mergeList(l1, l2 *ListNode) {
	var l1Tmp, l2Tmp *ListNode
	for l1 != nil && l2 != nil {
		// 保存下一个节点
		l1Tmp = l1.Next
		l2Tmp = l2.Next

		// 将 l2 的当前节点插入到 l1 后面
		l1.Next = l2
		l1 = l1Tmp

		// 将 l1 的下一个节点连接到 l2 后面
		l2.Next = l1
		l2 = l2Tmp
	}
}

/***** K 个一组翻转链表 *****/
// 给你一个链表，每 k 个节点一组进行翻转，请你返回翻转后的链表。
// k 是一个正整数，它的值小于或等于链表的长度。
// 如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
// 算法思路：每次找到 k 个节点，反转这 k 个节点，然后继续处理下一组
func reverseKGroup(head *ListNode, k int) *ListNode {
	dummyHead := &ListNode{Next: head} // 虚拟头节点
	lastEnd := dummyHead               // 上一组的尾节点

out:
	for head != nil {
		tail := head
		// 找到当前组的尾节点
		for i := 0; i < k-1; i++ {
			tail = tail.Next
			if tail == nil {
				// 如果剩余节点不足 k 个，保持原有顺序
				break out
			}
		}
		// 反转当前组，返回新的头节点和尾节点
		lastEnd.Next, lastEnd = myReverse(head, tail)
		head = lastEnd.Next // 移动到下一组的开始
	}
	return dummyHead.Next
}

// myReverse 反转从 head 到 tail 的链表（包括 tail）
// 返回反转后的头节点（原 tail）和尾节点（原 head）
func myReverse(head, tail *ListNode) (*ListNode, *ListNode) {
	end := tail.Next // 保存 tail 的下一个节点
	prev := end      // 反转后的最后一个节点应该指向 end

	p := head
	for p != end {
		next := p.Next
		p.Next = prev
		prev = p
		p = next
	}
	return tail, head // 返回新的头节点和尾节点
}

/***** 反转链表 II *****/
// 给你单链表的头指针 head 和两个整数 left 和 right ，其中 left <= right。
// 请你反转从位置 left 到位置 right 的链表节点，返回 反转后的链表 。
// 算法思路：
// 1. 找到需要反转的区间的前一个节点（preNode）和最后一个节点（endNode）
// 2. 使用 myReverse 反转区间内的节点
// 3. 连接反转后的链表
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	if head == nil || left <= 0 || left >= right {
		return head
	}

	// 找需要翻转的开始节点和结束节点
	right = right - left + 1 // 转换为需要遍历的节点数
	dummyHead := &ListNode{Next: head}
	preNode := dummyHead
	// 找到反转区间的前一个节点
	for left > 1 {
		left--
		preNode = preNode.Next
	}
	// 找到反转区间的最后一个节点
	endNode := preNode
	for right > 0 {
		right--
		endNode = endNode.Next
	}

	// 反转区间内的节点
	h, _ := myReverse(preNode.Next, endNode)
	preNode.Next = h // 连接反转后的链表

	return dummyHead.Next
}

/***** 链表中倒数第k个节点 *****/
// 输入一个链表，输出该链表中倒数第 k 个节点
// 算法思路：快慢指针
// 1. 快指针先走 k 步
// 2. 然后快慢指针同时移动，当快指针到达末尾时，慢指针就是倒数第 k 个节点
func getKthFromEnd(head *ListNode, k int) *ListNode {
	slow, fast := head, head
	// 快指针先走 k 步
	for k > 0 && fast != nil {
		fast = fast.Next
		k--
	}
	// 快慢指针同时移动
	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}
	return slow
}

/***** 链表中环的入口节点 *****/
// 给定一个链表，返回链表开始入环的第一个节点。如果链表无环，则返回 null
// 算法思路：快慢指针（Floyd 判圈算法）
// 1. 使用快慢指针找到相遇点
// 2. 将慢指针重置到 head，然后快慢指针同时移动，相遇点就是环的入口
// 原理：快指针走了 2k 步，慢指针走了 k 步
// 多出来的 k 步，就是 n 倍的环周长
// 假定环入口到相遇点距离为 m，那么：
// head 到环入口的距离是 k-m，相遇位置再走 k-m 步，便可以到达环入口
// https://leetcode-cn.com/problems/c32eOV/solution/tu-jie-kuai-man-zhi-zhen-ji-qiao-yuan-li-rdih/
func detectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}
	// 注意这里，slow 和 fast 初始要走相应的步数
	fast, slow := head.Next.Next, head.Next
	// 找到相遇点
	for fast != slow {
		slow = slow.Next
		if fast == nil || fast.Next == nil {
			return nil // 无环
		}
		fast = fast.Next.Next
	}
	// 快指针走了 2k 步，慢指针走了 k 步
	// 多出来的 k 步，就是 n 倍的环周长
	// 假定环入口到相遇点距离为：m，那么：
	// head 到环入口的距离是 k-m，相遇位置再走 k-m 步，便可以到达环入口
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
// 算法思路：
// 1. 遍历链表，找到合适的插入位置（p.Val <= x <= p.Next.Val 或 p.Next.Val < p.Val 表示到了循环的边界）
// 2. 如果找不到合适位置（所有节点值相同或 x 是最大/最小值），插入到最大值节点后面
func insert(aNode *Node, x int) *Node {
	xNode := &Node{Val: x}
	xNode.Next = xNode
	if aNode == nil {
		return xNode
	}

	maxNode := aNode // 记录最大值节点

	var p *Node
	for p != aNode {
		if p == nil {
			p = aNode
		}
		// 找到合适的插入位置：p.Val <= x <= p.Next.Val
		// 或者 p.Next.Val < p.Val（表示到了循环的边界，从最大值到最小值）
		if p.Val <= x && (p.Next.Val >= x || p.Next.Val < p.Val) {
			p.Next, xNode.Next = xNode, p.Next
			return aNode
		}
		// 更新最大值节点
		if maxNode.Val < p.Val {
			maxNode = p
		}
		p = p.Next
	}
	// 找不到合适位置时（所有节点值相同或 x 是最大/最小值），加到最大值节点后面
	maxNode.Next, xNode.Next = xNode, maxNode.Next
	return aNode
}

/***** 展平多级双向链表 *****/
// 多级双向链表中，除了指向下一个节点和前一个节点指针之外，它还有一个子链表指针，可能指向单独的双向链表
// 将这些子列表也扁平化，使所有节点出现在单级双链表中
// 算法思路：深度优先搜索（DFS）
// 按照 child 优先的顺序遍历链表，将节点依次连接到结果链表中
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
		// 保存下一个节点和子节点
		next := node.Next
		child := node.Child
		// 将当前节点连接到结果链表
		last.Next = node
		node.Prev = last
		last = last.Next

		// 先处理子节点，再处理下一个节点（深度优先）
		dfs(child)
		dfs(next)
		node.Child = nil // 清空子节点指针
	}
	dfs(root)
	dummyHead.Next.Prev = nil // 设置头节点的 Prev 为 nil
	return dummyHead.Next
}

/***** 旋转链表 *****/
// 给你一个链表的头节点 head ，旋转链表，将链表每个节点向右移动 k 个位置
// 算法思路：
// 1. 计算链表长度，并将链表首尾相连形成环
// 2. 找到新的尾节点（原链表的第 length - k%length 个节点）
// 3. 断开环，返回新的头节点
func rotateRight(head *ListNode, k int) *ListNode {
	if k == 0 || head == nil || head.Next == nil {
		return head
	}
	// 计算链表长度
	length := 1
	iter := head
	for iter.Next != nil {
		iter = iter.Next
		length++
	}
	// 计算需要移动的步数（实际移动步数 = k % length）
	add := length - k%length
	if add == length {
		return head // 不需要旋转
	}
	// 将链表首尾相连形成环
	iter.Next = head
	// 找到新的尾节点
	for add > 0 {
		iter = iter.Next
		add--
	}
	ret := iter.Next // 新的头节点
	iter.Next = nil  // 断开环
	return ret
}

// mergeList2 合并两个有序链表
// 算法思路：使用双指针，比较两个链表的节点值，将较小的节点连接到结果链表
func mergeList2(head1, head2 *ListNode) *ListNode {
	dummyHead := &ListNode{}
	p := dummyHead
	// 同时遍历两个链表
	for head1 != nil && head2 != nil {
		if head1.Val <= head2.Val {
			p.Next = head1
			head1 = head1.Next
		} else {
			p.Next = head2
			head2 = head2.Next
		}
		p = p.Next
	}
	// 连接剩余部分
	if head1 != nil {
		p.Next = head1
	} else if head2 != nil {
		p.Next = head2
	}
	return dummyHead.Next
}

/***** 链表排序2 *****/
// 给你链表的头结点 head ，请将其按升序排列并返回排序后的链表
// 算法思路：归并排序（自顶向下）
// 1. 使用快慢指针找到链表中点，将链表分成两部分
// 2. 递归排序两部分
// 3. 合并两个有序链表
func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// 找中点
	fast, slow := head.Next, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	head2 := slow.Next
	slow.Next = nil // 断开两部分
	// 递归排序两部分
	head1 := sortList(head)
	head2 = sortList(head2)
	// 合并两个有序链表
	return mergeList2(head1, head2)
}
