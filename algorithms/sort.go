package algorithms

import "sort"

/***** 在 D 天内送达包裹的能力 *****/
// 返回能在 days 天内将传送带上的所有包裹送达的船的最低运载能力。
func shipWithinDays(weights []int, days int) int {
	// 确定二分查找左右边界
	left, right := 0, 0
	for _, w := range weights {
		if w > left {
			left = w
		}
		right += w
	}
	return left + sort.Search(right-left, func(x int) bool {
		x += left
		day := 1 // 需要运送的天数
		sum := 0 // 当前这一天已经运送的包裹重量之和
		for _, w := range weights {
			if sum+w > x {
				day++
				sum = 0
			}
			sum += w
		}
		return day <= days
	})
}

/***** 数组中的逆序对 *****/
// 在数组中的两个数字，如果前面一个数字大于后面的数字，则这两个数字组成一个逆序对。
// 输入一个数组，求出这个数组中的逆序对的总数。
// 分解： 待排序的区间为 [l,r]，令 m = (l+r) / 2,
//       我们把 [l,r] 分成 [l,m] 和 [m+1,r]
// 解决： 使用归并排序递归地排序两个子序列
// 合并： 把两个已经排好序的子序列 [l,m] 和 [m+1,r] 合并起来
func reversePairs(nums []int) int {
	return mergeSort(nums, 0, len(nums)-1)
}

func mergeSort(nums []int, start, end int) int {
	if start >= end {
		return 0
	}
	mid := start + (end-start)/2 // 防止start和end相加引起的数组越界
	cnt := mergeSort(nums, start, mid) + mergeSort(nums, mid+1, end)
	// 左右分别是排好序的数组
	// cnt 是返回的逆序对的数量
	var tmp []int
	i, j := start, mid+1
	// i是左边数组的指针，j是右边数组的指针
	for i <= mid && j <= end { // 加判断防止越界
		if nums[i] <= nums[j] {
			tmp = append(tmp, nums[i]) // 将最小的元素放入tmp
			cnt += j - (mid + 1)
			// 当前右边数组被存入 tmp 的数量就是右边有几个元素小于左边数组的当前元素
			i++
		} else {
			tmp = append(tmp, nums[j]) // 将最小的元素放入tmp
			j++
		}
	}
	// 将左边数组剩余的加入
	for ; i <= mid; i++ {
		tmp = append(tmp, nums[i])
		cnt += end - (mid + 1) + 1
		// 右边数组全部被存入 tmp, 说明左边数组剩余元素都比右边数组中所有元素要大
	}
	// 将右边数组剩余的加入
	for ; j <= end; j++ {
		tmp = append(tmp, nums[j])
	}
	// 将排好序的 tmp 拷贝到当前数组片段中
	for i = start; i <= end; i++ {
		nums[i] = tmp[i-start]
	}
	return cnt
}

/***** 快速排序 *****/
// 最坏时间复杂度：O(n^2)
// 最好/平均时间复杂度：O(nlogn)
// 空间复杂度：O(logn)
func sortArray(nums []int) []int {

	var quickSort func(left, right int)
	var findPosition func(left, right int) int

	quickSort = func(left, right int) {
		if left >= right {
			return
		}
		pos := findPosition(left, right)
		quickSort(left, pos-1)
		quickSort(pos+1, right)
	}

	findPosition = func(left, right int) int {
		temp := nums[left]

		for left < right {
			for left < right && nums[right] >= temp {
				right--
			}
			nums[left] = nums[right]
			for left < right && nums[left] <= temp {
				left++
			}
			nums[right] = nums[left]
		}

		nums[left] = temp
		return left
	}

	quickSort(0, len(nums)-1)

	return nums
}

/***** 字典序排数 *****/
// 给定一个整数 n, 返回从 1 到 n 的字典顺序。
// 例如，给定 n =1 3，
// 返回 [1,10,11,12,13,2,3,4,5,6,7,8,9] 。
func lexicalOrder(n int) []int {
	if n <= 0 {
		return nil
	}
	res := make([]int, 0, n)

	var dfs func(int)
	dfs = func(i int) {
		if i > n {
			return
		}
		res = append(res, i)
		if 10*i <= n {
			// 优先10倍扩展
			dfs(10 * i)
		}
		// 当前尾数i是9, i+1已经被添加
		if (i+1)%10 == 0 {
			return
		}
		// 递增
		dfs(i + 1)
	}
	dfs(1)
	return res
}

/***** 数组归置 *****/
// 使得数组负数在左，正数在右，0在中间
func numsSwap(nums []int) {
	if len(nums) < 2 {
		return
	}
	left := 0
	right := len(nums) - 1
	for i, k := range nums {
		if k < 0 {
			nums[i], nums[left] = nums[left], nums[i]
			left++
		} else if k > 0 {
			nums[i], nums[right] = nums[right], nums[i]
			right--
		}
	}
}
