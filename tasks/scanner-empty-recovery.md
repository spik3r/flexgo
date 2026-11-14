# Recover safely from empty scanner-demo lists

Status: implemented and locally validated; independent exact-head review pending.
Owner: delegated Flexgo worker.
Branch: codex/scanner-empty-recovery.
Reused owned worktree: /tmp/flexgo-bubble-input-feedback.
Base: 23382b41e87d61d27337c5d586c8776e3cb94bcb.

Verified current main, no open PRs and no competing Flexgo claim in active
threads, task records or worktrees. Primary checkout and other worktrees remain
untouched. Disk has 30 GiB free. Applicable instructions and scanner guide were
read; no relevant repository AGENTS.md or .agents skill was found.

## Evidence and acceptance before implementation

Files Bottom sets filesCursor to len(files)-1; on an empty list that is -1.
Open checks only cursor < len(files), so End then Enter indexes -1 and panics.
History Bottom has the same negative cursor, and detailLines only clamps its
upper bound; repopulating history after empty End can panic while rendering.

- Reproduce both failures with synthetic key events and data, no file/scanner IO.
- Keep cursors valid for empty, repopulated and shortened lists. Empty Open must
  be a no-op; recovering data must select/open a valid item without a reset key.
- Clamp rendering projections without mutating stored data; reset history detail
  scroll only when a stale selection is reconciled or the list becomes empty.
- Show a bounded no-files message instead of an unexplained blank panel.
- Preserve normal navigation, open/close and existing populated fixture output.
- Run focused/full checks and independent exact-head review before one normal PR.
  Coordinator owns merge qualification; no dependent follow-on work.

No real scans, file reads, native terminal, clipboard, listener, denied cache or
process operation, Foundry, Barcode or deployment action is in scope.

## Implementation and validation

Both base reproductions failed with an index-out-of-range [-1] panic: empty Files
End then Enter in ScanScreen.Update, and empty History End then repopulation in
detailLines. Update now reconciles list cursors; empty Open is ignored. History
rendering and filetree rendering project valid selections without changing stored
models. The Files panel shows a clipped no-files message. Opened-file content is
kept until Close when its source list disappears.

Passed (exit 0): focused scanner race suite; full
`go test -mod=readonly -race -count=1 ./...`; `go vet -mod=readonly ./...`;
gofmt and whitespace checks. Existing goldens are unchanged. Tests include both
original failures, four repeated recovery cycles for each screen, shortened lists,
open-file retention, five stale cursor projections, seven empty view sizes and
remapped Open/Bottom bindings. Full output:
/tmp/flexgo-scanner-empty-recovery-race.log.
Go checks used GOPROXY=off, GOSUMDB=off, GOTOOLCHAIN=local and the existing
/tmp/codex-local-build-cache.

Local build/tidy remain unrun due to prior cache-write restrictions; hosted CI
must qualify them. Interactive terminal and real asynchronous data replacement
are unverified; the demo has no live scanner or reload pipeline. No denied action
was retried. No dependency or engine changes are included.
