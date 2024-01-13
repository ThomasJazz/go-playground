package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	maxConcurrentRequests = 100
	port                  = 80
)

type Payload struct {
	Time      time.Time `json:"time"`
	RequestId string    `json:"requestId"`
	Text      string    `json:"text"`
}

var currentRequests atomic.Int32

func main() {
	http.HandleFunc("/", handleRequest)
	log.Printf("Server starting on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if currentRequests.Load() >= maxConcurrentRequests {
			http.Error(w, "Server too busy, please try again later", http.StatusServiceUnavailable)
			return
		}
		currentRequests.Add(1)
		defer currentRequests.Add(-1)

		next.ServeHTTP(w, r)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Our Hello World Payload
	payload := Payload{
		Time:      time.Now(),
		RequestId: uuid.NewString(),
		Text:      "Hello World!",
	}

	// Return our payload
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
