package threads

import (
	"log"
	"runtime/debug"
)

func LogPanic() {
	if x := recover(); x != nil {
		// recovering from a panic; x contains whatever was passed to panic()
		log.Printf("panic: %s\n", debug.Stack())
	}
}
