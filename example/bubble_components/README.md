# Bubbles components in Flexgo

Run `go run ./example/bubble_components`.

The example combines a Bubbles text input, a scrollable viewport and an animated
spinner. All content is synthetic and local. It starts with the input focused:

- Type or paste into the input; `q` is text.
- Press Tab or Shift+Tab to switch between input and viewport.
- With the viewport focused, use Up/Down or Page Up/Page Down to scroll.
- Press Esc or Ctrl+C to quit from either focus.

At 20 columns and 5 rows, the footer still shows `Esc:quit Tab:focus`. Smaller
windows show the required size and, where space permits, the quit key. Editing,
paste, focus changes and scrolling pause while their controls are hidden; Esc
and Ctrl+C still quit. Resize back to continue at the same input cursor, focus
and scroll position (subject to viewport clamping). Spinner updates continue.
Clipboard or input-validation errors appear below the text input, within its
existing third row. Correcting the input or completing a successful paste clears
the feedback when Bubbles clears its error state. Long messages are clipped to
one row; they do not displace the viewport or quit controls.

The input cursor's command chain restarts when its focused control becomes
visible again. Spinner and input cursor updates use `tea.Cmd`.

## Copy the integration pattern

Keep component models on the application model. Route keys and paste by focus in
`Update`; pass spinner ticks and cursor messages to their components. Start both
the spinner and focused input cursor from `Init`. Return the command from `Focus` to
Bubble Tea when switching back from the viewport. Use the same visibility guard
for key routing and other component messages: paste and asynchronous input can
arrive outside the key handler. If cursor messages are paused while hidden,
restart the focused cursor command when the layout becomes usable again.

Window-size messages update the stored viewport dimensions so scrolling clamps to
the current visible height. The input width reserves two cells for the prompt
and one for the virtual cursor, so clipping cannot hide the cursor at the end
of a full input. The pinned Bubbles input does not recompute an interior
cursor's visible window on
`SetWidth`, so `setInputWidth` rebuilds it using public cursor methods and restores
the original cursor position.

Each leaf callback receives its content dimensions from Flexgo. Size a local copy
of the component for rendering, then use `flexgo.ClipView` to bound its output:

```go
{Flex: 1, View: flexgo.ClipView(func(w, h int) string {
    view := m.view
    view.SetWidth(w)
    view.SetHeight(h)
    return view.View()
})},
```

Rendering must not update the stored input value, cursor, focus or scroll offset.
`ClipView` bounds text; the application still chooses a layout whose fixed rows
fit. This example reserves three input rows and one footer row, then gives the
remaining height to the viewport. Render `input.Err` alongside `input.View()`;
the view string alone does not expose clipboard or validation failures. Keep
Bubbles as the owner of that error state rather than caching a second copy.

## Verification and limits

Synthetic tests cover focus/key/paste routing, cursor and scroll retention,
resizing, component startup, spinner updates, quit commands and terminal-cell
bounds. `FLEXGO_GOLDEN=1 go run ./example/bubble_components` renders a deterministic
snapshot without starting an interactive terminal or reading the clipboard.

Interactive terminal acceptance is not complete. Prior investigation of the
pinned Ultraviolet renderer found that some decomposed combining marks survive
Bubbles' logical value and View string but disappear from emitted terminal output.
The component tests verify the logical value; they do not establish terminal
rendering fidelity. No experimental renderer patch or dependency replacement is
included. Test your target terminal before relying on combining-mark display.
