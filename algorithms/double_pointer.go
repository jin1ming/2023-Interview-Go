package algorithms

// 上一次学习：2022.4.7，完成

import (
	"math"
	"sort"
)

/***** 两数之和 *****/
// 思路：
// - 通过额外的下标数组 indexes 来记录原始位置，并按对应 nums 值对 indexes 排序；
// - 随后在有序的“值序列”上使用左右指针收敛查找目标和；
// - 命中时返回原数组中的两个下标，未命中返回 nil。
func twoSum(nums []int, target int) []int {
	left, right := 0, len(nums)-1
	indexes := make([]int, len(nums))
	for i := range indexes {
		indexes[i] = i
	}
	// 对下标按数值从小到大排序，保证双指针可用
	sort.Slice(indexes, func(i, j int) bool {
		return nums[indexes[i]] < nums[indexes[j]]
	})
	for left < right {
		sum := nums[indexes[left]] + nums[indexes[right]]
		if sum == target {
			return []int{indexes[left], indexes[right]}
		} else if sum < target {
			// 当前和过小，左指针右移以增大和
			left++
		} else {
			// 当前和过大，右指针左移以减小和
			right--
		}
	}
	return nil
}

/***** 三数之和 *****/
// 给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，
// 使得 a + b + c = 0 ？请你找出所有和为 0 且不重复的三元组。
// 思路（经典双指针）：
// - 先整体排序，固定第一个数 nums[k]，在区间 (k, end] 上用左右指针找两数之和为 -nums[k]；
// - 为避免重复：当 nums[k] 与前一个相同则跳过；在命中三元组后，分别跳过左右侧的重复值；
// - 因为数组有序，若 nums[k] 已经 > 0，则后续三数之和不可能为 0，可直接返回。
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
			// 有序数组中首元素已大于 0，后续三数之和必然 > 0
			return res
		case k > 0 && nums[k-1] == nums[k]:
			// 跳过相同的起点，避免重复三元组
			continue
		default:
			ptrLeft = k + 1
			ptrRight = len(nums) - 1
			for ptrLeft < ptrRight {
				sum := nums[k] + nums[ptrLeft] + nums[ptrRight]
				if sum == 0 {
					r := []int{nums[k], nums[ptrLeft], nums[ptrRight]}
					res = append(res, r)
					// 跳过左侧重复值
					for ptrLeft < ptrRight && nums[ptrLeft] == nums[ptrLeft+1] {
						ptrLeft += 1
					}
					// 跳过右侧重复值
					for ptrLeft < ptrRight && nums[ptrRight] == nums[ptrRight-1] {
						ptrRight -= 1
					}
				}

				if sum > 0 {
					// 和偏大，右指针左移
					ptrRight -= 1
				} else {
					// 和偏小或等于 0 后需要继续尝试其他解，左指针右移
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
// 思路：
// - 先排序；固定第一个数 nums[k]，在后半段使用左右指针收敛；
// - 每次更新当前最优解 best（与 target 差值更小的和）；
// - 命中精确等于 target 时可直接返回。
func threeSumClosest(nums []int, target int) int {
	if len(nums) < 3 {
		return -1
	}

	sort.Ints(nums) // 必须先排序以使用双指针
	best := math.MaxInt32
	checkAndUpdate := func(value int) {
		if abs(value-target) < abs(best-target) {
			// 使用绝对值比较离 target 的更近程度
			best = value
		}
	}

	for k, v := range nums {
		if k > 0 && v == nums[k-1] {
			// 跳过相同的起点，避免重复计算
			continue
		}
		left, right := k+1, len(nums)-1
		for left < right {
			sum := v + nums[left] + nums[right]
			checkAndUpdate(sum)
			if sum == target {
				return target
			} else if sum < target {
				// 和偏小，左指针右移使和增大
				left++
			} else {
				// 和偏大，右指针左移使和减小
				right--
			}
		}
	}

	return best
}

// 计算整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}
