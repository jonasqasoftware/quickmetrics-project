package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

var static = fstest.MapFS{"index.html": {Data: []byte("<!doctype html><title>QuickMetrics</title>")}}

func TestHealthAndSecurityHeaders(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	New(static).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Header().Get("Content-Security-Policy") == "" || recorder.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing defensive headers: %v", recorder.Header())
	}
}

func TestAnalyzeContract(t *testing.T) {
	body := `{"runs":[{"id":"run-1","outcome":"passed","durationMs":125,"retries":0}],"thresholds":{"minimumPassRate":90,"maximumRetryRate":5,"maximumP95DurationMs":500}}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	New(static).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	for _, field := range []string{"totalRuns", "passRate", "failureRate", "retryRate", "averageDurationMs", "p95DurationMs", "decision"} {
		if _, exists := response[field]; !exists {
			t.Errorf("response missing %q", field)
		}
	}
}

func TestAnalyzeRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string
		status int
		code   string
	}{
		{name: "method", method: http.MethodGet, status: http.StatusMethodNotAllowed, code: "METHOD_NOT_ALLOWED"},
		{name: "malformed json", method: http.MethodPost, body: `{`, status: http.StatusBadRequest, code: "INVALID_JSON"},
		{name: "unknown field", method: http.MethodPost, body: `{"unexpected":true}`, status: http.StatusBadRequest, code: "INVALID_JSON"},
		{name: "invalid domain input", method: http.MethodPost, body: `{"runs":[],"thresholds":{"maximumP95DurationMs":1}}`, status: http.StatusUnprocessableEntity, code: "INVALID_INPUT"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/api/v1/analyze", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()
			New(static).ServeHTTP(recorder, request)
			if recorder.Code != test.status || !strings.Contains(recorder.Body.String(), test.code) {
				t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestStaticDashboardIsServed(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	New(static).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), "QuickMetrics") {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
