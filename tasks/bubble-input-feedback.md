# Show Bubbles input errors in the reserved status row

Status: implemented and locally validated; independent review pending.
Owner: delegated Flexgo worker.
Branch/worktree: `codex/bubble-input-feedback`, `/tmp/flexgo-bubble-input-feedback`.
Base: `7f3bb9f85982cc190fcfe22f25a2c1693d4952d4` (fresh remote main).

PR/issue queries returned no open work; current task listings, worktrees and
records showed no competing Flexgo owner. Primary changes and other worktrees
are preserved. Disk has 24 GiB free. Global instructions, CONTRIBUTING and the
example guide apply; no relevant repository AGENTS.md or .agents skill exists.

## Evidence and acceptance before implementation

Pinned Bubbles textinput stores clipboard-read failures in Err via pasteErrMsg;
validation errors also populate Err during ordinary input updates. The example
renders input.View without showing Err, so failures have no visible feedback.
The input already owns three rows but uses only a header and a text row.

- Reproduce missing feedback using a synthetic validator through Update, without
  reading or manipulating the system clipboard.
- Show a concise Error prefix and the component error on the existing third row.
  Use the leaf's existing ClipView bounds for long, multiline and styled errors.
- Keep input value/cursor/focus and viewport state unchanged by rendering. The
  component remains the sole owner of error state; successful edits/paste clear
  feedback when Bubbles clears Err, with no second error cache.
- Preserve error feedback across focus changes and hide/restore cycles; tiny
  layouts retain their current resize/quit hint and interaction guard.
- Keep viewport/footer placement, normal error-free goldens, renderer code and
  dependencies unchanged. Document error behavior and its qualification limits.
- Require focused/full local race tests, vet, independent exact-head review,
  normal PR and exact-head hosted CI before qualified merge.

No clipboard read, terminal acceptance, ps, blocked cache operation, experimental
patch, Barcode, Foundry or deployment action is in scope. Local build/tidy are
not retried; existing upstream combining-mark limitations remain documented.

## Implementation and local validation

The synthetic validator regression failed on the base because the error was not
visible. The input leaf now appends Error and the component detail in its existing
third row; the existing ClipView clips it. No second error state is introduced.
Four test groups cover validation/recovery, synthetic paste recovery, 15 bounded
error/size cases and repeated focus/tiny-window recovery. The guide documents the
behavior. No renderer, dependency or golden file changed.

Passed (exit 0): focused example race tests; full
`go test -mod=readonly -race -count=1 ./...`; `go vet -mod=readonly ./...`;
gofmt and `git diff --check`. Go checks used GOPROXY=off, GOSUMDB=off,
GOTOOLCHAIN=local and the existing /tmp/codex-local-build-cache. Full race output
is recorded in /tmp/flexgo-bubble-input-feedback-race.log.
Local build/tidy remain unrun because of previous cache-write restrictions;
normal hosted PR CI must qualify them. Clipboard failures were represented by a
synthetic stored component error; no native clipboard or interactive terminal
acceptance was performed. Existing combining-mark renderer limitations remain.
