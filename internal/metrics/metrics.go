package metrics

import (
	"errors"
	"math"
	"sort"
)

const maxRuns = 10000

type Run struct {
	ID         string `json:"id"`
	Outcome    string `json:"outcome"`
	DurationMS int64  `json:"durationMs"`
	Retries    int    `json:"retries"`
}

type Thresholds struct {
	MinimumPassRate      float64 `json:"minimumPassRate"`
	MaximumRetryRate     float64 `json:"maximumRetryRate"`
	MaximumP95DurationMS int64   `json:"maximumP95DurationMs"`
}

type Input struct {
	Runs       []Run      `json:"runs"`
	Thresholds Thresholds `json:"thresholds"`
}

type Decision struct {
	Status  string   `json:"status"`
	Reasons []string `json:"reasons"`
}

type Result struct {
	TotalRuns         int      `json:"totalRuns"`
	PassedRuns        int      `json:"passedRuns"`
	FailedRuns        int      `json:"failedRuns"`
	RetriedRuns       int      `json:"retriedRuns"`
	PassRate          float64  `json:"passRate"`
	FailureRate       float64  `json:"failureRate"`
	RetryRate         float64  `json:"retryRate"`
	AverageDurationMS float64  `json:"averageDurationMs"`
	P95DurationMS     int64    `json:"p95DurationMs"`
	Decision          Decision `json:"decision"`
}

func Analyze(input Input) (Result, error) {
	if len(input.Runs) == 0 {
		return Result{}, errors.New("at least one run is required")
	}
	if len(input.Runs) > maxRuns {
		return Result{}, errors.New("at most 10000 runs are accepted")
	}
	if input.Thresholds.MinimumPassRate < 0 || input.Thresholds.MinimumPassRate > 100 ||
		input.Thresholds.MaximumRetryRate < 0 || input.Thresholds.MaximumRetryRate > 100 ||
		input.Thresholds.MaximumP95DurationMS <= 0 {
		return Result{}, errors.New("thresholds must be valid percentages and a positive duration")
	}

	result := Result{TotalRuns: len(input.Runs)}
	durations := make([]int64, 0, len(input.Runs))
	var totalDuration int64
	seenIDs := make(map[string]struct{}, len(input.Runs))
	for _, run := range input.Runs {
		if run.ID == "" || run.DurationMS <= 0 || run.Retries < 0 {
			return Result{}, errors.New("each run requires an id, positive duration, and non-negative retries")
		}
		if _, exists := seenIDs[run.ID]; exists {
			return Result{}, errors.New("run ids must be unique")
		}
		seenIDs[run.ID] = struct{}{}
		switch run.Outcome {
		case "passed":
			result.PassedRuns++
		case "failed":
			result.FailedRuns++
		default:
			return Result{}, errors.New("outcome must be passed or failed")
		}
		if run.Retries > 0 {
			result.RetriedRuns++
		}
		totalDuration += run.DurationMS
		durations = append(durations, run.DurationMS)
	}

	result.PassRate = percentage(result.PassedRuns, result.TotalRuns)
	result.FailureRate = percentage(result.FailedRuns, result.TotalRuns)
	result.RetryRate = percentage(result.RetriedRuns, result.TotalRuns)
	result.AverageDurationMS = round(float64(totalDuration)/float64(result.TotalRuns), 2)
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	result.P95DurationMS = durations[int(math.Ceil(float64(len(durations))*0.95))-1]
	result.Decision = decide(result, input.Thresholds)
	return result, nil
}

func percentage(value, total int) float64 {
	return round(float64(value)*100/float64(total), 2)
}

func round(value float64, places int) float64 {
	power := math.Pow10(places)
	return math.Round(value*power) / power
}

func decide(result Result, thresholds Thresholds) Decision {
	reasons := make([]string, 0, 3)
	if result.PassRate < thresholds.MinimumPassRate {
		reasons = append(reasons, "pass rate is below the configured minimum")
	}
	if result.RetryRate > thresholds.MaximumRetryRate {
		reasons = append(reasons, "retry rate is above the configured maximum")
	}
	if result.P95DurationMS > thresholds.MaximumP95DurationMS {
		reasons = append(reasons, "p95 duration is above the configured maximum")
	}
	if len(reasons) > 0 {
		return Decision{Status: "investigate", Reasons: reasons}
	}
	return Decision{Status: "within-thresholds", Reasons: []string{"all calculated signals are within the configured thresholds"}}
}
