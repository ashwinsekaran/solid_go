// Package main demonstrates the fan-out/fan-in concurrency pattern applied
// to prime number discovery: a single stream of random integers is fanned
// out to one primeFinder goroutine per CPU core, and their outputs are
// fanned back into a single stream via fanIn.
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	done := make(chan int)
	defer close(done)

	// Single upstream source: an unbounded stream of random integers.
	randomNumFetcher := func() int { return rand.Intn(500000000) }
	randomIntStream := repeatFunc(done, randomNumFetcher)
	//primeStream := primeFinder(done, randomIntStream)

	// naive
	//for num := range take(done, primeStream, 10) {
	//	fmt.Println(num)
	//}

	// Fan-out: primality testing is CPU-bound (trial division), so we spin
	// up one primeFinder per core, all reading from the same
	// randomIntStream, to parallelize the work across CPUs.
	cpuCount := runtime.NumCPU()
	fmt.Println(cpuCount)

	primeFinderChannels := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primeFinderChannels[i] = primeFinder(done, randomIntStream)
	}

	// Fan-in: merge all primeFinder outputs into a single stream so the
	// caller can consume primes without caring how many workers produced
	// them.
	fannedInStream := fanIn(done, primeFinderChannels...)
	for num := range take(done, fannedInStream, 10) {
		fmt.Println(num)
	}

	fmt.Println(time.Since(start))
}

// repeatFunc returns a channel that continuously receives the result of
// calling fn, until done is closed.
func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()

	return stream
}

// take reads exactly n values from stream and forwards them, then closes
// its output channel. This lets the caller cap an otherwise-infinite
// stream to a fixed number of results.
func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()

	return taken
}

// primeFinder consumes integers from randomIntStream and forwards only the
// prime ones. isPrime uses trial division down to 2, which is intentionally
// naive (O(n)) so the workload is CPU-heavy enough to make fanning out
// across cores worthwhile.
func primeFinder(done <-chan int, randomIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randomIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()

	return primes
}

// fanIn merges multiple input channels into a single output channel. A
// WaitGroup tracks all transfer goroutines so the output channel is only
// closed once every input channel has been fully drained/closed.
func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	wg := sync.WaitGroup{}
	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}
	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}
