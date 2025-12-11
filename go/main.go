package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Each worker takes tasks from the tasks channel, "processes" them,
// and sends the result into the results channel.
func worker(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		// simulate some work
		time.Sleep(200 * time.Millisecond)
		results <- fmt.Sprintf("Worker %d processed %s", id, task)
	}
}

// This function writes every result received from the results channel
// into a text file. The file closes automatically because of defer.
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

	// start worker goroutines
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// send tasks
	for i := 1; i <= 10; i++ {
		tasks <- fmt.Sprintf("Task-%d", i)
	}
	close(tasks)

	// wait for workers, then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// write all results to file
	resultWriter("results_go.txt", results)
}
