# QuickMetrics

QuickMetrics is an educational Quality Engineering decision lab. It turns explicit test-run data into traceable signals, compares those signals with context-specific thresholds, and explains why a result needs investigation.

It addresses a common metrics failure: presenting percentages without their source, calculation, decision rule, or residual risk. Every displayed value is calculated by the Go API from visible input. The included sample is synthetic and is not production evidence.

## What it demonstrates

- deterministic calculation of pass, failure, retry, average-duration, and p95-duration signals;
- thresholds supplied by the user instead of universal quality targets;
- an API contract with strict JSON parsing and stable errors;
- risk-based unit, integration, contract, and accessibility checks;
- a dependency-free accessible dashboard embedded in the server binary;
- reproducible quality gates in GitHub Actions.

## Run locally

Requirements: Go 1.18 or newer.

```bash
go run ./cmd/server
```

Open <http://localhost:8080>. The default port can be changed with `PORT`.

Analyze the documented synthetic fixture directly:

```bash
curl -s http://localhost:8080/api/v1/analyze \
  -H 'content-type: application/json' \
  --data-binary @testdata/sample-runs.json
```

## API contract

| Method | Route | Purpose |
| --- | --- | --- |
| `GET` | `/health` | Process liveness |
| `POST` | `/api/v1/analyze` | Validate input, calculate signals, and compare thresholds |
| `GET` | `/` | Interactive dashboard |

Each run requires a unique `id`, an `outcome` of `passed` or `failed`, a positive `durationMs`, and a non-negative `retries` value. Percent thresholds use the inclusive range 0–100. Unknown JSON fields are rejected so contract drift remains visible.

The p95 calculation uses the nearest-rank method over run durations. `retryRate` is the percentage of runs that required one or more retries; it is a stability signal, not proof that a test is flaky. A repeated-test history would be required for that conclusion.

## Quality Strategy

Context, risks, test layers, evidence, and exit criteria are documented in [Quality Strategy](docs/QUALITY_STRATEGY.md). Metric semantics and decision constraints are documented in [Metric Decisions](docs/METRIC_DECISIONS.md).

Run the local gates:

```bash
gofmt -w .
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./...
```

The test suite covers exact calculations, nearest-rank p95, threshold boundaries, invalid datasets, duplicate IDs, malformed JSON, unknown fields, method restrictions, defensive headers, API response fields, static delivery, and key accessibility landmarks.

## Architecture

```text
Browser or API client
        |
        v
net/http boundary -> strict decoding -> metrics domain -> decision explanation
        |
        +-> embedded HTML, CSS, and JavaScript dashboard
```

The standard library keeps the service small and auditable. The domain calculator is independent of HTTP, while the dashboard and API ship as one binary. A database, framework, or frontend build chain would add cost without improving this lab's primary decision-making demonstration.

## Limitations and residual risks

- Input is evaluated in one request and is not persisted.
- The API has no authentication because it stores no data and is intended for local demonstration.
- Retry rate alone cannot identify flaky tests or root cause.
- Thresholds express the user's context; being within them is not proof of release readiness.
- There is no production telemetry, trend analysis, distributed tracing, DAST, or independent accessibility audit.
- Browser output encoding is implemented through DOM `textContent`; future rich rendering must preserve that boundary.

Do not submit confidential build data, customer information, credentials, or proprietary failure messages to an untrusted deployment.

## License

MIT — see [LICENSE](LICENSE).
