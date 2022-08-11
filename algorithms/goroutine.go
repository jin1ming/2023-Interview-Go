package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	gets()
}

/***交替并行打印hello-world***/
func helloWorld() {
	mu1, mu2 := &sync.Mutex{}, &sync.Mutex{}
	mu2.Lock()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			mu1.Lock()
			fmt.Print("hello")
			mu2.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			mu2.Lock()
			fmt.Println("world")
			mu1.Unlock()
		}
	}()
	wg.Wait()
}

/***多线程爬虫***/
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
	chBuffer := make(chan struct{}, 4)
	wg := sync.WaitGroup{}
	for i := range urls {
		chBuffer <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-chBuffer
				wg.Done()
			}()
			//resp, err := http.Get(url)
			//if err == nil {
			//	resp.Body.Close()
			//}
			//return err
			results[i] = urls[i]
		}(i)
	}
	wg.Wait()
	fmt.Println(results)
}

/***LRU Cache***/
func lru() {

}
