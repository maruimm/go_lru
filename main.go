package main

import (
	"context"
	"fmt"
	"github.com/maruimm/go_lru/lru_cache"
	"time"
)

func main() {
/*	c := lru_cache.NewLruCache()

	for i := 0 ; i < 60; i++ {
		ret, _ :=  c.Get(i)
		fmt.Printf("%+v\n",ret)
	}

	for i := 0 ; i < 60; i++ {
		ret, _ :=  c.Get(i)
		fmt.Printf("cache:%+v\n",ret)
	}*/

	ctx, CancelFunc := context.WithCancel(context.Background())

	c := lru_cache.NewCacheSvr(ctx,
		1024,
		1024 * time.Second,
		lru_cache.NewStorage())

	for i := 0 ; i < 40; i++ {
		ret, err :=  c.Get(fmt.Sprintf("%d",i))
		fmt.Printf("key:%d v:%+v, err:%+v\n",i,ret, err)
	}

	for i := 0 ; i < 10; i++ {
		ret, err :=  c.Get(fmt.Sprintf("%d",i))
		fmt.Printf("cache2 key:%d v:%+v, err:%+v\n",i,ret, err)
	}

	for i := 0 ; i < 30; i++ {
		ret, err :=  c.Get(fmt.Sprintf("%d",i))
		fmt.Printf("cache3 key:%d v:%+v, err:%+v\n",i,ret, err)
	}
	CancelFunc()

	time.Sleep(10*time.Second)

}