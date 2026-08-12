package web

import (
	"io/fs"
	"strings"
	"testing"
)

func TestDashboardContainsAccessibilityLandmarks(t *testing.T) {
	content, err := fs.ReadFile(Static(), "index.html")
	if err != nil {
		t.Fatal(err)
	}
	html := string(content)
	for _, marker := range []string{`lang="en"`, `<main`, `<nav aria-label=`, `class="skip-link"`, `aria-live="polite"`, `<label`} {
		if !strings.Contains(html, marker) {
			t.Errorf("dashboard missing accessibility marker %q", marker)
		}
	}
	if strings.Contains(html, "<script>") {
		t.Error("dashboard should not use inline scripts")
	}
}
