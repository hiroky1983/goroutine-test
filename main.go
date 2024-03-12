package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

// goの並列処理とchannelの学習
func main() {
	// learnGoRoutine()
	// learnTracer()
	// learnChannel()
	// learnChannel2()
	learnCloseChannel()
}

func learnGoRoutine() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("goroutine invoked")
	}()
	wg.Wait()
	fmt.Printf("num of goroutine: %d\n", runtime.NumGoroutine())
	fmt.Println("main function finished")
}

func learnTracer() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := trace.Start(f); err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()
	ctx, t := trace.NewTask(context.Background(), "main")
	defer t.End()
	fmt.Printf("The number of logical CPU Cores: %d\n", runtime.NumCPU())
	// task(ctx, "task1")
	// task(ctx, "task2")
	// task(ctx, "task3")
	var wg sync.WaitGroup
	wg.Add(3)
	go cTask(ctx, &wg, "task1")
	go cTask(ctx, &wg, "task2")
	go cTask(ctx, &wg, "task3")
	wg.Wait()
	fmt.Println("main function finished")
}

func task(ctx context.Context, name string) {
	defer trace.StartRegion(ctx, name).End()
	time.Sleep(time.Second)
	fmt.Println(name)
}

func cTask(ctx context.Context,wg *sync.WaitGroup, name string) {
	defer trace.StartRegion(ctx, name).End()
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Println(name)
}

func learnChannel() {
	// 読み取り専用チャネル
	// var ch <-chan int
	// 書き込み専用チャネル
	// var ch chan<- int
	// 双方向チャネル
	// var ch chan int
  // 読み込み、書き込みは矢印の向きで判断する
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- 10
		time.Sleep(500 * time.Millisecond)
	}()
	fmt.Println(<-ch)
	wg.Wait()
}

func learnChannel2() {
	ch := make(chan int)
	go func() {
		fmt.Println(<-ch)
	}()
	ch <- 10
	fmt.Printf("num of goroutine: %d\n", runtime.NumGoroutine())

	ch2 := make(chan int, 1)
	ch2 <- 2
	ch2 <- 3
	fmt.Println(<-ch2)
}

func learnCloseChannel() {
	ch1 := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(<-ch1)
	}()
	ch1 <- 10
	close(ch1)
	v, ok := <-ch1
	// fmt.Printf("%v, %v\n", v, ok)
	wg.Wait()

	// バッファ付きチャネル
	ch2 := make(chan int, 2)
	ch2 <- 1
	ch2 <- 2
	close(ch2)
	v, ok = <-ch2
	// 1, true
	fmt.Printf("%v, %v\n", v, ok)
	v, ok = <-ch2
	// 2, true
	fmt.Printf("%v, %v\n", v, ok)
	v, ok = <-ch2
	// 0, false
	fmt.Printf("%v, %v\n", v, ok)
	// バッファ付きチャネルの場合、closeしてもまだ読み込まれていないチャネルの値がある場合は中身は消えない
}
