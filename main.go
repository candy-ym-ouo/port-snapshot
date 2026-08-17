package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var version = "dev"

func main() {
	minPort, maxPort := flag.Int("min", 1, "minimum port (inclusive)"), flag.Int("max", 65535, "maximum port (inclusive)")
	output := flag.String("csv", "", "write results to this CSV file")
	flag.Parse()
	if *minPort < 1 || *maxPort > 65535 || *minPort > *maxPort {
		fmt.Fprintln(os.Stderr, "invalid port range: use 1<=min<=max<=65535")
		os.Exit(2)
	}

	entries, warnings, err := Collect(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	entries = FilterPorts(entries, *minPort, *maxPort)
	if *output != "" {
		if err := WriteCSV(*output, entries); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Println("PROTOCOL\tPORT\tPID\tPROCESS\tUSER")
		for _, e := range entries {
			fmt.Printf("%s\t%d\t%d\t%s\t%s\n", e.Protocol, e.Port, e.PID, e.Process, e.User)
		}
	}
	for _, warning := range warnings {
		fmt.Fprintln(os.Stderr, "warning:", warning)
	}
}

// PortEntry is one listening socket. PID, Process, or User may be zero/empty
// when the operating system denies access or the process exits during a scan.
type PortEntry struct {
	Protocol string
	Port     int
	PID      int
	Process  string
	User     string
}

func FilterPorts(entries []PortEntry, min, max int) []PortEntry {
	out := make([]PortEntry, 0, len(entries))
	for _, e := range entries {
		if e.Port >= min && e.Port <= max {
			out = append(out, e)
		}
	}
	return out
}

func WriteCSV(path string, entries []PortEntry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create CSV: %w", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.Write([]string{"protocol", "port", "pid", "process", "user"}); err != nil {
		return err
	}
	for _, e := range entries {
		if err := w.Write([]string{e.Protocol, strconv.Itoa(e.Port), strconv.Itoa(e.PID), e.Process, e.User}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

var errUnsupported = errors.New("port snapshot is unsupported on this platform")

// Collect obtains listening sockets and returns non-fatal per-process warnings.
func Collect(ctx context.Context) ([]PortEntry, []string, error) {
	return collectPlatform(ctx)
}

func normalizeProtocol(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
