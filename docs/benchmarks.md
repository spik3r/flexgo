# Render benchmark baseline

Run the public `Node.Render(w, h)` benchmarks from the repository root:

```sh
go test -run '^TestRenderBenchmarkFixtures$' .
go test -run '^$' -bench '^BenchmarkRender$' -benchmem -count=1 .
```

For a quick harness check, run one measured iteration per fixture:

```sh
go test -run '^$' -bench '^BenchmarkRender$' -benchtime=1x -benchmem -count=1 .
```

The one-iteration result confirms that each benchmark runs; it is too short
for a useful timing comparison.

The four fixed fixtures cover an 80×24 leaf, an 80×24 fixed/flex row, a
100×30 nested tree, and an 80×24 styled Unicode leaf. Each fixture is built
before timing. The timed operation includes layout, every leaf `View` callback,
Lip Gloss styling, and the final output string. It excludes tree construction,
output checks, terminal I/O, and data loading. Each callback uses fixed text
and performs no I/O.

The fixture test and each benchmark's setup check validate the tree, confirm
every leaf label, and check terminal width and height. Benchmarks retain their
output and report allocations. Compare `ns/op`, `B/op`, and `allocs/op` on the
same machine, Go version, and benchmark command before and after a change.
These numbers are a baseline for later Render work, not a performance target
or evidence of an improvement. They do not measure a live terminal, app state
updates, large trees, or every style and layout option.
