# Match scanner scroll state to its visible viewport

Status: claimed; implementation has not started.
Owner: Codex coordinator with a delegated builder.
Branch: codex/scanner-viewport-state.
Base: 443a6b406f7e15f3dcb9ba190c2665715ac8e3a3, freshly fetched main.

The owner authorized continued related-library improvements with cheaper agents,
review and tested merges. The coordinator selected this separately reproduced
scrolling defect after Foundry PR 184 and Flexgo PR 24 were merged and verified.
No Flexgo PR is open or overlapping claim active. Preserve the primary checkout's
unrelated README, examples and untracked artifacts.

## Observed defect and decision

For 39 log lines and ten visible rows, repeated PageDown stores offset 38 while
rendering projects offset 29 on a copy. Up changes stored offset to 37 but the
visible top line stays at 29. Several reverse keys are needed before anything
moves. Other Scan viewports, the opened-file viewer and History details share
the one-row update assumption. Reusing one viewport callback at different sizes
also mutates its captured scroll state during painting.

Plumb actual demo viewport dimensions from the root model into screen updates.
Reconcile the active offset before navigation and after resize or activation;
apply the requested movement from that legal position. Bottom uses the maximum
legal top row where that binding is already handled. Keep ten-row PageUp and
PageDown steps and all existing key ownership and file-list cursor semantics.
Do not add bindings or list paging.

Preserve the current top row across resize whenever legal, clamping only when
necessary. An unallocated viewport preserves its last logical position and
ignores scroll input until positive dimensions return. Inactive tabs retain
positions until activation, then reconcile before input. New Scan instances and
History activation must receive the current size without needing an intervening
View call. Modal input remains captive.

Rendering computes only a local projected offset. Neither model state nor a
closure's scroll state changes when painting, including negative or stale
provided positions. Repeated renders at different sizes must not alter later
output at the original size. Keep fixed-frame and no-shrink layout behavior;
this task does not promise that overconstrained chrome fits tiny terminals.

## Exact scope

Only demo/scanner/app.go, scan.go, history.go, widgets.go, a new
viewport_size.go helper, new viewport_state_test.go, existing widgets_test.go
if needed without weakening assertions, a behavior note in demo/scanner/README.md,
and this task. Keep core layout/ClipView, dependencies, example goldens, the
recent History sidebar window, root README, versions and release tags unchanged.
Any internal size helper must share the relevant construction constants and be
checked against actual leaf allocations, not a duplicate untested formula.

## Acceptance and validation

Commit this reviewed public-safe claim and record a genuine clean start before
implementation. Freeze regressions and prove failures on unchanged main.
At 39 lines and ten visible rows, PageDown and Bottom stop at offset 29; one Up
moves to 28 and visibly changes content. Verify through actual App updates and
WindowSize messages, without a required render between messages.

Cover visible heights one, two, ten and at least the content length; positive
and nonpositive widths; resize larger/smaller/zero/restore; tab and opened-file
transitions; History detail PageUp/PageDown; remapped keys and captive help.
Resizing ten visible rows to twenty clamps 29 to 19, then Up gives 18; shrinking
to five preserves 29 and Down gives 30. Keep independent tab positions. Preserve
selection/recovery and modal/navigation behavior already tested.

Compare actual layout callback dimensions with the update helper. A ten-row,
twenty-row, ten-row reuse of one callback must yield identical first and last
outputs with unchanged owned and captured offsets. Keep ordinary legal-position
render bytes and existing goldens unchanged. Do not read files or run a scanner;
use the demo's synthetic data and key messages only.

Run focused and full race tests, formatting, tidy with byte-identical dependency
files, vet and build. Report the actual local Go version and require hosted CI
at the module's Go version. A separate reviewer and the coordinator review the
frozen and committed source. The coordinator owns staging, publication and normal
merge after all gates pass. No other implementation begins before this merge.
