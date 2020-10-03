package main

import (
	"container/list"
	"context"
	"fmt"
	"sync"
)

type Lruer interface {
	Get(key interface{}) chan<- interface{}
	Close()
}

type Element struct {
	El  *list.Element
	Ext chan<- interface{}
}

type InitFunc func(recvGoroutines int) chan<- interface{}

type myLru struct {
	pool    map[uint32]Element
	lruList *list.List
	maxPool int
	init    InitFunc
}

func (l *myLru) Close() {
	for _, item := range l.pool {
		if item.Ext != nil {
			close(item.Ext)
		}
	}
}

func (l *myLru) Get(key interface{}) chan<- interface{} {

	k := key.(uint32)
	fmt.Printf("get key:%d\n", k)
	if value, ok := l.pool[k]; ok {
		el := value.El
		l.lruList.Remove(el)
		fmt.Printf("in cache:%d cache len:%d\n", k, l.lruList.Len())
		l.lruList.PushFront(el.Value.(uint32))
		return value.Ext
	}

	element := l.lruList.PushFront(k)

	stExt := l.init(1)

	l.pool[k] = Element{
		El:  element,
		Ext: stExt,
	}

	if l.lruList.Len() > l.maxPool {
		if endElement := l.lruList.Back(); endElement != nil {
			fmt.Printf("before del key list len:%d\n", l.lruList.Len())
			l.lruList.Remove(endElement)
			tmpKey := endElement.Value.(uint32)
			fmt.Printf("after del key list len:%d,tmpKey:%d\n", l.lruList.Len(), tmpKey)

			if tmpValue, ok := l.pool[tmpKey]; ok {
				close(tmpValue.Ext)
				delete(l.pool, tmpKey)
			}
		} else {
			fmt.Printf("error key endElement is nil\n")
		}
	}
	var n *list.Element
	fmt.Printf("list:")
	for e := l.lruList.Front(); e != nil; e = n {
		fmt.Printf("%d ", e.Value)
		n = e.Next()
	}
	fmt.Printf("\n")
	return stExt
}

func NewMyLru(maxPool int, myInit InitFunc) Lruer {

	l := &myLru{
		maxPool: maxPool,
		pool:    make(map[uint32]Element, maxPool),
		lruList: list.New(),
		init:    myInit,
	}
	return l
}

func initContext(ctx context.Context, wg *sync.WaitGroup) InitFunc {

	return func(recvGoroutinNumber int) chan<- interface{} {
		sendData := make(chan interface{}, 1024)
		wg.Add(recvGoroutinNumber)
		fmt.Printf("add times:%d\n", recvGoroutinNumber)
		for i := 0; i < recvGoroutinNumber; i++ {
			go func() {
				defer func(wg *sync.WaitGroup) {
					fmt.Printf("defer recv exit \n")
					wg.Done()
				}(wg)
			ELOOP:
				for {
					select {
					case <-ctx.Done():
						{
							fmt.Printf("unlimit loop exit\n")
							break ELOOP
						}
					case item, ok := <-sendData:
						{
							if ok {
								fmt.Printf("recv data:%v\n", item)
							} else {
								fmt.Printf("channel closed\n")
								break ELOOP
							}
						}
					}
				}
			}()
		}
		return sendData
	}
}

func main() {

	ctx := context.Background()

	ctx, cancelFunc := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	myInit := initContext(ctx, &wg)

	lruCache := NewMyLru(3000, myInit)

	for i := 0; i < 1000; i++ {
		//tmp := rand.Int() % 1000
		sendChannel := lruCache.Get(uint32(i))
		sendChannel <- i
	}
	cancelFunc()
	//lruCache.Close()

	wg.Wait()

}
