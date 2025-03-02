package thread

import (
	"fmt"
	"runtime/debug"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		println(fmt.Sprintf("recover stack: %s\n%s", p, string(debug.Stack())))
	}
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer Recover()

	fn()
}
