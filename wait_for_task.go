package main

import (
	"fmt"
	"math/rand"
	"time"
)

// the wait for task pattern is a foundational pattern used by larger patterns like pooling.”
func main() {
	// An unbuffered channel so there is a guarantee at the signaling level.
	// This is critically important for pooling, so I can add mechanics later if needed to allow for timeouts and cancellation
	ch := make(chan string)

	// a child goroutine waiting for a signal with data to perform work
	// Since the guarantee is at the signaling level, the child Goroutine doesn’t know how long it needs to wait.
	go func() {
		d := <-ch
		fmt.Println("child: recv'd signal:", d)
	}()

	// The parent Goroutine begins to prepare that work and finally signals the work to the child Goroutine
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	ch <- "data"
	fmt.Println("parent: sent signal")

	time.Sleep(time.Second)
}
