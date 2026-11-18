# Keep the selected History run inside its sidebar

Status: claimed; implementation has not started.
Owner: Codex coordinator with a delegated builder.
Branch: codex/history-sidebar-viewport.
Base: 83628794238856bd0f2e93c386c6ec2287b6dd0e, freshly fetched main.

The owner requested continued work across Foundry and its related UI libraries,
using cheaper subagents and reviewed, tested merges. The coordinator selected
this reproduced reference-app defect after merging the preceding Foundry fix.
No Flexgo PR is open and no existing claim overlaps this sidebar change.
Preserve the primary checkout's README, examples and untracked artifacts.

## Observed defect and decision

With the five supplied runs, History's sidebar paints 28 by 9 in a 28 by 5
allocation. The History page paints 80 by 19 in an 80 by 15 terminal. A selected
last run can lie outside the allocated rows. The callback renders every result;
Lipgloss Height sets a minimum rather than clipping those rows.

Use a bounded selection-following window for run rows. Keep the current
selection and detail viewport state; rendering and resizing must remain pure.
At tiny heights, decoration yields to the selected row. At narrow widths, keep
text cell-bounded with existing clipping helpers rather than a new ANSI parser.
When the normal list fits, preserve its complete rendered bytes.

## Exact scope

Only demo/scanner/history.go, new demo/scanner/history_sidebar_test.go, a small
behavior note in demo/scanner/README.md, and this task record. Keep navigation
bindings, Update's existing stale-selection reconciliation, detail scrolling,
core Render/ClipView and sizing contracts unchanged. No dependency, root README,
example golden, version, release, Foundry or consumer-pin changes. Scan Logs has
a separately observed scrolling defect; do not fold that repair into this task.

## Acceptance and validation

Commit this reviewed public-safe claim and acknowledge a genuine clean start
before implementation. Freeze literal regressions and reproduce their failures
on unchanged main. Nonpositive callback width or height returns empty text.
Positive sizes are bounded by their allocated cells. For heights one through
five, the projected selected run remains visible whenever its text can fit;
heading, blank lines and padding must not hide it. Cover first/middle/last runs,
Up/Down/Home/End navigation, resize and repeated renders, long/wide text and
supported ANSI styling. Do not promise arbitrary cursor controls are harmless.

Cover empty and shortened lists, with stale cursors projected safely without
mutating the model during rendering. The next Update keeps its existing
reconciliation and detail-reset rules. Preserve the normal 28 by 12 output
bytes, data and detail offset. The 80 by 15 page fits when its fixed frame fits;
this does not change the core no-shrink behavior for overconstrained frames.

Run focused and full race tests, gofmt, tidy with byte-identical go.mod/go.sum,
vet and build. Existing example goldens must stay unchanged. Report the actual
local Go version; require hosted CI at the module's Go version before merge.
A separate reviewer and the coordinator review the exact frozen commit.
The coordinator owns staging, publication and normal merge. No real scanner,
network/provider call, interactive app or home installation is part of testing.
