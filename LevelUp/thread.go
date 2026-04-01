package main

import (
	"fmt"
	"sync"
)

type Task struct {
	id int
	n  int
}

type Result struct {
	id     int
	factor int
}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		f := fib(task.n)
		results <- Result{id: task.id, factor: f}
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	numWorkers := 5
	numTasks := 20

	tasks := make(chan Task, numTasks)
	results := make(chan Result, numTasks)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	for i := 0; i < numTasks; i++ {
		tasks <- Task{id: i, n: i + 10}
	}
	close(tasks)

	wg.Wait()
	close(results)

	final := make(map[int]int)

	for r := range results {
		final[r.id] = r.factor
	}

	for i := 0; i < numTasks; i++ {
		fmt.Println(i, final[i])
	}
}