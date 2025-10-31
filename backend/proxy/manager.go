package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ProxyEventCallback is called when proxy is added
type ProxyEventCallback func(*Proxy)

// Manager handles proxy management and validation
type Manager struct {
	proxies        map[string]*Proxy
	healthCheck    map[string]ProxyHealth
	mu             sync.RWMutex
	onProxyAdded   ProxyEventCallback
	onProxyRemoved ProxyEventCallback
}

// NewManager creates a new proxy manager
func NewManager() *Manager {
	return &Manager{
		proxies:     make(map[string]*Proxy),
		healthCheck: make(map[string]ProxyHealth),
	}
}

// SetOnProxyAdded sets the callback for when a proxy is added
func (m *Manager) SetOnProxyAdded(callback ProxyEventCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onProxyAdded = callback
}

// SetOnProxyRemoved sets the callback for when a proxy is removed
func (m *Manager) SetOnProxyRemoved(callback ProxyEventCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onProxyRemoved = callback
}

// AddProxy adds a new proxy to the manager
func (m *Manager) AddProxy(proxyStr string) (*Proxy, error) {
	m.mu.Lock()

	proxy, err := ParseProxy(proxyStr)
	if err != nil {
		m.mu.Unlock()
		return nil, fmt.Errorf("failed to parse proxy: %w", err)
	}

	m.proxies[proxy.ID] = proxy

	// Call callback if set
	onProxyAdded := m.onProxyAdded
	m.mu.Unlock()

	if onProxyAdded != nil {
		onProxyAdded(proxy)
	}

	return proxy, nil
}

// ParseProxy parses a proxy string into a Proxy struct
func ParseProxy(proxyStr string) (*Proxy, error) {
	// Parse proxy URL
	parsedURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	proxy := &Proxy{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Original: proxyStr,
		Protocol: parsedURL.Scheme,
		Host:     parsedURL.Hostname(),
		Port:     parsedURL.Port(),
	}

	// Extract credentials if present
	if parsedURL.User != nil {
		proxy.Username = parsedURL.User.Username()
		password, _ := parsedURL.User.Password()
		proxy.Password = password
	}

	return proxy, nil
}

// ValidateProxy validates proxy connectivity
func (m *Manager) ValidateProxy(proxyID string) error {
	m.mu.RLock()
	proxy, exists := m.proxies[proxyID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("proxy not found: %s", proxyID)
	}

	// Test connectivity
	err := m.TestConnectivity(proxy)
	if err != nil {
		m.mu.Lock()
		m.healthCheck[proxyID] = ProxyHealth{
			Status:    Unhealthy,
			LastCheck: time.Now(),
			Error:     err.Error(),
		}
		m.mu.Unlock()
		return err
	}

	m.mu.Lock()
	m.healthCheck[proxyID] = ProxyHealth{
		Status:    Healthy,
		LastCheck: time.Now(),
	}
	m.mu.Unlock()

	return nil
}

// TestConnectivity tests proxy connectivity
func (m *Manager) TestConnectivity(proxy *Proxy) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create HTTP client with proxy
	proxyURL := &url.URL{
		Scheme: proxy.Protocol,
		Host:   fmt.Sprintf("%s:%s", proxy.Host, proxy.Port),
	}

	if proxy.Username != "" {
		proxyURL.User = url.UserPassword(proxy.Username, proxy.Password)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 5 * time.Second,
	}

	// Test with a simple request
	req, err := http.NewRequestWithContext(ctx, "GET", "http://www.google.com", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("proxy test failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("proxy returned error: %d", resp.StatusCode)
	}

	return nil
}

// GetProxy returns a proxy by ID
func (m *Manager) GetProxy(proxyID string) (*Proxy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	proxy, exists := m.proxies[proxyID]
	if !exists {
		return nil, fmt.Errorf("proxy not found: %s", proxyID)
	}

	return proxy, nil
}

// ListProxies returns all proxies
func (m *Manager) ListProxies() []*Proxy {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Proxy, 0, len(m.proxies))
	for _, proxy := range m.proxies {
		result = append(result, proxy)
	}

	return result
}

// GetProxyHealth returns health status of a proxy
func (m *Manager) GetProxyHealth(proxyID string) (ProxyHealth, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	health, exists := m.healthCheck[proxyID]
	return health, exists
}

// RemoveProxy removes a proxy
func (m *Manager) RemoveProxy(proxyID string) error {
	m.mu.Lock()

	proxy, exists := m.proxies[proxyID]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("proxy not found: %s", proxyID)
	}

	delete(m.proxies, proxyID)
	delete(m.healthCheck, proxyID)

	// Call callback if set
	onProxyRemoved := m.onProxyRemoved
	m.mu.Unlock()

	if onProxyRemoved != nil {
		onProxyRemoved(proxy)
	}

	return nil
}

// Proxy represents a proxy configuration
type Proxy struct {
	ID       string
	Original string
	Protocol string
	Host     string
	Port     string
	Username string
	Password string
}

// ProxyHealth represents proxy health status
type ProxyHealth struct {
	Status    HealthStatus
	LastCheck time.Time
	Latency   time.Duration
	Error     string
}

// HealthStatus represents the health status of a proxy
type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Unknown   HealthStatus = "unknown"
)

// FormatProxy formats a proxy for use in HTTP requests
func (p *Proxy) FormatProxy() string {
	if p.Username != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%s", p.Protocol, p.Username, p.Password, p.Host, p.Port)
	}
	return fmt.Sprintf("%s://%s:%s", p.Protocol, p.Host, p.Port)
}

// ImportProxiesFromFile imports proxies from a text file
func (m *Manager) ImportProxiesFromFile(filePath string) ([]*Proxy, error) {
	// TODO: Read file and parse each line as a proxy
	// For now, return empty
	return []*Proxy{}, nil
}
