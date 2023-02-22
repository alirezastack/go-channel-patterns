package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// The bounded work pooling pattern uses a pool of Goroutines
// to perform a FIXED amount of known work
func main() {
	// Right from the start, the function defines 2000 arbitrary pieces of work to perform.
	work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}

	// Then the GOMAXPROCS function is used to define the number of child Goroutines to use in the pool.
	g := runtime.GOMAXPROCS(0)

	// A WaitGroup is constructed to make sure the parent Goroutine can be told to wait until all 2000 pieces of work are completed.
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	// a pool of child Goroutines is created in the loop
	for c := 0; c < g; c++ {
		go func(child int) {
			// One change is the call to Done using defer() when each of the child Goroutines in the pool eventually terminate.
			// This will happen when all the work is completed and
			// this is how the pool will report back to the parent Goroutine they are aware they are not needed any longer.
			defer wg.Done()
			// they all wait on a receive-call using the for-range mechanics.
			for wrk := range ch {
				fmt.Println("child", child, "received work", wrk)
			}
			fmt.Println("child", child, "received shutdown signal")
		}(c)
	}

	// After the creation of the pool of child Goroutines,
	// A loop is executed by the parent Goroutine to start signaling work into the pool.
	for _, wrk := range work {
		ch <- wrk
	}

	// Once the last piece of work is signaled, the channel is closed.
	// Each of the child Goroutines will receive the closed signal once the signals in the buffer are emptied.
	close(ch)
	wg.Wait()

	time.Sleep(time.Second)
}
