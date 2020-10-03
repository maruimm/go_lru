package main

import (
	"fmt"
	"sync"
)

func main() {

	sendChannel := make(chan int)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("go goroutine start...\n")
		for item := range sendChannel {
			fmt.Printf("item:%v\n", item)
		}
		fmt.Printf("go goroutine exit...\n")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("go goroutine2 start...\n")
		for item := range sendChannel {
			fmt.Printf("item2:%v\n", item)
		}
		fmt.Printf("go goroutine2 exit...\n")
	}()


	sendChannel <- 1234
	close(sendChannel)

	wg.Wait()

}
