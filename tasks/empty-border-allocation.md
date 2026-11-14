# Preserve content space when an explicit border is empty

Status: local implementation for independent exact-head review.
Owner: delegated Flexgo worker.
Branch/worktree: `codex/empty-border-allocation`, `/tmp/flexgo-empty-border-allocation`.
Base: `ff40fbe5cd759d45dcc25b73fe98fc16556e706b` (verified remote main).
GitHub returned no open issues or PRs. Current worktrees and task records showed
no competing border-allocation claim. The existing Bubbles example worktree was
left alone. Primary changes and other worktrees were preserved. Disk had 37 GiB
free. Barcode #7 and its CI were not touched.

## Reproduction and scope

Node documents that ShowBorder=true with Lipgloss's zero Border renders no visible
frame. The pinned Lipgloss implementation returns content unchanged for that
border. Flexgo still reserved two columns and two rows: an 8x4 leaf received a
6x2 callback, with or without Debug enabled. The synthetic baseline regression
failed for both cases.

Distinguish the explicit-border selection (which still suppresses Debug) from a
nonempty frame (which needs reserved space and border painting). A zero border
now preserves the same allocation and output as an unframed node. HiddenBorder
contains space glyphs and still reserves its frame. Nonempty border sizing,
Debug sizing, public fields and content-overflow policy remain unchanged.

## Validation

Go checks use `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`:

- `go test -mod=readonly . -run TestEmptyExplicitBorderKeepsLeafAllocation -count=1`:
  failed before the fix for both Debug states, reproducing 6x2 instead of 8x4.
- `go test -mod=readonly . -run 'TestEmpty|TestNonemptyBorders' -count=1`:
  passed after the fix.
- `go test -mod=readonly -race ./... -count=1`: passed all packages, including
  unchanged example goldens, scanner regressions and responsive tests.
- `go vet -mod=readonly ./...`: passed.
- gofmt and diff whitespace checks passed.

Regression tests cover leaf callbacks, nested Row/Col containers, symmetric and
asymmetric margin/padding, fixed dimensions, automatic margins, zero-border
Debug precedence, NormalBorder, HiddenBorder, disabled borders and Debug-only
controls. All inputs are synthetic; no terminal session or scanner is launched.

Local build/tidy remain unrun due to the prior denied Go metadata-cache write.
Interactive terminal navigation/rendering and visual acceptance remain unrun.
No denied-action retry, cache workaround, listener, deployment or billing change
is included. Independent exact-head review precedes publication; qualified hosted checks
precede any merge.
