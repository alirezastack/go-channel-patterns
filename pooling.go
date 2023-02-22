package main

import (
	"fmt"
	"runtime"
	"time"
)

// The pooling pattern uses the wait_for_task.go pattern.
// The pooling pattern allows us to manage resource usage across a well-defined number of Goroutines.
// In Go, pooling is not needed for efficiency in CPU processing like at the operating system.
// It’s more important for efficiency in RESOURCE USAGE.
func main() {
	// It’s critically important that an unbuffered channel is used
	// because without the guarantee at the signaling level,
	// I can’t perform timeouts and cancellation on the send if needed at a later time.
	ch := make(chan string)

	// This line decides the number of child Goroutines the pool will contain.
	// The call to runtime.GOMAXPROCS is important in that it queries the runtime (when passing 0 as a parameter)
	// to the number of threads that exist for running Goroutines.
	// The number should always equal the number of cores/hardware_threads that are available to the program.
	// It represents the amount of CPU capacity available to the program.
	// When the size of the pool isn’t obvious, start with this number as a baseline.
	// It won’t be uncommon for this number to provide a reasonable performance benchmark.
	g := runtime.GOMAXPROCS(0)

	// The for loop creates the pool of child Goroutines where each child Goroutine sits
	// in a blocking receive call using the for/range mechanics for a channel.
	for c := 0; c < g; c++ {
		// A group of child Goroutines are created to service the same channel.
		// There is efficiency in this because the size of the pool dictates the amount of concurrent work happening at the same time.
		// If I have a pool of 16 Goroutines, that could represent 16 files being opened at any given time
		// or the amount of memory needed for 16 Goroutines to perform their work.
		go func(child int) {
			// The for range helps to minimize the amount of code I would otherwise need to receive a signal
			// and then shutdown once the channel is closed.
			// for { // infinite loop
			//   d, wd := <-ch
			//   if !wd {
			//      break
			//   }
			// }
			for d := range ch {
				// It must not matter which of the child Goroutines in the pool are chosen to receive a signal.
				// Depending on the amount of work being signaled,
				// it could be the same child Goroutines over and over while others are never selected.
				fmt.Printf("child %d: recv'd signal: %s\n", child, d)
			}
			fmt.Printf("child %d: recv'd shutdown signal\n", child)
		}(c)
	}

	const work = 100
	for w := 0; w < work; w++ {
		// As the channel ch is "unbuffered", the sender blocks until the receiver has received the value
		ch <- "data"
		fmt.Println("parent: sent signal:", w)
	}

	// close(ch) will cause the for loops to terminate and stop the program.
	// If the channel being used was a buffered channel,
	// data would flush out of the buffer first before the child Goroutines would receive the close signal.
	close(ch)
	fmt.Println("parent: sent shutdown signal")

	time.Sleep(time.Second)
}
