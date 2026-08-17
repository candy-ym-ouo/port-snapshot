//go:build darwin || freebsd || openbsd || netbsd || linux

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func collectPlatform(ctx context.Context) ([]PortEntry, []string, error) {
	if strings.HasPrefix(runtimeGOOS(), "linux") {
		return collectLinux(ctx)
	}
	cmd := exec.CommandContext(ctx, "lsof", "-nP", "-iTCP", "-sTCP:LISTEN", "-iUDP", "-F", "pcuPn")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("run lsof: %w (try installing lsof or running with sufficient privileges)", err)
	}
	return parseLsof(string(out)), nil, nil
}

// runtimeGOOS is a variable to keep platform selection testable without
// importing runtime in parser-focused tests.
var runtimeGOOS = func() string { return goos }

func parseLsof(data string) []PortEntry {
	var entries []PortEntry
	seen := make(map[string]struct{})
	var command, user, protocol string
	pid := 0
	for _, line := range strings.Split(data, "\n") {
		if line == "" {
			continue
		}
		switch line[0] {
		case 'p':
			pid, _ = strconv.Atoi(line[1:])
		case 'c':
			command = line[1:]
		case 'u':
			user = resolveUser(line[1:])
		case 'P':
			protocol = normalizeProtocol(line[1:])
		case 'n':
			port, ok := parseSocketName(line[1:])
			if ok && (protocol == "TCP" || protocol == "UDP") {
				entry := PortEntry{protocol, port, pid, command, user}
				key := fmt.Sprintf("%s/%d/%d", protocol, port, pid)
				if _, exists := seen[key]; !exists {
					entries = append(entries, entry)
					seen[key] = struct{}{}
				}
			}
		}
	}
	return entries
}

func resolveUser(value string) string {
	if uid, err := user.LookupId(value); err == nil {
		return uid.Username
	}
	return value
}

func parseSocketName(name string) (int, bool) {
	if strings.Contains(name, "->") {
		return 0, false
	}
	if i := strings.LastIndex(name, ":"); i >= 0 {
		portText := strings.TrimSpace(strings.TrimSuffix(name[i+1:], " (LISTEN)"))
		port, err := strconv.Atoi(portText)
		if err == nil {
			return port, true
		}
	}
	return 0, false
}

func collectLinux(ctx context.Context) ([]PortEntry, []string, error) {
	cmd := exec.CommandContext(ctx, "ss", "-ltnupH")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("run ss: %w (try running with sufficient privileges)", err)
	}
	return parseSS(string(out)), nil, nil
}

func parseSS(data string) []PortEntry {
	var entries []PortEntry
	s := bufio.NewScanner(strings.NewReader(data))
	for s.Scan() {
		f := strings.Fields(s.Text())
		if len(f) < 5 {
			continue
		}
		protocol := normalizeProtocol(f[0])
		addr := f[4]
		portText := addr[strings.LastIndex(addr, ":")+1:]
		port, err := strconv.Atoi(portText)
		if err != nil {
			continue
		}
		pid, process := parseSSProcess(s.Text())
		entries = append(entries, PortEntry{protocol, port, pid, process, userForPID(pid)})
	}
	return entries
}

var ssProcessPattern = regexp.MustCompile(`users:\(\(\"([^\"]+)\",pid=(\d+)`)

func parseSSProcess(line string) (int, string) {
	match := ssProcessPattern.FindStringSubmatch(line)
	if len(match) != 3 {
		return 0, ""
	}
	pid, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, match[1]
	}
	return pid, match[1]
}

func userForPID(pid int) string {
	if pid <= 0 {
		return ""
	}
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "status"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "Uid:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return ""
		}
		return resolveUser(fields[1])
	}
	return ""
}
