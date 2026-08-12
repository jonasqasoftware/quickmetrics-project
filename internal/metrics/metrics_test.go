package metrics

import (
	"reflect"
	"testing"
)

func TestAnalyzeCalculatesSignalsAndDecision(t *testing.T) {
	result, err := Analyze(Input{
		Runs: []Run{
			{ID: "1", Outcome: "passed", DurationMS: 100, Retries: 0},
			{ID: "2", Outcome: "passed", DurationMS: 200, Retries: 0},
			{ID: "3", Outcome: "passed", DurationMS: 300, Retries: 0},
			{ID: "4", Outcome: "failed", DurationMS: 400, Retries: 1},
		},
		Thresholds: Thresholds{MinimumPassRate: 80, MaximumRetryRate: 10, MaximumP95DurationMS: 350},
	})
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if result.PassRate != 75 || result.FailureRate != 25 || result.RetryRate != 25 {
		t.Fatalf("unexpected rates: %+v", result)
	}
	if result.AverageDurationMS != 250 || result.P95DurationMS != 400 {
		t.Fatalf("unexpected durations: %+v", result)
	}
	wantReasons := []string{
		"pass rate is below the configured minimum",
		"retry rate is above the configured maximum",
		"p95 duration is above the configured maximum",
	}
	if result.Decision.Status != "investigate" || !reflect.DeepEqual(result.Decision.Reasons, wantReasons) {
		t.Fatalf("unexpected decision: %+v", result.Decision)
	}
}

func TestAnalyzeReportsSignalsWithinThresholds(t *testing.T) {
	result, err := Analyze(Input{
		Runs:       []Run{{ID: "1", Outcome: "passed", DurationMS: 100, Retries: 0}},
		Thresholds: Thresholds{MinimumPassRate: 100, MaximumRetryRate: 0, MaximumP95DurationMS: 100},
	})
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if result.Decision.Status != "within-thresholds" {
		t.Fatalf("status = %q", result.Decision.Status)
	}
}

func TestAnalyzeRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input Input
	}{
		{name: "empty dataset", input: Input{}},
		{name: "invalid threshold", input: Input{Runs: []Run{{ID: "1", Outcome: "passed", DurationMS: 1}}, Thresholds: Thresholds{MinimumPassRate: 101, MaximumP95DurationMS: 1}}},
		{name: "missing id", input: Input{Runs: []Run{{Outcome: "passed", DurationMS: 1}}, Thresholds: Thresholds{MaximumP95DurationMS: 1}}},
		{name: "duplicate id", input: Input{Runs: []Run{{ID: "1", Outcome: "passed", DurationMS: 1}, {ID: "1", Outcome: "failed", DurationMS: 1}}, Thresholds: Thresholds{MaximumP95DurationMS: 1}}},
		{name: "invalid outcome", input: Input{Runs: []Run{{ID: "1", Outcome: "unknown", DurationMS: 1}}, Thresholds: Thresholds{MaximumP95DurationMS: 1}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Analyze(test.input); err == nil {
				t.Fatal("Analyze() error = nil, want validation error")
			}
		})
	}
}
