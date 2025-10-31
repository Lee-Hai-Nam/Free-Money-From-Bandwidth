package api

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func setWindowsAutoStart(enabled bool) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\\Microsoft\\Windows\\CurrentVersion\\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	appName := "BandwidthIncomeManager"
	if !enabled {
		// delete value if exists
		_ = k.DeleteValue(appName)
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get exe: %w", err)
	}
	// ensure absolute and quoted path
	exePath, _ = filepath.Abs(exePath)
	return k.SetStringValue(appName, fmt.Sprintf("\"%s\"", exePath))
}
