package algorithms

/***** Pow(x, n) *****/
func myPow(x float64, n int) float64 {
	if n >= 0 {
		return quickMul(x, n)
	}
	return 1.0 / quickMul(x, -n)
}

func quickMul(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	y := quickMul(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return y * y * x
}

/***** 跳跃游戏 II *****/
// 给定一个非负整数数组，你最初位于数组的第一个位置。
// 数组中的每个元素代表你在该位置可以跳跃的最大长度。
// 你的目标是使用最少地跳跃次数到达数组的最后一个位置。
// 假设你总是可以到达数组的最后一个位置。
func jump(nums []int) int {
	length := len(nums)
	end := 0
	maxPosition := 0
	steps := 0
	for i := 0; i < length-1; i++ {
		maxPosition = max(maxPosition, i+nums[i])
		// 当前可到达最远位置
		if i == end {
			// 已经到达可走的最远位置
			end = maxPosition
			steps++
		}
	}
	return steps
}

/***** 计数质数 *****/
// 统计所有小于非负整数 n 的质数的数量。
func countPrimes(n int) (cnt int) {
	isPrime := make([]bool, n)
	for i := range isPrime {
		isPrime[i] = true
	}
	for i := 2; i < n; i++ {
		if isPrime[i] {
			cnt++
			for j := 2 * i; j < n; j += i {
				isPrime[j] = false
			}
		}
	}
	return
}

/***** 用 Rand7() 实现 Rand10() *****/
// 方式一
func rand7() int {
	return 0
}
func rand10() int {
	for {
		row := rand7()
		col := rand7()
		idx := (row-1)*7 + col
		if idx <= 40 {
			return 1 + (idx-1)%10
		}
	}
}

// 方式二
func rand10B() int {
	a := rand5()
	b := rand2()
	if b == 1 {
		return a
	} else {
		return 5 + rand5()
	}
}
func rand2() int {
	t := rand7()
	for t == 7 {
		t = rand7()
	}
	return t % 2
}
func rand5() int {
	t := rand7()
	for t > 5 {
		t = rand7()
	}
	return t
}
