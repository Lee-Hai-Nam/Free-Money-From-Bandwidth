//go:build windows

package api

import (
	"os/exec"
	"syscall"
)

// hideConsoleWindow hides the spawned console window on Windows
func hideConsoleWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
