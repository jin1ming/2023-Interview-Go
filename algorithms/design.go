package algorithms

// 上一次学习：2022.4.7，完成

import (
	"math/rand"
	"time"
)

const maxLen = 50002

type StreamRank struct {
	data [maxLen]int
}

// Constructor4 /***** 数字流的秩 *****/
// 假设你正在读取一串整数。每隔一段时间，你希望能找出数字 x 的秩(小于或等于 x 的值的个数)。
// 请实现数据结构和算法来支持这些操作，也就是说：
// 实现 track(int x) 方法，每读入一个数字都会调用该方法；
// 实现 getRankOfNumber(int x) 方法，返回小于或等于 x 的值的个数。
// 思想：树状数组
func Constructor4() StreamRank {
	return StreamRank{}
}

func lowbit(x int) int {
	// 二进制里最低位1
	return x & -x
}

func (this *StreamRank) Track(x int) {
	for pos := x + 1; pos < maxLen; pos += lowbit(pos) {
		this.data[pos]++
	}
}

func (this *StreamRank) GetRankOfNumber(x int) int {
	ans := 0
	for pos := x + 1; pos > 0; pos -= lowbit(pos) {
		ans += this.data[pos]
	}
	return ans
}

type RandomizedSet struct {
	nums []int       // 用于随机访问
	M    map[int]int // 实现O(1)的访问和删除
	rand *rand.Rand  // 随机数种子
}

// NewRandomizedSet /***** 插入、删除和随机访问都是 O(1) 的容器 *****/
func NewRandomizedSet() RandomizedSet {
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

type Trie struct {
	children [26]*Trie
	isEnd    bool
}

// NewTrie /***** 前缀树/字典树 *****/
func NewTrie() Trie {
	return Trie{}
}

func (t *Trie) Insert(word string) {
	node := t
	for _, ch := range word {
		ch -= 'a'
		if node.children[ch] == nil {
			node.children[ch] = &Trie{}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *Trie) SearchPrefix(prefix string) *Trie {
	node := t
	for _, ch := range prefix {
		ch -= 'a'
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

func (t *Trie) Search(word string) bool {
	node := t.SearchPrefix(word)
	return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	return t.SearchPrefix(prefix) != nil
}
