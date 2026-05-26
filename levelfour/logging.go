package levelfour

import "net/http"

// Logger is an optional structured logging interface for the SDK.
// Implement this interface to receive debug and warning messages
// from the HTTP transport layer. Compatible with log/slog.Handler
// adapters.
//
// The SDK logs sparingly: request lifecycle at Debug, errors and
// retries at Warn. No logs are emitted by default (no-op).
type Logger interface {
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// WithLogger attaches an optional logger to the client.
// When set, the client logs outbound requests at Debug level and
// errors at Warn level. Pass nil to disable logging (the default).
func WithLogger(l Logger) ClientOption {
	return func(c *clientConfig) {
		c.logger = l
	}
}

type loggingTransport struct {
	base   http.RoundTripper
	logger Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logger.Debug("levelfour: request", "method", req.Method, "url", req.URL.String())
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		t.logger.Warn("levelfour: request failed", "method", req.Method, "url", req.URL.String(), "error", err.Error())
		return nil, err
	}
	return resp, nil
}
