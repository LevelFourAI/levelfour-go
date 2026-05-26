// Package levelfour contains the auto-generated request/response types and
// domain-specific sub-clients for the LevelFour API.
//
// Most users should import the [github.com/LevelFourAI/levelfour-go/levelfour]
// sub-package, which provides client construction, error helpers, and
// automatic retries:
//
//	import "github.com/LevelFourAI/levelfour-go/levelfour"
//
//	client, err := levelfour.NewClient("l4_live_...")
//	summary, err := client.Recommendations.GetSavingsByProvider(ctx)
//
// This root package is used for request/response struct types (e.g.,
// [ListRecommendationsRequest], [RecommendationSummary]) which are
// referenced in the sub-client method signatures.
package levelfour
