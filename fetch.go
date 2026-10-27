package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Proto      string
	Status     string
	Headers    http.Header
	Body       []byte
	Err        error
}

func fetchHedged(ctx context.Context, url string, ch chan<- Result) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- Result{URL: url, Err: err}
		return
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ch <- Result{URL: url, Err: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- Result{URL: url, Err: fmt.Errorf("bad status: %d", resp.StatusCode)}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- Result{URL: url, Err: err}
		return
	}
	ch <- Result{URL: url,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       body}
}
