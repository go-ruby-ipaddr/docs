# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference IPAddr workload, mirrored 1:1 by benchmarks/go/main.go. The same
# fixed IPv4 and IPv6 inputs drive both sides; run with the single argument
# `verify` to print the checkable outputs instead of timing (used to confirm the
# pure-Go library is byte-identical to MRI before any number is recorded).
require "ipaddr"
require_relative "_harness"

V4_NET  = "192.168.0.0/16"
V6_NET  = "2001:db8:85a3::8a2e:370:7334/64"
V6_ADDR = "2001:0db8:0000:0000:0000:ff00:0042:8329"
HOST    = "192.168.1.100"
MASKSRC = "192.168.15.7"

v4  = IPAddr.new(V4_NET)
v6  = IPAddr.new(V6_NET)
a6  = IPAddr.new(V6_ADDR)
mem = IPAddr.new(HOST)

if ARGV[0] == "verify"
  r = v4.to_range
  puts "parse-v4=#{IPAddr.new(V4_NET).to_s}"
  puts "parse-v6=#{IPAddr.new(V6_NET).to_s}"
  puts "include=#{v4.include?(mem)}"
  puts "to_s-v6=#{a6.to_s}"
  puts "to_range=#{r.first.to_s}..#{r.last.to_s}"
  puts "mask=#{IPAddr.new(MASKSRC).mask(24).to_s}"
  puts "netmask=#{v4.netmask}"
  exit
end

bench("parse-v4",  2000) { IPAddr.new(V4_NET) }
bench("parse-v6",  2000) { IPAddr.new(V6_NET) }
bench("include",   2000) { v4.include?(mem) }
bench("to_s-v6",   2000) { a6.to_s }
bench("to_range",  2000) { r = v4.to_range; r.first; r.last }
bench("mask",      2000) { IPAddr.new(MASKSRC).mask(24) }
