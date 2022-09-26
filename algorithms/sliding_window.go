package algorithms

import "math"

// 滑动窗口一定要 right 考虑为先！！！！！！！

/***** 无重复字符的最长子串 *****/
// 类型: 滑动窗口
func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}
	res, left := 0, 0
	M := make(map[byte]int)

	for right := 0; right < len(s); right++ {
		// 这里查的是和s[right]相同的字符，上一次出现的位置
		if pos, ok := M[s[right]]; ok {
			for left <= pos {
				// left 将自增，直到 c+1 的位置
				delete(M, s[left])
				left++
			}
		}
		M[s[right]] = right
		// 记录走过的位置
		if right-left > res {
			res = right - left
		}
	}
	return res + 1
}

/***** 找到字符串中所有字母异位词 *****/
// 滑动窗口
func findAnagrams2(s string, p string) []int {
	n, m := len(s), len(p)
	if n < m {
		return nil
	}

	var res []int
	cntS, cntP := [26]int{}, [26]int{}
	for i := 0; i < m; i++ {
		cntP[p[i]-'a']++
	}

	left, right := 0, 0
	// 右窗口开始不断向右移动
	for ; right < n; right++ {
		curRight := s[right] - 'a'
		// 将右窗口当前访问到的元素个数加1
		cntS[curRight]++
		// 当前窗口中 curRight 比 cntP 数组中对应元素的个数
		// 要多的时候就该移动左窗口指针
		for cntS[curRight] > cntP[curRight] {
			curLeft := s[left] - 'a'
			// 将左窗口当前访问到的元素个数减1
			cntS[curLeft]--
			left++
		}
		if right-left+1 == m {
			res = append(res, left)
		}
	}
	return res
}

/***** 和大于等于 target 的最短子数组 *****/
// 滑动窗口一定要 right 考虑为先！！！！！！！
func minSubArrayLen(target int, nums []int) int {
	res := math.MaxInt64
	left, sum := 0, 0
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		if sum >= target {
			for sum-nums[left] >= target {
				sum -= nums[left]
				left++
			}
			if right-left+1 < res {
				res = right - left
			}
		}
	}
	if res == math.MaxInt64 {
		return 0
	}
	return res
}

func minOperations(nums []int, x int) int {
	res := math.MaxInt64

	target := -x
	for _, n := range nums {
		target += n
	}
	if target == 0 {
		return len(nums)
	}
	if target < 0 {
		return -1
	}
	left := 0
	sum := 0
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		for sum > target {
			sum -= nums[left]
			left++
		}
		if sum == target {
			opTimes := left + len(nums) - right - 1
			if opTimes < res {
				res = opTimes
			}
		}
	}

	if res == math.MaxInt64 {
		return -1
	}
	return res
}

/***** 乘积小于 K 的子数组 *****/
// 滑动窗口一定要 right 考虑为先！！！！！！！
func numSubarrayProductLessThanK(nums []int, k int) int {
	left := 0
	sum, res := 1, 0

	// 每次循环找出满足条件的最大的 left，再将 right 加 1
	// 因为 nums 中每个数都大于等于 1
	// 所以每次 right 右移后，left 向左移动时不会满足条件
	for right := 0; right < len(nums); right++ {
		sum *= nums[right]
		for left <= right && sum >= k {
			sum /= nums[left]
			left++
		}
		res += right - left + 1
		// 关于 加 right - left + 1（子数组长度）:
		// 假设子数组 ABCD, right指向D, 此时统计的是D CD BCD ABCD
		// 统计的是包含right指向元素的子数组个数
		right++
	}
	return res
}

/***** 找到字符串中所有字母异位词 *****/
// 典型滑动窗口
func findAnagrams(s string, p string) []int {
	pl := len(p)
	sl := len(s)
	if pl > sl {
		return nil
	}
	var result []int

	m := make(map[byte]int)
	for i := 0; i < pl; i++ {
		m[p[i]]++
	}

	for i1 := 0; i1 < pl; i1++ {
		m[s[i1]]--
	}

out:
	for i := 0; i < sl-pl+1; i++ {
		if i > 0 {
			m[s[i-1]]++
			m[s[i+pl-1]]--
		}

		for _, v := range m {
			if v != 0 {
				continue out
			}
		}
		result = append(result, i)
	}
	return result
}

/***** 最小覆盖子串 *****/
// 给你一个字符串 s 、一个字符串 t 。返回 s 中涵盖 t 所有字符的最小子串。
// 如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。
func minWindow2(s string, t string) string {
	tMap := make(map[byte]int)
	sMap := make(map[byte]int)
	for _, c := range []byte(t) {
		tMap[c]++
	}

	res := ""
	left := 0

	isValid := func(ch byte) bool {
		for k, v := range tMap {
			if sMap[k] < v {
				return false
			}
		}
		return true
	}

	for right := 0; right < len(s); right++ {
		ch := s[right]
		if _, ok := tMap[ch]; !ok {
			continue
		}
		sMap[ch]++
		if isValid(ch) {
			for tMap[s[left]] == 0 || sMap[s[left]] > tMap[s[left]] {
				sMap[s[left]]--
				left++
			}
			if len(res) == 0 || len(res) > right-left+1 {
				res = s[left : right+1]
			}
			sMap[s[left]]--
			left++
		}
	}
	return res
}
