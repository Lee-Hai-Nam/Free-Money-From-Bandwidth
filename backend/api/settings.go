package api

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

type AppSettings struct {
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
		return &AppSettings{}, nil
	}
	var cfg AppSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &AppSettings{}, nil
	}
	return &cfg, nil
}

func (s *SettingsAPI) saveSettings(cfg *AppSettings) error {
	s.ensureDir()
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(s.settingsPath(), b, 0644)
}


