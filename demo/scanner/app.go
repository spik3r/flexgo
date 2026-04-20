package main

import (
	tea "charm.land/bubbletea/v2"

	"github.com/spik3r/flexgo"
)

// Screen is the top-level view discriminator.
type Screen int

const (
	ScreenLauncher Screen = iota
	ScreenScan
	ScreenHistory
)

// App is the root model. It owns the three screen models, the keymap,
// and the modal flag — nothing else. Screens store their own internal
// state; App routes messages and picks which screen to paint.
type App struct {
	width, height int
	ready         bool

	keys      KeyMap
	screen    Screen
	modalOpen bool

	launcher LauncherScreen
	scan     ScanScreen
	history  HistoryScreen
}

func NewApp() App {
	profile := DefaultProfile()
	return App{
		keys:     DefaultKeys(),
		screen:   ScreenLauncher,
		launcher: NewLauncherScreen(profile),
		scan:     NewScanScreen(profile),
		history:  NewHistoryScreen(),
	}
}

func (a App) Init() tea.Cmd { return nil }

// Update is the one and only dispatch point. Order: modal → global →
// screen. Modals are captive; global bindings can't be overridden by
// a screen; anything the screen doesn't handle is ignored.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height
		a.ready = true
		return a, nil

	case tea.KeyPressMsg:
		key := msg.String()

		// 1. Modal is captive. Only its own close/help bindings escape.
		if a.modalOpen {
			if Matches(key, a.keys.Close) || Matches(key, a.keys.Help) {
				a.modalOpen = false
			}
			return a, nil
		}

		// 2. Global bindings — available on every screen.
		switch {
		case Matches(key, a.keys.Quit):
			return a, tea.Quit
		case Matches(key, a.keys.Help):
			a.modalOpen = true
			return a, nil
		case Matches(key, a.keys.Launch):
			a.screen = ScreenLauncher
			return a, nil
		case Matches(key, a.keys.Scan):
			// Transition consumes the launcher's current profile.
			// One-way data flow — the scan screen doesn't reach back.
			a.scan = NewScanScreen(a.launcher.Profile())
			a.screen = ScreenScan
			return a, nil
		case Matches(key, a.keys.History):
			a.screen = ScreenHistory
			return a, nil
		}

		// 3. Screen-local dispatch.
		switch a.screen {
		case ScreenLauncher:
			a.launcher = a.launcher.Update(msg, a.keys)
		case ScreenScan:
			a.scan = a.scan.Update(msg, a.keys)
		case ScreenHistory:
			a.history = a.history.Update(msg, a.keys)
		}
	}
	return a, nil
}

// View always renders through pageShell, so the title bar, subheader
// and footer are identical across screens. Only the body changes —
// and when the modal is open, only the body is replaced. The chrome
// stays put, giving the user continuity.
func (a App) View() tea.View {
	if !a.ready {
		return tea.NewView("Loading...")
	}

	subtitle, screenHint, body := a.currentScreen()
	if a.modalOpen {
		screenHint = "esc / ? close help"
		body = buildKeymapCard(a.keys)
	}

	tree := pageShell(a.composedSubtitle(subtitle), a.screen, body, screenHint)
	v := tea.NewView(tree.Render(a.width, a.height))
	v.AltScreen = true
	return v
}

// composedSubtitle prefixes the screen's subtitle with its title,
// producing "Scanner · profile: default" for the subheader's left
// slot. Keeping the composition here (not on each screen) means the
// format change in one place lands everywhere.
func (a App) composedSubtitle(subtitle string) string {
	var title string
	switch a.screen {
	case ScreenLauncher:
		title = a.launcher.Title()
	case ScreenScan:
		title = a.scan.Title()
	case ScreenHistory:
		title = a.history.Title()
	}
	if subtitle == "" {
		return title
	}
	return title + "  ·  " + subtitle
}

// currentScreen extracts the three values the page shell needs from
// whichever screen is active. Adding a new screen is one new arm here
// + the Screen constant + field + KeyMap binding + screenTabs entry.
func (a App) currentScreen() (subtitle, screenHint string, body *flexgo.Node) {
	switch a.screen {
	case ScreenLauncher:
		return a.launcher.Subtitle(), a.launcher.Footer(), a.launcher.Body(a.width, a.height)
	case ScreenScan:
		return a.scan.Subtitle(), a.scan.Footer(), a.scan.Body(a.width, a.height)
	case ScreenHistory:
		return a.history.Subtitle(), a.history.Footer(), a.history.Body(a.width, a.height)
	}
	return "", "", &flexgo.Node{Flex: 1, Background: colBg}
}
