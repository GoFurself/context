package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// This program sends multiple requests to different urls concurrently and cancels the requests after one of them is done.
// It main purpose is to show how to use context to cancel requests and how to use WaitGroup to wait for all requests to finish.
// The idea can be used for example as a load balancer
func main() {

	urls := []string{
		"https://www.tumblr.com",
		"https://www.reddit.com",
		"https://www.snapchat.com",
		"https://www.whatsapp.com",
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.linkedin.com",
		"https://www.instagram.com",
		"https://www.pinterest.com",
	}

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			if err := DoGetRequestWithContext(ctx, urls[x]); err == nil {
				cancel()
				fmt.Println("Url: " + urls[x] + " got 200 status first. Cancelling other requests.")
			} else {
				fmt.Println("Error: ", err)
			}
		}(i)
	}
	wg.Wait()
}

func DoGetRequestWithContext(ctx context.Context, url string) error {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code is not 200. Status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}
