package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	id    int
	value int
}

type Node struct {
	jobs    chan Job
	results chan int
	next    *Node
}

func (n *Node) start(wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range n.jobs {
		res := process(job.value)
		if n.next != nil {
			n.next.jobs <- Job{id: job.id, value: res}
		} else {
			n.results <- res
		}
	}
	if n.next != nil {
		close(n.next.jobs)
	}
}

func process(x int) int {
	time.Sleep(10 * time.Millisecond)
	return x * x
}

func main() {
	stages := 4
	numJobs := 15

	results := make(chan int, numJobs)

	var first *Node
	var prev *Node

	for i := 0; i < stages; i++ {
		node := &Node{
			jobs:    make(chan Job, numJobs),
			results: results,
		}
		if first == nil {
			first = node
		}
		if prev != nil {
			prev.next = node
		}
		prev = node
	}

	var wg sync.WaitGroup

	curr := first
	for curr != nil {
		wg.Add(1)
		go curr.start(&wg)
		curr = curr.next
	}

	for i := 1; i <= numJobs; i++ {
		first.jobs <- Job{id: i, value: i}
	}
	close(first.jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r)
	}
}