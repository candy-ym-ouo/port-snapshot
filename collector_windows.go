//go:build windows

package main

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func collectPlatform(ctx context.Context) ([]PortEntry, []string, error) {
	out, err := exec.CommandContext(ctx, "netstat", "-ano").Output()
	if err != nil {
		return nil, nil, fmt.Errorf("run netstat: %w", err)
	}
	var entries []PortEntry
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) < 4 || (f[0] != "TCP" && f[0] != "UDP") {
			continue
		}
		if f[0] == "TCP" && !strings.EqualFold(f[3], "LISTENING") {
			continue
		}
		addr := f[1]
		i := strings.LastIndex(addr, ":")
		if i < 0 {
			continue
		}
		port, e := strconv.Atoi(addr[i+1:])
		if e != nil {
			continue
		}
		pid, _ := strconv.Atoi(f[len(f)-1])
		entries = append(entries, PortEntry{f[0], port, pid, "", ""})
	}
	return entries, []string{"Windows netstat does not reliably expose process names and users without elevated APIs"}, nil
}
