package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LevelFourAI/levelfour-go/levelfour"
)

func main() {
	// Uses LEVELFOUR_API_KEY environment variable by default.
	// Pass a key explicitly with: levelfour.NewClient("l4_live_...")
	client, err := levelfour.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	summary, err := client.Recommendations.GetSavingsByProvider(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Recommendations summary: %+v\n", summary)
}
