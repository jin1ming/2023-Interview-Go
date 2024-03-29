package algorithms

import "sort"

/***** 在排序数组中查找元素的第一个和最后一个位置 *****/
// 给定一个按照升序排列的整数数组 nums，和一个目标值 target。
// 找出给定目标值在数组中的开始位置和结束位置。
func searchRange(nums []int, target int) []int {
	leftmost := search(nums, target)
	if leftmost == len(nums) || nums[leftmost] != target {
		return []int{-1, -1}
	}
	rightmost := search(nums, target+1) - 1
	return []int{leftmost, rightmost}
}

/***** 二分搜索 *****/
func search(nums []int, target int) int {
	left, right := 0, len(nums)-1
	ans := len(nums)
	for left <= right {
		mid := int(uint(left+right) >> 1)
		if nums[mid] >= target {
			right = mid - 1
			ans = mid
		} else {
			left = mid + 1
		}
	}
	return ans
}

/***** 有效三角形的个数 *****/
func triangleNumber(nums []int) (ans int) {
	sort.Ints(nums)
	for i, v := range nums {
		for j := i + 1; j < len(nums); j++ {
			ans += sort.SearchInts(nums[j+1:], v+nums[j])
		}
	}
	return
}
