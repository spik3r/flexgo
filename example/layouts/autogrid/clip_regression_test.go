package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestBoundedClosesDiscardedStyle(t *testing.T) {
	text := "\x1b[31m" + ansi.SetHyperlink("https://example.test") + "first\nsecond" + ansi.ResetHyperlink() + ansi.ResetStyle
	out := label(text)(3, 1)
	if !strings.Contains(out, ansi.ResetStyle) || !strings.Contains(out, ansi.ResetHyperlink()) {
		t.Fatalf("discarded line held required ANSI closures: %q", out)
	}
}
