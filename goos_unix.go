//go:build darwin || freebsd || openbsd || netbsd || linux

package main

import "runtime"

var goos = runtime.GOOS
