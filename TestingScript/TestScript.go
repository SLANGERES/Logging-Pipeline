package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Logging struct {
	Id        string    `json:"id"`
	TimeStamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	MetaData  MetaData  `json:"metadata"`
}

type MetaData struct {
	SourceIp string `json:"source_ip"`
	Region   string `json:"region"`
}

func main() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("Program completed in %s\n", elapsed)
	}()

	// Test configuration
	service := []string{"Payment", "Internal", "External API", "Cloud", "Order"}
	levels := []string{"Info", "Err", "War"}
	url := "http://127.0.0.1:9090/log"
	numRequests := 100       // Total number of requests
	concurrency := 100        // Number of concurrent workers

	metadata := MetaData{
		SourceIp: "0.0.0.0",
		Region:   "Test_region",
	}

	rand.Seed(time.Now().UnixNano())

	// Channels for worker pool
	tasks := make(chan int, numRequests)
	var wg sync.WaitGroup

	// Start worker pool
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range tasks {
				randomService := service[rand.Intn(len(service))]
				randomLevel := levels[rand.Intn(len(levels))]

				logEntry := Logging{
					Id:        fmt.Sprintf("log-%d", i),
					TimeStamp: time.Now(),
					Service:   randomService,
					Level:     randomLevel,
					Message:   "Testing",
					MetaData:  metadata,
				}

				payloadBytes, _ := json.Marshal(logEntry)

				resp, err := http.Post(url, "application/json", bytes.NewReader(payloadBytes))
				if err != nil {
					fmt.Printf("Request %d failed: %v\n", i, err)
					continue
				}
				io.Copy(io.Discard, resp.Body) // Discard body to speed up
				resp.Body.Close()

				fmt.Printf("Request %d → Status: %d\n", i, resp.StatusCode)
			}
		}()
	}

	// Feed tasks
	for i := 1; i <= numRequests; i++ {
		tasks <- i
	}
	close(tasks)

	wg.Wait()
	fmt.Println("✅ All requests completed.")
}
