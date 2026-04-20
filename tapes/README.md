# tapes/

Source-of-truth scripts for the README/demo GIFs. Each `.tape` file is
a [vhs](https://github.com/charmbracelet/vhs) script — the rendered
GIF lands in `docs/` with the same base name.

## Install vhs

```bash
brew install vhs            # macOS
go install github.com/charmbracelet/vhs@latest   # any platform
```

`vhs` runs the recorded shell in a headless `ttyd` and captures
frames; it needs `ttyd` and `ffmpeg` on your `$PATH`. Both are
installed automatically with the brew formula; for `go install`
follow the upstream README.

## Regenerate everything

```bash
for tape in tapes/*.tape; do
  vhs "$tape"
done
```

Or just the one you changed:

```bash
vhs tapes/demo-scanner.tape
```

GIFs land in `docs/`. Commit both the `.tape` and the new `.gif`.

## Optimising file size

`vhs` output is reasonable but not minimal. After regenerating, run:

```bash
gifsicle -O3 -k 64 docs/<file>.gif -o docs/<file>.gif
```

`-O3` is the strongest lossless pass; `-k 64` caps the palette at 64
colours which is usually invisible for terminal content. Expect a
30–50% size cut.

If a GIF crosses ~5 MB even after `gifsicle`, shorten the recording
(more `Sleep` delete, less waiting) or switch its `Output` to
`docs/<file>.webp` — GitHub renders WebP inline now and it's much
smaller for terminal content.

## Conventions

- **Theme**: every tape uses `Catppuccin Mocha` for the surrounding
  terminal so the Charm-y dark vibe stays consistent. The flexgo
  apps paint their own palette (Tokyo Night Storm in the demo), so
  the terminal theme only affects the brief prompt at the start.
- **Padding**: `Set Padding 20` — gives the GIF a breathing-room
  gutter so it doesn't hug the README image bounds.
- **Sizing**:
  - Marquee (`demo-scanner`): `1200×760`.
  - Recipe demos: `1000×540`.
- **Length**: 8–18 seconds. Past 20s people scroll. Tighten `Sleep`
  blocks before adding more screens.
- **Hide the prompt**: every tape starts with a `Hide` / `clear` /
  `Show` block so the recording opens to a clean screen instead of
  a flickering `$ `.

## Adding a new tape

1. Copy the closest existing tape under a new name.
2. Adjust `Output`, `Width`, `Height`, and the `Type "go run …"`
   line.
3. Walk through the demo with `Type` / `Sleep` / `Tab` etc.
4. Run `vhs tapes/<your-tape>.tape` and check `docs/<your-tape>.gif`.
5. If size matters, run `gifsicle` (above).
6. Add an `![…](docs/<your-tape>.gif)` reference in the README.
