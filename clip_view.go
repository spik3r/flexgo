package flexgo

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// ClipView adapts a text callback to its allocated content area. It keeps the
// first h lines and first w terminal cells of each line, without wrapping,
// padding or ellipsis. Tabs become four spaces; CRLF and CR become newlines.
// Wide characters and grapheme clusters are kept whole. SGR styles and OSC-8
// hyperlinks carry across retained lines and close at each row boundary.
//
// A nil callback or a nonpositive dimension produces empty output without
// calling view. Otherwise view is called once with the supplied dimensions.
// The adapter retains no render state; view is responsible for its own state.
//
// Use ClipView with Node.View or NodeBuilder.View. It bounds callback text,
// not node borders, padding, gaps or overconstrained child allocations. Arbitrary
// terminal cursor-control sequences are outside its geometry guarantee.
func ClipView(view func(w, h int) string) func(w, h int) string {
	return func(w, h int) string {
		if view == nil || w <= 0 || h <= 0 {
			return ""
		}
		text := strings.ReplaceAll(view(w, h), "\t", "    ")
		text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
		var out strings.Builder
		var styles, link string
		row := 0
		for line := range strings.SplitSeq(text, "\n") {
			if row == h {
				break
			}
			if row > 0 {
				out.WriteByte('\n')
			}
			out.WriteString(styles)
			out.WriteString(link)
			// Truncate preserves escape sequences after the last visible cell,
			// including style changes that affect the following source line.
			line = ansi.Truncate(line, w, "")
			out.WriteString(line)
			styles, link = clipLineState(line, styles, link)
			if styles != "" {
				out.WriteString(ansi.ResetStyle)
			}
			if link != "" {
				out.WriteString(ansi.ResetHyperlink())
			}
			row++
		}
		return out.String()
	}
}

// Replay SGR changes instead of interpreting colour and attribute parameters.
// A standalone reset drops the history; OSC-8 always replaces the previous link.
func clipLineState(line, styles, link string) (string, string) {
	var state byte
	for len(line) > 0 {
		seq, _, n, next := ansi.DecodeSequence(line, state, nil)
		line, state = line[n:], next
		csi := strings.TrimPrefix(strings.TrimPrefix(seq, "\x1b["), "\x9b")
		if csi != seq && strings.HasSuffix(csi, "m") {
			if csi == "m" || csi == "0m" {
				styles = ""
			} else {
				styles += seq
			}
		}
		osc := strings.TrimPrefix(strings.TrimPrefix(seq, "\x1b]"), "\x9d")
		if osc != seq && strings.HasPrefix(osc, "8;") {
			link = seq
			body := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(osc, "\x1b\\"), "\a"), "\x9c")
			if _, uri, ok := strings.Cut(body[2:], ";"); ok && uri == "" {
				link = ""
			}
		}
	}
	return styles, link
}
