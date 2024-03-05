package main

import (
	"fmt"
	"time"
)

// goの並列処理とchannelの学習
func main() {
	fmt.Println("init")
	times := []int{1, 2, 3, 4, 5}
	for _, t := range times {
		go func() {
			fmt.Println(t)
		}()
	}
	go func () {
		fmt.Println("goroutine起動")
	}()
	fmt.Println("終了します")
	time.Sleep(1 * time.Second)
	fmt.Println("終了")

	// s1()
	s2()
	time.Sleep(2 * time.Second)
}

// result 
// init
// 終了します
// goroutine起動
// 5
// 2
// 3
// 4
// 1
// 終了
// chに2を送信: 1
// fmt.Printf("chに1を送信: %v\n", <-ch) が実行される前にch <- 1が実行されているため、deadlockが発生している為
func s1() {
	ch := make(chan int)
	go func() {
		ch <- 1
		fmt.Printf("chに1を送信: %v\n", <-ch)
	}()
	fmt.Printf("chに2を送信: %v\n", <-ch)

}

// result 
// init
// 終了します
// goroutine起動
// 4
// 5
// 3
// 1
// 2
// 終了
// fatal error: all goroutines are asleep - deadlock!

// goroutine 1 [chan send]:
// main.s2()
// 	/Users/yamadahiroki/myspace/goroutine-test/main.go:57 +0x77
// main.main()
// 	/Users/yamadahiroki/myspace/goroutine-test/main.go:25 +0x185

// goroutine 17 [chan send]:
// main.s2.func1()
// 	/Users/yamadahiroki/myspace/goroutine-test/main.go:54 +0x27
// created by main.s2 in goroutine 1
// 	/Users/yamadahiroki/myspace/goroutine-test/main.go:53 +0x66
// exit status 2
// ch <- 2が実行される前にch <- 1が実行されているため、deadlockが発生している為
func s2() {
	ch := make(chan int)
	go func() {
		ch <- 1
		fmt.Printf("chに1を送信: %v\n", <-ch)
	}()
	ch <- 2
	fmt.Printf("chに2を送信: %v\n", <-ch)
	}


