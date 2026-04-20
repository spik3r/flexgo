# layouts.Tabs interaction — h/l cycles through the tabs and the
# active-tab underline tracks the title (AlignSelf at work).
#
# Regenerate:  vhs tapes/layouts-tabs.tape

Output docs/layouts-tabs.gif

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

Type "go run ./example/layouts/tabs"
Enter
Sleep 2500ms

Type "l"
Sleep 1000ms
Type "l"
Sleep 1000ms
Type "h"
Sleep 1000ms
Type "h"
Sleep 1200ms

Type "q"
Sleep 400ms
