package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Loader manages loading and hot-reloading of app configurations
type Loader struct {
	configsDir string
	apps       map[string]*AppConfig
	mu         sync.RWMutex
}

// NewLoader creates a new config loader
func NewLoader(configsDir string) *Loader {
	return &Loader{
		configsDir: configsDir,
		apps:       make(map[string]*AppConfig),
	}
}

// LoadAppConfigs loads all app configurations from the configs directory
func (l *Loader) LoadAppConfigs() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	appsPath := filepath.Join(l.configsDir, "apps")
	files, err := filepath.Glob(filepath.Join(appsPath, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to glob config files: %w", err)
	}

	newApps := make(map[string]*AppConfig)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read config file %s: %w", file, err)
		}

		var config AppConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse config file %s: %w", file, err)
		}

		// Validate config against schema
		if err := l.ValidateConfig(&config); err != nil {
			return fmt.Errorf("invalid config file %s: %w", file, err)
		}

		newApps[config.AppID] = &config
	}

	l.apps = newApps
	return nil
}

// GetApps returns all loaded app configurations
func (l *Loader) GetApps() map[string]*AppConfig {
	l.mu.RLock()
	defer l.mu.RUnlock()

	result := make(map[string]*AppConfig)
	for id, config := range l.apps {
		result[id] = config
	}
	return result
}

// GetApp returns a specific app configuration by ID
func (l *Loader) GetApp(appID string) (*AppConfig, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	config, exists := l.apps[appID]
	if !exists {
		return nil, fmt.Errorf("app config not found: %s", appID)
	}

	return config, nil
}

// ValidateConfig validates an app configuration against the schema
func (l *Loader) ValidateConfig(config *AppConfig) error {
	if config.AppID == "" {
		return fmt.Errorf("app_id is required")
	}
	if config.Name == "" {
		return fmt.Errorf("name is required")
	}
	if config.DockerImage == "" {
		return fmt.Errorf("docker_image is required")
	}
	return nil
}

// WatchConfigs starts watching for config file changes (hot-reload)
func (l *Loader) WatchConfigs() error {
	// TODO: Implement file watcher for hot-reload
	// For now, just load once
	return l.LoadAppConfigs()
}

// AppConfig represents the structure of an app configuration
type AppConfig struct {
	AppID              string          `yaml:"app_id"`
	Name               string          `yaml:"name"`
	DockerImage        string          `yaml:"docker_image"`
	SupportedPlatforms []string        `yaml:"supported_platforms"`
	EnvironmentVars    []EnvVar        `yaml:"environment_vars"`
	Volumes            []Volume        `yaml:"volumes,omitempty"`
	NetworkMode        string          `yaml:"network_mode,omitempty"`
	ProxySupport       bool            `yaml:"proxy_support"`
	Ports              []string        `yaml:"ports,omitempty"`
	HealthCheck        *HealthCheck    `yaml:"health_check,omitempty"`
	EarningsAPI        *EarningsAPI    `yaml:"earnings_api,omitempty"`
	Description        string          `yaml:"description,omitempty"`
	ResourceLimits     *ResourceLimits `yaml:"resource_limits,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key         string `yaml:"key"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
}

// Volume represents a Docker volume mount
type Volume struct {
	Host      string `yaml:"host"`
	Container string `yaml:"container"`
	Type      string `yaml:"type,omitempty"` // bind, volume, tmpfs
}

// HealthCheck represents container health check configuration
type HealthCheck struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	Interval string `yaml:"interval,omitempty"`
}

// EarningsAPI represents earnings API configuration
type EarningsAPI struct {
	Endpoint   string            `yaml:"endpoint"`
	AuthMethod string            `yaml:"auth_method"`
	Headers    map[string]string `yaml:"headers,omitempty"`
}

// ResourceLimits represents container resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}
