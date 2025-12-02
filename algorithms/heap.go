package algorithms

import (
	"container/heap"
	"sort"
)

/***** 前K个高频元素 *****/
// 给定一个整数数组 nums 和整数 k，返回出现频率前 k 高的元素
// 算法思路：使用最小堆维护前 k 个高频元素
// 1. 统计每个元素的出现频率
// 2. 使用最小堆（大小为 k）维护频率最高的 k 个元素
// 3. 当堆大小超过 k 时，弹出频率最小的元素
func topKFrequent(nums []int, k int) []int {
	// 统计每个元素的出现频率
	occurrences := map[int]int{}
	for _, num := range nums {
		occurrences[num]++
	}
	h := &IHeap{}
	heap.Init(h) // 空的不需要也行
	// 遍历所有元素及其频率
	for key, value := range occurrences {
		heap.Push(h, [2]int{key, value}) // [元素值, 频率]
		// 如果堆大小超过 k，弹出频率最小的元素
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	// 从堆中取出结果（堆中剩余的是频率最高的 k 个元素）
	ret := make([]int, k)
	for i := 0; i < k; i++ {
		ret[k-i-1] = heap.Pop(h).([2]int)[0] // 倒序填充，因为堆顶是最小值
	}
	return ret
}

// IHeap 最小堆，用于存储 [元素值, 频率] 对
// 按照频率（第二个元素）进行排序，频率小的在堆顶
type IHeap [][2]int

func (h IHeap) Len() int           { return len(h) }
func (h IHeap) Less(i, j int) bool { return h[i][1] < h[j][1] } // 按频率升序
func (h IHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IHeap) Push(x interface{}) {
	*h = append(*h, x.([2]int))
}

func (h *IHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

/***** 课程表 III *****/
// 这里有 n 门不同的在线课程，他们按从 1 到 n 编号。每一门课程有一定的持续上课时间（课程时间）t 以及关闭时间第 d 天。
// 一门课要持续学习 t 天直到第 d 天时要完成，你将会从第 1 天开始。
// 给出 n 个在线课程用 (t, d) 对表示。你的任务是找出最多可以修几门课。
// 算法思路：贪心算法 + 最大堆
// 1. 按关闭时间排序
// 2. 使用最大堆维护已选课程的学习时间
// 3. 如果当前课程可以完成，直接加入；否则尝试替换堆顶（学习时间最长的课程）
type hp struct{ sort.IntSlice }

func (h hp) Less(i, j int) bool    { return h.IntSlice[i] > h.IntSlice[j] } // 最大堆（学习时间大的在堆顶）
func (h *hp) Push(v interface{})   { h.IntSlice = append(h.IntSlice, v.(int)) }
func (h *hp) Pop() (_ interface{}) { return }

func scheduleCourse(a [][]int) (ans int) {
	// 注意！需要按关闭时间排序！
	sort.Slice(a, func(i, j int) bool { return a[i][1] < a[j][1] })
	cur := 0 // 已学习时长
	h := hp{}
	for _, p := range a {
		t := p[0]          // 课程学习时间
		if cur+t <= p[1] { // 没有超期，直接学习
			cur += t
			heap.Push(&h, t)
			ans++
		} else if h.Len() > 0 && t < h.IntSlice[0] {
			// 该课程学习时间比之前的最长学习时间要短
			// 反悔策略：放弃之前学习时间最长的课程，改为学习当前课程
			// 这样可以腾出更多时间给后续课程
			cur += t - h.IntSlice[0] // 更新已学习时长
			h.IntSlice[0] = t
			heap.Fix(&h, 0) // 这样写比 Pop 后 Push 更高效
		}
	}
	return
}

/***** 数组中的第K个最大元素 *****/
// 在未排序的数组中找到第 k 个最大的元素
// 算法思路：使用最小堆维护前 k 个最大元素
// 维护一个大小为 k 的最小堆，堆顶是这 k 个元素中的最小值，也就是第 k 大的元素
func findKthLargest(nums []int, k int) int {
	h := &KHeap{}
	heap.Init(h)
	for _, v := range nums {
		if h.Len() == k {
			// 如果堆已满，只有当新元素大于堆顶（最小值）时才替换
			if v > h.nums[0] {
				h.nums[0] = v
				heap.Fix(h, 0) // 调整堆结构
			}
		} else {
			// 堆未满，直接加入
			heap.Push(h, v)
		}
	}
	// 堆顶就是第 k 大的元素
	return h.nums[0]
}

// KHeap 最小堆，用于存储数组中的元素
// 堆顶是最小值，用于维护前 k 个最大元素中的最小值
type KHeap struct {
	nums []int
}

func (kh *KHeap) Len() int {
	return len(kh.nums)
}

func (kh *KHeap) Swap(i, j int) {
	kh.nums[i], kh.nums[j] = kh.nums[j], kh.nums[i]
}

func (kh *KHeap) Push(val interface{}) {
	kh.nums = append(kh.nums, val.(int))
}

func (kh *KHeap) Less(i, j int) bool {
	return kh.nums[i] < kh.nums[j] // 最小堆
}

func (kh *KHeap) Pop() (_ interface{}) {
	return
}
