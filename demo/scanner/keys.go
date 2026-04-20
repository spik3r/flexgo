package main

import "slices"

// KeyMap is the single source of truth for every binding in the app.
// Screens and the root model consult it rather than comparing strings
// directly, so a rebind (or a future user-config layer) changes one
// file.
//
// The trade-off with centralising this: keys are decoupled from the
// feature code, so tests can drive any binding without knowing the
// string. The cost is one extra lookup per key press — negligible.
type KeyMap struct {
	// Global — active on every screen, outside modals.
	Quit    []string
	Help    []string
	Launch  []string // go to launcher
	Scan    []string // go to scan screen
	History []string // go to history screen

	// Main-axis navigation (used by tabs and the history list).
	Next []string
	Prev []string

	// Scrolling — used by every scrollable panel.
	Up       []string
	Down     []string
	PageUp   []string
	PageDown []string
	Top      []string
	Bottom   []string

	// File viewer (Files tab).
	Open []string

	// Modal / sub-screen exit.
	Close []string
}

func DefaultKeys() KeyMap {
	return KeyMap{
		Quit:    []string{"q", "ctrl+c"},
		Help:    []string{"?"},
		Launch:  []string{"1"},
		Scan:    []string{"2"},
		History: []string{"3"},

		Next: []string{"tab", "l", "right"},
		Prev: []string{"shift+tab", "h", "left"},

		Up:       []string{"k", "up"},
		Down:     []string{"j", "down"},
		PageUp:   []string{"pgup", "ctrl+b"},
		PageDown: []string{"pgdown", "ctrl+f"},
		Top:      []string{"g", "home"},
		Bottom:   []string{"G", "end"},

		Open: []string{"enter"},

		Close: []string{"esc"},
	}
}

// Matches reports whether the key pressed hits any of the bindings.
// Returns false for empty binding lists so callers can use zero-valued
// KeyMap fields safely.
func Matches(key string, bindings []string) bool {
	return slices.Contains(bindings, key)
}

// HelpEntry is one row in the keymap modal: (keys, description).
type HelpEntry struct {
	Keys string
	Desc string
}

// HelpEntries returns the full keymap rendered as (keys, description)
// pairs, organised by section. This is what the help modal displays.
func (k KeyMap) HelpEntries() []HelpEntry {
	return []HelpEntry{
		{format(k.Help), "toggle this help"},
		{format(k.Quit), "quit"},
		{"", ""},
		{format(k.Launch), "launcher / profile"},
		{format(k.Scan), "scan screen"},
		{format(k.History), "scan history"},
		{"", ""},
		{format(k.Next), "next tab / item"},
		{format(k.Prev), "prev tab / item"},
		{"", ""},
		{format(k.Up) + " / " + format(k.Down), "scroll line"},
		{format(k.PageUp) + " / " + format(k.PageDown), "scroll page"},
		{format(k.Top) + " / " + format(k.Bottom), "scroll to top / bottom"},
		{"", ""},
		{format(k.Open), "open (on a file)"},
		{format(k.Close), "close modal / back"},
	}
}

func format(bindings []string) string {
	if len(bindings) == 0 {
		return ""
	}
	return bindings[0] // first binding is the canonical label
}
