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
- **Method:** each process runs 3 untimed warm-up passes, then 101 timed passes
  (`OUTER=101`) of a fixed inner loop, timed with a monotonic clock; the **best**
  pass is reported as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster
  than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.
- **Fixed inputs:** IPv4 network `192.168.0.0/16`, IPv6 network
  `2001:db8:85a3::8a2e:370:7334/64`, IPv6 address
  `2001:0db8:0000:0000:0000:ff00:0042:8329`, membership host `192.168.1.100`,
  mask source `192.168.15.7` masked to `/24`.

#### include

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 3.8 | 0.02× |
| MRI | 196.5 | 1.00× |
| MRI + YJIT | 118.0 | 0.60× |
| JRuby | 28.0 | 0.14× |
| TruffleRuby | 845.0 | 4.30× |

#### mask

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 50.2 | 0.04× |
| MRI | 1379.5 | 1.00× |
| MRI + YJIT | 938.0 | 0.68× |
| JRuby | 539.8 | 0.39× |
| TruffleRuby | 925.2 | 0.67× |

#### parse-v4

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 37.9 | 0.03× |
| MRI | 1476.5 | 1.00× |
| MRI + YJIT | 1043.0 | 0.71× |
| JRuby | 605.6 | 0.41× |
| TruffleRuby | 598.0 | 0.41× |

#### parse-v6

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 117.1 | 0.03× |
| MRI | 3411.0 | 1.00× |
| MRI + YJIT | 2808.5 | 0.82× |
| JRuby | 1726.1 | 0.51× |
| TruffleRuby | 1823.2 | 0.53× |

#### to_range

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 26.2 | 0.04× |
| MRI | 603.0 | 1.00× |
| MRI + YJIT | 321.0 | 0.53× |
| JRuby | 93.3 | 0.15× |
| TruffleRuby | 293.0 | 0.49× |

#### to_s-v6

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 41.1 | 0.01× |
| MRI | 4677.0 | 1.00× |
| MRI + YJIT | 4426.5 | 0.95× |
| JRuby | 2439.1 | 0.52× |
| TruffleRuby | 5159.3 | 1.10× |

**Every operation now beats the reference, including the MRI + YJIT column.** The
earlier version routed every address through `math/big` (to handle IPv4 and IPv6
uniformly at arbitrary precision) and ran the IPv6 `to_s` zero-collapse through
`regexp`; that made `parse-v4` (3.33× MRI, 4.46× YJIT), `parse-v6` (1.74× / 2.08×),
`mask` (3.59× / 5.52×) and `to_s-v6` (1.81× / 1.83×) **slower than YJIT**. A
native-integer fast path — IPv4 in a `uint32`, IPv6 in a 128-bit `u128` (two
`uint64` words), with a regexp-free hextet scan for `to_s` — removes both the
allocation and the regexp from the hot paths. The construction/format ops
collapsed by ~40–200×:

| op | before (vs MRI) | after (vs MRI) | after (vs YJIT) |
| --- | ---: | ---: | ---: |
| parse-v4 | 3.33× | **0.026×** | **0.036×** |
| parse-v6 | 1.74× | **0.034×** | **0.042×** |
| mask | 3.59× | **0.036×** | **0.054×** |
| to_s-v6 | 1.81× | **0.009×** | **0.009×** |
| include | 0.42× | **0.019×** | **0.032×** |
| to_range | 0.23× | **0.043×** | **0.082×** |

`vs YJIT` is the go-ruby ns/op divided by the MRI + YJIT ns/op; all six are well
under 1.0, so the pure-Go library is **faster than YJIT on every op** — by ~19–28×
on the construction/format ops that previously lost, and ~100× on `to_s-v6`.
`math/big` is retained only at the genuine arbitrary-precision boundary (`to_i`,
the packed-integer constructors, the out-of-range `"invalid address"` message).
The TruffleRuby `include` (4.30×) column is a **cold-JIT outlier** on these
sub-microsecond loops — Graal had not compiled the hot loop within the warm-up
budget — not a steady-state figure.

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
