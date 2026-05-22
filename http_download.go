package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type bodyReadTimeoutError struct {
	timeout time.Duration
}

func (e *bodyReadTimeoutError) Error() string {
	return fmt.Sprintf("download timed out after %s while reading response body", e.timeout)
}

func (e *bodyReadTimeoutError) Timeout() bool {
	return true
}

func (e *bodyReadTimeoutError) Temporary() bool {
	return true
}

// fetchDownloadURL GETs a URL with short connect/response timeouts and a
// separate longer timeout for reading large response bodies.
func fetchDownloadURL(url string) ([]byte, error) {
	resp, err := downloadHTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := readBodyWithTimeout(resp.Body, downloadBodyTimeout)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return body, nil
}

func readBodyWithTimeout(body io.ReadCloser, timeout time.Duration) ([]byte, error) {
	type readResult struct {
		data []byte
		err  error
	}

	resultCh := make(chan readResult, 1)
	go func() {
		data, err := io.ReadAll(body)
		resultCh <- readResult{data: data, err: err}
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case result := <-resultCh:
		closeErr := body.Close()
		if result.err != nil {
			return nil, result.err
		}
		if closeErr != nil {
			return nil, closeErr
		}
		return result.data, nil
	case <-timer.C:
		body.Close()
		<-resultCh
		return nil, &bodyReadTimeoutError{timeout: timeout}
	}
}
