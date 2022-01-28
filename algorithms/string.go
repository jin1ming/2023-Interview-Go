package algorithms

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

/***** 字符串转换整数 (atoi) *****/
func myAtoi(str string) int {
	return convert(clean(str))
}

func clean(s string) (sign int, abs string) {
	// 先去除首尾空格
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	// 判断第一个字符
	switch s[0] {
	// 有效的
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		sign, abs = 1, s
	// 有效的，正号
	case '+':
		sign, abs = 1, s[1:]
	// 有效的，负号
	case '-':
		sign, abs = -1, s[1:]
	// 无效的，当空字符处理，并且直接返回
	default:
		abs = ""
		return
	}
	for i, b := range abs {
		// 遍历第一波处理过的字符，如果直到第i个位置有效，那就取s[:i]，
		// 从头到这个有效的字符，剩下的就不管了，也就是break掉
		// 比如 s=123abc，那么就取123，也就是s[:3]
		if b < '0' || '9' < b {
			abs = abs[:i]
			// 一定要break，因为后面的就没用了
			break
		}
	}
	return
}

// 接收的输入是已经处理过的纯数字
func convert(sign int, absStr string) int {
	absNum := 0
	for _, b := range absStr {
		// b - '0' ==> 得到这个字符类型的数字的真实数值的绝对值
		absNum = absNum*10 + int(b-'0')
		// 检查溢出
		switch {
		case sign == 1 && absNum > math.MaxInt32:
			return math.MaxInt32
		// 这里和正数不一样的是，必须和负号相乘，也就是变成负数，否则永远走不到里面
		case sign == -1 && absNum*sign < math.MinInt32:
			return math.MinInt32
		}
	}
	return sign * absNum
}

/***** 最长回文子串 *****/
// 中心扩展算法
func longestPalindrome(s string) string {
	if s == "" {
		return ""
	}
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		left1, right1 := expandAroundCenter(s, i, i)
		left2, right2 := expandAroundCenter(s, i, i+1)
		if right1-left1 > end-start {
			start, end = left1, right1
		}
		if right2-left2 > end-start {
			start, end = left2, right2
		}
	}
	return s[start : end+1]
}

func expandAroundCenter(s string, left, right int) (int, int) {
	for ; left >= 0 && right < len(s) && s[left] == s[right]; left, right = left-1, right+1 {
	}
	return left + 1, right - 1
}

/***** 反转每对括号间的子串 *****/
func reverseParentheses(s string) string {
	// 字符串无法直接修改，转换为byte slice
	brr := []byte(s)
	var stack []int
	for i := 0; i < len(brr); i++ {
		if brr[i] == '(' {
			// 遇到左括号，加入栈中
			stack = append(stack, i)
		} else if brr[i] == ')' {
			// 题目保证括号左右匹配，所以不用检验stack中是否有左括号
			lastIdx := stack[len(stack)-1]
			// 反转左括号位置+1到右括号位置-1之间的字符
			for lj, rj := lastIdx+1, i-1; lj < rj; lj, rj = lj+1, rj-1 {
				brr[lj], brr[rj] = brr[rj], brr[lj]
			}
			// 已匹配的左括号退栈
			stack = stack[:len(stack)-1]
		}
	}

	// 去掉所有括号字符
	sb := strings.Builder{}
	for i := 0; i < len(brr); i++ {
		if brr[i] != '(' && brr[i] != ')' {
			sb.WriteByte(brr[i])
		}
	}

	return sb.String()
}

/***** 二叉树的右视图 *****/
// 给定一个二叉树的 根节点 root，想象自己站在它的右侧，
// 按照从顶部到底部的顺序，返回从右侧所能看到的节点值。
func rightSideView(root *TreeNode) []int {
	var res []int
	if root == nil {
		return res
	}
	var queue []*TreeNode
	p := root
	res = append(res, root.Val)
	queue = append(queue, root)
	for len(queue) > 0 {
		qLen := len(queue)
		right := math.MinInt32
		// 保存每一层最右边的值
		for i := 0; i < qLen; i++ {
			p = queue[i]
			if p.Left != nil {
				queue = append(queue, p.Left)
				right = p.Left.Val
			}
			if p.Right != nil {
				queue = append(queue, p.Right)
				right = p.Right.Val
			}
		}
		if right != math.MinInt32 {
			res = append(res, right)
		}
		newQueue := make([]*TreeNode, len(queue)-qLen)
		copy(newQueue, queue[qLen:])
		queue = newQueue
	}
	return res
}

/***** 字符串解码 *****/
// 示例：
// 输入：s = "3[a2[c]]"
// 输出："accaccacc"
// 如果当前的字符为数位，解析出一个数字（连续的多个数位）并进栈
// 如果当前的字符为字母或者左括号，直接进栈
// 如果当前的字符为右括号，开始出栈，一直到左括号出栈
func decodeString(s string) string {
	var stack []string
	ptr := 0
	for ptr < len(s) {
		cur := s[ptr]
		if cur >= '0' && cur <= '9' {
			digits := getDigits(s, &ptr)
			stack = append(stack, digits)
		} else if (cur >= 'a' && cur <= 'z' || cur >= 'A' && cur <= 'Z') || cur == '[' {
			stack = append(stack, string(cur))
			ptr++
		} else {
			ptr++
			var sub []string
			for stack[len(stack)-1] != "[" {
				sub = append(sub, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			for i := 0; i < len(sub)/2; i++ {
				sub[i], sub[len(sub)-i-1] = sub[len(sub)-i-1], sub[i]
			}
			stack = stack[:len(stack)-1]
			repTime, _ := strconv.Atoi(stack[len(stack)-1])
			stack = stack[:len(stack)-1]
			t := strings.Repeat(getString(sub), repTime)
			stack = append(stack, t)
		}
	}
	return getString(stack)
}

func getDigits(s string, ptr *int) string {
	ret := ""
	for ; s[*ptr] >= '0' && s[*ptr] <= '9'; *ptr++ {
		ret += string(s[*ptr])
	}
	return ret
}

func getString(v []string) string {
	ret := ""
	for _, s := range v {
		ret += s
	}
	return ret
}

/***** 最多删除一个字符得到回文 *****/
func validPalindrome(s string) bool {
	var isValid func(s string, flag bool) bool
	isValid = func(s string, flag bool) bool {
		if len(s) < 2 {
			return true
		}
		i, j := 0, len(s)-1
		for i < j {
			if s[i] == s[j] {
				i++
				j--
			} else if !flag {
				return false
			} else {
				return isValid(s[i+1:j+1], false) ||
					isValid(s[i:j], false)
			}
		}
		return true
	}
	return isValid(s, true)
}

/***** 回文子字符串的个数 *****/
func countSubstrings(s string) int {
	res := 0
	for i := 0; i < len(s); i++ {
		left, right := i, i
		for left-1 >= 0 && right+1 < len(s) &&
			s[left-1] == s[right+1] {
			left--
			right++
		}
		res += (right-left)/2 + 1
		left, right = i-1, i
		for left >= 0 && right < len(s) &&
			s[left] == s[right] {
			left--
			right++
		}
		res += (right - left - 1) / 2
	}
	return res
}

/***** 含有所有字符的最短字符串 *****/
func minWindow(s string, t string) string {
	if len(t) > len(s) {
		return ""
	}
	count := 'z' - 'A' + 1
	// nums 用来存储哪些字母还不够
	nums := make([]int, count)
	// used 用来存储哪些字母出现在t中
	used := make([]bool, count)
	// status 表示还剩几个字母没满足条件
	status := 0
	res := ""
	// 初始化 nums、used、status
	for i := 0; i < len(t); i++ {
		k := int(t[i] - 'A')
		if nums[k] == 0 {
			status++
		}
		nums[k]--
		used[k] = true
	}

	left := 0
	for right := 0; right < len(s); right++ {
		k := int(s[right] - 'A')
		if used[k] == false {
			continue
		}
		nums[k]++
		if nums[k] == 0 {
			status--
		}
		if status == 0 {
			for !used[int(s[left]-'A')] || nums[int(s[left])-'A']-1 >= 0 {
				nums[int(s[left])-'A']--
				left++
			}
			if right-left+1 < len(res) || len(res) == 0 {
				res = s[left:right]
			}
		}
	}
	return res
}

/***** 字典序排数 *****/
// TODO: 还有别的方法
// 给你一个整数 n ，按字典序返回范围 [1, n] 内所有整数。
// 你必须设计一个时间复杂度为 O(n) 且使用 O(1) 额外空间的算法。
// 先排序，再比较首尾两个字符串
func longestCommonPrefix(strs []string) string {
	sort.Strings(strs)
	end := 0
	minLen := len(strs[0])
	if len(strs[len(strs)-1]) < minLen {
		minLen = len(strs[len(strs)-1])
	}
	for ; end < minLen; end++ {
		if strs[0][end] != strs[len(strs)-1][end] {
			break
		}
	}
	return strs[0][:end]
}
