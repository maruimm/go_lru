package lru_cache

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

type KeyToUint func(interface{}) uint32

type CacheSvrItem struct {
	taskChan   <-chan interface{}
	resultChan chan<- interface{}
	localCache LocalCache
}

func (svr *CacheSvrItem) start(ctx context.Context) <-chan interface{}{
	started := make(chan interface{})
	go func() {
		defer close(svr.resultChan)
		defer close(started)
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
			case started<- struct{}{}:

			}

		}
	}()
	return started
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
			return "", errors.New("send key timeout")
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
			return "", errors.New("get value timeout")
		}
	}
}

func NewCacheSvr(ctx context.Context,
	capNumber int,
	cacheExpire time.Duration,
	storage RemoteStorage,
	keyToUintFunc KeyToUint) LocalCache {
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
		started := svr.start(ctx)
		<-started
		chains = append(chains, Channels{
			taskChan:   taskChan,
			resultChan: resultChan,
		})
	}
	return &CacheSvr{
		svr: chains,
		keyToUintFunc: keyToUintFunc,
	}
}
