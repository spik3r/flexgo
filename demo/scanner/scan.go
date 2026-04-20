package main

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/spik3r/flexgo"
)

// ScanScreen is the main working view. It exposes Title/Subtitle/
// Footer/Body so the root App can plug it into the shared pageShell
// — the screen has no header/footer of its own, which is what keeps
// the chrome consistent across screens.
type ScanScreen struct {
	profile Profile

	activeTab int
	tabNames  []string

	progress Progress
	findings []Finding
	logs     []string
	files    []string

	// One viewport per tab, indexed by tab number.
	viewports [4]ViewportState

	// Files tab sub-state.
	filesCursor int
	openedPath  string
	openedBody  []string
}

func NewScanScreen(profile Profile) ScanScreen {
	return ScanScreen{
		profile:  profile,
		tabNames: []string{"Progress", "Findings", "Logs", "Files"},
		progress: SampleProgress(),
		findings: SampleFindings(),
		logs:     SampleLogs(),
		files:    SampleFiles(),
	}
}

func (s ScanScreen) Title() string { return "Scanner" }

func (s ScanScreen) Subtitle() string {
	if s.openedPath != "" {
		return "file: " + s.openedPath
	}
	return "profile: " + s.profile.Name
}

func (s ScanScreen) Footer() string {
	if s.openedPath != "" {
		return "viewing file  ·  j/k scroll  ·  esc close"
	}
	hints := []string{
		"progress  ·  j/k scroll  ·  tab next",
		"findings  ·  j/k scroll  ·  tab next",
		"logs      ·  j/k scroll  ·  pgup/pgdn  ·  tab next",
		"files     ·  j/k move  ·  enter open  ·  tab next",
	}
	return hints[s.activeTab]
}

// Update is called by App when the scan screen is active and no modal
// is open. It owns tab switching, scroll, and file-open transitions.
func (s ScanScreen) Update(msg tea.KeyPressMsg, keys KeyMap) ScanScreen {
	key := msg.String()

	// Opened-file viewer keys rebind Close to "back to Files tab".
	if s.openedPath != "" {
		switch {
		case Matches(key, keys.Close):
			s.openedPath = ""
			s.openedBody = nil
		case Matches(key, keys.Up):
			s.viewports[3].ScrollBy(-1, len(s.openedBody), 1)
		case Matches(key, keys.Down):
			s.viewports[3].ScrollBy(1, len(s.openedBody), 1)
		case Matches(key, keys.PageUp):
			s.viewports[3].ScrollBy(-10, len(s.openedBody), 1)
		case Matches(key, keys.PageDown):
			s.viewports[3].ScrollBy(10, len(s.openedBody), 1)
		}
		return s
	}

	switch {
	case Matches(key, keys.Next):
		s.activeTab = (s.activeTab + 1) % len(s.tabNames)
	case Matches(key, keys.Prev):
		s.activeTab = (s.activeTab - 1 + len(s.tabNames)) % len(s.tabNames)
	case Matches(key, keys.Up):
		if s.activeTab == 3 {
			if s.filesCursor > 0 {
				s.filesCursor--
			}
		} else {
			s.viewports[s.activeTab].ScrollBy(-1, s.contentLen(), 1)
		}
	case Matches(key, keys.Down):
		if s.activeTab == 3 {
			if s.filesCursor < len(s.files)-1 {
				s.filesCursor++
			}
		} else {
			s.viewports[s.activeTab].ScrollBy(1, s.contentLen(), 1)
		}
	case Matches(key, keys.PageUp):
		s.viewports[s.activeTab].ScrollBy(-10, s.contentLen(), 1)
	case Matches(key, keys.PageDown):
		s.viewports[s.activeTab].ScrollBy(10, s.contentLen(), 1)
	case Matches(key, keys.Top):
		s.viewports[s.activeTab].Offset = 0
		if s.activeTab == 3 {
			s.filesCursor = 0
		}
	case Matches(key, keys.Bottom):
		s.viewports[s.activeTab].Offset = s.contentLen()
		if s.activeTab == 3 {
			s.filesCursor = len(s.files) - 1
		}
	case Matches(key, keys.Open):
		if s.activeTab == 3 && s.filesCursor < len(s.files) {
			s.openedPath = strings.TrimSpace(s.files[s.filesCursor])
			s.openedBody = fakeFileBody(s.openedPath)
			s.viewports[3] = ViewportState{}
		}
	}
	return s
}

func (s ScanScreen) contentLen() int {
	switch s.activeTab {
	case 0:
		return len(renderProgress(s.progress))
	case 1:
		return len(renderFindings(s.findings))
	case 2:
		return len(s.logs)
	case 3:
		return len(s.files)
	}
	return 0
}

// Body builds the screen's content — everything between pageShell's
// chrome and footer. Split view normally; single-pane when a file is
// open. The chrome (title, subheader, hints) stays put either way.
func (s ScanScreen) Body(w, h int) *flexgo.Node {
	if s.openedPath != "" {
		viewerState := s.viewports[3]
		return &flexgo.Node{
			Flex:             1,
			Background:       colPanel,
			ShowBorder:       true,
			Border:           lipgloss.RoundedBorder(),
			BorderForeground: colAccent,
			View:             viewportView(&viewerState, s.openedBody),
		}
	}

	sidebar := &flexgo.Node{
		Width:      24,
		Background: colPanel,
		View:       s.sidebarView(),
	}

	tabs := tabStripNode(s.tabNames, s.activeTab)

	panel := &flexgo.Node{
		Flex:       1,
		Background: colBg,
		View:       s.activePanelView(),
	}

	main := &flexgo.Node{
		Flex:       1,
		Dir:        flexgo.Col,
		Background: colBg,
		Children:   []*flexgo.Node{tabs, panel},
	}

	return &flexgo.Node{
		Flex:       1,
		Dir:        flexgo.Row,
		Gap:        2,
		Background: colBg,
		Children:   []*flexgo.Node{sidebar, main},
	}
}

func (s ScanScreen) sidebarView() func(w, h int) string {
	return func(w, h int) string {
		rows := []string{
			headingStyle().Render("SCAN TARGET"),
			"",
			textStyle(colPanel).Render("root:    " + s.profile.Root),
			textStyle(colPanel).Render("profile: " + s.profile.Name),
			"",
			headingStyle().Render("INCLUDE"),
		}
		for _, inc := range s.profile.Include {
			rows = append(rows, textStyle(colPanel).Render("  "+inc))
		}
		rows = append(rows, "")
		rows = append(rows, headingStyle().Render("EXCLUDE"))
		for _, ex := range s.profile.Exclude {
			rows = append(rows, textStyle(colPanel).Render("  "+ex))
		}
		return panelPaint(rows, w, h, colPanel, 1, 2)
	}
}

func (s ScanScreen) activePanelView() func(w, h int) string {
	switch s.activeTab {
	case 0:
		return viewportView(&s.viewports[0], renderProgress(s.progress))
	case 1:
		return viewportView(&s.viewports[1], renderFindings(s.findings))
	case 2:
		return viewportView(&s.viewports[2], s.logs)
	case 3:
		return filetreeView(s.files, s.filesCursor)
	}
	return func(w, h int) string { return "" }
}

// tabStripNode builds the tab strip as a flexgo sub-tree rather than a
// single View callback, so every cell of the strip has a defined bg
// from the Node.Background fields — no transparent seams.
//
// Each tab is a flex-weighted leaf with its own Background. Inactive
// tabs use colBg, the active tab uses colAccent. A height-1 underline
// row below the titles uses colAccent as a visual anchor.
func tabStripNode(names []string, active int) *flexgo.Node {
	cells := make([]*flexgo.Node, 0, len(names))
	for i, name := range names {
		bg := colBg
		fg := colSubtle
		bold := false
		if i == active {
			bg = colAccent
			fg = colHighlight
			bold = true
		}
		label := name
		cells = append(cells, &flexgo.Node{
			Flex:       1,
			Background: bg,
			View:       tabTitleView(label, bg, fg, bold),
		})
	}
	return &flexgo.Node{
		Height:     2,
		Dir:        flexgo.Col,
		Background: colBg,
		Children: []*flexgo.Node{
			{Height: 1, Dir: flexgo.Row, Background: colBg, Children: cells},
			{Height: 1, Background: colBg, View: underlineView()},
		},
	}
}

func tabTitleView(text string, bg, fg color.Color, bold bool) func(w, h int) string {
	return func(w, h int) string {
		style := lipgloss.NewStyle().
			Width(w).Height(h).
			Background(bg).Foreground(fg).
			Align(lipgloss.Center, lipgloss.Center)
		if bold {
			style = style.Bold(true)
		}
		return style.Render(text)
	}
}

func underlineView() func(w, h int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Width(w).Height(h).
			Background(colBg).Foreground(colAccent).
			Render(strings.Repeat("─", w))
	}
}

func fakeFileBody(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return []string{"(no file)"}
	}
	return []string{
		"// " + path,
		"",
		"package example",
		"",
		"// Sample file content for the demo. Wire a real reader here",
		"// and dispatch via tea.Cmd — never block Update.",
		"",
		"func main() {",
		"    fmt.Println(\"hello\")",
		"}",
	}
}
