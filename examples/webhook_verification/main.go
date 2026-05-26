package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/LevelFourAI/levelfour-go/levelfour/webhooks"
)

func main() {
	verifier, err := webhooks.NewVerifier("whsec_your_secret_here")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		payload, err := verifier.Verify(r.Header, body)
		if err != nil {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		fmt.Printf("Received webhook: %v\n", payload)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
