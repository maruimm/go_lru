package main

import (
	"context"
	"fmt"
	"github.com/maruimm/go_lru/lru_cache"
	"time"
)

func main() {
	ctx, CancelFunc := context.WithCancel(context.Background())

	c := lru_cache.NewCacheSvr(ctx,
		10,
		1 * time.Second,
		lru_cache.NewStorage(),
		)

	for i := 0 ; i < 1024; i++ {
		go func() {
			for {
				for i := 0; i < 1024; i++ {
					ret, err := c.Get(fmt.Sprintf("%d", i))
					fmt.Printf("key:%d v:%+v, err:%+v\n", i, ret, err)
					time.Sleep(1*time.Second)
				}
			}
		}()
	}

	fmt.Printf("cache.....\n")

	for i := 0 ; i < 1000; i++ {
		ret, err :=  c.Get(fmt.Sprintf("%d",i))
		fmt.Printf("cache2 key:%d v:%+v, err:%+v\n",i,ret, err)
	}

	for i := 0 ; i < 3000; i++ {
		ret, err :=  c.Get(fmt.Sprintf("%d",i))
		fmt.Printf("cache3 key:%d v:%+v, err:%+v\n",i,ret, err)
	}
	CancelFunc()

	time.Sleep(10*time.Second)

}