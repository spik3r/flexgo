# Diagnose cyclic layout trees

Status: local commit for independent review

Base: local origin/main e616874. Isolated branch
`codex/layout-tree-validation` in `/tmp/flexgo-layout-tree-validation`.
Read-only live GitHub inspection on 2026-11-13 found no open PRs. Existing
worktrees and dirty primary checkout were preserved. The adaptive-key-hint
commit 160a0e8 was not changed, copied or published.

`Validate` now reports a cyclic Children reference with its ancestor path,
then continues checking siblings. `Inspect` marks that reference `<cycle>`.
Both track ancestors per traversal path, allowing repeated acyclic subtrees.
Tests cover a self-cycle with another invalid sibling, a longer cycle,
repeated deterministic inspection and shared subtrees with valid and invalid
nodes. Existing inspection output and layout goldens pass unchanged.

Validation on Go 1.27.0, GOPROXY=off, GOSUMDB=off, GOTOOLCHAIN=local,
GOCACHE=/tmp/codex-local-build-cache:

- `go test -mod=readonly ./...`: passed.
- `go test -mod=readonly -race ./...`: passed.
- `gofmt -l inspect.go inspect_cycles_test.go`: no output.
- `git diff --check`: passed.
- `go vet -mod=readonly ./...` and `go build -mod=readonly ./...` were
  invoked before the race run in the same shell. The shell finished at 0,
  but separate exit codes for vet/build were not captured, so those checks
  are not independently attested here. Go printed a blocked stat-cache write
  at `~/go/pkg/mod/cache/download/github.com/spik3r/flexgo/@v/` for e616874.
  That metadata write was not retried or escalated.

`Render` still assumes a valid tree; callers should validate generated trees
before rendering. This work does not bound arbitrarily deep acyclic graphs or
make concurrently mutated nodes safe. No dependencies changed, no `go mod tidy`
ran, and nothing was pushed. Independent review and any remaining build/vet
qualification should precede publication.
