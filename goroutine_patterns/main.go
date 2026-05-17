package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {

	//Fan in
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
		ch1 <- 3
		close(ch1)
	}()
	go func() {
		ch2 <- 1
		ch2 <- 2
		ch2 <- 3
		close(ch2)
	}()
	go func() {
		ch3 <- 1
		ch3 <- 2
		ch3 <- 3
		close(ch3)
	}()

	dd := merge(ch1, ch2, ch3)

	for d := range dd {
		fmt.Println(d)
	}

	//Worker Pool with fan out
	ctx := context.Background()
	a, err := processJobs(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)

	// Pipeline
	d := pipeline()

	for g := range d {
		fmt.Println(g)
	}
}

// Fan in method
func merge(ch1, ch2, ch3 <-chan int) <-chan int {
	wg := sync.WaitGroup{}
	outCh := make(chan int, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch1 {
			outCh <- v
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch2 {
			outCh <- v
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch3 {
			outCh <- v
		}
	}()

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh

}

// Worker pool with fan out
func processJobs(ctx context.Context, jobs []int, n int) ([]int, error) {
	wg := sync.WaitGroup{}

	jobCh := make(chan int, 10)
	resCh := make(chan int, len(jobs))

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				select {
				case resCh <- j * j:
				case <-ctx.Done():
				}
			}
		}()
	}

	go func() {
		for _, j := range jobs {
			select {
			case jobCh <- j:
			case <-ctx.Done():
			}
		}
		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(resCh)
	}()

	var res []int

	for r := range resCh {
		res = append(res, r)
	}

	return res, ctx.Err()
}

// Pipeline
func pipeline() <-chan int {

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)

	go func() {
		ch1 <- 1
		ch1 <- 2
		ch1 <- 3
		ch1 <- 4
		ch1 <- 5

		close(ch1)
	}()

	go func() {
		defer close(ch2)
		for j := range ch1 {
			ch2 <- j * j
		}
	}()

	go func() {
		defer close(ch3)
		for j := range ch2 {
			ch3 <- j * 2
		}
	}()

	return ch3
}
