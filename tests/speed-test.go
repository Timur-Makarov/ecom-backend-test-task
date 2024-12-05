package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	const (
		numRequests = 1000
		concurrency = 50
		url         = "http://localhost:8080/counter/1"
	)

	var wg sync.WaitGroup
	startTime := time.Now()

	// Channel to limit concurrency
	sem := make(chan struct{}, concurrency)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			resp, err := http.Post(url, "", nil)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()
		}()
	}

	wg.Wait()
	duration := time.Since(startTime)
	fmt.Printf("Made %d requests in %v\n", numRequests, duration)
}
