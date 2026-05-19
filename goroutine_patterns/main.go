package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	// ── Fan-in ────────────────────────────────────────────────────────────────
	// Three producers each send values on their own channel.
	// merge() funnels all three into one channel so the consumer has a single
	// range loop instead of three separate ones.
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() { ch1 <- 1; ch1 <- 2; ch1 <- 3; close(ch1) }()
	go func() { ch2 <- 1; ch2 <- 2; ch2 <- 3; close(ch2) }()
	go func() { ch3 <- 1; ch3 <- 2; ch3 <- 3; close(ch3) }()

	dd := merge(ch1, ch2, ch3)
	for d := range dd { // blocks until merge closes outCh
		fmt.Println(d)
	}

	// ── Fan-out / Worker Pool ─────────────────────────────────────────────────
	// processJobs distributes 10 jobs across 3 worker goroutines (fan-out).
	// Each worker squares its job value and sends the result to resCh (fan-in back).
	ctx := context.Background()
	a, err := processJobs(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)

	// ── Pipeline ──────────────────────────────────────────────────────────────
	// pipeline() chains three goroutine stages:
	//   stage 1 → produces raw ints
	//   stage 2 → squares them
	//   stage 3 → doubles the squares
	// Each stage is connected by a channel; closing one channel terminates the next.
	d := pipeline()
	for g := range d {
		fmt.Println(g)
	}
}

// merge fans-in three channels into one output channel.
// A WaitGroup tracks when every drain goroutine finishes so outCh is closed exactly once.
func merge(ch1, ch2, ch3 <-chan int) <-chan int {
	wg := sync.WaitGroup{}
	outCh := make(chan int, 1) // buffer of 1 keeps producers from blocking while draining

	// One goroutine per input channel; each drains its channel into outCh.
	drain := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			outCh <- v
		}
	}

	wg.Add(1)
	go drain(ch1)
	wg.Add(1)
	go drain(ch2)
	wg.Add(1)
	go drain(ch3)

	// A separate goroutine waits for all drains to finish, then closes outCh.
	// This signals the consumer's range loop to exit.
	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

// processJobs fans work out to n workers and collects results.
// Context cancellation is propagated into every send/receive via select.
func processJobs(ctx context.Context, jobs []int, n int) ([]int, error) {
	wg := sync.WaitGroup{}

	jobCh := make(chan int, 10)        // buffered so the producer rarely blocks
	resCh := make(chan int, len(jobs)) // large enough to hold all results without blocking workers

	// Start n worker goroutines — this is the fan-out.
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh { // blocks until a job arrives or jobCh is closed
				select {
				case resCh <- j * j: // fan-in: send result back
				case <-ctx.Done():   // abort if context is cancelled
				}
			}
		}()
	}

	// Producer: feeds jobs into jobCh, then closes it so workers know there's no more work.
	go func() {
		for _, j := range jobs {
			select {
			case jobCh <- j:
			case <-ctx.Done():
			}
		}
		close(jobCh) // signals workers to exit their range loop
	}()

	// Closer: waits for all workers to finish, then closes resCh so the collector can range over it.
	go func() {
		wg.Wait()
		close(resCh)
	}()

	// Collect all results from resCh (blocks until resCh is closed).
	var res []int
	for r := range resCh {
		res = append(res, r)
	}

	return res, ctx.Err() // non-nil if context was cancelled
}

// pipeline chains three transformation stages via channels.
// Each stage reads from the previous channel, transforms the value, and writes to the next.
func pipeline() <-chan int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)

	// Stage 1: emit raw integers 1..5.
	go func() {
		for i := 1; i <= 5; i++ {
			ch1 <- i
		}
		close(ch1) // signals stage 2 to stop
	}()

	// Stage 2: square each value from stage 1.
	go func() {
		defer close(ch2) // signals stage 3 to stop when stage 1 is exhausted
		for j := range ch1 {
			ch2 <- j * j
		}
	}()

	// Stage 3: double each squared value from stage 2.
	go func() {
		defer close(ch3) // signals the final consumer to stop
		for j := range ch2 {
			ch3 <- j * 2
		}
	}()

	return ch3 // caller ranges over the final stage's channel
}
