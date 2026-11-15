# Keep debug labels within their reserved row

Status: local implementation complete; frozen for independent review.
Owner: delegated Flexgo contributor, coordinated by Codex under the owner's
related-repository maintenance request.
Base: `0f0f49435b4c0f0c368c9eba028685d2aa46f54f` (verified remote main).
Live PRs and issues are empty. The earlier border task is merged in PR 22 and
left debug label wrapping for a separate change. Preserve that record and all
unrelated worktrees. The primary checkout's README and examples are another lane;
use a separate clean worktree for this task.

## Reproduction and decision

A valid debug node rendered at 10 by 5 with a long ASCII Name gives its bounded
leaf an 8 by 2 callback but paints 10 by 10. An empty Name at 6 by 5 paints six
rows because the dimension label wraps. Unicode, ANSI and explicit line breaks
also add rows. The short-label control fits. All callbacks return exactly their
allocated rectangle; generated metadata alone exceeds its one reserved row.

Keep the existing callback budget. Clip only the complete generated debug label
with existing ClipView semantics at width w and height 1, then keep ordinary bold
styling and padding. Retain the first logical line's cell prefix without an
ellipsis. Long or multiline names may omit the dimension suffix. Normalize tabs
and CR/CRLF as ClipView already does, keep graphemes whole and close SGR/OSC8.
At nonpositive label widths, display no metadata text. Keep the stored Name and
Inspect output unchanged. Do not create another text parser.

Measuring and reserving every wrapped label row would change callback dimensions
and let arbitrary diagnostic names consume the body. Prefer one bounded row,
matching the renderer's existing reservation. This does not change allowed user
View or overconstrained-child overflow, no-shrink policy or frame arithmetic.
Arbitrary cursor-control sequences remain outside ClipView's geometry guarantee;
add no new security, constant-time or input-size claim.

## Exact scope

Only `render.go`, `render_edge_test.go`, `node.go`'s Name field comment, and this
`tasks/bounded-debug-labels.md` note. Clarify the display limit without changing
public types or signatures. No ClipView implementation change is needed.
Do not edit README, examples, prior task records, dependencies, goldens, versions,
release tags, Foundry or consumer pins. Stop for scope review if a golden changes.

## Acceptance and validation

Acknowledge a genuine clean start at this committed claim before edits. Freeze
literal public Render regressions and reproduce their failures on old rendering
before changing production code. Cover long and spaced ASCII, narrow dimensions
with empty Name, wide/combining/grapheme names, SGR/OSC8, LF/CRLF/CR and tabs.
Assert expected visible prefixes, omitted later lines and style closure, exact
callback sizes and bounded physical output when the frame itself fits.

Keep short-label bytes and bounded body sentinel text unchanged. Cover repeated
renders/resizes and preserved Name/styles/Inspect, zero/one content widths,
nonnegative tiny callbacks, margins/padding, fixed sizing and nested Row/Col.
Do not promise a fixed frame fits an allocation smaller than the frame.
Explicit zero, ordinary and custom-wide borders still suppress Debug; disabled
Debug and unwrapped overflowing user View controls remain unchanged.

Run focused regressions, gofmt, tidy with byte-identical go.mod/go.sum, vet,
build and the full uncached race suite with unchanged example goldens. Report
local Go 1.25.9 honestly; require hosted Go 1.25.5 CI before normal merge.
Use owned temporary HOME and the dedicated batch cache. Run no interactive app,
scanner, provider or real-home installation.

Freeze the final diff and observed checks for two independent reviewers. The
coordinator owns staging, commit, push, PR and normal merge after exact-head
checks. No action outside these paths is authorized by this claim.

## Implementation and observed validation

The generated debug label now passes through existing ClipView at the current
label width and one row before its ordinary bold styling. Frame arithmetic,
callback dimensions, caller content and explicit-border precedence are unchanged.
The Name field comment states the cropped display; stored Name and Inspect stay
unchanged. No ClipView implementation or additional text parser changed.

Frozen public Render tests failed on the old renderer for long and narrow
labels, Unicode/ANSI/line-break names, nested layouts, spacing and zero-width
metadata. Exact short-label bytes and body/overconstrained-child overflow
controls already passed. After the change, all cases pass, including SGR/OSC8
closure before the body, grapheme prefixes, tabs, fixed sizing, tiny callbacks,
repeated renders/resizes and unchanged styles/Inspect output.

Observed local checks with Go 1.25.9 and cached dependencies:

- Focused debug-label regressions with the race detector: passed.
- `gofmt -l .`: no output.
- `go mod tidy`: passed; go.mod and go.sum stayed byte-identical.
- `go vet ./...` and `go build ./...`: passed.
- `go test -race -count=1 ./...`: passed, including unchanged example goldens.
- Working and staged whitespace checks: passed; nothing was staged.

The local toolchain was Go 1.25.9. Hosted Go 1.25.5 and independent exact-head
review remain required before normal merge. No interactive application,
external scanner/service or real-home installation was launched. No dependency,
golden, README, example, prior task, Foundry or consumer pin changed.
