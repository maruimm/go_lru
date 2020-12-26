package lru_cache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

type LocalCache interface { //仅暴露一个Get接口就可以了
	Get(interface{}) (interface{} ,error)
}


type entry struct{
	updateTime time.Time
	data      interface{}
	el* list.Element
}


type LruCache struct {
	lk *sync.RWMutex
	ll *list.List //list里面存key的列表,
	pool map[interface{}]entry
	capNumber int
	storage RemoteStorage
	cacheExpire time.Duration
	rchan chan interface{}
}


func (pCache* LruCache)insert(key interface{}, val interface{}) error {
	pCache.lk.Lock()
	defer pCache.lk.Unlock()
	el := pCache.ll.PushFront(key)
	pCache.pool[key] = entry{
		updateTime: time.Now(),
		data: val,
		el:el,
	}
	return nil
}

func (pCache *LruCache) refreshCoroutineStart() {
	go func() {
		for {
			select {
				case key,ok :=<- pCache.rchan: {
					if !ok {
						return
					}
					val, err := pCache.storage.Get(key)
					if err != nil {
						return
					}
					_ = pCache.insert(key, val)
				}
			}
		}
	}()
}


func (pCache* LruCache)needRefresh(key interface{}) error {
	select {
	case pCache.rchan <- key:
	case <-time.After(5*time.Microsecond): {
		return errors.New("needRefresh time out")
	}
	}
	return nil
}


func (pCache* LruCache)update(key interface{}, val entry) error {

	if time.Now().Sub(val.updateTime) > pCache.cacheExpire {
		err := pCache.needRefresh(key)
		if err != nil {
			fmt.Printf("error:%+v\n",err)
			return err
		}
	} else { //还未过期之前移动到list头部
		pCache.lk.Lock()
		defer pCache.lk.Unlock()
		pCache.ll.MoveToFront(val.el)
	}
	return nil
}

func (pCache* LruCache)trim() error{ //超过最大容量时删除最老的
	pCache.lk.Lock()
	defer pCache.lk.Unlock()
	if pCache.ll.Len() > pCache.capNumber {
		trimCount := pCache.ll.Len() - pCache.capNumber
		for i := 0; i < trimCount; i++ {
			el := pCache.ll.Back()
			key := el.Value
			delete(pCache.pool,key)
			pCache.ll.Remove(el)
		}
	}
	return nil
}

func (pCache* LruCache)get(key interface{}) (entry ,bool) {

	pCache.lk.RLock()
	defer pCache.lk.RUnlock()
	val,ok := pCache.pool[key]
	return val, ok
}

func (pCache* LruCache)Get(key interface{}) (interface{} ,error) {

	val, ok := pCache.get(key)
	if !ok  {//没找到
		originVal, err := pCache.storage.Get(key)
		if err != nil {
			return "" , errors.New("no data")
		}
		_ = pCache.insert(key, originVal)
		_ = pCache.trim()
		return originVal,nil
	} else {
		_ = pCache.update(key , val)
	}
	return val.data, nil //找到了

}


func NewLruCache(capNumber int,
	cacheExpire time.Duration,
	storage RemoteStorage) LocalCache{
	lru := &LruCache{
		ll :     list.New(),
		pool :   make(map[interface{}]entry),
		capNumber :    capNumber,
		storage: storage,
		cacheExpire:cacheExpire,
		lk: new(sync.RWMutex),
		rchan: make(chan interface{}, 1024),
	}
	lru.refreshCoroutineStart()
	return lru
}
