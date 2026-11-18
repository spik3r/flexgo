# Keep the reference Launcher usable in smaller terminals

Status: claimed; implementation not started.
Owner: Codex coordinator with a delegated builder.
Branch: codex/compact-launcher.
Base: 7c9f78c9658797518d1dc5644bddc0707016c847, freshly fetched main.

The owner authorized related UI improvements with cheaper agents, independent
review and tested merges. No Flexgo PR is open. Preserve unrelated work in the
primary checkout. This is a new bounded reference-demo requirement: the existing
docs describe a fixed centered card, not a core sizing guarantee.

## Requirement and decision

The actual default App Launcher paints 80x32 for an 80x24 terminal and 40x55 for
40x16. Keep the help-closed Launcher within both requested dimensions. All five
fields (Name, Root, Include, Exclude, Deep) must remain reachable and the selected
field visible through the existing cursor keys. Preserve current profile data,
cursor, screen transitions, modal capture, remapped keys and synthetic Scan
handoff. Text, Enter and Space remain inert; Deep is not a toggle or editor.

Use a compact demo-only presentation when the allocated body cannot contain the
existing card. Prefer a title and five bounded rows over new scroll state. If a
smaller body needs a window, derive it from cursor and allocation while painting;
never mutate profile or cursor. Selected values may clip by terminal cells.
Keep a visible short start hint and existing global navigation. Bound narrow
shell text only where needed, preserving fitting shared chrome bytes.

Preserve the current fitting Launcher output byte-for-byte at 100x40 and120x40.
Painting is deterministic and leaves state unchanged. Resize100x40 to80x24 to
40x16 and back retains profile, cursor and active screen. Help activation and
close remain reachable; the fixed help card content is explicitly separate from
this requirement. Below the fixed shell chrome, existing overflow limitations
remain; do not claim arbitrary tiny-size support.

## Scope

Only demo/scanner/launcher.go, app.go if needed for actual body allocation,
widgets.go if needed for bounded shell text, new launcher_compact_test.go,
existing widgets_test.go if needed without weakening tests, a behavior note in
demo/scanner/README.md, and this task. No core layout/render/ClipView changes,
Scan or History body changes, key bindings, editing, dependencies, root README,
example goldens, versions or releases. Do not run a scanner or read target files.

## Validation and handoff

Commit this reviewed public-safe claim and record a genuine clean start before
implementation. Freeze regressions and prove failures on unchanged production.
Use actual App WindowSize and key messages, ANSI-aware per-line width and frame
height, all five cursor positions, resize round trips, deterministic output,
profile handoff, modal/global routing and existing rebindings. Capture fitting
bytes before edits and compare after. Check a custom profile with long and wide
values without interpreting control strings as markup or adding a new policy.

Run focused tests and full race tests, vet, build, gofmt and tidy with dependency
files unchanged. Record actual Go version. A separate reviewer reviews the frozen
source and immutable commit; hosted CI must pass before the coordinator merges.
The coordinator owns named staging, commit, push and normal merge. No next
implementation starts before this merge.
