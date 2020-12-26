package lru_cache

import (
	"math/rand"
	"time"
)

type RemoteStorage interface {
	Get(interface{}) (interface{} ,error)
}


type myStorage struct {
	channel Channels
}

func (s *myStorage) Get(key interface{}) (interface{} ,error) {
	time.Sleep(1*time.Microsecond)
	return rand.Int(),nil
}

func NewStorage() RemoteStorage {
	taskChan := make(chan interface{}, 1024)
	resultChan := make(chan interface{}, 1024)
	channel :=  Channels{
		taskChan:taskChan,
		resultChan:resultChan,
	}
	return &myStorage{
		channel:channel,
	}
}
