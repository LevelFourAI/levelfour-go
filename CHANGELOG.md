# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.1.0...v0.2.0) (2026-09-20)


### Features

* sync generated Go SDK ([85c8cd4](https://github.com/LevelFourAI/levelfour-go/commit/85c8cd425eb96d4a05c32371cccead9ea7313590))
* sync generated Go SDK from levelfour-api ([fca5ce9](https://github.com/LevelFourAI/levelfour-go/commit/fca5ce900dcacafd21f6c16b6e67935986ba1381))


### Bug Fixes

* let a feature cut a minor version while pre-1.0 ([789c51c](https://github.com/LevelFourAI/levelfour-go/commit/789c51ca871d31f2980745698831fa257f19c326))


### Miscellaneous

* repository hygiene for a public SDK ([16df3c6](https://github.com/LevelFourAI/levelfour-go/commit/16df3c62ce6b808f15853403aae513dc063bba68))
* repository hygiene for a public SDK ([0a43289](https://github.com/LevelFourAI/levelfour-go/commit/0a43289800599371721983b593294bd49402dd1a))

## [0.1.0](https://github.com/LevelFourAI/levelfour-go/releases/tag/v0.1.0) (2026-05-26)

Initial public release of the LevelFour Go SDK.

### Services

* `client.APIKeys`: list, create, revoke, rotate API keys.
* `client.Accounts`: manage connected cloud accounts and integration installations.
* `client.Audit`: top-level realized-savings audit detail.
* `client.Auth`: identity and organization context (`GetWhoami`).
* `client.Costs`: cost summary, breakdown, daily and monthly aggregates, per-provider drill-downs.
* `client.Health`: API readiness probes.
* `client.Providers`: list connected providers.
* `client.Recommendations`: cost-optimization recommendations, list, detail, per-provider drill-downs, in-progress queue.
* `client.Recommendations.Audit`: realized-savings detail per provider.
* `client.Webhooks`: register, list, and delete webhook endpoints.

### Client features

* Typed errors (`*BadRequestError`, `*NotFoundError`, `*TooManyRequestsError`, etc.) inspectable with `errors.As`.
* Automatic retries on `408`, `409`, `429`, and `5xx` responses (configurable via `WithMaxRetries`, `WithNoRetries`).
* Pagination iterators on every list method (`page.Iterator()`, `levelfour.CollectAll`).
* Webhook signature verification (`webhooks.NewVerifier`, HMAC-SHA256).
* Per-request overrides via the `option` package.
* Single hand-written wrapper around the Fern-generated base client. OpenTelemetry instrumentation is left to the caller via `WithHTTPClient(otelhttp.NewTransport(...))`.
