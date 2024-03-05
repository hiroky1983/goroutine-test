package main

import (
	"fmt"
	"time"
)

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
}