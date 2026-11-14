# Opt-in bounded leaf callbacks

Status: implemented locally for independent exact-head review.
Owner: delegated Flexgo worker.
Branch/worktree: `codex/clip-view`, `/tmp/flexgo-clip-view`.
Base: `248c5529b4605146ff482e3d24d34b2d4dfe0b9e` (fresh remote main).
No open GitHub PRs/issues or competing callback-bounding claim were found.
The Bubbles example worktree and all primary changes are preserved. Disk has
36 GiB free. Barcode #7 and Foundry are outside this task.

## Design and acceptance

The public Render contract allows overflowing callbacks. Three existing examples
independently clip their own leaf output. Add an opt-in callback adapter rather
than an engine sizing or overflow policy; keep existing Render behavior intact.
Independent read-only design review accepted the adapter before implementation.
The review required per-line ANSI closure, retained cross-line styling, unchanged
plain text, and explicit limits for frames and cursor-control programs.

- `ClipView(view func(w, h int) string)` returns a callback of the same signature.
- Nil callbacks and nonpositive allocations produce empty output without calling
  the original callback; otherwise call it once with the original allocation.
- Crop text to the first h explicit lines and first w terminal cells per line,
  without wrapping, ellipsis or padding; preserve wide and combining characters.
- Normalize tabs to four spaces, CRLF and CR to newlines before measuring.
- Preserve ANSI styles and OSC-8 links and close discarded trailing state so a
  clipped leaf cannot style its siblings. This is for text callbacks, not arbitrary
  terminal cursor-control programs or a security filter.
- Test direct callbacks and full Row/Col compositions with borders, padding,
  sibling markers and repeated resizes. Unwrapped callbacks still overflow.
- Replace duplicate clipping in both AutoGrid examples and the responsive example;
  keep their goldens stable and document opt-in use and allocation limits.

Local permitted checks: focused tests, full uncached race suite, vet, gofmt and
whitespace checks. Local build/tidy are not rerun because of the earlier denied
Go metadata-cache write. Independent exact-head review and qualified hosted CI
must pass before merge. Interactive terminal acceptance remains unverified.

## Reproduction and implementation

Before the change, the AutoGrid helper returned an open red SGR/OSC-8 hyperlink
when truncating `first\nsecond` to one row: the discarded row held both closing
sequences. The synthetic `TestBoundedClosesDiscardedStyle` failed on the old helper.
A read-only Go overlay of the original responsive source also reproduced
`fit([]string{"first\nsecond", "third"}, 5, 1)` returning two rows; the new
`TestLeafClippingCountsEmbeddedLines` failed against that original source.

`ClipView` now bounds normalized text with the pinned ANSI grapheme-aware truncator.
It closes each emitted row and replays SGR changes and the current OSC-8 link on
following retained rows. Standalone SGR resets discard accumulated style history. Each
invocation owns its state; callback state remains the caller's responsibility.
Three examples share the adapter. Existing renderer code and golden files are
unchanged. Documentation distinguishes leaf bounds from node sizing policy.

## Validation

Environment: `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`.

- Both baseline regressions failed as expected, then passed with the adapter.
- Focused callback/example tests passed: `go test -mod=readonly .
  ./example/layouts/autogrid ./example/builder/layouts/autogrid ./example/responsive
  -run 'TestClipView|ExampleClipView|TestBoundedCloses|TestLeafClipping|TestResizeBounds'
  -count=1`.
- `go test -mod=readonly -race -count=1 ./...` passed, including all unchanged
  example goldens, scanner and layout tests.
- `go test -mod=readonly -race . -run 'TestClipView|ExampleClipView' -count=1`
  passed after adding explicit closed-style/closed-link and ST-terminated link
  cases to the ANSI table.
- `go vet -mod=readonly ./...`, gofmt and whitespace checks passed.
- Initial development compilation errors (decoder API name and an unused import)
  were corrected before passing checks.

Local build/tidy and interactive terminal acceptance remain unrun. No browser,
listener, deployment or Barcode CI operation was attempted. A read-only `ps`
process-status command was denied by the sandbox (`operation not permitted`);
it was not retried. Test completion was obtained through its existing session.
Independent review and hosted CI remain the publication/merge gates.
