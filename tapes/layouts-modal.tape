# layouts.Modal — page content with a bordered modal opened on top.
# Press Space to open the modal, Escape to close.
#
# Regenerate:  vhs tapes/layouts-modal.tape

Output docs/layouts-modal.gif

Set FontSize 14
Set Width 1000
Set Height 540
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms

Hide
Type "clear"
Enter
Show

Type "go run ./example/layouts/modal"
Enter
Sleep 2500ms

# Open the modal.
Space
Sleep 2500ms

# Close it.
Escape
Sleep 1500ms

Type "q"
Sleep 400ms
