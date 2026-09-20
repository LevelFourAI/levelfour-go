# LevelFour Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/LevelFourAI/levelfour-go/levelfour.svg)](https://pkg.go.dev/github.com/LevelFourAI/levelfour-go/levelfour)
[![CI](https://github.com/LevelFourAI/levelfour-go/actions/workflows/ci.yml/badge.svg)](https://github.com/LevelFourAI/levelfour-go/actions/workflows/ci.yml)

The official Go SDK for the [LevelFour](https://levelfour.ai) API.

> **Status: v0.x.** The public API is still stabilizing. Minor versions may include breaking changes until v1.0. Pin to an exact version in production and review the [CHANGELOG](CHANGELOG.md) before upgrading.

## Requirements

Go 1.24 or later. We test against the three most recent Go minor releases and drop support for older versions when they reach end-of-life.

## Documentation

API reference is available at [pkg.go.dev](https://pkg.go.dev/github.com/LevelFourAI/levelfour-go).

## Installation

```bash
# x-release-please-start-version
go get github.com/LevelFourAI/levelfour-go@v0.2.0
# x-release-please-end
```

## Usage

Set `LEVELFOUR_API_KEY` in your environment, or pass the key directly:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LevelFourAI/levelfour-go/levelfour"
)

func main() {
	// Uses LEVELFOUR_API_KEY env var. Pass a key explicitly:
	// levelfour.NewClient("l4_live_...")
	client, err := levelfour.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	summary, err := client.Recommendations.GetSavingsByProvider(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(summary)
}
```

## Error Handling

Use the convenience functions in the `levelfour` package to check API error types:

```go
_, err := client.Recommendations.Get(ctx, id)
if err != nil {
	if levelfour.IsNotFound(err) {
		// handle 404
	}
	if levelfour.IsRateLimited(err) {
		// back off and retry
	}
	// extract structured error details
	fmt.Println(levelfour.ErrorCode(err))    // e.g. "NOT_FOUND"
	fmt.Println(levelfour.ErrorMessage(err)) // e.g. "resource not found"
	fmt.Println(levelfour.ErrorBody(err))    // raw response body
}
```

Available checks: `IsBadRequest` (400), `IsUnauthorized` (401), `IsForbidden` (403), `IsNotFound` (404), `IsTimeout` (408), `IsConflict` (409), `IsRateLimited` (429), `IsValidationError` (422), `IsServerError` (5xx).

## Configuration

```go
client, err := levelfour.NewClient("l4_live_...",
	levelfour.WithBaseURL("https://custom.api.com"),
	levelfour.WithMaxRetries(5),
	levelfour.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
	levelfour.WithLogger(slog.Default()),
)
```

| Option | Default | Description |
|--------|---------|-------------|
| `WithBaseURL` | `https://api.levelfour.ai` | API base URL |
| `WithMaxRetries` | `2` | Retry attempts for 429, 408, 409, 5xx |
| `WithNoRetries` | - | Disable retries entirely |
| `WithHTTPClient` | 30s timeout | Custom `*http.Client` |
| `WithLogger` | none | Structured logging (Debug/Warn) |

## Retries

All methods (including `Execute`, `Get`, `Post`, and `RawRequest`) automatically retry failed requests with exponential backoff, jitter, and `Retry-After` header support. By default, retries up to 2 times on:

- `429` Too Many Requests (respects `Retry-After` header)
- `408` Request Timeout
- `409` Conflict
- `5xx` Server Errors

```go
client, err := levelfour.NewClient("l4_live_...", levelfour.WithMaxRetries(5))
// or disable entirely:
client, err := levelfour.NewClient("l4_live_...", levelfour.WithNoRetries())
```

## Timeouts

The default HTTP timeout is 30 seconds. Override with a custom `http.Client`:

```go
client, err := levelfour.NewClient("l4_live_...",
	levelfour.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
)
```

For per-request timeouts, use `context.WithTimeout`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
result, err := client.Recommendations.GetSavingsByProvider(ctx)
```

## Pagination

Paginated endpoints return a `*levelfour.Page` with an auto-paging iterator:

```go
page, err := client.Recommendations.List(ctx, &levelfour.ListRecommendationsRequest{})
if err != nil {
	log.Fatal(err)
}

iter := page.Iterator()
for iter.Next(ctx) {
	item := iter.Current()
	fmt.Println(item)
}
if err := iter.Err(); err != nil {
	log.Fatal(err)
}
```

Or collect all results at once:

```go
items, err := levelfour.CollectAll(ctx, page)
```

<details>
<summary>Raw HTTP Requests</summary>

For undocumented or beta endpoints:

```go
var result map[string]interface{}
err := client.Get(ctx, "/api/v1/custom-endpoint", &result)
err := client.Post(ctx, "/api/v1/custom", body, &result)
```

For full response access (headers, status):

```go
resp, err := client.RawRequest(ctx, http.MethodGet, "/api/v1/health")
defer resp.Body.Close()
fmt.Println(resp.Header.Get("X-Request-Id"))
```

</details>

## Webhooks

**Verify incoming events** using the `levelfour/webhooks` package (HMAC-SHA256):

```go
import "github.com/LevelFourAI/levelfour-go/levelfour/webhooks"

verifier, err := webhooks.NewVerifier("whsec_...")
if err != nil {
	log.Fatal(err)
}

payload, err := verifier.Verify(r.Header, body)
if err != nil {
	http.Error(w, "invalid signature", http.StatusUnauthorized)
	return
}
```

**Manage webhook endpoints** using the generated `client.Webhooks` sub-client:

```go
endpoints, err := client.Webhooks.List(ctx)

_, err = client.Webhooks.Register(ctx, &levelfour.RegisterEndpointRequest{
	Url: "https://yourapp.com/webhooks",
})

err = client.Webhooks.Delete(ctx, endpointID)
```

<details>
<summary>HTTP Middleware</summary>

```go
import "github.com/LevelFourAI/levelfour-go/middleware"

mux := http.NewServeMux()
handler := middleware.WithLevelFour(mux, middleware.Config{
	APIKey: "l4_live_...",
})

// In handlers:
client := middleware.ClientFromContext(r.Context())
```

</details>

<details>
<summary>OpenTelemetry</summary>

The SDK accepts a custom `*http.Client` via `WithHTTPClient`, so you can add [OpenTelemetry](https://opentelemetry.io/) tracing using the standard [`otelhttp`](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp) transport:

```go
import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

client, _ := levelfour.NewClient("l4_live_...",
	levelfour.WithHTTPClient(&http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}),
)
```

This automatically creates spans for every HTTP request, injects W3C `traceparent` headers, and records `http.client.request.duration` metrics.

</details>

## Versioning

This SDK follows [Semantic Versioning](https://semver.org/). While on `v0.x`, minor versions may include breaking changes. After `v1.0.0`, breaking changes will only ship in major version bumps.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## Security

To report vulnerabilities, see [SECURITY.md](SECURITY.md).

## License

Apache-2.0 - see [LICENSE](LICENSE).
