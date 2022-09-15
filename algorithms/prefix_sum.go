package algorithms

/***** 航班预订统计 *****/
// 核心思想：差分数组
// 这里有 n 个航班，它们分别从 1 到 n 进行编号。
// 有一份航班预订表 bookings ，表中第 i 条预订记录 bookings[i] = [firsti, lasti, seatsi]
// 意味着在从 firsti 到 lasti （包含 firsti 和 lasti ）的 每个航班 上预订了 seatsi 个座位。
// 请你返回一个长度为 n 的数组 answer，里面的元素是每个航班预定的座位总数。
func corpFlightBookings(bookings [][]int, n int) []int {
	nums := make([]int, n+1)
	for _, book := range bookings {
		nums[book[0]] += book[2]
		if book[1] != n {
			nums[book[1]+1] -= book[2]
		}
	}
	for i := 2; i <= n; i++ {
		nums[i] += nums[i-1]
	}
	return nums[1:]
}

type Dis struct {
	left     int
	right    int
	leftNum  int
	rightNum int
}

/***** 相同元素的间隔之和 *****/
// 给你一个下标从 0 开始、由 n 个整数组成的数组 arr 。
// arr 中两个元素的 间隔 定义为它们下标之间的 绝对差 。更正式地，arr[i] 和 arr[j] 之间的间隔是 |i - j| 。
// 返回一个长度为 n 的数组 intervals ，其中 intervals[i] 是 arr[i] 和 arr 中每个相同元素（与 arr[i] 的值相同）的 间隔之和 。
// 注意：|x| 是 x 的绝对值。
func getDistances(arr []int) []int64 {
	closedDis := make([]Dis, len(arr)) // 为每个点创建个记录左右前缀和的结构
	walked := make(map[int]int)        // 上一个拥有该值的坐标
	for k, v := range arr {
		// 从左到右扫描，记录每个节点左侧的必要信息
		if _, ok := walked[v]; ok {
			// 记录与左边最近相同值节点的距离
			closedDis[k].left = k - walked[v]
			// 记录左边有多少个有相同值得点
			closedDis[k].leftNum = closedDis[walked[v]].leftNum + 1
		}
		// 最近经过v是在k处
		walked[v] = k
	}
	// 清除walk记录
	for k := range walked {
		delete(walked, k)
	}
	for k := len(arr) - 1; k >= 0; k-- {
		// 从右到左扫描，记录每个节点右侧的必要信息
		v := arr[k]
		if _, ok := walked[v]; ok {
			// 记录与右边最近相同值节点的距离
			closedDis[k].right = walked[v] - k
			// 记录右左边有多少个有相同值得点
			closedDis[k].rightNum = closedDis[walked[v]].rightNum + 1
		}
		walked[v] = k
	}

	res := make([]int64, len(arr))
	for i := range arr {
		// 从左向右遍历，建立左侧的前缀和（距离和）
		closeLD := closedDis[i].left
		if closeLD != 0 {
			closedDis[i].left = closedDis[i-closeLD].left + (closedDis[i-closeLD].leftNum+1)*closeLD
		}
		res[i] += int64(closedDis[i].left)
	}

	for i := len(arr) - 1; i >= 0; i-- {
		// 从右向左遍历，建立右侧的前缀和（距离和）
		closeRD := closedDis[i].right
		if closeRD != 0 {
			closedDis[i].right = closedDis[i+closeRD].right + (closedDis[i+closeRD].rightNum+1)*closeRD
		}
		res[i] += int64(closedDis[i].right)
	}
	return res
}
