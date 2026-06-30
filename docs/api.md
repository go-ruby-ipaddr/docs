# Usage & API

The public API lives at the module root (`github.com/go-ruby-ipaddr/ipaddr`). It is **Ruby-shaped but
Go-idiomatic**: the surface mirrors MRI's `IPAddr`, while following Go conventions —
an explicit `error` where Ruby raises, value types, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-ipaddr/ipaddr`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-ipaddr/ipaddr
```

## Worked example

```go
package main

import (
	"fmt"

	github.com/go-ruby-ipaddr/ipaddr
)

func main() {
	net, _ := ipaddr.New("192.168.1.5/24")
	net.ToS()   // "192.168.1.0"   (masked, like MRI)
	net.Cidr()  // "192.168.1.0/24"

	in, _ := net.Include("192.168.1.99") // true

	lo, hi, _ := net.ToRange()           // 192.168.1.0 .. 192.168.1.255
	_ = fmt.Sprint
}
```

## Shape

```go
// Construction.
func New(s string) (*IPAddr, error)                       // IPAddr.new(s)
func NewFamily(s string, family Family) (*IPAddr, error)   // IPAddr.new(s, family)
func NewFromInt(addr *big.Int, family Family) (*IPAddr, error)
func NewNtoh(addr []byte) (*IPAddr, error)                 // IPAddr.new_ntoh
func Ntop(addr []byte) (string, error)                     // IPAddr.ntop

// Strings.
func (ip *IPAddr) ToS() string        // to_s   (compact)
func (ip *IPAddr) ToString() string   // to_string (canonical)
func (ip *IPAddr) Cidr() string       // cidr
func (ip *IPAddr) Inspect() string    // inspect
func (ip *IPAddr) Netmask() string    // netmask

// Masking / prefix.
func (ip *IPAddr) Mask(prefixlen string) (*IPAddr, error) // mask("24"|"255.255.255.0")
func (ip *IPAddr) MaskLen(prefixlen int) (*IPAddr, error)
func (ip *IPAddr) Prefix() int
func (ip *IPAddr) SetPrefix(prefix int) (*IPAddr, error)  // prefix=

// Membership / range.
func (ip *IPAddr) Include(other any) (bool, error)        // include? / ===
func (ip *IPAddr) ToRange() (lo, hi *IPAddr, err error)   // to_range
func (ip *IPAddr) Each(fn func(*IPAddr) error) error

// Bitwise / arithmetic (other: *IPAddr | string | int | int64 | uint64 | *big.Int).
func (ip *IPAddr) And(other any) (*IPAddr, error)         // &
func (ip *IPAddr) Or(other any) (*IPAddr, error)          // |
func (ip *IPAddr) Xor(other any) (*IPAddr, error)         // ^ (extension)
func (ip *IPAddr) Not() (*IPAddr, error)                  // ~
func (ip *IPAddr) Add(offset int64) (*IPAddr, error)      // +
func (ip *IPAddr) Sub(offset int64) (*IPAddr, error)      // -
func (ip *IPAddr) Succ() (*IPAddr, error)                 // succ

// Comparison.
func (ip *IPAddr) Cmp(other any) (int, bool)              // <=>  (ok=false ~> nil)
func (ip *IPAddr) Eql(other any) bool                     // ==
func (ip *IPAddr) Hash() uint64

// Predicates.
func (ip *IPAddr) Ipv4() bool
func (ip *IPAddr) Ipv6() bool
func (ip *IPAddr) Loopback() bool
func (ip *IPAddr) Private() bool
func (ip *IPAddr) LinkLocal() bool
func (ip *IPAddr) Multicast() bool        // extension: MRI has no multicast?
func (ip *IPAddr) IsIpv4Mapped() bool     // ipv4_mapped?
func (ip *IPAddr) IsIpv4Compat() bool     // ipv4_compat?

// Conversions.
func (ip *IPAddr) Native() (*IPAddr, error)
func (ip *IPAddr) Ipv4Mapped() (*IPAddr, error)
func (ip *IPAddr) Ipv4Compat() (*IPAddr, error)
func (ip *IPAddr) HtonString() ([]byte, error)  // hton
func (ip *IPAddr) Family() Family
func (ip *IPAddr) ToI() *big.Int

// Errors (Error is the base; InvalidPrefixError is an InvalidAddressError in MRI).
type Error struct{ Msg string }
type InvalidAddressError struct{ Msg string }
type InvalidPrefixError struct{ Msg string }
type AddressFamilyError struct{ Msg string }
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
corpus through both the system `ruby` and this library and compares the results
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so the
cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-ipaddr/ipaddr` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
