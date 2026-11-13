# Keep tab labels within the header strip

Status: local implementation for independent review

Base: verified remote main 915fb99631b0fdde7cc63c2bcbfa44d5e8160f61.
Branch/worktree: `codex/bounded-tab-headers`, `/tmp/flexgo-bounded-tab-headers`.
Live PR inspection found no open Flexgo PRs. Existing worktrees and primary
checkout edits were preserved. No key-hint or cycle-diagnostic work is included.

## Scope

Long or multiline tab labels previously wrapped beyond their one-row title
slot, and the underline could extend beyond its column. A 12x6 regression
example rendered 31x13 before this change, displacing the active panel.
Header callbacks now truncate by terminal cell width with an ellipsis, retain
only the first title line and bound underlines without adding an ellipsis.
Zero/nonpositive callback allocations return empty output. The two-row strip,
node shape, customization fields and normal-label goldens remain unchanged.
No general overflow policy or layout-engine refactor is introduced.

## Validation

Go 1.27.0 with GOPROXY=off, GOSUMDB=off, GOTOOLCHAIN=local and
GOCACHE=/tmp/codex-local-build-cache:

- Before implementation, focused tests reproduced overflow for ASCII,
  multiline, wide Unicode and styled labels, plus zero-allocation callbacks.
- `go test -mod=readonly ./layouts . -count=1`: passed, including root example
  goldens without updates.
- `go test -mod=readonly -race ./layouts . -count=1`: passed after the final
  title/underline assertions. Cases include 1–40 column allocations, exact
  visible truncation, CR/LF labels, panel position and callback dimensions.
- Formatting and diff whitespace checks passed.

No dependency change, listener, module-cache write request or go mod tidy was
used. Separate vet/build and interactive terminal checks are not claimed;
those remain for hosted CI and review. Panels must still bound their own
output, and a viewport shorter than the fixed two-row strip is not made into
a scrolling layout. No publication occurred at this review checkpoint.

## Independent review correction

Independent review found that first-line truncation could discard SGR resets
or OSC-8 hyperlink terminators on later lines. Multiline ANSI-bearing labels
now close both style and hyperlink state after their visible ellipsis, before
padding, neighboring tabs or panel output. Plain labels gain no escape bytes;
single-line labels are unchanged. Raw-escape tests cover SGR, OSC-8 with BEL
and ST terminators, combined styling/links and LF/CRLF/CR separators at full
and narrow widths. They assert exact closing bytes, including placement before
other content rather than relying only on stripped display width.

After this correction, `go test -mod=readonly -race ./layouts . -count=1`
passed again, including existing example goldens. Formatting and diff checks
passed. The earlier separate vet/build and interactive-terminal limitations
remain. This correction is local and requires independent re-review.

## Hosted module-consistency correction

PR #10 CI run 43422932546 passed formatting but refused module consistency:
this recipe directly imports the already pinned x/ansi v0.11.7, so its require
entry must be direct. Applied the exact two-line go.mod move printed by hosted
tidy. No version, go.sum byte, production source or test changed in this
correction. No local go mod tidy or denied metadata-cache write was retried.
The earlier source review remains tied to 667e63a; fresh CI on the corrected
head is required before merge.
