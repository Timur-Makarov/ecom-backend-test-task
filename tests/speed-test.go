package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	const runs = 150
	var timeTook time.Duration

	for i := 0; i < runs; i++ {
		duration := runTest()
		timeTook += duration
	}

	log.Printf("%d runs took %v s \n", runs, timeTook.Seconds())
	log.Printf("One run took %v on avarage \n", timeTook/runs)
}

func runTest() time.Duration {
	const (
		numRequests = 100
		concurrency = 10
		url         = "http://localhost:8080/banners/1/stats"
	)

	var wg sync.WaitGroup
	startTime := time.Now()

	sem := make(chan struct{}, concurrency)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			req, err := http.NewRequest("PUT", url, nil)
			if err != nil {
				log.Println("failed to create request:", err)
				return
			}

			_, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Println("failed to send request:", err)
				return
			}
		}()
	}

	wg.Wait()
	duration := time.Since(startTime)
	fmt.Printf("Made %d requests in %v\n", numRequests, duration)
	return duration
}
