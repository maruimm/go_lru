package main

import (
	"context"
	"fmt"
	"github.com/maruimm/go_lru/lru_cache"
	"sync"
	"time"
)

func main() {
	ctx, CancelFunc := context.WithCancel(context.Background())

	c := lru_cache.NewCacheSvr(ctx,
		10,
		1024 * time.Second,
		lru_cache.NewStorage(),
		)

	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0 ; i < 20; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				ret, err := c.Get(fmt.Sprintf("%d", i))
				fmt.Printf("key:%d v:%+v, err:%+v\n", i, ret, err)
				//time.Sleep(1*time.Second)
			}
		}()
	}
	wg.Wait()

	go func() {
		for i := 10; i < 20; i++ {
			ret, err := c.Get(fmt.Sprintf("%d", i))
			fmt.Printf("cache key:%d v:%+v, err:%+v\n", i, ret, err)
			//time.Sleep(1*time.Second)
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
	CancelFunc()

}