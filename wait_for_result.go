package main

import (
	"math/rand"
	"time"
)

func main() {
	// type: unbuffered channel with open state
	// guarantees at the signaling level (means the sending Goroutine wants a guarantee that the signal being sent has been received)
	ch := make(chan string)

	// sending goroutine
	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		ch <- "data payload"
		println("child: sent signal")
	}()

	// this blocks the main goroutine
	// parent Goroutine waits to receive a signal with string data
	// the amount of time the parent Goroutine will need to wait is unknown. It’s the unknown latency cost of this type of channel
	d := <-ch
	println("parent: recvd signal:", d)

	time.Sleep(time.Second)
	println("---------------------------------------")
}
