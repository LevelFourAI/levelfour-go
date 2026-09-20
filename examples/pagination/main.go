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

	ctx := context.Background()

	// Auto-paging iterator: handles page fetching automatically.
	page, err := client.Recommendations.List(ctx, &levelfour.ListRecommendationsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	iter := page.Iterator()
	for iter.Next(ctx) {
		item := iter.Current()
		fmt.Printf("Recommendation: %s\n", item.RecommendationID)
	}
	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}

	// Or use CollectAll to get everything in one call.
	page2, err := client.Recommendations.List(ctx, &levelfour.ListRecommendationsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	items, err := levelfour.CollectAll(ctx, page2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total recommendations: %d\n", len(items))
}
