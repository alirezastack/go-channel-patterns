package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// The retry timeout pattern is great when I have to ping something (like a database)
// which might fail, but I don’t want to fail immediately.
// I want to retry for a specified amount of time before I fail.
func main() {
	fmt.Println("call retryTimeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	retryTimeout(ctx, time.Second*10, DBHealthCheck)
}

// retryTimeout The function takes a context for the amount of time
// the function should attempt to perform work unsuccessfully
// It also takes a retryInterval that specifies how long to wait between attempts,
// and finally a check() function to execute.
// This function is coded by the caller for the specific work (like pinging the database) that needs to be performed and could fail.
func retryTimeout(ctx context.Context, retryInterval time.Duration, check func(ctx context.Context) error) {
	// The core of the function runs in an endless loop
	for {
		fmt.Println("performing user check call...")
		// The first step in the loop is to run the check function passing in the context
		// so the caller’s function can also respect the context.
		// If that doesn't fail, the function returns that life is good.
		// If it fails, the code goes on to the next step.
		if err := check(ctx); err == nil {
			fmt.Println("work finished successfully")
			return
		}

		fmt.Println("check if timeout has expired")
		// the context is checked to see if the amount of time given has expired.
		// If it has, the function returns the timeout error,
		// else it continues to the next step which is to create a timer value.
		if ctx.Err() != nil {
			fmt.Println("time expired 1:", ctx.Err())
			return
		}

		fmt.Println("wait", retryInterval, "before trying again")
		// The time value is set to the retryInterval.
		// The timer could be created above the for loop and reused,
		// which would be good if this function was going to be running a lot.
		// To simplify the code, a new timer is created every time
		t := time.NewTimer(retryInterval)

		// The last step is to block on a select statement waiting to receive one of the two signals
		select {
		case <-ctx.Done(): // if context expires
			fmt.Println("time expired 2:", ctx.Err())
			t.Stop()
			return
		case <-t.C: // if retryInterval expires, the loop is restarted and the process runs again.
			fmt.Println("retry again")
		}
	}
}

func DBHealthCheck(ctx context.Context) error {
	// assume ping is always ok
	//return nil

	// if error happens in DB health check
	return errors.New("ping error")
}
