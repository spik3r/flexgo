package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

const (
	launcherCardWidth  = 60
	launcherCardHeight = 20
	// Zero-height spacer leaves still paint one row above and below the card.
	launcherCardBodyHeight = launcherCardHeight + 2
)

// LauncherScreen presents a scan profile and a cursor over its fields.
// It uses a centred card when space allows and compact rows otherwise.
type LauncherScreen struct {
	profile Profile
	cursor  int
	fields  []string
}

func NewLauncherScreen(profile Profile) LauncherScreen {
	return LauncherScreen{
		profile: profile,
		fields:  []string{"Name", "Root", "Include", "Exclude", "Deep"},
	}
}

func (l LauncherScreen) Title() string    { return "Launcher" }
func (l LauncherScreen) Subtitle() string { return "configure a scan" }
func (l LauncherScreen) Footer() string {
	return "profile editor  ·  j/k or tab to move cursor  ·  2 start scan"
}

func (l LauncherScreen) Update(msg tea.KeyPressMsg, keys KeyMap) LauncherScreen {
	key := msg.String()
	switch {
	case Matches(key, keys.Up), Matches(key, keys.Prev):
		if l.cursor > 0 {
			l.cursor--
		}
	case Matches(key, keys.Down), Matches(key, keys.Next):
		if l.cursor < len(l.fields)-1 {
			l.cursor++
		}
	}
	return l
}

func (l LauncherScreen) Profile() Profile { return l.profile }

// Body uses its allocated body dimensions to choose a centred card or
// compact rows. Both presentations leave profile and cursor state unchanged.
func (l LauncherScreen) Body(w, h int) *flexgo.Node {
	if compactLauncher(w, h) {
		return &flexgo.Node{Flex: 1, Background: colPanel, View: l.compactView()}
	}
	card := &flexgo.Node{
		Width:            launcherCardWidth,
		Height:           launcherCardHeight,
		ShowBorder:       true,
		Border:           lipgloss.RoundedBorder(),
		BorderForeground: colAccent,
		Background:       colPanel,
		Dir:              flexgo.Col,
		Children: []*flexgo.Node{
			{Height: 2, Background: colAccent, View: cardTitleBar("Scan profile")},
			{Flex: 1, Background: colPanel, View: l.cardBody()},
			{Height: 2, Background: colPanel, View: l.cardFooter()},
		},
	}
	return centered(card, launcherCardHeight)
}

func compactLauncher(w, h int) bool {
	return w < launcherCardWidth || h < launcherCardBodyHeight
}

func (l LauncherScreen) fieldValues() []string {
	return []string{
		l.profile.Name,
		l.profile.Root,
		strings.Join(l.profile.Include, ", "),
		strings.Join(l.profile.Exclude, ", "),
		boolLabel(l.profile.Deep),
	}
}

func (l LauncherScreen) compactView() func(int, int) string {
	values := l.fieldValues()
	return func(w, h int) string {
		if w <= 0 || h <= 0 {
			return ""
		}
		rows := make([]string, h)
		for i := range rows {
			rows[i] = padLineBg("", w, colPanel, colText)
		}
		first, visible := 0, h
		if h > len(l.fields) {
			title := flexgo.ClipView(func(int, int) string { return headingStyle().Render("Scan profile") })(w, 1)
			rows[0] = padLineBg(title, w, colPanel, colText)
			first, visible = 1, h-1
		}
		cursor := max(0, min(l.cursor, len(l.fields)-1))
		offset := max(0, cursor-visible+1)
		for i := offset; i < min(offset+visible, len(l.fields)); i++ {
			bg, fg := colPanel, colText
			if i == cursor {
				bg, fg = colSelected, colHighlight
			}
			line := l.fields[i] + ": " + values[i]
			clipped := flexgo.ClipView(func(int, int) string { return line })(w, 1)
			rows[first+i-offset] = padLineBg(clipped, w, bg, fg)
		}
		return strings.Join(rows, "\n")
	}
}

func (l LauncherScreen) cardBody() func(int, int) string {
	values := l.fieldValues()
	return func(w, h int) string {
		rows := make([]string, len(l.fields))
		for i, name := range l.fields {
			bg := colPanel
			fg := colText
			if i == l.cursor {
				bg = colSelected
				fg = colHighlight
			}
			rowStyle := lipgloss.NewStyle().Background(bg).Foreground(fg)
			label := rowStyle.Width(12).Render(name + ":")
			val := rowStyle.Render(values[i])
			line := lipgloss.JoinHorizontal(lipgloss.Top, label, val)
			rows[i] = rowStyle.Width(w).Padding(0, 2).Render(line)
		}
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(colPanel).
			Padding(1, 0).
			Render(strings.Join(rows, "\n"))
	}
}

func (l LauncherScreen) cardFooter() func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(colPanel).
			Foreground(colSubtle).
			Padding(0, 2).
			Align(lipgloss.Left, lipgloss.Center).
			Render("press 2 to start scanning with this profile")
	}
}

func boolLabel(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
