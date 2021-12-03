package shutdownManager

import "sync"

var (
	mut   sync.Mutex
	funcs = make([]func(), 0, 5)
)

func Initiate() {
	// TODO
}

func Register(f func()) {
	mut.Lock()
	funcs = append(funcs, f)
	mut.Unlock()
}
