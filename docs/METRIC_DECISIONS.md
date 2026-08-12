# Metric Decisions

## Signals

- **Pass rate:** passed runs divided by total runs.
- **Failure rate:** failed runs divided by total runs.
- **Retry rate:** runs with one or more retries divided by total runs.
- **Average duration:** arithmetic mean of run duration in milliseconds.
- **P95 duration:** nearest-rank 95th percentile after sorting durations.

Rates and average duration are rounded to two decimal places. Counts and p95 duration remain integers.

## Interpretation

The API returns `investigate` when any calculated signal breaches a supplied threshold. Otherwise it returns `within-thresholds`. These labels communicate comparison results only; they do not approve a release.

Pass rate does not express severity. A single failed critical scenario may matter more than many passed low-risk scenarios. Retry rate is a stability prompt, not a flaky-test diagnosis. Duration depends on environment and workload. Teams should combine these signals with defect impact, risk, change context, and qualitative evidence.

## Deliberate exclusions

Defect density, escaped defects, mean time to repair, coverage, and code complexity are excluded because the current input contract cannot calculate them honestly. They should be added only with defined source systems, ownership, formulas, and decision use.
