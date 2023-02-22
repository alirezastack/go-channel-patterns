package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	// If the cancel function is not called at least once, there will be a memory leak!
	defer cancel()

	// creating a buffered channel of 1
	ch := make(chan string, 1)

	// a child Goroutine is created to perform some I/O bound work
	go func() {
		// Seed() will randomize Intn() output
		rand.Seed(time.Now().Unix())
		randomIOTime := rand.Intn(200)
		fmt.Println("workload io time:", randomIOTime, "ms")

		// a random Sleep call is made to simulate blocking work that can’t be directly cancelled.
		// That work can take up to 200 milliseconds to finish.
		// There is a 50-millisecond difference between the timeout and the amount of time the work could take.
		time.Sleep(time.Duration(randomIOTime) * time.Millisecond)
		ch <- "data"
	}()

	// the parent Goroutine blocks in a select statement waiting on two signals.
	// if the parent Goroutine receives the timeout signal, it walks away.
	// In this situation, it can’t inform the child Goroutine that it won’t be around to receive its signal.
	// This is why it’s so important for the work channel to be a buffer of 1 (ch channel).
	// The child Goroutine needs to be able to send its signal, either the parent Goroutine is around to receive it or not.
	// If a non-buffered channel is used, the child Goroutine will block forever and become a memory leak.
	select {
	// The first case represents the child Goroutine finishing the work on time and the result being received.
	// That is what I want
	case d := <-ch:
		fmt.Println("work complete", d)
	// The second case represents a timeout from the Context.
	// This means the work didn't finish within the 150 millisecond time limit.
	case <-ctx.Done():
		fmt.Println("work cancelled")
	}

	time.Sleep(time.Second)
	fmt.Println("---------------------------------")
}
