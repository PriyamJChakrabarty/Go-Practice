package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID int
}

type Output struct {
	ID    int
	Value int
	Err   error
}

func generator(n int) <-chan Task {
	out := make(chan Task)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- Task{ID: i}
		}
	}()
	return out
}

func processor(ctx context.Context, task Task) (int, error) {
	delay := time.Duration(rand.Intn(300)) * time.Millisecond
	select {
	case <-time.After(delay):
		if rand.Float32() < 0.3 {
			return 0, fmt.Errorf("failed task %d", task.ID)
		}
		return task.ID * 10, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func worker(ctx context.Context, id int, tasks <-chan Task, results chan<- Output, retries int) {
	for task := range tasks {
		var val int
		var err error

		for attempt := 0; attempt <= retries; attempt++ {
			taskCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
			val, err = processor(taskCtx, task)
			cancel()

			if err == nil {
				break
			}
		}

		results <- Output{
			ID:    task.ID,
			Value: val,
			Err:   err,
		}
	}
}

func fanOut(ctx context.Context, in <-chan Task, workerCount int, retries int) <-chan Output {
	out := make(chan Output)
	var wg sync.WaitGroup

	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, in, out, retries)
		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func fanIn(channels ...<-chan Output) <-chan Output {
	var wg sync.WaitGroup
	out := make(chan Output)

	output := func(c <-chan Output) {
		defer wg.Done()
		for val := range c {
			out <- val
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	tasks := generator(20)

	stage1 := fanOut(ctx, tasks, 4, 2)

	final := fanIn(stage1)

	success := 0
	fail := 0

	for res := range final {
		if res.Err != nil {
			fmt.Printf("Task %d failed after retries\n", res.ID)
			fail++
		} else {
			fmt.Printf("Task %d success: %d\n", res.ID, res.Value)
			success++
		}
	}

	fmt.Println("\nSummary:")
	fmt.Println("Success:", success)
	fmt.Println("Failed:", fail)
}