package main

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// UI primitives the screens compose. Two rules shape this file:
//
//  1. Every View callback paints its full (w, h) region with a defined
//     background colour. Terminals (and tmux in particular) leave
//     unpainted cells transparent, so a hole in one line of output
//     shows the wallpaper through. The fixBg helper below handles this.
//
//  2. Layout composition uses flexgo *Node trees wherever possible
//     rather than ad-hoc lipgloss JoinHorizontal/JoinVertical. That
//     way flexgo's ambient-background propagation keeps backdrops
//     consistent and we get alignment for free.

// Palette kept in one place so the demo has a consistent look and you
// can see how a Theme abstraction would fit. Values are lifted from
// Tokyo Night Storm (https://github.com/folke/tokyonight.nvim) —
// bg is #24283b, the accent blue is #7aa2f7, etc. Keeping the hex
// literal here makes it obvious where to swap a palette.
var (
	colBg        = lipgloss.Color("#24283b") // storm bg
	colPanel     = lipgloss.Color("#1f2335") // bg_dark — sidebars, cards
	colAccent    = lipgloss.Color("#1f3d5f") // muted blue — header, active
	colMuted     = lipgloss.Color("#3b4261") // fg_gutter — dividers, inactive
	colText      = lipgloss.Color("#a9b1d6") // fg_dark — body text, dimmer than fg
	colSubtle    = lipgloss.Color("#737aa2") // dark5 — subtitles, even dimmer
	colHighlight = lipgloss.Color("#c9a45a") // dimmed yellow — emphasis / brand
	colSelected  = lipgloss.Color("#394b70") // blue7 — selection bg, softer
	colDanger    = lipgloss.Color("#c75c78") // muted red
	colWarn      = lipgloss.Color("#c98356") // muted orange
	colOk        = lipgloss.Color("#85b569") // muted green
)

// --- page shell -----------------------------------------------------

// The global help hint, shown in the always-on footer. Changes here
// propagate to every screen — that's the point of the shell.
const globalNavHint = "? help  ·  1/2/3 screens  ·  q quit"

// The screen entries for the subheader's flex-between indicator.
// Order matches the Screen enum.
var screenTabs = []struct {
	key    string
	label  string
	screen Screen
}{
	{"1", "Launcher", ScreenLauncher},
	{"2", "Scan", ScreenScan},
	{"3", "History", ScreenHistory},
}

// pageShell wraps a body in the shared chrome: a 1-cell outer frame
// in the page bg, then a branded title bar, subheader with flex-
// between subtitle + screen indicator, a margined body slot, and a
// two-row footer carrying screen-specific + global keymap hints.
//
// Screens just provide Subtitle + Footer + Body. Everything else —
// title, screen indicator, global hints, padding — lives here, so
// "consistent look" stays consistent across screens.
func pageShell(subtitle string, active Screen, body *flexgo.Node, screenHint string) *flexgo.Node {
	chrome := &flexgo.Node{
		Flex:       1,
		Dir:        flexgo.Col,
		Background: colBg,
		Children: []*flexgo.Node{
			// Brand title bar — three rows of colAccent, "FlexGo" centred.
			{Height: 3, Background: colAccent, View: brandTitleView()},

			// Subheader — subtitle on left, screen indicator on right.
			{Height: 1, Background: colPanel, View: subheaderView(subtitle, active)},

			// Divider — a faint horizontal rule to separate chrome from body.
			{Height: 1, Background: colBg, View: dividerView()},

			// Body slot, margined with explicit spacer nodes rather than
			// the Node.Paddings field. Paddings + a full-width bg relies
			// on lipgloss padding every cell with the given background,
			// and multi-line styled content breaks that guarantee — the
			// trailing padding shows through to the terminal default.
			// Solid Flex-spacer children dodge the issue entirely.
			{
				Flex:       1,
				Dir:        flexgo.Row,
				Background: colBg,
				Children: []*flexgo.Node{
					{Width: 2, Background: colBg, View: solidView(colBg)},
					body,
					{Width: 2, Background: colBg, View: solidView(colBg)},
				},
			},

			// Spacer row between body and footer.
			{Height: 1, Background: colBg, View: solidView(colBg)},

			// Footer row 1 — screen-specific hint.
			{Height: 1, Background: colBg, View: screenHintView(screenHint)},

			// Footer row 2 — global hints; always the same.
			{Height: 1, Background: colPanel, View: globalHintView()},
		},
	}

	// Outer 1-cell frame so the chrome doesn't touch the terminal
	// edges — gives the whole UI a small breathing margin.
	return &flexgo.Node{
		Dir:        flexgo.Col,
		Background: colBg,
		Children: []*flexgo.Node{
			{Height: 1, Background: colBg, View: solidView(colBg)},
			{
				Flex:       1,
				Dir:        flexgo.Row,
				Background: colBg,
				Children: []*flexgo.Node{
					{Width: 1, Background: colBg, View: solidView(colBg)},
					chrome,
					{Width: 1, Background: colBg, View: solidView(colBg)},
				},
			},
			{Height: 1, Background: colBg, View: solidView(colBg)},
		},
	}
}

// brandTitleView paints "FlexGo" centred both axes in a 3-row bar.
// The leading ⬢ + spacing gives it a brand feel without needing a
// real logo asset.
func brandTitleView() func(w, h int) string {
	return func(w, h int) string {
		barBg := lipgloss.NewStyle().Background(colAccent)
		brand := barBg.
			Foreground(colHighlight).
			Bold(true).
			Italic(true).
			Render("⬢  FlexGo")
		return barBg.
			Width(w).Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(brand)
	}
}

// subheaderView renders the contextual bar just under the title: the
// screen's subtitle on the left, and a "1 Launcher · 2 Scan ·
// 3 History" indicator on the right with the active entry highlighted.
// The two sides are separated by a flex spacer so the indicator hugs
// the right edge at any terminal width — equivalent to CSS
// justify-content: space-between.
func subheaderView(subtitle string, active Screen) func(w, h int) string {
	return func(w, h int) string {
		barBg := lipgloss.NewStyle().Background(colPanel)

		left := barBg.
			Foreground(colText).
			Italic(true).
			Padding(0, 2).
			Render(subtitle)

		rendered := make([]string, 0, len(screenTabs))
		for _, it := range screenTabs {
			entry := it.key + " " + it.label
			if it.screen == active {
				rendered = append(rendered, barBg.Foreground(colHighlight).Bold(true).Render(entry))
			} else {
				rendered = append(rendered, barBg.Foreground(colMuted).Render(entry))
			}
		}
		sep := barBg.Foreground(colMuted).Render("  ·  ")
		right := barBg.Padding(0, 2).Render(strings.Join(rendered, sep))

		gapWidth := w - lipgloss.Width(left) - lipgloss.Width(right)
		if gapWidth < 0 {
			gapWidth = 0
		}
		gap := barBg.Width(gapWidth).Render("")
		row := lipgloss.JoinHorizontal(lipgloss.Top, left, gap, right)
		return barBg.Width(w).Height(h).Render(row)
	}
}

// dividerView renders a single row of box-drawing rule in colMuted on
// the page background — the visual seam between chrome and body.
func dividerView() func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(colBg).
			Foreground(colMuted).
			Render(strings.Repeat("─", w))
	}
}

// screenHintView is the first footer row: the active screen's own
// keymap hint, in muted text on the page bg.
func screenHintView(hint string) func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(colBg).
			Foreground(colSubtle).
			Padding(0, 2).
			Render(hint)
	}
}

// globalHintView is the last row of the footer: the fixed global
// bindings. Centred to look like a status line; also doubles as a
// visual anchor so the user knows the screen "ends" there.
func globalHintView() func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(colPanel).
			Foreground(colSubtle).
			Align(lipgloss.Center, lipgloss.Center).
			Render(globalNavHint)
	}
}

// solidView paints a uniform bg fill. Used as the explicit spacer
// between body and footer so no cell is left unpainted.
func solidView(bg color.Color) func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(bg).
			Render("")
	}
}

// --- viewport -------------------------------------------------------

// ViewportState holds scroll position for one scrollable panel. The
// screen owns the state; viewportView paints it.
type ViewportState struct {
	Offset int
}

// ScrollBy applies a delta to the offset, clamping so the user can't
// scroll past either end of content. The visible-row count is needed
// to know where "bottom" is; it's passed in because only the screen
// knows the current allocated height.
func (v *ViewportState) ScrollBy(delta, contentLines, visibleLines int) {
	v.Offset += delta
	maxOffset := contentLines - visibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}
	if v.Offset < 0 {
		v.Offset = 0
	}
	if v.Offset > maxOffset {
		v.Offset = maxOffset
	}
}

// viewportView renders `lines` into a (w, h) box starting at
// state.Offset with a scroll indicator on the right. Every cell is
// painted with colBg — crucially including trailing whitespace on
// each line, so terminal cells don't show through.
//
// A production build should swap this for bubbles/viewport, which
// handles word wrap, horizontal scroll, and mouse events.
func viewportView(state *ViewportState, lines []string) func(w, h int) string {
	return func(w, h int) string {
		if h <= 0 || w <= 0 {
			return ""
		}
		maxOffset := len(lines) - h
		if maxOffset < 0 {
			maxOffset = 0
		}
		if state.Offset > maxOffset {
			state.Offset = maxOffset
		}
		contentW := w - 1 // reserve one column for the scroll indicator

		rendered := make([]string, h)
		for i := 0; i < h; i++ {
			src := ""
			if idx := state.Offset + i; idx >= 0 && idx < len(lines) {
				src = lines[idx]
			}
			rendered[i] = padLineBg(src, contentW, colBg, colText)
		}
		body := strings.Join(rendered, "\n")

		indicator := scrollIndicator(state.Offset, maxOffset, h)
		return lipgloss.JoinHorizontal(lipgloss.Top, body, indicator)
	}
}

func scrollIndicator(offset, maxOffset, h int) string {
	bg := lipgloss.NewStyle().Background(colBg)
	if maxOffset == 0 {
		cells := make([]string, h)
		for i := range cells {
			cells[i] = bg.Foreground(colMuted).Render("│")
		}
		return strings.Join(cells, "\n")
	}
	pos := offset * (h - 1) / maxOffset
	out := make([]string, h)
	for i := range out {
		if i == pos {
			out[i] = bg.Foreground(colAccent).Render("█")
		} else {
			out[i] = bg.Foreground(colMuted).Render("│")
		}
	}
	return strings.Join(out, "\n")
}

// --- filetree (stub) ------------------------------------------------

// filetreeView renders `lines` with one row highlighted. The sidebar
// filetree is a list with a cursor — enough to demonstrate the
// pattern; swap for a real tree widget when you need expand/collapse.
func filetreeView(lines []string, cursor int) func(w, h int) string {
	return func(w, h int) string {
		if h <= 0 {
			return ""
		}
		offset := 0
		if cursor >= h {
			offset = cursor - h + 1
		}
		rendered := make([]string, h)
		for i := 0; i < h; i++ {
			idx := offset + i
			var line string
			if idx < len(lines) {
				line = lines[idx]
			}
			bg := colBg
			fg := colText
			if idx == cursor && idx < len(lines) {
				bg = colSelected
				fg = colHighlight
			}
			rendered[i] = padLineBg(line, w, bg, fg)
		}
		return strings.Join(rendered, "\n")
	}
}

// --- keymap modal card ----------------------------------------------

// buildKeymapCard returns the keymap modal as a body-friendly node.
// Sits inside pageShell's body slot, so the shared title and footer
// stay visible while help is open.
//
// Centring uses explicit Flex:1 spacer nodes (each with solidView)
// rather than Justify/Align or auto-margins. Both of those approaches
// rely on lipgloss to paint trailing padding cells with the correct
// background, and multi-line styled content can leave the right-edge
// cells transparent — which in tmux shows wallpaper through the
// "background". Explicit spacers make each cell its own bg-painted
// leaf, so there's nothing for the terminal to leave unpainted.
func buildKeymapCard(keys KeyMap) *flexgo.Node {
	entries := keys.HelpEntries()
	body := func(w, h int) string {
		rows := make([]string, 0, len(entries))
		bg := lipgloss.NewStyle().Background(colPanel)
		for _, e := range entries {
			if e.Keys == "" && e.Desc == "" {
				rows = append(rows, bg.Width(w-4).Render(""))
				continue
			}
			k := bg.Width(18).Foreground(colAccent).Bold(true).Render(e.Keys)
			d := bg.Foreground(colText).Render(e.Desc)
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, k, d))
		}
		return bg.
			Width(w).Height(h).
			Padding(1, 2).
			Render(strings.Join(rows, "\n"))
	}

	card := &flexgo.Node{
		Width:            56,
		Height:           24,
		ShowBorder:       true,
		Border:           lipgloss.RoundedBorder(),
		BorderForeground: colAccent,
		Background:       colPanel,
		Dir:              flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 2, Background: colAccent, View: cardTitleBar("Keybindings")},
			{Flex: 1, Background: colPanel, View: body},
		},
	}
	return centered(card, 24)
}

// centered wraps a fixed-size node in Flex-spacer rows and columns so
// it renders in the middle of whatever slot the parent gives it. The
// cardHeight parameter is the node's fixed Height — needed because the
// middle row has to declare a height for the Col-distribute to keep
// the card at its natural size.
//
// Structure:
//
//	[Flex:1 spacer row]
//	[Row: Flex:1 spacer | card | Flex:1 spacer]  (height = cardHeight)
//	[Flex:1 spacer row]
//
// Every spacer is a leaf with Background:colBg and a solid View, so no
// cell is ever unpainted.
func centered(card *flexgo.Node, cardHeight int) *flexgo.Node {
	middleRow := &flexgo.Node{
		Height:     cardHeight,
		Dir:        flexgo.Row,
		Background: colBg,
		Children: []*flexgo.Node{
			{Flex: 1, Background: colBg, View: solidView(colBg)},
			card,
			{Flex: 1, Background: colBg, View: solidView(colBg)},
		},
	}
	return &flexgo.Node{
		Flex:       1,
		Dir:        flexgo.Col,
		Background: colBg,
		Children: []*flexgo.Node{
			{Flex: 1, Background: colBg, View: solidView(colBg)},
			middleRow,
			{Flex: 1, Background: colBg, View: solidView(colBg)},
		},
	}
}

// cardTitleBar is the short header strip used inside centered cards
// (the keymap modal, the launcher). Keeping it a single function keeps
// the two cards visually consistent.
func cardTitleBar(title string) func(w, h int) string {
	return func(w, h int) string {
		barBg := lipgloss.NewStyle().Background(colAccent)
		text := barBg.
			Foreground(colHighlight).Bold(true).
			Padding(0, 2).
			Render(title)
		gapWidth := w - lipgloss.Width(text)
		if gapWidth < 0 {
			gapWidth = 0
		}
		gap := barBg.Width(gapWidth).Render("")
		row := lipgloss.JoinHorizontal(lipgloss.Top, text, gap)
		return barBg.Width(w).Height(h).Render(row)
	}
}

// --- small style helpers --------------------------------------------

// headingStyle is the faded-bold section header used inside panels
// with colPanel as their backdrop.
func headingStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(colPanel).
		Foreground(colSubtle).
		Bold(true)
}

// textStyle is the default body-text style for a given bg. Foreground
// is colText; callers paint on whichever bg the parent container uses.
func textStyle(bg color.Color) lipgloss.Style {
	return lipgloss.NewStyle().Background(bg).Foreground(colText)
}

// panelPaint renders a slice of lines into a (w, h) box with a fixed
// bg, applying padding and filling every cell. Used by sidebars and
// cards where the backdrop must be solid.
func panelPaint(lines []string, w, h int, bg color.Color, padY, padX int) string {
	return lipgloss.NewStyle().
		Width(w).Height(h).
		Background(bg).
		Padding(padY, padX).
		Render(strings.Join(lines, "\n"))
}

// --- text helpers ---------------------------------------------------

// padLineBg renders a single line's worth of content into exactly w
// columns with a specific bg. Any trailing space is filled with the
// same bg — the explicit fix for the "wallpaper shows through the
// right-hand gap" problem.
func padLineBg(line string, w int, bg, fg color.Color) string {
	if w <= 0 {
		return ""
	}
	trimmed := truncateToWidth(line, w)
	contentW := lipgloss.Width(trimmed)
	style := lipgloss.NewStyle().Background(bg).Foreground(fg)
	if contentW >= w {
		return style.Render(trimmed)
	}
	pad := style.Width(w - contentW).Render("")
	return lipgloss.JoinHorizontal(lipgloss.Top, style.Render(trimmed), pad)
}

func truncateToWidth(s string, w int) string {
	if lipgloss.Width(s) <= w {
		return s
	}
	// Crude: chop runes until it fits. For ANSI-bearing strings this
	// is lossy, but the content we feed it is plain text lines.
	r := []rune(s)
	for len(r) > 0 && lipgloss.Width(string(r)) > w {
		r = r[:len(r)-1]
	}
	return string(r)
}

// --- findings / progress formatters --------------------------------

func renderProgress(p Progress) []string {
	pct := 0
	if p.FilesTotal > 0 {
		pct = (p.FilesScanned * 100) / p.FilesTotal
	}
	bar := progressBar(pct, 40)
	elapsed := "—"
	if !p.Started.IsZero() {
		elapsed = "47s"
	}
	return []string{
		"status:    " + p.Status.String(),
		"elapsed:   " + elapsed,
		fmt.Sprintf("progress:  %d / %d  (%d%%)", p.FilesScanned, p.FilesTotal, pct),
		bar,
		"",
		"current:   " + p.CurrentPath,
	}
}

func progressBar(pct, width int) string {
	if width < 4 {
		width = 4
	}
	filled := (width * pct) / 100
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return lipgloss.NewStyle().Foreground(colOk).Background(colBg).Render("[" + bar + "]")
}

func renderFindings(findings []Finding) []string {
	out := make([]string, 0, len(findings)+1)
	header := lipgloss.NewStyle().Bold(true).Background(colBg).Foreground(colSubtle).
		Render("severity  location                               message")
	out = append(out, header)
	for _, f := range findings {
		sev := severityTag(f.Severity)
		loc := fmt.Sprintf("%s:%d", f.Path, f.Line)
		loc = padRight(loc, 38)
		line := lipgloss.NewStyle().Background(colBg).Foreground(colText).Render("  " + loc + "  " + f.Message)
		out = append(out, lipgloss.JoinHorizontal(lipgloss.Top, sev, line))
	}
	return out
}

func severityTag(s Severity) string {
	var bg color.Color
	switch s {
	case SevHigh:
		bg = colDanger
	case SevMed:
		bg = colWarn
	default:
		bg = colMuted
	}
	return lipgloss.NewStyle().
		Background(bg).
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Padding(0, 1).
		Render(s.Label())
}

func padRight(s string, w int) string {
	if lipgloss.Width(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-lipgloss.Width(s))
}
