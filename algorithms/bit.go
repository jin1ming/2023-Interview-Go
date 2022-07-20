package algorithms

import "math/bits"

// 上一次学习：2022.7.20

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
	// 先对x、y进行异或，找出二进制下不同的值
	// 不同的即为1， 统计1出现的次数即可
	return bits.OnesCount(uint(x ^ y))
}

/***** 两整数之和 *****/
// 不允许使用 +、-
func getSum(a int, b int) int {
	// 异或+与运算:时间复杂度O(logSum) | 空间复杂度O(1)
	// a + b 的问题拆分为 (a 和 b 的无进位结果) + (a 和 b 的进位结果)
	// 最后这个 + 再用无进位结果（异或）来表示
	// a, b = a^b, (a&b)<<1
	for b != 0 {
		a, b = a^b, (a&b)<<1
	}
	return a
}
