package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/LevelFourAI/levelfour-go/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		client := middleware.ClientFromContext(r.Context())
		if client == nil {
			http.Error(w, "no client", http.StatusInternalServerError)
			return
		}

		summary, err := client.Recommendations.GetSavingsByProvider(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Summary: %+v\n", summary)
	})

	// WithLevelFour injects a configured client into every request's context.
	// Uses LEVELFOUR_API_KEY environment variable.
	handler := middleware.WithLevelFour(mux, middleware.Config{
		APIKey: "", // reads from LEVELFOUR_API_KEY env var
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
