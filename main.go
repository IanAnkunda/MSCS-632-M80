package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Each worker pulls tasks from the channel, "processes" them, and sends results back
func worker(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		time.Sleep(200 * time.Millisecond) // simulate work
		results <- fmt.Sprintf("Worker %d processed %s", id, task)
	}
}

// Writes results to a file as they come in
func resultWriter(filename string, results <-chan string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for line := range results {
		file.WriteString(line + "\n")
	}
}

func main() {
	tasks := make(chan string, 10)
	results := make(chan string, 10)

	var wg sync.WaitGroup

	// start workers
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// send tasks
	for i := 1; i <= 10; i++ {
		tasks <- fmt.Sprintf("Task-%d", i)
	}
	close(tasks)

	// close results channel after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// write everything to file
	resultWriter("results_go.txt")
}
