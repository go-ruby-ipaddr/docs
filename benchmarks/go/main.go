// SPDX-License-Identifier: BSD-3-Clause
//
// Pure-Go IPAddr workload, mirroring benchmarks/ruby/ipaddr.rb 1:1 through the
// go-ruby-ipaddr/ipaddr Go API. Same fixed IPv4 and IPv6 inputs, same iteration
// counts. Run with the single argument `verify` to print the checkable outputs
// (compared byte-for-byte against MRI before any timing).
package main

import (
	"fmt"
	"os"

	"github.com/go-ruby-ipaddr/ipaddr"
)

const (
	v4Net   = "192.168.0.0/16"
	v6Net   = "2001:db8:85a3::8a2e:370:7334/64"
	v6Addr  = "2001:0db8:0000:0000:0000:ff00:0042:8329"
	host    = "192.168.1.100"
	maskSrc = "192.168.15.7"
)

func must(ip *ipaddr.IPAddr, err error) *ipaddr.IPAddr {
	if err != nil {
		panic(err)
	}
	return ip
}

// mask24 masks to a /24 prefix, mirroring Ruby's IPAddr#mask(24).
func mask24(ip *ipaddr.IPAddr) *ipaddr.IPAddr { return must(ip.Mask("24")) }

func main() {
	v4 := must(ipaddr.New(v4Net))
	a6 := must(ipaddr.New(v6Addr))
	mem := must(ipaddr.New(host))

	if len(os.Args) > 1 && os.Args[1] == "verify" {
		lo, hi, err := v4.ToRange()
		if err != nil {
			panic(err)
		}
		inc, err := v4.Include(mem)
		if err != nil {
			panic(err)
		}
		fmt.Printf("parse-v4=%s\n", must(ipaddr.New(v4Net)).ToS())
		fmt.Printf("parse-v6=%s\n", must(ipaddr.New(v6Net)).ToS())
		fmt.Printf("include=%t\n", inc)
		fmt.Printf("to_s-v6=%s\n", a6.ToS())
		fmt.Printf("to_range=%s..%s\n", lo.ToS(), hi.ToS())
		fmt.Printf("mask=%s\n", mask24(must(ipaddr.New(maskSrc))).ToS())
		fmt.Printf("netmask=%s\n", v4.Netmask())
		return
	}

	bench("parse-v4", 2000, func() { sink = must(ipaddr.New(v4Net)) })
	bench("parse-v6", 2000, func() { sink = must(ipaddr.New(v6Net)) })
	bench("include", 2000, func() {
		inc, err := v4.Include(mem)
		if err != nil {
			panic(err)
		}
		sink = inc
	})
	bench("to_s-v6", 2000, func() { sink = a6.ToS() })
	bench("to_range", 2000, func() {
		lo, hi, err := v4.ToRange()
		if err != nil {
			panic(err)
		}
		sink = lo
		sink = hi
	})
	bench("mask", 2000, func() { sink = mask24(must(ipaddr.New(maskSrc))) })
}
