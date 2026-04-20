package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// LauncherScreen is the wizard-style profile editor. Its body is a
// centred card (auto-margins all four sides) on a colBg backdrop, so
// the page shell's header/footer sit flush and the card floats in
// the middle without touching them.
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

// Body returns the centred profile card. Centring is delegated to the
// `centered` helper in widgets.go, which builds an explicit Flex-
// spacer frame around the card — the only approach that guarantees
// every cell around the card is painted with colBg. See that helper
// for the reasoning.
func (l LauncherScreen) Body(w, h int) *flexgo.Node {
	card := &flexgo.Node{
		Width:            60,
		Height:           20,
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
	return centered(card, 20)
}

func (l LauncherScreen) cardBody() func(int, int) string {
	values := []string{
		l.profile.Name,
		l.profile.Root,
		strings.Join(l.profile.Include, ", "),
		strings.Join(l.profile.Exclude, ", "),
		boolLabel(l.profile.Deep),
	}
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
