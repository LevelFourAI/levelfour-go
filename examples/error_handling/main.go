package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LevelFourAI/levelfour-go/levelfour"
)

func main() {
	client, err := levelfour.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Recommendations.Get(context.Background(), "nonexistent-id")
	if err != nil {
		if levelfour.IsNotFound(err) {
			fmt.Println("Recommendation not found")
			return
		}

		if levelfour.IsRateLimited(err) {
			fmt.Println("Rate limited, retry later")
			return
		}

		fmt.Printf("API error (status %d): %v\n", levelfour.StatusCode(err), err)
		fmt.Printf("Error code: %s\n", levelfour.ErrorCode(err))
		fmt.Printf("Error message: %s\n", levelfour.ErrorMessage(err))
		fmt.Printf("Raw body: %s\n", levelfour.ErrorBody(err))
		return
	}
}
