//go:build darwin || freebsd || openbsd || netbsd || linux

package main

import "testing"

func TestParseLsof(t *testing.T) {
	entries := parseLsof("p123\ncserver\nuroot\nPTCP\nn*:8080\nPUDP\nn*:5353\n")
	if len(entries) != 2 || entries[0].PID != 123 || entries[0].Port != 8080 || entries[0].Protocol != "TCP" || entries[1].Protocol != "UDP" {
		t.Fatalf("unexpected: %#v", entries)
	}
}

func TestParseSSProcess(t *testing.T) {
	pid, process := parseSSProcess(`tcp LISTEN 0 4096 *:8080 *:* users:(("server",pid=123,fd=3))`)
	if pid != 123 || process != "server" {
		t.Fatalf("unexpected: pid=%d process=%q", pid, process)
	}
}
