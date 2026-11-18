# Make all reference help actions reachable in smaller terminals

Status: claimed; implementation not started.
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
