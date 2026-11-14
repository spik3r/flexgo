# Bound modal titles and form labels

Status: local implementation for independent review; not published.

Owner: delegated Flexgo/barcode worker. Branch: `codex/bounded-recipe-labels`.
Worktree: `/tmp/flexgo-form-field-snapshot`.
Base: `10e821478b1d91e92a4f35028bbc8fe07ae51d8f` (fresh remote main).
Open PR query returned none. Existing worktrees and task records showed no
competing label/modal claim; primary checkout changes were preserved.

## Scope and reproduction

A modal reserving a single title row grew from 14x8 to 14x10–12 when the title
contained explicit newlines or wrapped. A two-row form grew from 14x2 to 14x6
with a multiline label, displacing the next field. Both synthetic regression
tests failed before the implementation.

Static labels now truncate each explicit line by terminal-cell width with an
ellipsis. Modal titles use one line; form labels retain up to the allocated
row height. Excess lines add an ellipsis and close ANSI styles and hyperlinks
whose original closing sequences may have been discarded. Empty allocations
produce no label. Alignment, form colon suffixes, custom multi-row forms and
normal examples remain supported. This does not change user body/field callbacks,
core overflow policy, fixed-size overconstraint behavior or the prior Tabs work.

## Validation

All Go commands used `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local`
and `GOCACHE=/tmp/codex-local-build-cache`:

- `go test -mod=readonly ./layouts -run 'TestModalTitleKeepsBodyInItsRow|TestFormLabelsKeepRowsAndFieldsAligned' -count=1`:
  failed both regressions on the original implementation as expected.
- `go test -mod=readonly ./layouts -count=1`: passed after the initial fix.
- `go test -mod=readonly -race ./layouts . -count=1`: passed after final code/tests.
- `go test -mod=readonly -race ./... -count=1`: passed all packages, including
  unchanged example golden snapshots and responsive/autogrid tests.
- `go vet -mod=readonly ./...`: passed, exit 0.
- `go build -mod=readonly ./...`: exit 0, but emitted a blocked Go stat-cache
  write warning for `/Users/kaitait/go/pkg/mod/cache/download/github.com/spik3r/flexgo/@v/v0.1.1-0.20261113121353-10e821478b1d.info*.tmp`.
  This is not a clean build check. No retry, escalation or alternate cache was
  used after that warning. Hosted build verification remains required.
- gofmt and diff whitespace checks passed.

Tests cover CR/LF/CRLF, CJK and combining graphemes, narrow/empty/negative callback
allocations, explicit two-row labels, field/body placement, SGR and OSC-8 BEL/ST
closures in raw output. Interactive terminal behavior is unverified. Local
`go mod tidy` was deliberately not run because of the prior denied metadata-cache
operation. No dependency manifests, versions, sums or golden files changed.
No listener, deployment, push or merge was performed for this assignment.

## Review correction: tab expansion

Independent review identified that ANSI width measurement treats tabs as zero
cells while Lipgloss replaces each tab with four spaces before wrapping. A new
synthetic regression reproduced callback and full-layout overflow with `A\tB`,
leading/repeated tabs, CJK, SGR, OSC-8 and multiline input at narrow widths.
The test failed on 5d9367108c071f4054e83e1b9248ba2fc1906833 before the correction.

`boundedLabel` now performs Lipgloss's four-space tab expansion before measuring
and truncating. Tests cover callback bounds, full modal body/form row placement,
and exact visible expansion when enough width exists. No change to Tabs or core
rendering is included.

With the same offline/read-only Go flags and existing build cache listed above,
`go test -mod=readonly -race ./... -count=1` and `go vet -mod=readonly ./...`
both passed after the correction. Formatting/diff checks passed. Disk had 39 GiB
available before validation. Build/tidy were not rerun, respecting the prior
blocked metadata-cache operation. Interactive terminal checks remain unrun.
The revised exact head requires independent review before publication.
