package httpapi

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/jonasqasoftware/quickmetrics-project/internal/metrics"
)

const maxBodyBytes = 1 << 20

type API struct {
	static fs.FS
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(static fs.FS) http.Handler {
	api := API{static: static}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", api.health)
	mux.HandleFunc("/api/v1/analyze", api.analyze)
	mux.Handle("/", http.FileServer(http.FS(static)))
	return securityHeaders(requestLog(mux))
}

func (api API) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "use GET for this endpoint")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (api API) analyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "use POST for this endpoint")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input metrics.Input
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must match the documented JSON contract")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must contain one JSON object")
		return
	}
	result, err := metrics.Analyze(input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "INVALID_INPUT", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; connect-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration_ms=%d", r.Method, r.URL.Path, time.Since(started).Milliseconds())
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Code: code, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("response_encode_error=%q", err)
	}
}
