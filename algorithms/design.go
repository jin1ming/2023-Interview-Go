package algorithms

// 上一次学习：2022.4.7，完成

import (
	"math/rand"
	"time"
)

const maxLen = 50002

// StreamRank 数字流的秩数据结构
// 使用树状数组（Binary Indexed Tree, BIT）实现
type StreamRank struct {
	data [maxLen]int // 树状数组，用于快速统计前缀和
}

// Constructor4 /***** 数字流的秩 *****/
// 假设你正在读取一串整数。每隔一段时间，你希望能找出数字 x 的秩(小于或等于 x 的值的个数)。
// 请实现数据结构和算法来支持这些操作，也就是说：
// 实现 track(int x) 方法，每读入一个数字都会调用该方法；
// 实现 getRankOfNumber(int x) 方法，返回小于或等于 x 的值的个数。
// 思想：树状数组，支持 O(log n) 的单点更新和前缀和查询
func Constructor4() StreamRank {
	return StreamRank{}
}

// lowbit 返回 x 的二进制表示中最低位的 1 所对应的值
// 例如：lowbit(6) = lowbit(110) = 2
// 利用补码的性质：x & -x 可以得到最低位的 1
func lowbit(x int) int {
	// 二进制里最低位1
	return x & -x
}

// Track 记录一个数字 x
// 算法思路：树状数组的单点更新
// 将 x+1 作为索引（+1 是为了处理 0 的情况），沿着树状数组向上更新
func (this *StreamRank) Track(x int) {
	// 从 x+1 位置开始，沿着树状数组向上更新
	// pos += lowbit(pos) 是树状数组的标准更新方式
	for pos := x + 1; pos < maxLen; pos += lowbit(pos) {
		this.data[pos]++
	}
}

// GetRankOfNumber 返回小于或等于 x 的值的个数（即 x 的秩）
// 算法思路：树状数组的前缀和查询
// 查询 [1, x+1] 的前缀和，即小于等于 x 的数字个数
func (this *StreamRank) GetRankOfNumber(x int) int {
	ans := 0
	// 从 x+1 位置开始，沿着树状数组向下累加
	// pos -= lowbit(pos) 是树状数组的标准查询方式
	for pos := x + 1; pos > 0; pos -= lowbit(pos) {
		ans += this.data[pos]
	}
	return ans
}

// RandomizedSet 支持 O(1) 插入、删除和随机访问的数据结构
// 算法思路：使用数组 + 哈希表的组合
// - 数组用于随机访问（O(1)）
// - 哈希表存储值到数组索引的映射（O(1) 查找和更新）
type RandomizedSet struct {
	nums []int       // 用于随机访问，存储所有元素
	M    map[int]int // 实现O(1)的访问和删除，key: 元素值, value: 在数组中的索引
	rand *rand.Rand  // 随机数种子，用于生成随机索引
}

// NewRandomizedSet /***** 插入、删除和随机访问都是 O(1) 的容器 *****/
// 初始化一个 RandomizedSet
func NewRandomizedSet() RandomizedSet {
	r := RandomizedSet{M: make(map[int]int),
		rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
	return r
}

// Insert 插入一个元素
// 时间复杂度：O(1)
func (r *RandomizedSet) Insert(val int) bool {
	// 如果元素已存在，返回 false
	if _, ok := r.M[val]; ok {
		return false
	}
	// 将元素添加到数组末尾，并在哈希表中记录索引
	r.M[val] = len(r.nums)
	r.nums = append(r.nums, val)
	return true
}

// Remove 删除一个元素
// 算法思路：为了保持 O(1) 删除，将待删除元素与数组末尾元素交换，然后删除末尾
// 时间复杂度：O(1)
func (r *RandomizedSet) Remove(val int) bool {
	// 如果元素不存在，返回 false
	if _, ok := r.M[val]; !ok {
		return false
	}
	// 将当前值和数组末尾值互换位置再删除
	index := r.M[val]          // 获取待删除元素的索引
	t := r.nums[len(r.nums)-1] // 获取数组末尾元素
	r.nums[index] = t          // 将末尾元素移到待删除位置
	r.M[t] = index             // 更新末尾元素在哈希表中的索引

	delete(r.M, val)                // 从哈希表中删除
	r.nums = r.nums[:len(r.nums)-1] // 从数组中删除末尾元素
	return true
}

// GetRandom 随机返回一个元素
// 时间复杂度：O(1)
func (r *RandomizedSet) GetRandom() int {
	index := r.rand.Intn(len(r.nums))
	return r.nums[index]
}

// Trie 前缀树（字典树）数据结构
// 用于高效存储和检索字符串集合
// 每个节点代表一个字符，从根到节点的路径表示一个字符串
type Trie struct {
	children [26]*Trie // 26 个小写字母的子节点
	isEnd    bool      // 标记当前节点是否是某个单词的结尾
}

// NewTrie /***** 前缀树/字典树 *****/
// 初始化一个空的前缀树
func NewTrie() Trie {
	return Trie{}
}

// Insert 插入一个单词到前缀树中
// 算法思路：从根节点开始，沿着单词的每个字符向下遍历
// 如果路径不存在则创建，最后标记单词结尾
// 时间复杂度：O(m)，m 为单词长度
func (t *Trie) Insert(word string) {
	node := t
	for _, ch := range word {
		ch -= 'a' // 将字符转换为 0-25 的索引
		// 如果当前字符对应的子节点不存在，创建新节点
		if node.children[ch] == nil {
			node.children[ch] = &Trie{}
		}
		// 移动到下一个节点
		node = node.children[ch]
	}
	// 标记单词结尾
	node.isEnd = true
}

// SearchPrefix 查找前缀对应的节点
// 如果前缀存在，返回最后一个字符对应的节点；否则返回 nil
// 时间复杂度：O(m)，m 为前缀长度
func (t *Trie) SearchPrefix(prefix string) *Trie {
	node := t
	for _, ch := range prefix {
		ch -= 'a' // 将字符转换为 0-25 的索引
		// 如果路径不存在，返回 nil
		if node.children[ch] == nil {
			return nil
		}
		// 移动到下一个节点
		node = node.children[ch]
	}
	return node
}

// Search 查找单词是否在前缀树中
// 算法思路：先查找前缀，然后检查最后一个节点是否标记为单词结尾
// 时间复杂度：O(m)，m 为单词长度
func (t *Trie) Search(word string) bool {
	node := t.SearchPrefix(word)
	return node != nil && node.isEnd
}

// StartsWith 检查前缀树中是否有以给定前缀开头的单词
// 算法思路：查找前缀对应的节点是否存在即可
// 时间复杂度：O(m)，m 为前缀长度
func (t *Trie) StartsWith(prefix string) bool {
	return t.SearchPrefix(prefix) != nil
}
