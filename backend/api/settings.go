package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type AppSettings struct {
	AutoStart  bool `json:"auto_start"`
	ShowInTray bool `json:"show_in_tray"`
}

type SettingsAPI struct {
	ctx     context.Context
	baseDir string
}

func NewSettingsAPI(baseDir string) *SettingsAPI {
	return &SettingsAPI{baseDir: baseDir}
}

func (s *SettingsAPI) OnStartup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SettingsAPI) settingsPath() string {
	return filepath.Join(s.baseDir, "data", "settings.json")
}

func (s *SettingsAPI) ensureDir() {
	_ = os.MkdirAll(filepath.Join(s.baseDir, "data"), 0755)
}

func (s *SettingsAPI) GetSettings() (*AppSettings, error) {
	s.ensureDir()
	path := s.settingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		// defaults
		return &AppSettings{AutoStart: false, ShowInTray: true}, nil
	}
	var cfg AppSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &AppSettings{AutoStart: false, ShowInTray: true}, nil
	}
	return &cfg, nil
}

func (s *SettingsAPI) saveSettings(cfg *AppSettings) error {
	s.ensureDir()
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(s.settingsPath(), b, 0644)
}

func (s *SettingsAPI) SetAutoStart(enabled bool) (bool, error) {
	cfg, _ := s.GetSettings()
	cfg.AutoStart = enabled
	if err := s.saveSettings(cfg); err != nil {
		return false, err
	}
	if runtime.GOOS == "windows" {
		if err := setWindowsAutoStart(enabled); err != nil {
			return false, fmt.Errorf("autostart failed: %w", err)
		}
	}
	return true, nil
}

func (s *SettingsAPI) SetShowInTray(enabled bool) (bool, error) {
	cfg, _ := s.GetSettings()
	cfg.ShowInTray = enabled
	if err := s.saveSettings(cfg); err != nil {
		return false, err
	}
	return true, nil
}
