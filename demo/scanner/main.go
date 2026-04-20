// demo/scanner is a reference TUI app built on flexgo. It wires up
// three screens (Launcher, Scan, History), a keymap modal, four
// tabbed sub-panels on the Scan screen, and centralised key dispatch.
//
// It exists to show how the pieces fit together once an app grows
// beyond one screen: state split by screen, one root model that owns
// routing, a single KeyMap consulted by everyone. Read demo/scanner/
// README.md for the architectural tour.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		// Render the Scan screen via the page shell at a fixed size.
		app := NewApp()
		app.width, app.height = 100, 32
		app.ready = true
		app.screen = ScreenScan
		subtitle, footer, body := app.currentScreen()
		tree := pageShell(app.composedSubtitle(subtitle), app.screen, body, footer)
		fmt.Print(tree.Render(app.width, app.height))
		return
	}
	if _, err := tea.NewProgram(NewApp()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
