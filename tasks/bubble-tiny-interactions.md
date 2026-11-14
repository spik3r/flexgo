# Suspend interactions while Bubbles controls are hidden

Status: implemented locally for independent exact-head review.
Owner: delegated Flexgo worker.
Branch/worktree: `codex/bubble-tiny-interactions`,
`/tmp/flexgo-bubble-tiny-interactions`.
Base: `41a77f1b9cf8d7ce0611b5e0e5f1155c9f121792` (fresh remote main).

Current PR/issue queries were empty. Worktrees, task records and active task
listings showed no competing Flexgo owner. Primary changes and previous
worktrees are preserved. Disk has 33 GiB available. No relevant repository
AGENTS.md or .agents skill was found; global instructions and CONTRIBUTING apply.

## Scope and acceptance before implementation

The Bubbles example's tiny-window View hides both components while Update still
routes keys, paste and viewport scrolling. The responsive reference already
suspends navigation while unusable. Apply that application-level pattern to the
Bubbles example, without renderer or dependency changes.

- At widths below 20 or heights below 5, ignore component interaction messages:
  edits, deletion, paste, focus switching, scrolling and delayed component input.
- Preserve logical input, cursor position, selected focus and viewport offset
  across ignored messages. Resize may still clamp scroll to the new dimensions.
- Esc/Ctrl+C quit at any size; window-size messages always apply; spinner ticks
  continue. No clipboard commands may start from hidden input.
- Pause other component updates while hidden and restart the focused input's
  cursor command when transitioning back to usable dimensions. Avoid restarting
  it on every ordinary resize or for viewport focus.
- Show the required size and quit hint where the tiny allocation allows it.
- Test direct synthetic events, repeated hide/restore cycles, command routing,
  actual subsequent editing/scrolling and unchanged normal example goldens.

No terminal acceptance, ps, blocked cache operation, listener, experimental
upstream patch, Barcode or Foundry action is in scope. Local build/tidy remain
unrun; require permitted local race/vet checks, independent exact-head review,
normal PR and green exact-head hosted CI before a qualified merge.

## Reproduction and result

Before the fix, `TestTinyWindowIgnoresHiddenInteractions` reproduced eight
violations at 19x8: insertion, deletion, Tab/Shift+Tab focus changes, synthetic
paste, and viewport scrolling changed state behind the resize hint or issued
commands. The same regression now passes.

The shared usability predicate gates key routing after global quit keys and
also gates other component messages. Resize and spinner handling remain active.
Hidden cursor messages stop scheduling successors, so the focused cursor is
restarted only on an unusable-to-usable transition. Logical focus stays intact.
The tiny view now names the required 20x5 size and shows Esc:quit when it fits.
No normal-layout, renderer, dependency or golden bytes changed.

## Validation

Go checks use `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`.

- Baseline hidden-interaction test: failed before implementation as expected.
- `go test -mod=readonly -race ./example/bubble_components -count=1`: passed.
- `go test -mod=readonly -race -count=1 ./...`: passed all packages and unchanged
  example goldens.
- `go vet -mod=readonly ./...`, gofmt and diff whitespace checks: passed.

New tests cover 90 hidden-event cases (five allocation shapes, both focus states,
nine event types), quit and spinner continuity, three hide/restore cycles for
each focus, visible cursor recovery, no redundant cursor restarts, subsequent
editing/scrolling, and bounded actionable hints including partial-line cutoffs.
Commands that could read the clipboard are never executed. Cursor restart checks
inspect returned commands and raw rendered cursor cells; no terminal is started.

Local build/tidy and interactive terminal acceptance remain unrun. Existing
upstream combining-mark output loss remains documented and unchanged. No denied
operation was retried. Independent exact-head review and exact-head hosted CI
must qualify the result before merge.
