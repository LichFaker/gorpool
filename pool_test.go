package gorpool

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	p := GetPool()
	p.Start()
	var seed int32
	go func() {
		for i := 0; i < 100; i++ {
			p.Submit(func() {
				fmt.Println("work id:", atomic.AddInt32(&seed, 1))
			})
		}
	}()
	
	time.Sleep(2 * time.Second)
	
	if seed != 100 {
		t.Errorf("seed expected 100, got:%d", seed)
	}
	
}
