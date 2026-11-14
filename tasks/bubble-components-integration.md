# Qualify the planned Bubbles component example

Status: integrated locally for independent exact-head review.
Owner: delegated Flexgo worker.
Branch/worktree: `codex/bubble-components-integration`,
`/tmp/flexgo-bubble-components-integration`.
Base: `80904e5e5c32dcfdfadb9ccf73537359556f0348` (fresh remote main).
Source example: existing local `490f9079465c55ce0340f65422c06780a2252dfc`.

## Ownership and scope

The repository plan lists text input, viewport and spinner integration as an open
example item. GitHub has no open PRs or issues. Current task/thread listings,
worktrees and task records show no active competing Bubbles owner. The original
example worktree is clean and remains unchanged. Its prior task record reports
local independent review, but no publication. This task carries that work forward
in a new branch rather than creating a second example. Primary changes and all
other worktrees are preserved. Disk has 35 GiB free.

The earlier dependency tidy attempt hit an uncached protected module cache path;
that operation is not retried. The previous task also traced lost combining marks
in terminal output to pinned Ultraviolet after the component View string. Its
experimental dependency patch remains unadopted and outside this task. No live
terminal claim is inferred from logical input tests or snapshots.

## Acceptance before implementation

- Integrate the existing text input, viewport and spinner example with current
  main and use `ClipView` for bounded leaf output.
- Keep component models in application state; sizing/rendering local copies must
  not change input value, cursor position, focus or viewport offset.
- Preserve state through repeated normal/tiny/normal resizes and clamp scroll
  when the viewport grows. Input remains editable after resize.
- Tab/Shift+Tab switch input/viewport focus; keys and synthetic paste reach only
  the focused component; q stays text and Esc/Ctrl+C quit from either focus.
- Spinner ticks continue in either focus. No tests access the system clipboard.
- Test terminal-cell bounds, Unicode logical values, narrow cursor windows,
  footer retention, embedded content and nonpositive allocations.
- Register the deterministic example golden, document copyable integration and
  known terminal limitations, and mark only this planned example item complete.
- Obtain independent exact-head review and full local race/vet checks; require
  hosted tidy/build/race CI and an exact unchanged base before merge.

Barcode #7, Foundry, deployment, upstream patches, blocked cache operations,
listeners and the denied ps command are outside this task.

## Implementation and validation

Carried the source example, prior regression tests and pinned Bubbles v2.1.0
checksums forward from `490f907`. Each leaf now uses current `ClipView`; the
application still owns component state, resize allocation and message routing.
A new minimum-layout test failed against the original example: at 20x5, its
footer ended with `Esc` and did not show the complete quit/focus controls. The
footer now puts `Esc:quit Tab:focus` before optional focus detail; the test passes
for input and viewport focus.

Five additional test groups cover the minimum footer, initial component commands,
repeated normal/tiny/zero resizes followed by cursor-positioned editing, long
styled/tabbed viewport content, and quitting from viewport focus. The existing
nine test groups cover input/paste focus, viewport scrolling/clamping, spinner
updates, logical Unicode input, interior cursor windows and stable rendering.
No test reads the system clipboard or starts a terminal.

Using `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`:

- `go test -mod=readonly ./example/bubble_components -count=1`: passed before
  integration and after the callback/footer changes.
- `go test -mod=readonly ./example/bubble_components
  -run TestMinimumLayoutKeepsControlsVisible -count=1`: failed on the prior footer
  as expected, then passed with the fix in the package/full suites.
- `go test -mod=readonly . -run '^TestExampleGolden/bubble_components$' -update
  -count=1`: passed; registered a deterministic 24-row snapshot.
- `go test -mod=readonly -race -count=1 ./...`: passed all packages, including the
  new example and all existing snapshots.
- `go vet -mod=readonly ./...` and gofmt passed.
- Source/documentation whitespace checks pass. The new golden intentionally
  contains trailing cell padding; its exact bytes are checked by the golden test.

Local build/tidy and interactive terminal acceptance remain unrun. No denied
operation was retried. Hosted dependency tidy/build/race and independent review
remain required gates before merge. The example guide records the existing
upstream combining-mark limitation instead of claiming terminal fidelity.

## Independent review correction

The first exact-head review of `24db9be` reproduced an omitted cursor-cell
reservation. At width 20 with 18 characters and the cursor at the end, the raw
row retained cursor style escapes but no painted cursor space. The new
`TestInputCursorSurvivesBoundsAndResize` failed on this case before correction.
Both stored resize sizing and the render-copy sizing now reserve two prompt
cells plus one virtual-cursor cell. The guide and width assertion agree.

The new test checks one painted cursor cell, logical value and cursor position
across 144 combinations of width, text length, start/interior/end position and
shrink/grow transitions. There are now 15 example test groups. After the fix,
`go test -mod=readonly -race ./example/bubble_components -count=1`, the example
golden check/regeneration, full `go test -mod=readonly -race -count=1 ./...`,
vet, gofmt and source/documentation whitespace checks passed. Regeneration left
the new example golden byte-for-byte unchanged. Revised exact-head review is
required before publication; terminal-rendering acceptance remains unverified.
