package lru_cache

import (
	"fmt"
	//"math/rand"
	"sync"
	"time"
)

type operatorType int

const (
	_ operatorType = iota
	NeedUpdateCache
	DontNeedUpdateCache
)

type RemoteStorage interface {
	Get(interface{}) (interface{} ,operatorType,error)
}

type Rsp struct {
	Result interface{}
	Err error
}

type call []chan Rsp

type myStorage struct {
	mux *sync.Mutex
	filter map[interface{}]call
}

func(s *myStorage) get(key interface{}) (interface{} ,error) {

	valChan := make(chan Rsp)
	go func () {
		defer close(valChan)
		time.Sleep(1 * time.Second)
		//return rand.Int(),nil
		valChan <- Rsp{
			fmt.Sprintf("%+v-%+v", key.(string), time.Now().Unix()),
			nil,
		}
	}()
	rsp := <- valChan
	return rsp.Result, rsp.Err
}

func (s *myStorage) Get(key interface{}) (interface{} ,operatorType ,error) {
	s.mux.Lock()
	if calls, ok := s.filter[key]; ok {
		call := make(chan Rsp)
		calls = append(calls, call)
		s.filter[key] = calls
		s.mux.Unlock()
		select {
			case ret := <-call:
				fmt.Printf("flight succ:key:%+v,val:%+v\n",key, ret)
				return ret.Result, DontNeedUpdateCache,ret.Err
		}
	} else {
		s.filter[key] = nil
		s.mux.Unlock()
		result, err := s.get(key)
		fmt.Printf("real read from remote:key:%+v,val:%+v\n",key, result)
		s.mux.Lock()
		if calls, ok := s.filter[key]; ok {
			for _,call := range calls {
				call <- Rsp {
					Result:result,
					Err:err,
				}
				close(call)
			}
		}
		delete(s.filter,key)
		s.mux.Unlock()
		return result , NeedUpdateCache, err

	}
}

func NewStorage() RemoteStorage {

	return &myStorage{
		filter:make(map[interface{}]call),
		mux: new(sync.Mutex),
	}
}
