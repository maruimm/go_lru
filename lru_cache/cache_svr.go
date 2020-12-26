package lru_cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)




type CacheSvr struct {
	localCache LocalCache
	ctx context.Context
}

func (pCache *CacheSvr) Get(key interface{}) (interface{}, error) {
	resultChan := make(chan interface{})
	go func(ctx context.Context) {
		defer close(resultChan)
		val, err := pCache.localCache.Get(key)
		if err != nil {
			fmt.Printf("svr local cache failed...\n")
			return
		}
		for {
			select {
				case <- ctx.Done():
					fmt.Printf("svr ctx done\n")
					return
				case resultChan <- val:
					return
			}
		}
	}(pCache.ctx)
	result , ok:= <- resultChan
	if !ok {
		return "", errors.New("result has closed")
	}
	return result, nil
}

func NewCacheSvr(ctx context.Context,
	capNumber int,
	cacheExpire time.Duration,
	storage RemoteStorage,
	) LocalCache {

	return &CacheSvr{
		localCache: NewLruCache(
			capNumber,
			cacheExpire,
			storage,
		),
		ctx: ctx,
	}
}
