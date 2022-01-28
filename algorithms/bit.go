package algorithms

/***** 只出现一次的数字 *****/
// 给定一个非空整数数组，除了某个元素只出现一次以外，
// 其余每个元素均出现两次。找出那个只出现了一次的元素。
func singleNumber(nums []int) int {
	single := 0
	for _, num := range nums {
		// 出现两次的数字会抵消为 0
		// 最后剩的就是最终的结果
		single ^= num
	}
	return single
}

/***** 汉明距离 *****/
func hammingDistance(x int, y int) int {
	xor := x ^ y
	res := 0
	for xor != 0 {
		if xor%2 == 1 {
			res++
		}
		xor >>= 1
	}
	return res
}

/***** 两整数之和 *****/
// 不允许使用 +、-
func getSum(a, b int) int {
	for a != 0 {
		temp := a ^ b
		a = (a & b) << 1
		b = temp
	}
	return b
}
