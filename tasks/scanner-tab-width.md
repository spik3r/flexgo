# Keep scanner-demo tabbed rows within their allocation

Status: local implementation for independent exact-head review.
Owner: delegated Flexgo/barcode worker.
Branch/worktree: `codex/scanner-tab-width`, `/tmp/flexgo-scanner-tab-width`.
Base: `7741ed91ad56b487b7f9afe7e68e00cc7d87a75b` (verified remote main).
Open PR query returned none; current worktrees and task records showed no
competing scanner-row claim. Primary changes and other worktrees were preserved.
Disk had 37 GiB available before validation. Barcode #7 remains untouched.

## Reproduction and change

A synthetic regression demonstrated that padLineBg("A\tB", 3, ...) rendered
seven columns instead of three. Repeated tabs also widened file-tree rows and
pushed viewport scroll indicators beyond their allocated columns. Width
measurement counted tabs as zero cells, while Lipgloss later expanded them to
four spaces; background padding then added further surplus width.

Normalize tabs to the same four spaces before the existing clipping and padding
calculations. This affects the scanner demo's viewport and file-tree rows only.
The core library, recipes, key hints, dependencies, existing truncation policy,
scroll state and row highlighting are unchanged. This demo displays synthetic
content; no real scanner or repository inspection runs during the tests.

## Validation

Go checks use `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`:

- `go test -mod=readonly ./demo/scanner -run 'TestPadLineBgMeasuresExpandedTabs|TestScannerRowsKeepTabContentInsideAllocation' -count=1`:
  failed on the original implementation, reproducing helper and full-widget
  width failures.
- `go test -mod=readonly -race ./... -count=1`: passed all packages, including
  scanner regressions, existing example golden snapshots and responsive tests.
- `go vet -mod=readonly ./...`: passed.
- gofmt and diff whitespace checks passed.

Regression coverage checks exact expanded/clipped text, leading/repeated tabs,
CJK/combining text, plain text, empty helper allocations and viewport/file-tree
sizes at widths 1–12, heights 1–4 and multiple cursor/scroll offsets. It checks
that the scroll indicator remains in the last column.

Local build/tidy are deliberately unrun because of the prior denied Go
metadata-cache write. No retry, workaround, listener, billing change or deployment
is included. Interactive terminal navigation, rendering and visual acceptance
remain unrun. Hosted checks remain required after independent review/publication.
