# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1](https://github.com/LevelFourAI/levelfour-go/compare/v0.1.0...v0.1.1) (2026-05-26)


### Features

* initial public release of the LevelFour Go SDK ([1b9eebe](https://github.com/LevelFourAI/levelfour-go/commit/1b9eebec24eed0a734a19d27e8c4a5412828b0ac))

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
