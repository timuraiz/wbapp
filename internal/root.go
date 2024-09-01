package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

const RetryCount = 3

// backoff calculates the retry backoff duration based on the number of retries.
func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

// shouldRetry determines if the request should be retried based on the error and response.
func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	if resp != nil && (resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout) {
		return true
	}
	return false
}

// retryableTransport is a custom RoundTripper that handles retries for transient errors.
type retryableTransport struct {
	transport http.RoundTripper
}

// RoundTrip executes the request and retries on failure.
func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	var resp *http.Response
	var err error
	retries := 0

	for retries < RetryCount {
		resp, err = t.transport.RoundTrip(req)
		if !shouldRetry(err, resp) {
			break
		}

		// Wait before retrying
		time.Sleep(backoff(retries))

		// Consume any response body to reuse the connection
		drainBody(resp)

		// Recreate the request body for the retry
		if req.Body != nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		retries++
	}

	return resp, err
}

// drainBody reads and discards the response body.
func drainBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

// NewRetryableClient creates an HTTP client with retryable transport.
func NewRetryableClient() *http.Client {
	return &http.Client{
		Transport: &retryableTransport{
			transport: http.DefaultTransport,
		},
	}
}

// WbScanner defines a struct for making HTTP requests.
type WbScanner struct {
	host   string
	client *http.Client
}

// NewWbScanner creates a new WbScanner with a customizable HTTP client.
func NewWbScanner(host string, client *http.Client) *WbScanner {
	if client == nil {
		client = NewRetryableClient() // Use the retryable client if none is provided
	}
	return &WbScanner{host: host, client: client}
}

// Request sends an HTTP request with retry logic.
func (wb *WbScanner) Request(route string, method string, payload bytes.Buffer, queryParams url.Values, headers map[string]string) ([]byte, error) {
	fullURL := wb.host + route
	if method == http.MethodGet && queryParams != nil {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryParams.Encode())
	}

	var req *http.Request
	var err error
	if method == http.MethodGet {
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		req, err = http.NewRequest(method, fullURL, &payload)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := wb.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
