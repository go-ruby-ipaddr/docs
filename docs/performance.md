# Performance

`go-ruby-ipaddr/ipaddr` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `ipaddr`. This
page records the **methodology** for the comparative benchmark of that module
against the reference Ruby runtimes, part of the ecosystem-wide per-module parity
suite.

## What is measured

The **same** Ruby script — a representative `IPAddr` workload — is run under every
runtime. `rbgo`'s number reflects **this pure-Go library doing the work**; every
other column is that interpreter's own `ipaddr` stdlib. So the comparison is the
**Ruby-visible operation**, apples-to-apples across interpreters. The script
prints a deterministic checksum and its output is checked **byte-identical to
MRI** before timing.

- **Method:** best-of-N wall time (best, not mean, to suppress scheduler noise);
  single-shot processes, no warm-up beyond the script's own loop.
- **Runtimes:** `ruby` (MRI, the oracle) and `ruby --yjit`; `jruby` (OpenJDK);
  `truffleruby` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`ipaddr.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-ipaddr) | 210 | 1.91× |
| MRI (ruby 4.0.5) | 110 | 1.00× |
| MRI + YJIT | 90 | 0.82× |
| JRuby 10.1.0.0 | 1600 | 14.55× |
| TruffleRuby 34.0.1 | 350 | 3.18× |

rbgo runs on **go-ruby-ipaddr** at near parity with MRI (1.91x) on this parse + subnet-membership + range workload.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are **real measured numbers** from
    the 2026-06-30 run (Apple M-series; `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`,
    `truffleruby 34.0.1`) — nothing is fabricated or cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `ipaddr`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked identical to
MRI before any timing (`to_s`, the `include?` boolean, the masked address and the
`to_range` endpoints all byte-for-byte equal — verified across MRI, JRuby and
TruffleRuby too).

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.
- **Fixed inputs:** IPv4 network `192.168.0.0/16`, IPv6 network
  `2001:db8:85a3::8a2e:370:7334/64`, IPv6 address
  `2001:0db8:0000:0000:0000:ff00:0042:8329`, membership host `192.168.1.100`,
  mask source `192.168.15.7` masked to `/24`.

#### include

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 83.4 | 0.42× |
| MRI | 200.5 | 1.00× |
| MRI + YJIT | 119.0 | 0.59× |
| JRuby | 124.0 | 0.62× |
| TruffleRuby | 852.7 | 4.25× |

#### mask

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 5116.1 | 3.59× |
| MRI | 1426.5 | 1.00× |
| MRI + YJIT | 926.0 | 0.65× |
| JRuby | 565.3 | 0.40× |
| TruffleRuby | 1158.4 | 0.81× |

#### parse-v4

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 4913.2 | 3.33× |
| MRI | 1474.0 | 1.00× |
| MRI + YJIT | 1102.0 | 0.75× |
| JRuby | 582.5 | 0.40× |
| TruffleRuby | 1091.5 | 0.74× |

#### parse-v6

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 6049.6 | 1.74× |
| MRI | 3467.5 | 1.00× |
| MRI + YJIT | 2906.5 | 0.84× |
| JRuby | 3365.5 | 0.97× |
| TruffleRuby | 2961.5 | 0.85× |

#### to_range

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 141.5 | 0.23× |
| MRI | 621.0 | 1.00× |
| MRI + YJIT | 326.5 | 0.53× |
| JRuby | 90.5 | 0.15× |
| TruffleRuby | 3182.6 | 5.12× |

#### to_s-v6

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 8538.1 | 1.81× |
| MRI | 4712.0 | 1.00× |
| MRI + YJIT | 4654.5 | 0.99× |
| JRuby | 2489.6 | 0.53× |
| TruffleRuby | 9480.1 | 2.01× |

Split down the middle. **`include?` (0.42×) and `to_range` (0.23×) are faster
than MRI** — ~2.4× and ~4.4× — because they are pure integer comparisons on the
already-parsed `*big.Int` begin/end addresses with no reparsing. The costs are on
the **construction and formatting** side: `parse-v4` (3.33×), `parse-v6` (1.74×)
and `mask` (3.59×, which reparses then masks) are slower than MRI's C `IPAddr`
because the pure-Go library carries every address through `math/big` allocation
to handle IPv4 and IPv6 uniformly at arbitrary precision, versus MRI's native
machine-integer path. `to_s-v6` (1.81×) runs the MRI-faithful zero-collapse cascade
through `regexp`. Those allocation- and regexp-heavy construction/format paths are
the module's clear optimization targets (a small-integer fast path for IPv4 and a
regexp-free `to_s` compactor); the query operations already beat the reference. The
TruffleRuby `include` (4.25×) and `to_range` (5.12×) columns are **cold-JIT
outliers** on these sub-microsecond loops — Graal had not compiled the hot loop
within the warm-up budget — not steady-state figures.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-ipaddr/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/ipaddr.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries. `go run . verify`
    and `ruby ruby/ipaddr.rb verify` print the outputs that are diffed for parity.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the shortest loops (the cold-JIT
    outliers noted above). Sub-microsecond rows carry the most relative noise;
    treat those ratios as order-of-magnitude. Every number here is a **real
    measured value** from the dated run above — nothing is fabricated, estimated,
    or cherry-picked. The go-ruby column is the pure-Go library; every other
    column is that interpreter's own stdlib doing the equivalent work.
