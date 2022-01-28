package algorithms

import (
	"math/rand"
	"time"
)

type RandomizedSet struct {
	nums []int       // 用于随机访问
	M    map[int]int // 实现O(1)的访问和删除
	rand *rand.Rand  // 随机数种子
}

// Constructor /***** 插入、删除和随机访问都是 O(1) 的容器 *****/
func Constructor() RandomizedSet {
	r := RandomizedSet{M: make(map[int]int),
		rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
	return r
}

func (r *RandomizedSet) Insert(val int) bool {
	if _, ok := r.M[val]; ok {
		return false
	}
	r.M[val] = len(r.nums)
	r.nums = append(r.nums, val)
	return true
}

func (r *RandomizedSet) Remove(val int) bool {
	if _, ok := r.M[val]; !ok {
		return false
	}
	// 将当前值和数组末尾值互换位置再删除
	index := r.M[val]
	t := r.nums[len(r.nums)-1]
	r.nums[index] = t
	r.M[t] = index

	delete(r.M, val)
	r.nums = r.nums[:len(r.nums)-1]
	return true
}

func (r *RandomizedSet) GetRandom() int {
	index := r.rand.Intn(len(r.nums))
	return r.nums[index]
}
