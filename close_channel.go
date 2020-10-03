package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 10 * time.Second)

	go func(ctx context.Context) {
		ELOOP:
		for {
			select {
				case <-ctx.Done(): {
					fmt.Printf("done reason:%s\n",ctx.Err())
					break ELOOP
				}
			}
		}
	}(ctx)

	ctx2, _ := context.WithTimeout(ctx, 20 * time.Second)
	go func(ctx context.Context) {
	ELOOP:
		for {
			select {
			case <-ctx.Done(): {
				fmt.Printf("done reason 2:%s\n",ctx.Err())
				break ELOOP
			}
			}
		}
	}(ctx2)
	time.Sleep(5 * time.Second)
	cancelFunc()
	time.Sleep(5 * time.Second)
}
