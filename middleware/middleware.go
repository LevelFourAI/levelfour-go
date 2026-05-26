package middleware

import (
	"context"
	"net/http"

	"github.com/LevelFourAI/levelfour-go/levelfour"
)

type contextKey struct{}

// Config configures the LevelFour middleware.
type Config struct {
	APIKey     string
	BaseURL    string
	MaxRetries uint
}

// WithLevelFour wraps an http.Handler, injecting a configured LevelFour client
// into every request's context. Retrieve it with ClientFromContext.
// Panics at startup if the API key is invalid.
func WithLevelFour(next http.Handler, cfg Config) http.Handler {
	var opts []levelfour.ClientOption
	if cfg.BaseURL != "" {
		opts = append(opts, levelfour.WithBaseURL(cfg.BaseURL))
	}
	if cfg.MaxRetries > 0 {
		opts = append(opts, levelfour.WithMaxRetries(cfg.MaxRetries))
	}

	client, err := levelfour.NewClient(cfg.APIKey, opts...)
	if err != nil {
		panic("levelfour middleware: " + err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKey{}, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ClientFromContext retrieves the LevelFour client from the request context.
// Returns nil if the middleware was not applied.
func ClientFromContext(ctx context.Context) *levelfour.Client {
	client, _ := ctx.Value(contextKey{}).(*levelfour.Client)
	return client
}
