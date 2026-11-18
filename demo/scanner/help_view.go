package main

import (
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/spik3r/flexgo"
)

const (
	helpCardWidth  = 56
	helpCardHeight = 24
	// The centred frame needs one painted spacer row above and below the card.
	helpCardBodyHeight = helpCardHeight + 2
)

type helpLayout struct {
	compact bool
	size    viewportSize
	rows    []string
}

func (a App) helpLayout() helpLayout {
	body := pageBodySize(a.width, a.height)
	if body.width >= helpCardWidth && body.height >= helpCardBodyHeight {
		// Remapped labels can make the fixed card taller than its declared height.
		rendered := buildKeymapCard(a.keys).Render(body.width, body.height)
		if lipgloss.Width(rendered) <= body.width && lipgloss.Height(rendered) <= body.height {
			return helpLayout{}
		}
	}
	size := body.inset(0, 1)
	return helpLayout{compact: true, size: size, rows: helpWrappedRows(a.keys, max(0, size.width-1))}
}

func helpWrappedRows(keys KeyMap, width int) []string {
	if width <= 0 {
		return nil
	}
	rows := []string{}
	for _, entry := range keys.HelpEntries() {
		if entry.Keys == "" && entry.Desc == "" {
			continue
		}
		line := entry.Desc
		if entry.Keys != "" {
			line = entry.Keys + " · " + line
		}
		line = strings.ReplaceAll(line, "\t", "    ")
		line = strings.ReplaceAll(strings.ReplaceAll(line, "\r\n", "\n"), "\r", "\n")
		rows = append(rows, strings.Split(ansi.Wrap(line, width, ""), "\n")...)
	}
	return rows
}

func (a App) reconcileHelp() App {
	layout := a.helpLayout()
	if layout.compact {
		layout.size.scroll(&a.help, 0, len(layout.rows))
	} else {
		a.help.Offset = 0
	}
	return a
}

func (a App) updateHelp(key string) App {
	a = a.reconcileHelp()
	layout := a.helpLayout()
	if !layout.compact || !layout.size.positive() {
		return a
	}
	switch {
	case Matches(key, a.keys.Up):
		layout.size.scroll(&a.help, -1, len(layout.rows))
	case Matches(key, a.keys.Down):
		layout.size.scroll(&a.help, 1, len(layout.rows))
	case Matches(key, a.keys.PageUp):
		layout.size.scroll(&a.help, -10, len(layout.rows))
	case Matches(key, a.keys.PageDown):
		layout.size.scroll(&a.help, 10, len(layout.rows))
	case Matches(key, a.keys.Top):
		a.help.Offset = 0
	case Matches(key, a.keys.Bottom):
		a.help.Offset = max(0, len(layout.rows)-layout.size.height)
	}
	return a
}

func (a App) helpBody() *flexgo.Node {
	layout := a.helpLayout()
	if !layout.compact {
		return buildKeymapCard(a.keys)
	}
	state := a.help
	title := flexgo.ClipView(func(w, h int) string {
		return padLineBg("Keybindings", w, colPanel, colHighlight)
	})
	return &flexgo.Node{
		Flex:       1,
		Dir:        flexgo.Col,
		Background: colBg,
		Children: []*flexgo.Node{
			{Height: 1, Background: colPanel, View: title},
			{Flex: 1, Background: colBg, View: viewportView(&state, layout.rows)},
		},
	}
}

func usableHelpKey(bindings []string, blocked [][]string, limit int) string {
	for _, key := range bindings {
		if key == "" || strings.ContainsAny(key, "\r\n\t") || ansi.StringWidth(key) > limit {
			continue
		}
		usable := true
		for _, group := range blocked {
			if Matches(key, group) {
				usable = false
				break
			}
		}
		if usable {
			return key
		}
	}
	return ""
}

func (a App) helpHint() string {
	layout := a.helpLayout()
	defaults := DefaultKeys()
	if !layout.compact && slices.Equal(a.keys.Close, defaults.Close) && slices.Equal(a.keys.Help, defaults.Help) {
		return "esc / ? close help"
	}
	width := pageBodySize(a.width, a.height).width
	groups := []string{}
	closeKey := usableHelpKey(a.keys.Close, nil, width-len(" close"))
	if closeKey == "" {
		closeKey = usableHelpKey(a.keys.Help, nil, width-len(" close"))
	}
	if closeKey != "" {
		groups = append(groups, closeKey+" close")
	}
	blocked := [][]string{a.keys.Close, a.keys.Help}
	pairs := []struct {
		first, second []string
		label         string
	}{
		{a.keys.Up, a.keys.Down, "scroll"},
		{a.keys.PageUp, a.keys.PageDown, "page"},
		{a.keys.Top, a.keys.Bottom, "ends"},
	}
	for _, pair := range pairs {
		aliases := []string{}
		for _, bindings := range [][]string{pair.first, pair.second} {
			if key := usableHelpKey(bindings, blocked, width-len(pair.label)-1); key != "" {
				aliases = append(aliases, key)
			}
			blocked = append(blocked, bindings)
		}
		if len(aliases) > 0 {
			groups = append(groups, strings.Join(aliases, "/")+" "+pair.label)
		}
	}
	hint := ""
	for _, group := range groups {
		next := group
		if hint != "" {
			next = hint + " · " + group
		}
		if ansi.StringWidth(next) <= width {
			hint = next
		}
	}
	return hint
}
