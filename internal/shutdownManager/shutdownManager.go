package shutdownManager

import "sync"

var (
	mut sync.Mutex
	funcs = make([]func() error, 0)
)

func Initiate() {
	// TODO
}

func Register(f func() error) {
	mut.Lock()
	defer mut.Unlock()

	funcs = append(funcs, f)
}
