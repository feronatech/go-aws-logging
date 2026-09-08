# Golang AWS Logging Utility Library

[![Static Badge](https://img.shields.io/badge/aws-SDK_v2-orange)](go.mod) [![Static Badge](https://img.shields.io/badge/Go-1.25%2B-%2300ADD8?logo=go)](go.mod) [![Go Reference](https://img.shields.io/badge/Reference-%2300ADD8?logo=go&logoColor=white)](https://pkg.go.dev/github.com/feronatech/go-aws-logging) [![GitHub License](https://img.shields.io/github/license/feronatech/go-aws-logging?label=License&color=purple)](LICENSE)



Structured JSON logging for Go services - built with AWS Lambda in mind, but cloud-agnostic by design.

Every log line is a single JSON object with consistent top-level fields (`timestamp`, `level`, `message`, `service`, `region`, `environment`), an arbitrary `context` bag for your own fields, and for errors, an automatically extracted `error` object containing the error name, type and stack trace.

The core `logging` package has no cloud SDK dependencies, so you can use it anywhere (Lambda, ECS, Cloud Run, Kubernetes, or a plain binary). An optional `logging/aws` subpackage lets you construct a logger directly from an `aws.Config`.

---

## Installation

```bash
go get github.com/feronatech/go-aws-logging
```

## Import

```go
import "github.com/feronatech/go-aws-logging/logging"
```

Optional AWS helpers:

```go
import "github.com/feronatech/go-aws-logging/logging/aws"
```

---

## Quick start

```go
package main

import (
	"errors"

	"github.com/feronatech/go-aws-logging/logging"
)

func main() {
	log := logging.NewLogger("us-east-1", "prod", "app name")

	log.Info("bla bla natural language message", map[string]any{
		"extra1": 456,
		"extra2": "Blabla",
	})

	log.Error("bla bla natural language message", map[string]any{
		"extra1": 456,
	}, errors.New("something went wrong"))
}
```

Output:

```json
{"timestamp":"2026-08-07T12:34:56.030Z","level":"INFO","message":"bla bla natural language message","service":"app name","region":"us-east-1","environment":"prod","context":{"extra1":456,"extra2":"Blabla"}}
```

---

## Log output format

Every entry shares the same envelope:

| Field | Description |
| --- | --- |
| `timestamp` | ISO-8601 UTC timestamp with milliseconds |
| `level` | `DEBUG`, `INFO`, `WARN` or `ERROR` |
| `message` | The natural-language message you passed in |
| `service` | The application name the logger was created with |
| `region` | The region the logger was created with |
| `environment` | The environment the logger was created with |
| `context` | The `extra` map you passed in |
| `error` | Only present on `Error`/`ErrorCtx`: see below |

### Info / Debug / Warn

```json
{
  "timestamp": "2026-08-07T12:34:56.030Z",
  "level": "INFO",
  "message": "bla bla natural language message",
  "service": "app name",
  "region": "us-east-1",
  "environment": "prod",
  "context": { "extra1": 456, "extra2": "Blabla" }
}
```

### Error

```json
{
  "timestamp": "2026-08-07T12:34:56.030Z",
  "level": "ERROR",
  "message": "bla bla natural language message",
  "service": "app name",
  "region": "us-east-1",
  "environment": "prod",
  "context": { "extra1": 456, "extra2": "Blabla" },
  "error": {
    "name": "Validation failed",
    "type": "errors.ValidationError",
    "stack": []
  }
}
```

The `error` object is derived from the `error` value you pass in; you don't need to format anything yourself:

- **`name`**: the error message (`err.Error()`)
- **`type`**: the concrete Go type of the error, including its package
- **`stack`**: the stack trace, when one is available on the error

---

## Creating a logger

There are several constructors, depending on where your configuration lives.

### `NewLogger(region, env, app)`

The most explicit option: pass the values directly.

```go
log := logging.NewLogger("us-east-1", "staging", "order-api")
```

| Parameter | Description |
| --- | --- |
| `region` | Deployment region (e.g. `us-east-1`). Not AWS-specific, any region identifier works. |
| `env` | Environment name (e.g. `dev`, `staging`, `prod`). Emitted as `environment`. |
| `app` | Application or service name. Emitted as `service`. |

### `FromContext(ctx)`

Reads region, environment and application from the given context, and keeps that context as the logger's context.

```go
func handler(ctx context.Context, event MyEvent) error {
	log := logging.FromContext(ctx)
	log.Info("handling event", nil)
	return nil
}
```

### `FromOptions(options)`

For full control, pass a `logging.LoggingOptions` value.

```go
log := logging.FromOptions(logging.LoggingOptions{
	Region: "ap-southeast-2",
	Env:    "prod",
	App:    "inventory-worker",
})
```

### `aws.FromConfig(cfg)` (optional)

If you already have an `aws.Config` in hand, the `logging/aws` subpackage can derive the logger for you:

```go
import (
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awslogging "github.com/feronatech/go-aws-logging/logging/aws"
)

cfg, err := awsconfig.LoadDefaultConfig(ctx)
if err != nil {
	return err
}

log := awslogging.FromConfig(cfg)
log.Info("logger initialised from aws.Config", nil)
```

`aws.FromConfig` returns the same `logging.Logger` as the other constructors, so it's a drop-in replacement.

> Keeping the AWS integration in a separate package means the core `logging` package stays free of AWS SDK dependencies.

---

## Logging

Each level comes in two flavours: a plain variant that uses the logger's own context, and a `Ctx` variant that takes a context for that single call.

```go
log.Debug("cache lookup", map[string]any{"key": cacheKey})
log.Info("user created", map[string]any{"userId": id})
log.Warn("retrying request", map[string]any{"attempt": 3})
log.Error("payment failed", map[string]any{"orderId": orderID}, err)
```

Pass `nil` when you have no extra fields:

```go
log.Info("service starting", nil)
```

### The `extra` map

Everything you put in `extra` ends up under the `context` key in the JSON output; keeping your custom fields cleanly separated from the log envelope:

```go
log.Info("bla bla natural language message", map[string]any{
	"extra1": 456,
	"extra2": "Blabla",
})
```

```json
{"...":"...","context":{"extra1":456,"extra2":"Blabla"}}
```

Values can be any JSON-serialisable type: numbers, strings, booleans, slices, maps or structs.

### Errors

`Error` and `ErrorCtx` take the error as their **last** argument, after the extra map:

```go
if err := validate(order); err != nil {
	log.Error("order validation failed", map[string]any{
		"orderId": order.ID,
	}, err)
	return err
}
```

The library inspects the error and extracts its name, type and stack trace into the `error` field automatically.

### Levels at a glance

| Method | Context variant | Typical use |
| --- | --- | --- |
| `Debug` | `DebugCtx` | Verbose diagnostics, usually disabled in production |
| `Info` | `InfoCtx` | Normal operational events |
| `Warn` | `WarnCtx` | Unexpected but recoverable conditions |
| `Error` | `ErrorCtx` | Failures that need attention (takes an `error`) |

---

## Contexts

A logger always carries a context. It defaults to `context.Background()`, and you can replace it whenever you have a more meaningful one, typically at the start of a function or handler.

### `SetContext`

`SetContext` sets the context the logger remembers and uses for subsequent calls. It returns a `Logger`, so it chains nicely:

```go
func processOrder(ctx context.Context, log logging.Logger, order Order) error {
	log = log.SetContext(ctx)

	log.Info("processing order", map[string]any{"orderId": order.ID})
	// ...
	log.Info("order processed", nil)
	return nil
}
```

This is the usual pattern in a Lambda handler or HTTP middleware: set the real request context once at the top, then log normally throughout the scope.

### `Ctx` variants

The `Ctx` variants let you supply a context for a single call, without changing the logger's remembered context. Useful when you're passing a logger around and only occasionally have a different context at hand:

```go
log.DebugCtx(ctx, "cache lookup", nil)
log.InfoCtx(ctx, "user created", map[string]any{"userId": id})
log.WarnCtx(ctx, "retrying request", map[string]any{"attempt": 3})
log.ErrorCtx(ctx, "payment failed", map[string]any{"orderId": orderID}, err)
```

---

## Usage with AWS Lambda

```go
package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/feronatech/go-aws-logging/logging"
)

type Request struct {
	OrderID string `json:"orderId"`
}

func handler(ctx context.Context, req Request) error {
	log := logging.FromContext(ctx)

	log.Info("processing order", map[string]any{"orderId": req.OrderID})

	if req.OrderID == "" {
		err := errors.New("missing order id")
		log.Error("cannot process request", nil, err)
		return err
	}

	log.Debug("order processed successfully", map[string]any{
		"orderId":    req.OrderID,
		"durationMs": 42,
	})
	return nil
}

func main() {
	lambda.Start(handler)
}
```

---

## Usage outside AWS

Because the constructors take plain strings (or a `LoggingOptions` struct), the same logger works anywhere:

```go
log := logging.FromOptions(logging.LoggingOptions{
	Region: "europe-west4", // GCP region
	Env:    "prod",
	App:    "billing-service",
})

log.Info("running on Cloud Run", nil)
```

Or read straight from environment variables:

```go
log := logging.NewLogger(
	os.Getenv("REGION"),
	os.Getenv("ENVIRONMENT"),
	os.Getenv("APP_NAME"),
)
```

---

## Package layout

| Package | Purpose |
| --- | --- |
| `logging` | Core `Logger` type and cloud-agnostic constructors |
| `logging/aws` | AWS-specific helper: `FromConfig(aws.Config)` |

---

## API reference

```go
// Constructors
func NewLogger(region, env, app string) Logger
func FromContext(ctx context.Context) Logger
func FromOptions(options LoggingOptions) Logger

// logging/aws
func FromConfig(cfg aws.Config) logging.Logger

// Logger
type Logger interface {
	// SetContext replaces the context the logger uses for subsequent calls.
	// Defaults to context.Background().
	SetContext(ctx context.Context) Logger

	Debug(message string, extra map[string]any)
	DebugCtx(ctx context.Context, message string, extra map[string]any)

	Info(message string, extra map[string]any)
	InfoCtx(ctx context.Context, message string, extra map[string]any)

	Warn(message string, extra map[string]any)
	WarnCtx(ctx context.Context, message string, extra map[string]any)

	Error(message string, extra map[string]any, err error)
	ErrorCtx(ctx context.Context, message string, extra map[string]any, err error)
}
```

---

## Contributing

Issues and pull requests are welcome. Please open an issue first for larger changes so we can discuss the approach.

## License

See [LICENSE](LICENSE) for details.
