# go-ruby-ipaddr documentation

**Ruby's `IPAddr` in pure Go — parsing, masking, predicates & operators, MRI byte-exact, no cgo.**

`go-ruby-ipaddr/ipaddr` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's [`IPAddr`](https://docs.ruby-lang.org/en/master/IPAddr.html),
matching reference Ruby (MRI). The module path is `github.com/go-ruby-ipaddr/ipaddr`.

It is the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other way:
this library has **no dependency on the Ruby runtime**.

!!! success "Status: complete"
    **Complete — MRI byte-exact.** Faithful port of Ruby's `IPAddr`: **parsing** (prefixlen / netmask / bracketed IPv6 / `%zone` / embedded IPv4), **formatting** (`to_s` / `to_string` / `cidr` / `inspect`), **masking**, **membership** (`include?` / `to_range`), the **bitwise / arithmetic** operators, **comparison**, the **predicates** and **conversions**, all in `math/big.Int`. Validated by a **differential oracle** against the system `ruby -ripaddr` — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
net, _ := ipaddr.New("192.168.1.5/24")
net.ToS()   // "192.168.1.0"   (masked, like MRI)
net.Cidr()  // "192.168.1.0/24"

in, _ := net.Include("192.168.1.99") // true

lo, hi, _ := net.ToRange()           // 192.168.1.0 .. 192.168.1.255
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`ipaddr`](https://github.com/go-ruby-ipaddr/ipaddr) | the library — Ruby's IPAddr in pure Go |
| [`docs`](https://github.com/go-ruby-ipaddr/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-ipaddr.github.io`](https://github.com/go-ruby-ipaddr/go-ruby-ipaddr.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-ipaddr/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby, validated by a differential
  oracle against the `ruby` binary.
- **Standalone & reusable.** A standalone module bound by `rbgo`; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-ipaddr/ipaddr](https://github.com/go-ruby-ipaddr/ipaddr).
