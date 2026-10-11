package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func Con(url string) time.Duration{
	start := time.Now()

	res, err := http.Get(url)

	if (err != nil) {
		return 0
	}

	defer res.Body.Close()

	return time.Since(start).Round(time.Millisecond)
}

func main(){
	url := "https://go.dev/"
	workers := 2
	times := 10
	jobs := make(chan string)
	results := make(chan time.Duration)

	var wg sync.WaitGroup

	start := time.Now()

	go func(){
		for i := 0; i < times; i++ {
			jobs <- url
		}
		close(jobs)
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			for job := range jobs {
				results <- Con(job)
			}
		}()
	}

	go func(){
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}
	fmt.Printf("Time since app started: %v", time.Since(start).Round(time.Millisecond))
}
