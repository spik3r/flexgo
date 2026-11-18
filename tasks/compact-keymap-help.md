# Make all reference help actions reachable in smaller terminals

Status: implementation and local checks complete; awaiting independent review and CI.
Owner: Codex coordinator with delegated builder.
Branch: codex/compact-keymap-help.
Base: 178499e606146c0af5e94f1f5f19376dc9c115f3, freshly fetched main.

The owner authorized continued UI improvements with cheaper agents, independent
review and tested merges. No Flexgo PR is open. Foundry PR 185 and Flexgo PR 26
are merged with exact reviewed trees verified. Preserve shared checkout work.

## Requirement and decision

The actual help-open App returns 80 by 36 at an 80 by 24 terminal and 40 by 60
at 40 by 16. Existing modal routing ignores scroll keys; the final actions
cannot be brought into the allocated frame. This is a new bounded demo
requirement, not a change to core fixed-size or no-shrink behavior.

Keep help-open App frames within both target sizes. Every existing nonblank
HelpEntries key label and description must be reachable, including the final
Close action. There are twelve actions and four blank separators. Use a compact
terminal-cell-wrapped help viewport when the existing card cannot fit. Keep the
initial full-card ANSI bytes unchanged at 100 by 40 and 120 by 40. Count and
paint the same wrapped rows; do not truncate away descriptions or long labels.

Use current remappable Up, Down, PageUp, PageDown, Top and Bottom only for modal
navigation. Keep ten-row page steps consistent with existing demo scrolling.
Close and Help retain first priority. Other keys, including global quit, screen
switching, item/tab/file actions and text, remain captive. Opening starts at
zero; reopening resets. Resize retains a legal top position and clamps only
when needed. Nonpositive allocations preserve intent and ignore navigation;
painting stays pure. Existing underlying active-screen resize reconciliation
may still legally clamp its viewport; preserve selection, profile and cursor.

Derive compact navigation/close hints from current usable bindings. A Close or
Help collision wins over scrolling, so a colliding canonical alias must not be
advertised as usable navigation; use another existing safe alias or omit that
hint. Empty bindings must not invent keys. Existing HelpEntries may retain raw
configured labels. No new bindings or editing. Below fixed shell chrome, retain
the documented tiny-size limitation; do not claim arbitrary terminal support.

## Exact scope

Only demo/scanner/app.go, widgets.go, new help_view.go if a small demo-local
state/wrapping helper makes the code clearer, new help_compact_test.go,
demo/scanner/README.md and this task. Existing tests remain unchanged. Keep
Launcher, Scan and History body implementations, core renderer, ClipView,
KeyMap definitions, dependencies, root README, example goldens and releases
unchanged. Reuse merged narrow chrome rather than adding a competing fix.
No scanner, target file, provider, network or URI action.

## Validation and handoff

Commit this reviewed public-safe claim and record a genuine clean start before
edits. Freeze regressions and prove failure on unchanged production. Capture
full fitting help bytes before edits. Use actual App WindowSize/key messages
and ANSI-aware per-line bounds; collect every label/description across scroll
frames rather than checking only the first frame. Cover each underlying screen,
last Close reachability, line/page/top/bottom keys, remaps and collisions,
empty and long/wide bindings, wrapping, resize round trips, open/reset, no
render needed between messages, pure paint and no leakage. Verify update size
against actual compact callback allocation and test card-transition boundaries.
Preserve existing compact Launcher, viewport, recovery and golden tests.

Run focused and full race tests, vet, build, gofmt and tidy with dependency files
unchanged. Record the actual local Go version and require hosted CI. Independent
review covers frozen source and immutable commit. The coordinator owns named
staging, commit, push and normal merge after all gates. No next implementation
starts before that merge.

## Observed implementation and checks

The coordinator recorded a genuine clean start at claim commit
`368b4daeee45b89ac5dd7e55a66e24f3ffe9c64e` before implementation edits.
Frozen actual-App regressions failed against unchanged production: help returned
36 rows at 80 by 24 and 60 at 40 by 16, final entry text could not be brought into
view, modal scroll keys did nothing, and remapped close hints stayed hardcoded.
Long/wide Next and Prev labels also expanded the full card at larger terminals.
Six full ANSI baselines were captured before edits at 100 by 40 and 120 by 40,
covering each underlying screen; all six hashes still match exactly.

App now owns a help offset. Current modal scroll keys move only help, while
Close and Help keep first priority and other keys stay captive. The full card is
used only when its actual rendered extent fits the current body allocation and
labels. Otherwise a compact title and wrapped viewport show the same row list
used to count scroll limits. Wrapping uses the existing ANSI-aware terminal-cell
helper. Compact hints choose current aliases that survive dispatch priority,
prefer a fitting close key, and omit empty or colliding aliases.

Opening resets to the first row; resize preserves a legal top row and clamps only
when needed. Nonpositive allocations, including MinInt sizes, preserve position
and ignore navigation. Rendering and callback reuse leave owned and captured
positions unchanged. Existing underlying screen resize reconciliation remains in
place. Actual callback allocation tests cover the compact geometry; transition
checks cover body widths 54 through 58 and 74 and heights 22 through 28, including
the centred frame's spacer rows. Long labels select the compact path even at
100 by 40 and 120 by 40 rather than overflowing the fixed key column.

With default bindings, all twelve entry labels and descriptions, including Close,
were collected across scroll frames on Launcher, Scan and History at both target
sizes. Tested safe remaps and long/wide labels also retain text reachability.
No new navigation is invented when bindings are empty. Private actual-App ANSI
and plain frames at both sizes show the initial and final help pages. At 40 by 16
the final page has offset seven and includes the Close description.

Local validation used Go 1.27.4 on darwin/arm64, an isolated cache, GOWORK and GOENV
off, offline dependencies and readonly module inputs. Observed final checks
passed:

- `go test -count=1 ./demo/scanner`
- `go test -race -count=1 ./...`, including unchanged example golden tests
- `go vet ./...`
- `go build ./...`
- `gofmt -l .` returned no paths
- `go mod tidy` left go.mod and go.sum byte-identical

Existing tests were not edited. Launcher, Scan and History bodies, core,
ClipView, KeyMap definitions, dependency files and goldens remain unchanged.
Independent review and hosted CI remain required; the coordinator owns staging,
publication and normal merge after those gates pass.
