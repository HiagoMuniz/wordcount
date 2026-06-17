package main

import (
	"bufio"
	"os"
)

// RunSequential processes the file and counts word frequencies sequentially.
func RunSequential(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	freqs := make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := TokenizeLine(scanner.Text())
		for _, w := range words {
			freqs[w]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return freqs, nil
}
