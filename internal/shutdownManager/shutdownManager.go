package shutdownManager

import (
	"os"
	"sync"
)

var (
	mut   sync.Mutex
	funcs = make([]func(*sync.WaitGroup), 0, 5)
)

func Initiate() {
	mut.Lock()
	defer mut.Unlock()

	wg := &sync.WaitGroup{}
	wg.Add(len(funcs))
	for _, f := range funcs {
		go f(wg)
	}
	wg.Wait()
	os.Exit(0)
}

func Register(f func(wg *sync.WaitGroup)) {
	mut.Lock()
	funcs = append(funcs, f)
	mut.Unlock()
}
