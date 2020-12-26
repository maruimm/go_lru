package main

import (
	"context"
	"fmt"
	"github.com/maruimm/go_lru/lru_cache"
	"strconv"
	"time"
)

func main() {
	ctx, CancelFunc := context.WithCancel(context.Background())

	c := lru_cache.NewCacheSvr(ctx,
		10,
		1024 * time.Second,
		lru_cache.NewStorage(),
		func(val interface{}) uint32 {
			str := val.(string)
			ret, err := strconv.Atoi(str)
			if err != nil {
				return 0
			}
			return uint32(ret)
		})

	for i := 0 ; i < 1024*10; i++ {
		go func(i int) {
			ret, err := c.Get(fmt.Sprintf("%d", i))
			fmt.Printf("key:%d v:%+v, err:%+v\n", i, ret, err)
		}(i)
	}

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