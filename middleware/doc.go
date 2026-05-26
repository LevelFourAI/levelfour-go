// Package middleware provides HTTP middleware for injecting a LevelFour
// client into the request context.
//
// Use [WithLevelFour] to wrap an http.Handler and [ClientFromContext] to
// retrieve the client in downstream handlers:
//
//	handler := middleware.WithLevelFour(mux, middleware.Config{
//	    APIKey: "l4_live_...",
//	})
//
//	func myHandler(w http.ResponseWriter, r *http.Request) {
//	    client := middleware.ClientFromContext(r.Context())
//	    // use client...
//	}
package middleware
