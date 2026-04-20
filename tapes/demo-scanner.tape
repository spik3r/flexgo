# Marquee GIF for the README — walks through the scanner demo:
# launcher → scan screen with each tab → keymap modal → history.
#
# Regenerate:  vhs tapes/demo-scanner.tape

Output docs/demo-scanner.gif

Set FontSize 14
Set Width 1200
Set Height 760
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms

# Clear the prompt before recording so it doesn't flicker.
Hide
Type "clear"
Enter
Show

Type "go run ./demo/scanner"
Enter
Sleep 2500ms

# Launcher: cursor through a couple of fields.
Down
Sleep 600ms
Down
Sleep 600ms
Up
Sleep 600ms

# Switch to scan screen.
Type "2"
Sleep 1500ms

# Tab through Progress → Findings → Logs → Files.
Tab
Sleep 1200ms
Tab
Sleep 1200ms
Tab
Sleep 1200ms

# Open the keymap modal.
Type "?"
Sleep 2500ms
Escape
Sleep 600ms

# Jump to history and browse.
Type "3"
Sleep 1000ms
Down
Sleep 600ms
Down
Sleep 600ms
Down
Sleep 800ms

# Quit cleanly.
Type "q"
Sleep 500ms
