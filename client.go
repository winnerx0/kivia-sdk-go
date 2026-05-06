package kiviasdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// NewClient initializes a new Kivia SDK client with the provided API key.
// It configures the default production base URL and the underlying HTTP client.
func NewClient(apiKey string) *Client {
	client := &Client{
		BaseURL:    "https://kivia.winnerx0.dev/api/v1",
		ApiKey:     apiKey,
		HttpClient: &http.Client{},
	}

	return client
}

// NewLog returns an HTTP middleware handler that automatically tracks
// and logs request metrics such as latency, status code, and target path.
//
// Example:
//
//	kiviaClient := kiviasdk.NewClient("API_KEY")
//	mux := http.NewServeMux()
//	loggedHandler := kiviaClient.NewLog(mux)
//	http.ListenAndServe(":8080", loggedHandler)
func (c *Client) NewLog(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// start time for the request
		start := time.Now()

		rw := &ResponseWriteWrapper{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			statusCode:     200,
		}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		go c.sendLog(r, rw, latency)
	}
}

// sendLog is called asynchronously to send the collected log data to the Kivia backend.
// It builds the Log entry and executes a POST request against the observability API.
func (c *Client) sendLog(r *http.Request, rw *ResponseWriteWrapper, latency time.Duration) {
    url := fmt.Sprintf("%s/logs/create", c.BaseURL)
    logEntry := &Log{
        Path:      r.URL.String(),
        Status:    rw.statusCode,
        IPAddress: r.RemoteAddr,
        Timestamp: time.Now(),
        Latency:   int(latency.Milliseconds()),
    }

    bodyBytes, err := json.Marshal(logEntry)
    if err != nil {
        log.Printf("kiviasdk: failed to marshal log: %v", err)
        return
    }

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
    if err != nil {
        log.Printf("kiviasdk: failed to create request: %v", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-kivia-api-key", c.ApiKey)

    res, err := c.HttpClient.Do(req)
    if err != nil {
        log.Printf("kiviasdk: failed to send log: %v", err)
        return
    }
    defer res.Body.Close()

    if res.StatusCode != 201 {
        body, _ := io.ReadAll(res.Body)
        log.Printf("kiviasdk: log rejected (%d): %s", res.StatusCode, string(body))
    }
}
