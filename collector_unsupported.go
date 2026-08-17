//go:build !darwin && !freebsd && !openbsd && !netbsd && !linux && !windows

package main

import "context"

func collectPlatform(context.Context) ([]PortEntry, []string, error) {
	return nil, nil, errUnsupported
}
