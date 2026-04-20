# Contributing to flexgo

Welcome. This is the orientation doc for anyone working on flexgo —
yourself a few months from now, a colleague, an open-source
contributor. Everything you need to be productive in an hour.

If something here is unclear or out of date, fix it. The doc lives
in the repo so it doesn't drift.

---

## What flexgo is

A Go library for building flexible, responsive terminal UIs (TUIs).
It wraps the [Charm](https://charm.sh) ecosystem (BubbleTea,
Lipgloss) and exposes a tree-based layout API modelled after CSS
Flexbox. Public types are `Node`, `Direction`, `Justify`, `Align`,
`Spacing`. Recipes for common shapes (Dashboard, Form, Tabs, …) live
in the `flexgo/layouts` sub-package.

Module path: `github.com/spik3r/flexgo`.

---

## Repo layout

```
flexgo/
├── node.go               core Node struct + Direction/Justify/Align/Spacing
├── render.go             Render entrypoint + container/leaf rendering
├── layout.go             distribute() — main-axis size partitioning
├── render_helpers.go     join, padding, margin helpers
├── align.go              cross-axis alignment
├── builder.go            NodeBuilder fluent API
├── inspect.go            Validate(root) + Inspect(root) diagnostics
├── *_test.go             unit tests + golden test runner
├── testdata/             *.golden snapshots for example outputs
│
├── layouts/              recipe sub-package — one file per recipe
│
├── example/              runnable examples; each main.go also doubles
│   ├── basics/           as a golden-test fixture (FLEXGO_GOLDEN=1)
│   ├── builder/
│   ├── layouts/
│   └── …
│
├── demo/scanner/         larger reference app — see its own README
│
├── tapes/                vhs scripts for the README GIFs
├── docs/                 generated GIFs
│
├── todo.md               outstanding work, prioritised
├── CLAUDE.md             notes for AI assistants working here
└── CONTRIBUTING.md       this file
```

---

## Prerequisites

- **Go 1.25.5+** (the version pinned in `go.mod`). Check with `go version`.
- **git** for the obvious reasons.

Optional but useful:

- **[vhs](https://github.com/charmbracelet/vhs)** to regenerate the
  README GIFs. `brew install vhs` on macOS.
- **`gifsicle`** to shrink the regenerated GIFs.
  `brew install gifsicle`.

That's it — no Node, no Make, no Docker.

---

## Getting started

```bash
git clone git@github.com:spik3r/flexgo.git
cd flexgo
go test ./...
```

Tests should be green on a fresh clone. If they're not, file an
issue — that's a CI gap, not your problem.

### Run an example

Every directory under `example/` is a runnable BubbleTea program:

```bash
go run ./example/layouts/dashboard
go run ./example/basics/justify
go run ./example/dynamic
```

`q` or `ctrl+c` quits.

### Run the reference app

```bash
go run ./demo/scanner
```

This is the larger demo wiring three screens, tabs, modals, and a
centralised keymap. Read [`demo/scanner/README.md`](demo/scanner/README.md)
for the architectural tour.

---

## The CI checks (run them locally)

CI runs the same six commands on every push/PR. Run them locally
before pushing and you'll never be the one to break the build:

```bash
gofmt -l .                    # nothing should print
go mod tidy                   # then check `git diff` is empty
go vet ./...
go build ./...
go test -race -count=1 ./...
```

If any of those fail and the failure is not obvious, open an issue.
The CI workflow is at [`.github/workflows/ci.yml`](.github/workflows/ci.yml).

---

## Code conventions

The codebase is small and opinionated. A few rules to keep in mind:

### Comments

Default to writing none. Only add a comment when the *why* is
non-obvious — a hidden constraint, a subtle invariant, a workaround,
or behaviour that would surprise a reader.

Don't explain *what* the code does — well-named identifiers do that.
Don't reference the current task or PR ("added for the X flow") —
that belongs in the commit message.

Functions exposed publicly (capitalised names) get a doc comment
because they show up in pkg.go.dev — write the comment for someone
who doesn't know the codebase.

### Patterns to match

- **Struct literals over builders** for internal use.
  `&flexgo.Node{Dir: flexgo.Row, Children: …}` reads cleaner than
  the builder chain. The `NodeBuilder` exists for users who prefer
  the fluent style; it's not the canonical idiom.
- **Value receivers** for screen models in `demo/scanner/` —
  `Update` returns the new state, parent stores it. Keeps everything
  functional and snapshot-able.
- **Side effects via `tea.Cmd`**. Anything that reads a file, runs a
  process, or talks to the network goes through `tea.Cmd` so the
  Update loop never blocks.
- **Don't pre-abstract**. Three similar lines are better than a
  premature helper. Wait for the third repetition.
- **Don't add error handling for things that can't happen.** Internal
  code trusts internal invariants; only validate at boundaries.

### Naming

- Files are lower-case with underscores (`render_helpers.go`,
  `inspect_test.go`).
- Test files mirror the source name (`render.go` → `render_test.go`).
- Examples nest under their group: `example/<group>/<name>/main.go`.

---

## Adding a new layout recipe

This is the most common contribution. The standing checklist:

1. New file `layouts/<recipe>.go` with a function returning
   `*flexgo.Node`. Every dimension/style overridable on the returned
   node — no internal-only secrets.
2. Doc comment on the function with a `// Customize (top 3 overrides):`
   block. Other recipes show the format.
3. Test in `layouts/layouts_test.go` exercising the basic shape and
   any clamp/branch logic.
4. Runnable example at `example/layouts/<recipe>/main.go` with
   BubbleTea wrapping. Use the existing recipe examples as templates.
5. Equivalent built-via-`NodeBuilder` example at
   `example/builder/layouts/<recipe>/main.go`.
6. Register both example paths in `golden_test.go`'s
   `exampleGoldens` slice.
7. Generate goldens: `go test -run TestExampleGolden -update .`
8. Run the full test suite and the CI commands above.

---

## Adding a new example

Lighter-weight than a recipe:

1. New directory `example/<group>/<name>/` with a `main.go`.
2. The `main` function honours `FLEXGO_GOLDEN=1`: when set, render
   the tree at a fixed size and print to stdout, then exit. Every
   example does this — copy the pattern.
3. Add the path to `exampleGoldens` in `golden_test.go`.
4. Generate the golden: `go test -run TestExampleGolden -update .`
5. Verify with `go test -count=1 .`.

---

## Regenerating the README GIFs

Source `.tape` files in `tapes/`, output `.gif` files in `docs/`.

```bash
vhs tapes/demo-scanner.tape         # one tape
for tape in tapes/*.tape; do        # everything
  vhs "$tape"
done
gifsicle -O3 -k 64 docs/<file>.gif -o docs/<file>.gif   # shrink
```

Full guide: [`tapes/README.md`](tapes/README.md).

---

## Releases & versioning

flexgo follows [SemVer](https://semver.org). Tags are
`vMAJOR.MINOR.PATCH` — the `v` prefix is mandatory for Go modules.

### Where we are

- **`v0.x.y`** — current. Pre-1.0 means breaking changes are allowed
  between **minor** versions (e.g. `v0.1.0 → v0.2.0` may break API).
  Patch bumps (`v0.1.0 → v0.1.1`) are bug fixes only.
- **`v1.x.y`** — when we commit to API stability. After `v1.0.0`,
  breaking changes require a `v2` module path migration
  (`github.com/spik3r/flexgo/v2`), so we don't go there lightly.

### When to bump what (pre-1.0)

| Change | Bump |
|---|---|
| Bug fix, no API change | patch (`v0.1.0` → `v0.1.1`) |
| New feature, additive | minor (`v0.1.1` → `v0.2.0`) |
| Breaking API change | minor (`v0.2.0` → `v0.3.0`) |
| New layout recipe | minor |
| New example | none — examples aren't part of the public API |

### Cutting a release

1. Make sure `main` is green: all CI checks pass.
2. Update any version-mentioning docs (rare; usually nothing).
3. Tag the commit:
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0

   - <highlight>
   - <highlight>"
   git push origin main
   git push origin v0.2.0
   ```
4. The Go module proxy picks the tag up automatically; within a few
   minutes it's available for `go get github.com/spik3r/flexgo@v0.2.0`
   and indexed at `https://pkg.go.dev/github.com/spik3r/flexgo`.
5. Optional: cut a GitHub Release from the tag with the same
   highlights as the tag annotation. Helps people skim the changelog.

### What goes in the release notes

Three sections, in this order:

- **Breaking changes** — the headline. List every renamed/removed
  identifier and the migration. Skip the section if there are none.
- **New** — features and recipes added. One line each.
- **Fixed** — bug fixes. One line each.

Keep prose minimal. The audience is busy.

### Yanking a bad release

If `v0.2.0` ships broken:

1. Tag a fix as `v0.2.1` immediately.
2. Mark `v0.2.0` retracted in `go.mod`:
   ```
   retract v0.2.0  // <reason>
   ```
   Commit, push, and tag the next release. The proxy honours
   `retract` and steers users away from the bad version.

Avoid force-pushing or deleting tags — the Go proxy may have already
cached them, which causes confusing "version not found" errors for
users.

---

## Asking for help

Open an issue with:

- What you tried.
- What you expected.
- What happened (with the `Inspect(root)` output if it's a layout
  question — that prints the tree shape).
- Output of `go version`.

The maintainer's response time isn't always fast; PRs that include a
failing test are picked up first.

---

## Where the next work is

[`todo.md`](todo.md) is the prioritised backlog. The "Suggested
order of work" section at the bottom is the recommended pick-up
order. Items marked 🏛️ are large enough to want a design discussion
in an issue first; smaller items (🐞 / 🧹 / 🧪) are fair game to
just send a PR.
