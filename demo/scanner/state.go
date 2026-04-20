package main

import (
	"fmt"
	"time"
)

// Domain types. Kept small and separate so the screens consume values,
// not ad-hoc strings — this is what lets you swap real data in later
// without touching the layout code.

type ScanStatus int

const (
	ScanIdle ScanStatus = iota
	ScanRunning
	ScanComplete
)

func (s ScanStatus) String() string {
	switch s {
	case ScanIdle:
		return "idle"
	case ScanRunning:
		return "running"
	case ScanComplete:
		return "complete"
	}
	return "?"
}

type Progress struct {
	Status       ScanStatus
	FilesScanned int
	FilesTotal   int
	Started      time.Time
	CurrentPath  string
}

type Severity int

const (
	SevLow Severity = iota
	SevMed
	SevHigh
)

func (s Severity) Label() string {
	return []string{"low", "med", "high"}[s]
}

type Finding struct {
	Severity Severity
	Path     string
	Line     int
	Message  string
}

type ScanResult struct {
	ID       string
	Started  time.Time
	Duration time.Duration
	Profile  string
	Findings int
}

// Profile is the editable scan configuration the launcher wizard
// collects before a scan starts.
type Profile struct {
	Name    string
	Root    string
	Include []string
	Exclude []string
	Deep    bool
}

func DefaultProfile() Profile {
	return Profile{
		Name:    "default",
		Root:    "./src",
		Include: []string{"*.go", "*.ts", "*.py"},
		Exclude: []string{"vendor", "node_modules"},
		Deep:    true,
	}
}

// Sample data — stand-ins for what a real scanner would produce. In a
// real app these come from commands dispatched by tea.Cmd; here they
// seed the screens so the layout is visible without wiring a backend.

func SampleProgress() Progress {
	return Progress{
		Status:       ScanRunning,
		FilesScanned: 847,
		FilesTotal:   1204,
		Started:      time.Now().Add(-47 * time.Second),
		CurrentPath:  "src/auth/session_test.go",
	}
}

func SampleFindings() []Finding {
	return []Finding{
		{SevHigh, "src/auth/session.go", 142, "credential logged at INFO level"},
		{SevHigh, "src/payments/stripe.go", 301, "webhook signature not verified"},
		{SevMed, "src/api/handler.go", 88, "request body read without size limit"},
		{SevMed, "src/db/query.go", 44, "SQL built via fmt.Sprintf"},
		{SevLow, "src/util/uuid.go", 12, "math/rand used for token generation"},
		{SevLow, "src/util/slug.go", 7, "unicode normalisation missing"},
	}
}

func SampleLogs() []string {
	lines := []string{
		"scanner started profile=default root=./src",
		"discover: 1204 files queued",
		"worker-0 claim: src/main.go",
		"worker-1 claim: src/app.go",
		"worker-2 claim: src/auth/session.go",
		"finding src/auth/session.go:142 severity=high",
		"worker-0 done: src/main.go (0 findings)",
		"worker-3 claim: src/payments/stripe.go",
		"finding src/payments/stripe.go:301 severity=high",
	}
	for i := 10; i < 40; i++ {
		lines = append(lines, fmt.Sprintf("worker-%d claim: src/pkg/file_%03d.go", i%4, i))
	}
	return lines
}

func SampleFiles() []string {
	return []string{
		"src/",
		"  auth/",
		"    session.go",
		"    session_test.go",
		"    tokens.go",
		"  payments/",
		"    stripe.go",
		"    paypal.go",
		"  api/",
		"    handler.go",
		"    middleware.go",
		"  db/",
		"    query.go",
		"    migrate.go",
		"  util/",
		"    uuid.go",
		"    slug.go",
		"  main.go",
		"  app.go",
	}
}

func SampleHistory() []ScanResult {
	now := time.Now()
	return []ScanResult{
		{"2026-04-20-a", now.Add(-3 * time.Hour), 92 * time.Second, "default", 12},
		{"2026-04-19-b", now.Add(-26 * time.Hour), 88 * time.Second, "default", 9},
		{"2026-04-19-a", now.Add(-31 * time.Hour), 104 * time.Second, "deep", 21},
		{"2026-04-18-a", now.Add(-50 * time.Hour), 67 * time.Second, "fast", 4},
		{"2026-04-17-a", now.Add(-74 * time.Hour), 99 * time.Second, "default", 11},
	}
}
