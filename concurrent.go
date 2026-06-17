package main

import (
	"bufio"
	"os"
	"sync"
)

// RunConcurrent processes the file and counts word frequencies concurrently using a worker pool.
func RunConcurrent(filename string, numWorkers int) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Channel to distribute lines to workers
	linesChan := make(chan string, 2048)
	// Channel to collect partial maps from workers
	partialMapsChan := make(chan map[string]int, numWorkers)

	var wg sync.WaitGroup
	var scanErr error
	var scanErrMu sync.Mutex

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localFreqs := make(map[string]int)
			for line := range linesChan {
				words := TokenizeLine(line)
				for _, w := range words {
					localFreqs[w]++
				}
			}
			partialMapsChan <- localFreqs
		}()
	}

	// Producer: read lines and send to linesChan
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			linesChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			scanErrMu.Lock()
			scanErr = err
			scanErrMu.Unlock()
		}
		close(linesChan)
	}()

	// Wait for workers to finish and then close the partial maps channel
	go func() {
		wg.Wait()
		close(partialMapsChan)
	}()

	// Merge partial maps
	finalFreqs := make(map[string]int)
	for localMap := range partialMapsChan {
		for k, v := range localMap {
			finalFreqs[k] += v
		}
	}

	scanErrMu.Lock()
	err = scanErr
	scanErrMu.Unlock()
	if err != nil {
		return nil, err
	}

	return finalFreqs, nil
}
