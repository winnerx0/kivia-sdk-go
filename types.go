package kiviasdk

import (
	"bytes"
	"net/http"
	"time"
)

// Log represents a single request lifecycle metric captured by the SDK.
type Log struct {

	Id string `json:"id"`

	Path string `json:"path"`

	Status int `json:"status"`

	IPAddress string `json:"ip_address,omitempty"`

	Timestamp time.Time `json:"timestamp"`

	Latency int `json:"latency"`
}

// Client is the primary Kivia SDK client instance, holding configuration and connection states.
type Client struct {
	BaseURL    string
	ApiKey     string
	HttpClient *http.Client
}
// ResponseWriteWrapper intercepts status codes and body lengths to facilitate latency and logging analytics.
type ResponseWriteWrapper struct {
	http.ResponseWriter
	body *bytes.Buffer
	statusCode int
}
