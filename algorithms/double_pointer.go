package algorithms

// 上一次学习：2022.4.7，完成

import (
	"math"
	"sort"
)

/***** 两数之和 *****/
func twoSum(nums []int, target int) []int {
	left, right := 0, len(nums)-1
	indexes := make([]int, len(nums))
	for i := range indexes {
		indexes[i] = i
	}
	sort.Slice(indexes, func(i, j int) bool {
		return nums[indexes[i]] < nums[indexes[j]]
	})
	for left < right {
		sum := nums[indexes[left]] + nums[indexes[right]]
		if sum == target {
			return []int{indexes[left], indexes[right]}
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
	if len(nums) < 3 {
		return -1
	}

	sort.Ints(nums) // 注意不要忘记排序
	best := math.MaxInt32
	checkAndUpdate := func(value int) {
		if abs(value-target) < abs(best-target) {
			// 注意不要忘记和上次比较时取差值的绝对值
			best = value
		}
	}

	for k, v := range nums {
		if k > 0 && v == nums[k-1] {
			continue
		}
		left, right := k+1, len(nums)-1
		for left < right {
			sum := v + nums[left] + nums[right]
			checkAndUpdate(sum)
			if sum == target {
				return target
			} else if sum < target {
				left++
			} else {
				right--
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
