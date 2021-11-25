package panicHandler

import "fmt"

func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Printf("recovered panic: %v\n", r)
	}
}
