package lru_cache

import (
	"math/rand"
	"time"
)

type RemoteStorage interface {
	Get(interface{}) (interface{} ,error)
}


type myStorage struct {

}

func (s *myStorage) Get(key interface{}) (interface{} ,error) {
	time.Sleep(1*time.Second)
	return rand.Int(),nil
}

func NewStorage() RemoteStorage {

	return &myStorage{

	}
}
