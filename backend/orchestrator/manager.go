package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"bandwidth-income-manager/backend/docker"
)

// Manager manages multiple Docker hosts (local + remote)
type Manager struct {
	devices       map[string]*Device
	dockerClients map[string]*docker.Client
	mu            sync.RWMutex
}

// NewManager creates a new orchestrator manager
func NewManager() *Manager {
	return &Manager{
		devices:       make(map[string]*Device),
		dockerClients: make(map[string]*docker.Client),
	}
}

// AddDevice adds a new device to manage
func (m *Manager) AddDevice(device *Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Test connection
	client, err := docker.NewDockerClient(device.Host)
	if err != nil {
		return fmt.Errorf("failed to connect to device: %w", err)
	}

	if err := client.TestConnection(); err != nil {
		return fmt.Errorf("device connection test failed: %w", err)
	}

	m.devices[device.ID] = device
	m.dockerClients[device.ID] = client
	return nil
}

// RemoveDevice removes a device
func (m *Manager) RemoveDevice(deviceID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.devices, deviceID)
	delete(m.dockerClients, deviceID)
	return nil
}

// GetDevice returns a device by ID
func (m *Manager) GetDevice(deviceID string) (*Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, exists := m.devices[deviceID]
	if !exists {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	return device, nil
}

// ListDevices returns all devices
func (m *Manager) ListDevices() []*Device {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Device, 0, len(m.devices))
	for _, device := range m.devices {
		result = append(result, device)
	}

	return result
}

// DeployApp deploys an app to a specific device
func (m *Manager) DeployApp(deviceID string, appID string, config *AppDeploymentConfig) error {
	m.mu.RLock()
	client, exists := m.dockerClients[deviceID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	// TODO: Implement app deployment logic
	_ = appID
	_ = config
	_ = client

	return nil
}

// GetDeviceStatus gets the status of a device
func (m *Manager) GetDeviceStatus(deviceID string) (*DeviceStatus, error) {
	m.mu.RLock()
	client, exists := m.dockerClients[deviceID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	// Test connection
	err := client.TestConnection()
	status := DeviceStatus{
		ID:       deviceID,
		Online:   err == nil,
		LastSeen: time.Now(),
	}

	if err != nil {
		status.Error = err.Error()
	}

	return &status, nil
}

// Device represents a managed device
type Device struct {
	ID     string
	Name   string
	Host   string
	Auth   *AuthCredentials
	Status string
}

// DeviceStatus represents the status of a device
type DeviceStatus struct {
	ID       string
	Online   bool
	LastSeen time.Time
	Error    string
}

// AuthCredentials represents authentication credentials for a device
type AuthCredentials struct {
	Type     string // "none", "ssh", "tcp_tls"
	Username string
	Password string
	KeyPath  string
}

// AppDeploymentConfig represents configuration for deploying an app
type AppDeploymentConfig struct {
	AppID          string
	Config         map[string]string
	ProxyID        string
	ResourceLimits map[string]string
}
