// Package levelfour provides the entry point for the LevelFour Go SDK.
//
// Use [NewClient] to create a configured API client:
//
//	client, err := levelfour.NewClient("l4_live_...")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	summary, err := client.Recommendations.GetSavingsByProvider(ctx)
//
// The client validates API key format, sets default timeouts, and
// configures automatic retries with exponential backoff.
//
// Use the error-checking functions ([IsNotFound], [IsRateLimited], etc.)
// to handle API errors:
//
//	if levelfour.IsNotFound(err) {
//	    // handle 404
//	}
//
// Sub-packages:
//   - [github.com/LevelFourAI/levelfour-go/levelfour/webhooks] — webhook signature verification
//   - [github.com/LevelFourAI/levelfour-go/middleware] — HTTP middleware for context injection
package levelfour
