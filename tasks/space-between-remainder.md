# Distribute every spare cell between children

Status: local implementation for independent exact-head review.
Owner: delegated Flexgo worker.
Branch/worktree: `codex/space-between-remainder`, `/tmp/flexgo-space-between-remainder`.
Base: `db35c1a570068b7cb655f3f553a2bc7824aa5b66` (verified remote main).
No open PR or competing claim was found. The older gap-justification worktree
is clean and its fe131ec commit is already an ancestor of main; it covered
center/end gap accounting and two-child space-between cases, not this remainder.
Other worktrees and primary changes were preserved. Disk had 37 GiB available.
Barcode #7 was not touched or rerun.

## Evidence and scope

A public Render regression failed on both axes: three one-cell children in eight
cells with Gap=1 and JustifySpaceBetween produced `A  B  C `, leaving the last
child one cell short of the far edge. Integer division in join discarded the
remainder and the parent rendered it as trailing padding.

Distribute leftover cells among earlier inter-child gaps, like the existing
allocator's integer tie-breaking. This keeps the first and last children at the
edges when space permits, with gap sizes differing by at most one cell. Configured
Gap remains a minimum. Single/empty child lists and overconstrained layouts keep
their existing behavior. Only the space-between join branch changes; no layout
engine, public type, dependency, user content or truncation policy changes.

## Validation

Checks use `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`:

- `go test -mod=readonly . -run TestSpaceBetweenKeepsLastChildAtFarEdge -count=1`:
  failed both row and column before the fix as expected.
- `go test -mod=readonly . -run TestSpaceBetween -count=1`: passed all three
  regression groups after the fix.
- `go test -mod=readonly -race ./... -count=1`: final run passed all packages,
  including unchanged example goldens, scanner tests and responsive tests.
- `go vet -mod=readonly ./...`: passed.
- gofmt and diff whitespace checks passed.

Tests cover public Render placement on both axes, three/four children, divisible
and nondivisible spare space, zero-sized gaps, minimum gaps, insufficient space,
empty/single/two-child controls, cross-axis dimensions and painted gap rows.
An initial full-suite run failed only the new paint test's overly literal ANSI
comparison (blank fill adds redundant style sequences); the assertion was fixed
to check painted spaces on each affected row. No production change was needed.

Local build/tidy remain unrun due to the prior denied metadata-cache write.
Interactive terminal navigation and visual acceptance are unrun. No denied-action
retry, cache workaround, listener, deployment or billing changes are included.
Hosted checks remain required after independent review and publication.
