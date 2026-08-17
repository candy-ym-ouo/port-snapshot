//go:build darwin || freebsd || openbsd || netbsd || linux

package main

import "testing"

func TestParseLsof(t *testing.T) {
	entries := parseLsof("p123\ncserver\nuroot\nnTCP *:8080 (LISTEN)\n")
	if len(entries) != 1 || entries[0].PID != 123 || entries[0].Port != 8080 || entries[0].Protocol != "TCP" {
		t.Fatalf("unexpected: %#v", entries)
	}
}
