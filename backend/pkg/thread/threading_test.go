package thread

import (
	"testing"
	"time"
)

func TestRecover(t *testing.T) {
	called := false
	cleanup := func() {
		called = true
	}

	Recover(cleanup)

	if !called {
		t.Error("cleanup function not called")
	}
}

func TestGoSafe(t *testing.T) {
	panicCalled := make(chan bool)

	fn := func() {
		panic("test panic")
	}

	go func() {
		RunSafe(fn)
		panicCalled <- true
	}()

	select {
	case <-time.After(time.Second):
		t.Error("panic not recovered")
	case <-panicCalled:
		t.Log("panic recovered")
	}

}

func TestRunSafe(t *testing.T) {
	called := false

	fn := func() {
		called = true
	}

	RunSafe(fn)

	if !called {
		t.Error("function not called")
	}
}
