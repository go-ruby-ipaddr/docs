<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-ipaddr` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-ipaddr/ipaddr`
library** against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby,
TruffleRuby). It measures the **library primitive** through its Go API, isolated
from the rbgo interpreter, so the numbers answer: *is the pure-Go implementation
as fast as the reference runtime's own `ipaddr`?*

## Layout

- `go/`            — self-contained Go driver; `go.mod` pins the **published**
  library by pseudo-version (no `replace`).
- `ruby/ipaddr.rb` — the equivalent workload; `ruby/_harness.rb` is the shared
  timer.
- `run.sh`         — runs every available runtime and prints one Markdown table
  per sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build **identical inputs** (the
same fixed IPv4/IPv6 addresses) and their outputs are checked identical to MRI
before timing — run `go run . verify` (in `go/`) and `ruby ruby/ipaddr.rb verify`
and diff the two. Results are published, dated, in `../docs/performance.md`.

## Operations

`parse-v4` / `parse-v6` (`IPAddr.new`), `include` (subnet membership test),
`to_s-v6` (compact IPv6 formatting), `to_range` (network → first/last address),
and `mask` (mask an address to a `/24`). Fixed inputs: `192.168.0.0/16`,
`2001:db8:85a3::8a2e:370:7334/64`, `2001:0db8:0000:0000:0000:ff00:0042:8329`,
membership host `192.168.1.100`, mask source `192.168.15.7`.
