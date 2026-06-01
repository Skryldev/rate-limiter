package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/didip/tollbooth/v8"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

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
	message := Message{
		Status: "Request limit exceeded",
		Body:   "limit exceeded",
	}
	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		log.Fatal(err)
	}

	//tlbthLimiter := tollbooth.NewLimiter(1, nil)
	tlbthLimiter := tollbooth.NewLimiter(5, nil)
	tlbthLimiter.SetMessageContentType("application/json")
	tlbthLimiter.SetMessage(string(jsonMessage))
	http.Handle("/ping", tollbooth.LimitFuncHandler(tlbthLimiter, http.HandlerFunc(endPointHandler)))
	log.Println("✅ Server started successfully on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
