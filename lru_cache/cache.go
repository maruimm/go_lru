package lru_cache

import (
	"container/list"
	"errors"
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
	ll *list.List //list里面存key的列表,
	pool map[interface{}]entry
	capNumber int
	storage RemoteStorage
	cacheExpire time.Duration
}


func (pCache* LruCache)insert(key interface{}, val interface{}) error {
	el := pCache.ll.PushFront(key)
	pCache.pool[key] = entry{
		updateTime: time.Now(),
		data: val,
		el:el,
	}
	return nil
}

func (pCache* LruCache)update(key interface{}, val entry) error {

	if time.Now().Sub(val.updateTime) > pCache.cacheExpire {
		//远端访问,过期了
		val, err := pCache.storage.Get(key)
		if err != nil {
			return err
		}
		return pCache.insert(key, val)
	} else { //还未过期之前移动到list头部
		pCache.ll.MoveToFront(val.el)
	}
	return nil
}

func (pCache* LruCache)trim() error{ //超过最大容量时删除最老的
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

func (pCache* LruCache)get(key interface{}) (interface{} ,error) {

	val,ok := pCache.pool[key]
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

func (pCache* LruCache)Get(key interface{}) (interface{} ,error) {
	return pCache.get(key)
}


func NewLruCache(capNumber int,
	cacheExpire time.Duration,
	storage RemoteStorage) LocalCache{
	return &LruCache{
		ll :     list.New(),
		pool :   make(map[interface{}]entry),
		capNumber :    capNumber,
		storage: storage,
		cacheExpire:cacheExpire,
	}
}
