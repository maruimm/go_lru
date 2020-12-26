package lru_cache

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

type KeyToUint func(interface{}) uint32

type CacheSvrItem struct {
	taskChan   <-chan interface{}
	resultChan chan<- interface{}
	localCache LocalCache
}

func (svr *CacheSvrItem) start(ctx context.Context) {
	go func() {
		defer close(svr.resultChan)
	LOOP:
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("svr exit...\n")
				break LOOP
			case key, ok := <-svr.taskChan:
				{
					if !ok {
						fmt.Printf("taskChan exit..\n")
						break LOOP
					}
					val, err := svr.localCache.Get(key)
					if err != nil {
						fmt.Printf("svr local cache failed...\n")
						break
					}
					svr.resultChan <- val
				}
			}
		}
	}()
}

type Channels struct {
	taskChan   chan<- interface{}
	resultChan <-chan interface{}
}

type CacheSvr struct {
	svr           []Channels
	keyToUintFunc KeyToUint
}

func (pCache *CacheSvr) Get(key interface{}) (interface{}, error) {
	index := pCache.keyToUintFunc(key)
	index = index % uint32(len(pCache.svr))
	select {
	case pCache.svr[index].taskChan <- key:
	case <-time.After(5 * time.Microsecond): {
			return "", errors.New("get cache value timeout")
		}
	}
	select {
	case val, ok := <-pCache.svr[index].resultChan:
		{
			if !ok {
				fmt.Printf("resultChan has closed:%+v\n", ok)
				return "", errors.New("resultChan has closed")
			}
			return val, nil
		}
	case <-time.After(5 * time.Microsecond): {
			return "", errors.New("get cache value timeout")
		}
	}
}

func NewCacheSvr(ctx context.Context,
	capNumber int,
	cacheExpire time.Duration,
	storage RemoteStorage) LocalCache {
	var chains []Channels
	numCPU := runtime.NumCPU()
	if numCPU == 0 {
		numCPU = 4
	}
	for i := 0; i < numCPU; i++ {
		taskChan := make(chan interface{})
		resultChan := make(chan interface{})
		svr := &CacheSvrItem{
			localCache: NewLruCache(
				capNumber,
				cacheExpire,
				storage,
			),
			taskChan:   taskChan,
			resultChan: resultChan,
		}
		svr.start(ctx)
		chains = append(chains, Channels{
			taskChan:   taskChan,
			resultChan: resultChan,
		})
	}
	return &CacheSvr{
		svr: chains,
		keyToUintFunc: func(val interface{}) uint32 {
			str := val.(string)
			ret, err := strconv.Atoi(str)
			if err != nil {
				return 0
			}
			return uint32(ret)
		},
	}
}
