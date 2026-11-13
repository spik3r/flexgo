package layouts

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Static recipe labels must not wrap into space reserved for other content.
func boundedLabel(text string, w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	lines := strings.Split(strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n"), "\n")
	if len(lines) > h {
		lines = lines[:h]
		lines[h-1] += "…"
		// A discarded line may hold the closing style or hyperlink sequence.
		if ansi.Strip(text) != text {
			lines[h-1] += ansi.ResetStyle + ansi.ResetHyperlink()
		}
	}
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], w, "…")
	}
	return strings.Join(lines, "\n")
}
