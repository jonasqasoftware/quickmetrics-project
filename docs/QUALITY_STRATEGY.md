# Quality Strategy

## Context

QuickMetrics is a local educational application that calculates a small set of test-execution signals. Its quality risk is not financial loss or production availability; it is misleading a reviewer with an incorrect calculation, an unexplained threshold, fabricated evidence, or inaccessible presentation.

## Risk-to-test matrix

| Priority | Risk | Control | Evidence | Residual risk |
| --- | --- | --- | --- | --- |
| Critical | A calculated signal is wrong | Pure domain calculator and exact-value tests | Unit tests cover rates, averages, and nearest-rank p95 | Tests use selected datasets, not property-based generation |
| Critical | A metric is presented without provenance | API accepts explicit runs; samples are labeled synthetic | Contract test and documentation review | Users can still provide poor-quality source data |
| High | A universal threshold is implied | Thresholds are required request inputs | Boundary and decision tests | A user may choose thresholds without business context |
| High | Invalid data creates plausible output | Strict validation, unique IDs, capped dataset, stable 4xx errors | Table-driven domain and HTTP negative tests | Semantic duplication across external systems is not detectable |
| High | API contract drifts silently | Unknown fields rejected and response fields asserted | HTTP contract test | No version negotiation beyond the `/v1` path |
| Medium | Dashboard excludes keyboard or assistive-tech users | Semantic landmarks, labels, focus styles, live results, reduced motion | Static accessibility checks and manual keyboard review | No independent screen-reader audit |
| Medium | Browser renders input as executable markup | Dynamic values use `textContent`; restrictive CSP | Code review and defensive-header test | Future UI changes could weaken the boundary |
| Medium | A large request exhausts resources | 1 MB HTTP body limit and 10,000-run domain limit | Invalid-input coverage | No distributed rate limiting or authentication |

## Test levels

1. Domain unit tests validate calculations and decisions without HTTP.
2. HTTP integration tests validate parsing, status codes, headers, errors, and response shape.
3. Embedded-asset tests validate delivery and key accessibility landmarks.
4. CI runs formatting verification, `go vet`, race-enabled tests, coverage reporting, build, and secret-marker scanning.

The same calculation is not duplicated through a browser E2E test because the domain and API levels cover it faster and more precisely. A small manual dashboard check verifies responsive layout, keyboard flow, and error/result presentation.

## Decision criteria

A change is ready for review when formatting, vet, race-enabled tests, build, and secret scan pass; calculation changes include exact expected-value tests; and documentation continues to distinguish synthetic samples from evidence.

Coverage is reported as diagnostic evidence. It is not treated as proof of quality or optimized by removing meaningful branches.
