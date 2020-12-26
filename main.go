package main

import (
	"fmt"
	"github.com/maruimm/go_lru/lru_cache"
)

func main() {
	c := lru_cache.NewLruCache()

	for i := 0 ; i < 60; i++ {
		ret, _ :=  c.Get(i)
		fmt.Printf("%+v\n",ret)
	}

	for i := 0 ; i < 60; i++ {
		ret, _ :=  c.Get(i)
		fmt.Printf("cache:%+v\n",ret)
	}

}