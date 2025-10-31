//go:build !windows

package docker

import "os/exec"

// hideConsoleWindow is a no-op on Unix-like systems
func hideConsoleWindow(cmd *exec.Cmd) {}
