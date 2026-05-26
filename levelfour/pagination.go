package levelfour

import "context"

// CollectAll iterates through all pages and returns every result.
// This is a convenience wrapper around the [PageIterator] for the
// common case of collecting all items into a single slice.
//
//	page, err := client.Recommendations.List(ctx, request)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	items, err := levelfour.CollectAll(ctx, page)
func CollectAll[Cursor comparable, T any, R any](ctx context.Context, page *Page[Cursor, T, R]) ([]T, error) {
	var results []T
	iter := page.Iterator()
	for iter.Next(ctx) {
		results = append(results, iter.Current())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
