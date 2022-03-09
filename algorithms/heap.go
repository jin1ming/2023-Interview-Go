package algorithms

import (
	"container/heap"
	"sort"
)

/***** 前K个高频元素 *****/
func topKFrequent(nums []int, k int) []int {
	occurrences := map[int]int{}
	for _, num := range nums {
		occurrences[num]++
	}
	h := &IHeap{}
	heap.Init(h) // 空的不需要也行
	for key, value := range occurrences {
		heap.Push(h, [2]int{key, value})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	ret := make([]int, k)
	for i := 0; i < k; i++ {
		ret[k-i-1] = heap.Pop(h).([2]int)[0]
	}
	return ret
}

type IHeap [][2]int

func (h IHeap) Len() int           { return len(h) }
func (h IHeap) Less(i, j int) bool { return h[i][1] < h[j][1] }
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
type hp struct{ sort.IntSlice }

func (h hp) Less(i, j int) bool    { return h.IntSlice[i] > h.IntSlice[j] } // 最大堆
func (h *hp) Push(v interface{})   { h.IntSlice = append(h.IntSlice, v.(int)) }
func (h *hp) Pop() (_ interface{}) { return }

// 这里有 n 门不同的在线课程，他们按从 1 到 n 编号。每一门课程有一定的持续上课时间（课程时间）t 以及关闭时间第 d 天。
// 一门课要持续学习 t 天直到第 d 天时要完成，你将会从第 1 天开始。
// 给出 n 个在线课程用 (t, d) 对表示。你的任务是找出最多可以修几门课。
func scheduleCourse(a [][]int) (ans int) {
	// 注意！需要按关闭时间排序！
	sort.Slice(a, func(i, j int) bool { return a[i][1] < a[j][1] })
	cur := 0 // 已学习时长
	h := hp{}
	for _, p := range a {
		if t := p[0]; cur+t <= p[1] { // 没有超期，直接学习
			cur += t
			heap.Push(&h, t)
			ans++
		} else if h.Len() > 0 && t < h.IntSlice[0] { // 该课程学习时间比之前的最长学习时间要短
			cur += t - h.IntSlice[0] // 反悔，放弃之前的最长学习时间的课程，改为学习该课程
			h.IntSlice[0] = t
			heap.Fix(&h, 0) // 这样写比 Pop 后 Push 更高效
		}
	}
	return
}

/***** 数组中的第K个最大元素 *****/
// 维护一个最小堆，这个最小堆保证存放的是数组中的最大值
func findKthLargest(nums []int, k int) int {
	h := &KHeap{}
	heap.Init(h)
	for _, v := range nums {
		if h.Len() == k {
			if v > h.nums[0] {
				h.nums[0] = v
				heap.Fix(h, 0)
			}
		} else {
			heap.Push(h, v)
		}
	}
	return h.nums[0]
}

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
	return kh.nums[i] < kh.nums[j]
}

func (kh *KHeap) Pop() (_ interface{}) {
	return
}
