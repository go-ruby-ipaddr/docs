# Roadmap

`go-ruby-ipaddr/ipaddr` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's IPAddr — the deterministic,
interpreter-independent slice — is **complete**.

| Stage | What | Status |
| --- | --- | --- |
| Parsing | `"a.b.c.d"`, `"addr/prefixlen"`, `"addr/netmask"`, bracketed IPv6 (`"[::1]/64"`), `%zone` identifiers and embedded IPv4 tails (`"::ffff:1.2.3.4"`), with MRI's ambiguous-zero-fill and out-of-range octet rejections. | **Done** |
| Formatting | `to_s` (compact, with MRI's exact zero-run collapsing and the `::a.b.c.d` / `::ffff:a.b.c.d` rewrites), `to_string` (canonical expanded), `cidr`, `inspect` (`#<IPAddr: IPv4:…/mask>`) and `netmask`. | **Done** |
| Masking & membership | `mask(prefixlen|netmask)`, `prefix` / `prefix=`, masked construction with non-contiguous-netmask rejection; `include?` / `===`, `to_range`, and an idiomatic `Each` over the range. | **Done** |
| Bitwise & arithmetic | `&`, `|`, `~`, `+`, `-`, `succ`, and `Xor` (the natural completion of the `&`/`|` set), all coercing strings / integers / `IPAddr`, carried in `math/big.Int` over the full 128-bit space. | **Done** |
| Comparison, predicates & conversions | `<=>` (Comparable), `==`, `Hash`; `ipv4?` / `ipv6?` / `loopback?` / `private?` / `link_local?` / `ipv4_mapped?` / `ipv4_compat?`; `native` / `ipv4_mapped` / `ipv4_compat` / `hton` / `ntop` / `new_ntoh` / `family` / `to_i`; the `InvalidAddressError` / `InvalidPrefixError` / `AddressFamilyError` family. | **Done** |
| Differential oracle & coverage | A wide corpus parsed and formatted here and compared to the system `ruby -ripaddr` — strings, predicates, operators, `to_range`, conversions, error class + message. 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Conformance targets MRI's behaviour;
  differences across MRI releases are matched to the reference used by the
  differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
