# Bubbles components in Flexgo

Run `go run ./example/bubble_components`.

The example combines a Bubbles text input, a scrollable viewport and an animated
spinner. All content is synthetic and local. It starts with the input focused:

- Type or paste into the input; `q` is text.
- Press Tab or Shift+Tab to switch between input and viewport.
- With the viewport focused, use Up/Down or Page Up/Page Down to scroll.
- Press Esc or Ctrl+C to quit from either focus.

At 20 columns and 5 rows, the footer still shows `Esc:quit Tab:focus`. Smaller
windows show a resize hint while retaining the component state. Resize back to
continue editing or scrolling. Spinner and input cursor updates use `tea.Cmd`.

## Copy the integration pattern

Keep component models on the application model. Route keys and paste by focus in
`Update`; pass spinner ticks and cursor messages to their components. Start both
the spinner and focused input cursor from `Init`. Return the command from `Focus` to
Bubble Tea when switching back from the viewport.

Window-size messages update the stored viewport dimensions so scrolling clamps to
the current visible height. The input width excludes its two-cell prompt. The
pinned Bubbles input does not recompute an interior cursor's visible window on
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
remaining height to the viewport.

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
