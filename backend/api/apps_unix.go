//go:build !windows

package api

import "os/exec"

// hideConsoleWindow is a no-op on Unix-like systems
func hideConsoleWindow(cmd *exec.Cmd) {}
