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
// The structure can be used as a load balancer
func main() {

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

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

	for i := 0; i < len(urls); i++ {

		wg.Add(1)
		go func(x int) {
			// If the context is cancelled, the request will be cancelled as well
			RequestWithContext(ctx, &wg, urls[x], cancel)
		}(i)
	}

	wg.Wait()
}

func RequestWithContext(ctx context.Context, wg *sync.WaitGroup, url string, cancel context.CancelFunc) error {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	duration := time.Since(start)

	defer resp.Body.Close()
	defer cancel()
	fmt.Println("Url: "+url+" got response status:", resp.Status, "in", duration)

	return nil
}
