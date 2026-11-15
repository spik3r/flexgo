# Reserve the painted width of custom border sides

Status: reviewed implementation claim; work has not started.
Owner: delegated Flexgo contributor, coordinated by Codex under the owner's
related-repository maintenance request.
Base: `d2698de4e213073eece4e27856efa2bdaf53b554` (verified remote main).
Live GitHub PRs and issues are empty. The earlier empty-border task is merged
in PR 14 and explicitly excluded nonempty border sizing. Preserve that record.
The primary checkout's uncommitted README and examples belong to another lane;
all implementation uses a separate worktree.

## Reproduction and scope

With NormalBorder's left and right sides changed to `界`, a valid node rendered
at 10 by 5 gives its bounded leaf an 8 by 3 callback and paints 12 columns.
One wide side and one ASCII side paint 11 columns. The frame fits when the leaf
gets 6 or 7 columns, respectively. This is an avoidable frame allocation error,
not content overflow or a change to the documented no-shrink layout policy.

Reserve the horizontal width used by pinned Lipgloss v2.0.3's border painter.
Use side-only border measurements and substitute a space only for an empty side.
A vertical pattern alternates glyphs; reserve its widest painted size rather
than summing the whole string. Do not count corner width as added side columns
or count a wide horizontal glyph as multiple physical rows.

Keep the existing two frame rows, ordinary and hidden borders, zero-border
allocation and Debug precedence, and debug-only sizing. Nonempty zero-width
side content is distinct from an empty side's space fallback. Add no clipping,
new layout phase, API, parser, dependency, version constant or release tag.
No Foundry or consumer pin changes are part of this task.

Exact paths: `render.go`, `render_edge_test.go`, and this task note. Follow the
existing plain task-note convention; add no task framework or Project Memory
store. Do not edit README, examples, golden files or the earlier task records.

## Acceptance and validation

Freeze public Render regressions and show failure on unchanged rendering code
before implementation. Use literal callback and output dimensions for wide and
asymmetric sides, alternating patterns, empty-side fallback, ordinary/hidden
borders, wide top/corner controls, zero-border and Debug precedence. Cover a
nested Row/Col layout and spacing with bounded leaf output. Tiny allocations
must keep callback sizes nonnegative, without promising a frame fits when its
own fixed width exceeds the allocation. Repeated renders must preserve styles.

Run focused regressions, gofmt, tidy with unchanged go.mod/go.sum, vet, build
and the full race suite with unchanged example goldens. Local Go 1.25.9 is
available; report that patch version honestly and require the hosted Go 1.25.5
CI job to pass before merge. Use owned temporary homes and the dedicated batch
cache; launch no interactive terminal, scanner or external service.

Freeze the final diff and observed checks for two independent reviews. The
coordinator owns staging, commit, publication and merge after exact-head checks.
No work outside the stated paths is authorized by this claim.
