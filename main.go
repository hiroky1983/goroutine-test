package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"time"
)

const bufferSize = 3

// goの並列処理とchannelの学習
func main() {
	// learnGoRoutine()
	// learnTracer()
	// learnChannel()
	// learnChannel2()
	// learnCloseChannel()
	// learnSelect()
	// learnSelectDefault()
	// learnSelectContinuous()
	// noMutex()
	// mutex()
	// rwMutex()
	// atomicfunc()
	// learningContext()
	// learningContext2()
	learnContextWithDeadline()
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

	ch3 := generateCountStream()
	for v := range ch3 {
		fmt.Println(v)
	}

	nCh := make(chan struct{})
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("goroutine %v started\n", i)
			<-nCh
			fmt.Println(i)
		}(i)
	}
	time.Sleep(2 * time.Second)
	close(nCh)
	fmt.Println("closed channel")

	wg.Wait()
	fmt.Println("all goroutine finished")
}

func generateCountStream() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i <= 5; i++ {
			ch <- i
		}
	}()
	return ch
}

func learnSelect() {
	ch1 := make(chan string,1)
	ch2 := make(chan string,1)
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(2)
		go func() {
			defer wg.Done()
			time.Sleep(500 * time.Millisecond)
			ch1 <- "A"
	}()
		go func() {
			defer wg.Done()
			time.Sleep(800 * time.Millisecond)
			ch2 <- "B"
	}()
loop:
	for ch1 != nil || ch2 != nil {
		select {
		case <-ctx.Done():
			fmt.Println("timeout!")
			break loop
		case v := <-ch1:
			fmt.Println(v)
			ch1 = nil
		case v := <-ch2:
			fmt.Println(v)
			ch2 = nil
		}
	}
	wg.Wait()
	fmt.Println("all goroutine finished")
}

func learnSelectDefault() {
	var wg sync.WaitGroup
	ch := make(chan string, bufferSize)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < bufferSize; i++ {
			time.Sleep(1000 * time.Millisecond)
			ch <- "hello"
		}
		close(ch)
	}()
		for i := 0; i < bufferSize; i++ {
			select {
			case m := <-ch:
				fmt.Println(m)
			default:
				fmt.Println("no value")
			}
			time.Sleep(1500 * time.Millisecond)
		}
}

const bufsize = 5

func countProducer(wg *sync.WaitGroup, ch chan<- int, size, sleep int) {
	defer wg.Done()
	defer close(ch)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		ch <- i
	}
}

func countConsumer(ctx context.Context, wg *sync.WaitGroup, ch1 <-chan int, ch2 <-chan int) {
	defer wg.Done()
// loop:
	for ch1 != nil || ch2 != nil {
		select {
			case <-ctx.Done():
				fmt.Println(ctx.Err())
				// break loop
				for ch1 != nil || ch2 != nil {
					select {
					case v, ok := <-ch1:
						if !ok {
							ch1 = nil
							break
						}
						fmt.Printf("ch1 %v\n",v)
					case v, ok := <-ch2:
						if !ok {
							ch2 = nil
							break
						}
						fmt.Printf("ch2 %v\n",v)
					}
				}
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					break
				}
				fmt.Printf("ch1 %v\n",v)
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					break
				}
				fmt.Printf("ch2 %v\n",v)
		}
	}
}

func learnSelectContinuous() {
	ch1 := make(chan int, bufsize)
	ch2 := make(chan int, bufsize)
	var wg sync.WaitGroup
	ctx , cancel := context.WithTimeout(context.Background(), 180 * time.Millisecond)
	defer cancel()

	wg.Add(3)
	go countProducer(&wg, ch1, bufsize, 50)
	go countProducer(&wg, ch2, bufsize, 500)
	go countConsumer(ctx, &wg, ch1, ch2)
	wg.Wait()
	fmt.Println("all goroutine finished")
}

func noMutex() {
	var wg sync.WaitGroup
	var i int
	wg.Add(2)
	go func() {
		defer wg.Done()
		i++
	}()
	go func() {
		defer wg.Done()
		i++
	}()
	// 複数のgoroutineが同じ変数にアクセスしたときに上記だと排他制御がないため、iの値が不定になる
	// 基本fmt.Println(i)の結果は2になるが、たまに1になることがある
	// go run -race main.goで実行するとデッドロックが発生する
	// mutexを使うことで排他制御を行う
	wg.Wait()
	fmt.Println(i)
}

func mutex() {
	var wg sync.WaitGroup
	// データ競合を防ぐためにmutexを使う
	var mu sync.Mutex
	var i int
	wg.Add(2)
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		i++
	}()
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		i++
	}()
	wg.Wait()
	fmt.Println(i)
}

func rwMutex() {
	var wg sync.WaitGroup
	var rwMu sync.RWMutex
	var c int

	wg.Add(4)
	go write(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	wg.Wait()
	fmt.Println("finished")
}

func read (mu *sync.RWMutex, wg *sync.WaitGroup, c *int) {
	defer wg.Done()
	time.Sleep(10 * time.Millisecond)
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("read lock")
	fmt.Println(*c)
	time.Sleep(1 * time.Second)
	fmt.Println("read unlock")
}

func write (mu *sync.RWMutex, wg *sync.WaitGroup, c *int) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("write lock")
	*c +=1
	time.Sleep(1 * time.Second)
	fmt.Println("write unlock")
}

func atomicfunc() {
	var wg sync.WaitGroup
	var c int64

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				atomic.AddInt64(&c, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println(c)
	fmt.Println("finished")
}

func learningContext() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	wg.Add(3)
	go subTask(ctx, &wg, "a")
	go subTask(ctx, &wg, "b")
	go subTask(ctx, &wg, "c")
	wg.Wait()

}

func subTask(ctx context.Context, wg *sync.WaitGroup, id string) {
	defer wg.Done()
	t := time.NewTicker(500 * time.Millisecond)
	select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		case <-t.C:
			t.Stop()
			fmt.Println(id)
		}
}

func learningContext2() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := criticalTask(ctx)
		if err != nil {
			fmt.Printf("critical task cancelled due to %v\n", err)
			cancel()
			return
		}
		fmt.Println("success", v)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := normalTask(ctx)
		if err != nil {
			fmt.Printf("normal task cancelled due to %v\n", err)
			cancel()
			return
		}
		fmt.Println("success", v)
	}()
	wg.Wait()
}

func criticalTask(ctx context.Context) (string, error) {
  ctx, cancel := context.WithTimeout(ctx, 1200 * time.Millisecond)
	t := time.NewTicker(1000 * time.Millisecond)
	defer cancel()
	select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-t.C:
			t.Stop()
		}
		return "A", nil
}

func normalTask(ctx context.Context) (string, error) {
	t := time.NewTicker(3000 * time.Millisecond)
	select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-t.C:
			t.Stop()
	}
	return "B", nil
}

func learnContextWithDeadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20 * time.Millisecond))
	defer cancel()
	ch := subTaskWithDeadline(ctx)
	v, ok := <-ch
	if ok {
		fmt.Println(v)
	}
	fmt.Println("subtask finished")

}

func subTaskWithDeadline(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		deadline, ok := ctx.Deadline()
		if ok {
			if deadline.Sub(time.Now().Add(30 * time.Millisecond)) < 0 {
				fmt.Println("impossible to meet deadline")
				return
			}
		}
		time.Sleep(30 * time.Millisecond)
		ch <- "hello"
	}()
	return ch
}
