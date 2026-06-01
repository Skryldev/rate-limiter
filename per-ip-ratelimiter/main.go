package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func endPointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hello World!",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		return
	}
}
func main() {
	cm := newCleanupManager(5*time.Minute, 10*time.Minute)
	cm.Start()
	defer cm.Stop()

	http.Handle("/ping", perClientRateLimiter()(http.HandlerFunc(endPointHandler)))
	log.Println("✅ Server started successfully on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
