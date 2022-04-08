package algorithms

// 上一次学习：2022.4.7，完成

import (
	"math/rand"
	"time"
)

type LRUCache struct {
	size       int
	capacity   int
	cache      map[int]*DLinkedNode
	head, tail *DLinkedNode
}

type DLinkedNode struct {
	key, value int
	prev, next *DLinkedNode
}

func initDLinkedNode(key, value int) *DLinkedNode {
	return &DLinkedNode{
		key:   key,
		value: value,
	}
}

// Constructor2 /***** LRU 缓存 *****/
func Constructor2(capacity int) LRUCache {
	l := LRUCache{
		cache:    map[int]*DLinkedNode{},
		head:     initDLinkedNode(0, 0),
		tail:     initDLinkedNode(0, 0),
		capacity: capacity,
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

func (r *LRUCache) Get(key int) int {
	if _, ok := r.cache[key]; !ok {
		return -1
	}
	node := r.cache[key]
	r.moveToHead(node)
	return node.value
}

func (r *LRUCache) Put(key int, value int) {
	if _, ok := r.cache[key]; !ok {
		node := initDLinkedNode(key, value)
		r.cache[key] = node
		r.addToHead(node)
		r.size++
		if r.size > r.capacity {
			removed := r.removeTail()
			delete(r.cache, removed.key)
			r.size--
		}
	} else {
		node := r.cache[key]
		node.value = value
		r.moveToHead(node)
	}
}

func (r *LRUCache) addToHead(node *DLinkedNode) {
	node.prev = r.head
	node.next = r.head.next
	r.head.next.prev = node
	r.head.next = node
}

func (r *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (r *LRUCache) moveToHead(node *DLinkedNode) {
	r.removeNode(node)
	r.addToHead(node)
}

func (r *LRUCache) removeTail() *DLinkedNode {
	node := r.tail.prev
	r.removeNode(node)
	return node
}

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
