package gorpool

import (
	"fmt"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	p := GetPool()
	p.Start()
	defer p.Destroy()
	var seed int32
	go func() {
		for i := 0; i < 100; i++ {
			p.Submit(CreateTask(func(i ...interface{}) bool {
				fmt.Println("work id:", atomic.AddInt32(&seed, 1))
				tt := reflect.TypeOf(i)
				fmt.Println("got:", i[0], tt)
				return true
			}, i))
		}
	}()
	
	time.Sleep(2 * time.Second)
	
	if seed != 100 {
		t.Errorf("seed expected 100, got:%d", seed)
	}
	
}
