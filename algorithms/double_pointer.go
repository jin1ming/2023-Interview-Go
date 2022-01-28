package algorithms

import (
	"math"
	"sort"
)

/***** 两数之和 *****/
func twoSum(nums []int, target int) []int {
	left, right := 0, len(nums)-1
	indexs := make([]int, len(nums))
	for i, _ := range indexs {
		indexs[i] = i
	}
	sort.Slice(indexs, func(i, j int) bool {
		return nums[indexs[i]] < nums[indexs[j]]
	})
	for left < right {
		sum := nums[indexs[left]] + nums[indexs[right]]
		if sum == target {
			return []int{indexs[left], indexs[right]}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	return nil
}

/***** 三数之和 *****/
// 给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，
// 使得 a + b + c = 0 ？请你找出所有和为 0 且不重复的三元组。
func threeSum(nums []int) [][]int {
	if len(nums) < 3 {
		return [][]int{}
	}

	sort.Ints(nums)
	var res [][]int

	var ptrLeft, ptrRight int
	for k, _ := range nums {
		switch {
		case nums[k] > 0:
			return res
		case k > 0 && nums[k-1] == nums[k]:
			continue
		default:
			ptrLeft = k + 1
			ptrRight = len(nums) - 1
			for ptrLeft < ptrRight {
				sum := nums[k] + nums[ptrLeft] + nums[ptrRight]
				if sum == 0 {
					r := []int{nums[k], nums[ptrLeft], nums[ptrRight]}
					res = append(res, r)
					for ptrLeft < ptrRight && nums[ptrLeft] == nums[ptrLeft+1] {
						ptrLeft += 1
					}
					for ptrLeft < ptrRight && nums[ptrRight] == nums[ptrRight-1] {
						ptrRight -= 1
					}
				}

				if sum > 0 {
					ptrRight -= 1
				} else {
					ptrLeft += 1
				}
			}
		}
	}
	return res
}

/***** 最接近的三数之和 *****/
// 给定一个包括 n 个整数的数组 nums 和 一个目标值 target。
// 找出 nums 中的三个整数，使得它们的和与 target 最接近。
// 返回这三个数的和。假定每组输入只存在唯一答案。
func threeSumClosest(nums []int, target int) int {
	sort.Ints(nums)
	var (
		n    = len(nums)
		best = math.MaxInt32
	)

	// 根据差值的绝对值来更新答案
	update := func(cur int) {
		if abs(cur-target) < abs(best-target) {
			best = cur
		}
	}

	// 枚举 a
	for i := 0; i < n; i++ {
		// 保证和上一次枚举的元素不相等
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		// 使用双指针枚举 b 和 c
		j, k := i+1, n-1
		for j < k {
			sum := nums[i] + nums[j] + nums[k]
			// 如果和为 target 直接返回答案
			if sum == target {
				return target
			}
			update(sum)
			if sum > target {
				// 如果和大于 target，移动 c 对应的指针
				k0 := k - 1
				// 移动到下一个不相等的元素
				for j < k0 && nums[k0] == nums[k] {
					k0--
				}
				k = k0
			} else {
				// 如果和小于 target，移动 b 对应的指针
				j0 := j + 1
				// 移动到下一个不相等的元素
				for j0 < k && nums[j0] == nums[j] {
					j0++
				}
				j = j0
			}
		}
	}
	return best
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}
