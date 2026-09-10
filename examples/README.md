# Examples

Runnable examples for the LevelFour Go SDK.

## Setup

Set your API key:

```bash
export LEVELFOUR_API_KEY="l4_live_..."
```

## Running

Each example is a standalone Go module. Run from within the repository:

```bash
cd examples/<name>
go run .
```

Examples use `replace` directives to reference the local SDK source.
This means they must be run from inside the cloned repository, because they
cannot be copied standalone without updating the `go.mod` to reference
a published SDK version.

## Available Examples

| Example | Description |
|---------|-------------|
| [client](client/) | Basic client setup and API call |
| [error_handling](error_handling/) | Error predicates, status codes, and structured error extraction |
| [pagination](pagination/) | Auto-paging iterator and `CollectAll` helper |
| [middleware](middleware/) | HTTP middleware with `ClientFromContext` |
| [webhook_verification](webhook_verification/) | HMAC-SHA256 webhook signature verification |
