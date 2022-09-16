package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	sif := make(chan os.Signal, 1)
	signal.Notify(sif, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		listen(ctx)
		wg.Done()
	}()

	f := <-sif
	if f == os.Interrupt {
		cancel()
	}
	wg.Wait()

	fmt.Println("goodbye")
}

func listen(ctx context.Context) {
	res := make(chan bool, 1)
	defer close(res)

	res <- true
	for {
		var num int
		select {
		case <-ctx.Done():
			fmt.Println("\nВыхожу из программы")
			return
		case <-res:
			go func(res chan bool) {
				defer func() {
					res <- true
				}()
				fmt.Println("Введите натуральное число:")
				_, err := fmt.Scan(&num)
				if err != nil {
					return
				}
				if num > 0 {
					square := num * num
					fmt.Println(square)
				}
			}(res)
		}
	}
}
