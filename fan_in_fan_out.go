package main

import (
	"fmt"
	"math/rand"
	"time"
)

// the fan out/in uses the wait_for_result.go pattern
// the idea of this pattern is to create a Goroutine for each
// individual piece of work that is pending and can be done concurrently.
// You MUST remember, a fan out is dangerous in a running service
// since the number of child Goroutines I create for the fan are a multiplier.
// If I have a service handling 50k requests on 50 thousand Goroutines,
// and I decide to use a fan out pattern of 10 child Goroutines for some requests,
// in a worse case scenario I would be talking 500k Goroutines existing at the same time.
// Depending on the resources those child Goroutines needed,
// I might not have them available at that scale and the back pressure could bring the service down.
func main() {
	// number of child goroutines
	children := 2000
	// A buffered channel of 2000 is constructed, one for each child Goroutine being created.
	// creating a buffered channel since there is only one receiver, and
	// it’s not important to have a guarantee at the signaling level.
	// That will only create extra latency.
	ch := make(chan string, children)

	for c := 0; c < children; c++ {
		// in a loop, 2000 child Goroutines are created, and they are off to do their work.
		go func(child int) {
			// A random sleep is used to simulate the work and the unknown amount of time it takes to get the work done.
			// The key is that the order of the work is undefined, out of order, execution which also changes each time the program runs.
			// If this is not acceptable, I can’t use concurrency
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			ch <- fmt.Sprintf("data-%d", child)
			fmt.Println("child", child, "sent signal.")
		}(c)
	} // Once all the Goroutines are created, the parent Goroutine waits in a receive-loop.

	// the idea is to move the guarantee to know when all the signals have been received.
	// This will reduce the cost of latency from the channels.
	// That will be done with a counter that is decremented for each received signal until it reaches zero.
	for children > 0 {
		// Eventually as data is signaled into the buffered channel,
		// the parent Goroutine will pick up the data and eventually all the work is received.
		data := <-ch
		children--
		fmt.Println("received:", data)
	}

	time.Sleep(time.Second)
	fmt.Println("--------------------------------")
}
