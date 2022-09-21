package algorithms

import (
	"errors"
	"fmt"
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

func NewLRUCache(capacity int) LRUCache {
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

func (c *LRUCache) Get(key int) int {
	if _, ok := c.cache[key]; !ok {
		return -1
	}
	node := c.cache[key]
	c.moveToHead(node)
	return node.value
}

func (c *LRUCache) Put(key int, value int) {
	if _, ok := c.cache[key]; !ok {
		node := initDLinkedNode(key, value)
		c.cache[key] = node
		c.addToHead(node)
		c.size++
		if c.size > c.capacity {
			removed := c.removeTail()
			delete(c.cache, removed.key)
			c.size--
		}
	} else {
		node := c.cache[key]
		node.value = value
		c.moveToHead(node)
	}
}

func (c *LRUCache) addToHead(node *DLinkedNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

func (c *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache) moveToHead(node *DLinkedNode) {
	c.removeNode(node)
	c.addToHead(node)
}

func (c *LRUCache) removeTail() *DLinkedNode {
	node := c.tail.prev
	c.removeNode(node)
	return node
}

// 带超时的lru
type object struct {
	time    time.Time
	timeout time.Duration
	obj     interface{}
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
	go gc()
}

func gc() {
	for {
		for k, v := range cache {
			if v.time.Add(v.timeout).Before(time.Now()) {
				delete(cache, k)
			}
			time.Sleep(time.Microsecond)
		}
		time.Sleep(time.Second)
	}
}

func Set(key string, obj interface{}, timeout time.Duration) {
	cache[key] = &object{time.Now(), timeout, obj}
}

func Get(key string) (obj interface{}, err error) {
	c, ok := cache[key]
	if ok {
		now := time.Now()
		if c.time.Add(c.timeout).After(now) {
			c.time = now
			return c.obj, nil
		}
		delete(cache, key)
		return nil, ErrTimeOut
	}
	return nil, ErrKeyNotFound
}

func Delete(key string) {
	delete(cache, key)
}

func HasKey(key string) bool {
	_, ok := cache[key]
	return ok
}

func main() {
	test := "test"
	Set("test", test, time.Duration(10*time.Second))
	obj, err := Get("test")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("The first result is :", obj)
	}
	type p struct {
		a, b, c, d, e int
	}
	struct2 := p{56, 7, 8, 9, 9}
	Set("name", struct2, time.Duration(10*time.Second))
	resultstruct2, err2 := Get("name")
	if err != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("The second result is :", resultstruct2)
	}
	Delete("name")
	resultstruct3, _ := Get("name")
	if resultstruct3 == nil {
		fmt.Println("Delete was success")
	} else {
		fmt.Println("Delete was error")
	}

	isornot := HasKey("test")
	if isornot {
		fmt.Println("test exsited is true")
	} else {
		fmt.Println("test exsited is false")
	}
	isornot2 := HasKey("name")
	if isornot2 {
		fmt.Println("name exsited is true")
	} else {
		fmt.Println("name exsited is false")
	}

}
