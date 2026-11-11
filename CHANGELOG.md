# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.8.0...v0.9.0) (2026-11-11)


### Features

* sync generated Go SDK ([#265](https://github.com/LevelFourAI/levelfour-go/issues/265)) ([dce4698](https://github.com/LevelFourAI/levelfour-go/commit/dce46984564e34386a24e65aac54d275c7e05781))

## [0.8.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.7.0...v0.8.0) (2026-11-10)


### Features

* sync generated Go SDK ([#260](https://github.com/LevelFourAI/levelfour-go/issues/260)) ([affba8e](https://github.com/LevelFourAI/levelfour-go/commit/affba8eceb84c81a46d3f9ef022150f8f42f1994))

## [0.7.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.6.0...v0.7.0) (2026-11-07)


### Features

* sync generated Go SDK ([#258](https://github.com/LevelFourAI/levelfour-go/issues/258)) ([0c56610](https://github.com/LevelFourAI/levelfour-go/commit/0c566108ed5be65998694ca49fad32f81968783e))

## [0.6.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.5.0...v0.6.0) (2026-11-06)


### Features

* sync generated Go SDK ([#256](https://github.com/LevelFourAI/levelfour-go/issues/256)) ([6186b55](https://github.com/LevelFourAI/levelfour-go/commit/6186b550e4dc32b5e0205aa7eb1736ca23e360b7))

## [0.5.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.4.1...v0.5.0) (2026-11-05)


### Features

* sync generated Go SDK ([#254](https://github.com/LevelFourAI/levelfour-go/issues/254)) ([446e6d4](https://github.com/LevelFourAI/levelfour-go/commit/446e6d42f60af0caaf4d4336921eeb657485af61))

## [0.4.1](https://github.com/LevelFourAI/levelfour-go/compare/v0.4.0...v0.4.1) (2026-11-04)


### Miscellaneous

* block internal working documents from the repository ([#73](https://github.com/LevelFourAI/levelfour-go/issues/73)) ([aa2efd2](https://github.com/LevelFourAI/levelfour-go/commit/aa2efd2e21231f4bd4dd3d107aac524b30e0d655))

## [0.4.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.3.0...v0.4.0) (2026-11-04)


### Features

* sync generated Go SDK ([#247](https://github.com/LevelFourAI/levelfour-go/issues/247)) ([a238475](https://github.com/LevelFourAI/levelfour-go/commit/a2384757d9238810b289033a95534e43e3474a63))

## [0.3.0](https://github.com/LevelFourAI/levelfour-go/compare/v0.2.1...v0.3.0) (2026-11-04)


### Features

* sync generated Go SDK ([#245](https://github.com/LevelFourAI/levelfour-go/issues/245)) ([7f8b5d5](https://github.com/LevelFourAI/levelfour-go/commit/7f8b5d5aa91aa80b35b7dbce2430d7b38fe54ca7))

## [0.2.1](https://github.com/LevelFourAI/levelfour-go/compare/v0.2.0...v0.2.1) (2026-11-04)


### Bug Fixes

* **license:** use the canonical Apache 2.0 text and name LevelFour Inc. in NOTICE ([#212](https://github.com/LevelFourAI/levelfour-go/issues/212)) ([647186b](https://github.com/LevelFourAI/levelfour-go/commit/647186bf7c427a44bd2fb9c9b599618639e9e266))


### Documentation

* show how to create, change, pause and delete a cost alert ([af3aacf](https://github.com/LevelFourAI/levelfour-go/commit/af3aacfb858022e230f32634b40ebb2c49a6dea7))
* show how to create, change, pause and delete a cost alert ([a8bbe80](https://github.com/LevelFourAI/levelfour-go/commit/a8bbe80378eb5f335bc8286e968e7278efde3447))

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
