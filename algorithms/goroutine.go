package algorithms

import (
	"fmt"
	"strconv"
	"sync"
)

// main 函数示例（注释掉以避免与其他文件冲突）
// func main() {
// 	gets()
// }

/***交替并行打印hello-world***/
// 使用两个互斥锁实现两个 goroutine 交替打印 "hello" 和 "world"
// 算法思路：使用互斥锁控制执行顺序，通过锁的获取和释放实现交替执行
func helloWorld() {
	mu1, mu2 := &sync.Mutex{}, &sync.Mutex{} // 两个互斥锁
	mu2.Lock()                               // 先锁定 mu2，让第二个 goroutine 等待
	wg := sync.WaitGroup{}
	wg.Add(2)
	// goroutine 1: 打印 "hello"
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			mu1.Lock() // 获取 mu1，开始打印
			fmt.Print("hello")
			mu2.Unlock() // 释放 mu2，让第二个 goroutine 可以打印
		}
	}()
	// goroutine 2: 打印 "world"
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			mu2.Lock() // 等待 mu2 被释放（由第一个 goroutine 释放）
			fmt.Println("world")
			mu1.Unlock() // 释放 mu1，让第一个 goroutine 可以继续
		}
	}()
	wg.Wait() // 等待两个 goroutine 完成
}

/***多线程爬虫***/
// 使用 goroutine 并发处理多个 URL，通过带缓冲的 channel 控制并发数量
// 算法思路：使用信号量模式（Semaphore Pattern）限制并发数
// chBuffer 作为信号量，容量为 4，表示最多同时有 4 个 goroutine 在执行
func gets() {
	//var urls = []string{
	//	"http://www.golang.org/",
	//	"http://www.google.com/",
	//}
	var urls []string
	for i := 0; i < 100; i++ {
		urls = append(urls, strconv.Itoa(i))
	}
	results := make([]string, 100)
	// 带缓冲的 channel，容量为 4，用于控制并发数（信号量）
	chBuffer := make(chan struct{}, 4)
	wg := sync.WaitGroup{}
	for i := range urls {
		chBuffer <- struct{}{} // 获取信号量（如果 channel 已满，这里会阻塞）
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-chBuffer // 释放信号量
				wg.Done()
			}()
			// 实际爬虫代码示例：
			//resp, err := http.Get(url)
			//if err == nil {
			//	resp.Body.Close()
			//}
			//return err
			results[i] = urls[i]
		}(i)
	}
	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println(results)
}
