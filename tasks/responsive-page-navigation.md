# Add page and boundary navigation to the responsive example

Status: implemented and locally validated; independent exact-head review pending.
Owner: delegated Flexgo worker.
Branch: codex/responsive-page-navigation.
Reused owned worktree: /tmp/flexgo-bubble-input-feedback.
Base: 7ac999a8577d607586954ff343b9676986eda762.

Fresh remote main matches merged #19. Open PRs/issues are empty; current task
records, worktree inventory and active threads show no competing Flexgo owner.
Disk has 32 GiB free. Primary checkout and other worktrees are preserved.
Applicable global instructions, CLAUDE.md, CONTRIBUTING.md and responsive guide
were read in this lane; no repository AGENTS.md or relevant .agents skill exists.

## Evidence and acceptance before implementation

The example has 12 records with 33 detail rows each, but only Up/Down and j/k
movement. Page Up/Down and Home/End are ignored. Reproduce that missing behavior
with synthetic key messages, then add a bounded active-pane navigation contract:

- Page Up/Down move by the visible content height (terminal height minus four
  reserved rows); Home/End move to the first/last record or detail position.
- Clamp all movement. Detail navigation starts from the currently visible
  projection after resize, preserving stored scroll intent until an input occurs.
- Only a changed list selection clears its detail offset. Boundary no-ops retain
  selection, detail position and focus. Empty/tiny views ignore movement.
- Keep selection visible and output bounded. Show keys and quit/back controls at
  every usable breakpoint; update the intentional help-only golden change.
- Run focused/full permitted checks and independent exact-head review before
  publishing one normal PR. No merge or dependent follow-on work.

No dependency, engine, native terminal, clipboard, listener, denied process/cache
operation, Barcode, Foundry or deployment change is in scope.

## Implementation and validation

The base failed the new Page Down test (selection stayed 0 instead of moving to
5 at 40x9). The active-pane movement branch now handles existing one-row keys,
page keys and Home/End with one clamp. List selection changes clear detail scroll;
list boundary no-ops keep it. Help exposes the new keys at compact and wide sizes.
Only the two help rows changed in the responsive golden.

Passed (exit 0): focused responsive race tests; focused golden regeneration;
full `go test -mod=readonly -race -count=1 ./...`;
`go vet -mod=readonly ./...`; gofmt; source/docs whitespace checks.
The golden contains intentional terminal-cell trailing padding; its whitespace
check disables blank-at-eol for that file only. Full race output is at
/tmp/flexgo-responsive-page-navigation-race.log. Go checks use GOPROXY=off,
GOSUMDB=off, GOTOOLCHAIN=local and /tmp/codex-local-build-cache.

Regression coverage exercises page/boundary sequences at four widths, visible
projection after resize, remembered scroll at list boundaries, empty/short/tiny
content, repeated recovery and complete bounded help at minimum/wide sizes.
Existing single-row navigation and all other golden tests still pass.

Local build/tidy remain unrun due to prior cache-write restrictions; hosted CI
must qualify them. Synthetic key dispatch is tested; interactive terminal
acceptance has not been performed. No denied action was retried.
