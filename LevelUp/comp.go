package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	Value int
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, rateLimiter <-chan time.Time) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			<-rateLimiter

			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))

			result := Result{
				JobID: job.ID,
				Value: job.ID * job.ID,
			}

			fmt.Printf("Worker %d processed Job %d\n", id, job.ID)

			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numWorkers := 5
	numJobs := 20

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rateLimiter := time.Tick(100 * time.Millisecond)

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg, rateLimiter)
	}

	go func() {
		for i := 1; i <= numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	finalResults := make(map[int]int)
	var mu sync.Mutex

	var collectorWG sync.WaitGroup
	collectorWG.Add(1)

	go func() {
		defer collectorWG.Done()
		for res := range results {
			mu.Lock()
			finalResults[res.JobID] = res.Value
			mu.Unlock()
		}
	}()

	collectorWG.Wait()

	fmt.Println("\nFinal Results:")
	for k, v := range finalResults {
		fmt.Printf("Job %d => %d\n", k, v)
	}
}