package algorithms

import (
	"errors"
	"time"
)

// LRUCache LRU（Least Recently Used）缓存实现
// 算法思路：使用双向链表 + 哈希表
// - 双向链表：维护访问顺序，头部是最新访问的，尾部是最久未访问的
// - 哈希表：实现 O(1) 的查找
// 时间复杂度：Get O(1), Put O(1)
type LRUCache struct {
	size       int                  // 当前缓存大小
	capacity   int                  // 缓存容量
	cache      map[int]*DLinkedNode // 哈希表，key -> 节点指针
	head, tail *DLinkedNode         // 双向链表的虚拟头尾节点
}

// DLinkedNode 双向链表节点
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

// NewLRUCache 创建 LRU 缓存
func NewLRUCache(capacity int) LRUCache {
	l := LRUCache{
		cache:    map[int]*DLinkedNode{},
		head:     initDLinkedNode(0, 0), // 虚拟头节点
		tail:     initDLinkedNode(0, 0), // 虚拟尾节点
		capacity: capacity,
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

// Get 获取缓存值
// 如果 key 存在，将节点移到头部（标记为最近使用），返回 value
// 如果 key 不存在，返回 -1
func (c *LRUCache) Get(key int) int {
	if _, ok := c.cache[key]; !ok {
		return -1
	}
	node := c.cache[key]
	c.moveToHead(node) // 移到头部，标记为最近使用
	return node.value
}

// Put 添加或更新缓存
// 如果 key 不存在，创建新节点并添加到头部；如果容量超限，删除尾部节点
// 如果 key 存在，更新值并移到头部
func (c *LRUCache) Put(key int, value int) {
	if _, ok := c.cache[key]; !ok {
		// key 不存在，创建新节点
		node := initDLinkedNode(key, value)
		c.cache[key] = node
		c.addToHead(node) // 添加到头部
		c.size++
		// 如果容量超限，删除尾部节点（最久未使用的）
		if c.size > c.capacity {
			removed := c.removeTail()
			delete(c.cache, removed.key)
			c.size--
		}
	} else {
		// key 存在，更新值并移到头部
		node := c.cache[key]
		node.value = value
		c.moveToHead(node)
	}
}

// addToHead 将节点添加到头部（最近使用）
func (c *LRUCache) addToHead(node *DLinkedNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

// removeNode 从双向链表中删除节点
func (c *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead 将节点移到头部（标记为最近使用）
func (c *LRUCache) moveToHead(node *DLinkedNode) {
	c.removeNode(node)
	c.addToHead(node)
}

// removeTail 删除尾部节点（最久未使用的）
func (c *LRUCache) removeTail() *DLinkedNode {
	node := c.tail.prev
	c.removeNode(node)
	return node
}

// 带超时的 LRU 缓存实现
// 使用全局 map 存储缓存项，每个缓存项都有超时时间
// 通过后台 goroutine 定期清理过期项

// object 缓存对象，包含值、创建时间和超时时间
type object struct {
	time    time.Time     // 创建时间
	timeout time.Duration // 超时时间
	obj     interface{}   // 缓存的值
}

type cacheTable map[string]*object

var (
	cache            cacheTable
	ErrTimeOut       = errors.New("The cache has been timeout.")
	ErrKeyNotFound   = errors.New("The key was not found.")
	ErrTypeAssertion = errors.New("Type assertion error.")
)

func init() {
	cache = make(cacheTable, 1000)
	go gc() // 启动后台清理 goroutine
}

// gc 垃圾回收函数，定期清理过期的缓存项
// 算法思路：每隔 1 秒遍历一次缓存，删除已过期的项
func gc() {
	for {
		for k, v := range cache {
			// 如果当前时间超过了创建时间 + 超时时间，删除该项
			if v.time.Add(v.timeout).Before(time.Now()) {
				delete(cache, k)
			}
			time.Sleep(time.Microsecond)
		}
		time.Sleep(time.Second) // 每隔 1 秒执行一次清理
	}
}

// Set 设置缓存项
// key: 缓存键
// obj: 缓存值
// timeout: 超时时间
func Set(key string, obj interface{}, timeout time.Duration) {
	cache[key] = &object{time.Now(), timeout, obj}
}

// Get 获取缓存项
// 如果 key 存在且未过期，更新访问时间并返回值
// 如果 key 不存在或已过期，返回错误
func Get(key string) (obj interface{}, err error) {
	c, ok := cache[key]
	if ok {
		now := time.Now()
		// 检查是否过期
		if c.time.Add(c.timeout).After(now) {
			c.time = now // 更新访问时间（类似 LRU 的最近使用）
			return c.obj, nil
		}
		// 已过期，删除并返回错误
		delete(cache, key)
		return nil, ErrTimeOut
	}
	return nil, ErrKeyNotFound
}

// Delete 删除缓存项
func Delete(key string) {
	delete(cache, key)
}

// HasKey 检查 key 是否存在
func HasKey(key string) bool {
	_, ok := cache[key]
	return ok
}

// main 函数示例（注释掉以避免与其他文件冲突）
// // func main() {
// 	test := "test"
// 	Set("test", test, time.Duration(10*time.Second))
// 	obj, err := Get("test")
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("The first result is :", obj)
// 	}
// 	type p struct {
// 		a, b, c, d, e int
// 	}
// 	struct2 := p{56, 7, 8, 9, 9}
// 	Set("name", struct2, time.Duration(10*time.Second))
// 	resultstruct2, err2 := Get("name")
// 	if err != nil {
// 		fmt.Println(err2)
// 	} else {
// 		fmt.Println("The second result is :", resultstruct2)
// 	}
// 	Delete("name")
// 	resultstruct3, _ := Get("name")
// 	if resultstruct3 == nil {
// 		fmt.Println("Delete was success")
// 	} else {
// 		fmt.Println("Delete was error")
// 	}
//
// 	isornot := HasKey("test")
// 	if isornot {
// 		fmt.Println("test exsited is true")
// 	} else {
// 		fmt.Println("test exsited is false")
// 	}
// 	isornot2 := HasKey("name")
// 	if isornot2 {
// 		fmt.Println("name exsited is true")
// 	} else {
// 		fmt.Println("name exsited is false")
// 	}
// }
