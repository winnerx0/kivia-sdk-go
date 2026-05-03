# Kivia Go SDK

The Kivia Go SDK is the official Go middleware client for integrating with the Kivia observability platform. It provides seamless tracking of request metrics, including latencies, target paths, request statuses, and IP addresses.

## Installation

```bash
go get github.com/kivia-observe/kivia-sdk-go
```

## Quick Start (with standard `net/http`)

Here is a quick example of how you can wrap your existing `net/http` handlers to capture analytics automatically:

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/kivia-observe/kivia-sdk-go"
)

func main() {
	// Initialize the Kivia client with your unique project API key.
	// You can get this key from your Kivia dashboard.
	kiviaClient := kiviasdk.NewClient("YOUR_KIVIA_API_KEY")

	// Your standard application logic
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Kivia Go SDK!"))
	})

	// Wrap your entire application router inside the Kivia middleware logger
	loggedHandler := kiviaClient.NewLog(mux)

	fmt.Println("Server running on :8080...")
	http.ListenAndServe(":8080", loggedHandler)
}
```

## Architecture

This SDK exports the following:
* `NewClient(apiKey)`: Prepares the backend target routing and authentication payload.
* `client.NewLog(http.Handler)`: Native `net/http` compatible middleware. Wraps around any multiplexer (e.g. standard library `ServeMux`, `gorilla/mux`, `chi`) and automatically traps the execution time to log analytics asynchronously without holding up the server responses!
