# Keep responsive record selection visible

Status: implemented and locally validated; independent exact-head review pending.
Owner: delegated Flexgo worker.
Branch: codex/responsive-record-rows.
Reused clean owned worktree: /tmp/flexgo-bubble-input-feedback.
Base: 8661f7e2d3340fa0f2acd7c9a4f8a8ce43b05e3a.

Current remote main, open PR/issue queries, worktrees, task records and active
thread listings show no competing Flexgo work. The primary checkout and other
worktrees remain untouched. Disk has 17 GiB free; no installs or new checkout.
Read global instructions, CLAUDE.md, CONTRIBUTING.md and the responsive guide.
No repository AGENTS.md or applicable .agents skill was found.

## Evidence and acceptance before implementation

The list derives its visible start from selected record and allocated row count.
It appends each title as one entry, then clips the joined text. A title with an
embedded line break can occupy multiple rows and hide the selection that the
row calculation intended to keep visible.

- Reproduce a missing selected record with synthetic multiline titles.
- Limit each record title and its selection prefix to one terminal row using
  existing ClipView behavior, including tabs, CRLF, wide text and ANSI closure.
- Keep selection visible across compact/wide/tiny resize and navigation cycles;
  retain stored titles, detail content, selection and scroll state.
- Keep normal fixture output, renderer and dependencies unchanged.
- Run focused and aggregate permitted tests, independent exact-head review and
  open a normal PR if accepted. No dependent PR or merge is part of this task.

No native terminal, clipboard, listener, denied process/cache operation,
experimental patch, Barcode, Foundry or deployment action is in scope.

## Implementation and validation

The base regression failed: a two-row list returned the first title and its
continuation instead of the selected second record. Each prefixed title is now
clipped to one row before rows are joined. Nonpositive callback allocations
return empty output. The guide documents the first-line projection.

Passed (exit 0): focused responsive race tests; full
`go test -mod=readonly -race -count=1 ./...`; `go vet -mod=readonly ./...`;
gofmt and `git diff --check`. Tests cover four multiline forms, 20 title/width
cases with raw ANSI/OSC-8 closure checks, nonpositive callback areas, and 48
selection/size combinations plus detail/back/tiny recovery. Normal example
goldens passed unchanged. No renderer, dependency or golden modifications.
Go checks used GOPROXY=off, GOSUMDB=off, GOTOOLCHAIN=local and the existing
/tmp/codex-local-build-cache. Full output: /tmp/flexgo-responsive-record-rows-race.log.

Local build/tidy remain unrun because of prior cache-write restrictions; normal
hosted PR CI must qualify them. Interactive terminal acceptance was not performed.
The latest disk check showed 21 GiB free. No denied operation was retried.
