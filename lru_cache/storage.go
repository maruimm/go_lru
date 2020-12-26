package lru_cache

import (
	"errors"
	"math/rand"
)

type RemoteStorage interface {
	Get(interface{}) (interface{} ,error)
}


type myStorage struct {
	channel Channels
}

func (s *myStorage) Get(key interface{}) (interface{} ,error) {

	select {
		case s.channel.taskChan <- key:
		case result := <- s.channel.resultChan:
			return result, nil
		default:
			return "", errors.New("get storage failed")
	}

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
