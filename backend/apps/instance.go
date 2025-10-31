package apps

import (
	"fmt"
	"sync"
)

// AppInstance represents a container instance of an app (with or without proxy)
type AppInstance struct {
	InstanceID  string            // Unique instance ID
	AppID       string            // App identifier
	ProxyID     string            // Proxy ID (empty string for local)
	ContainerID string            // Docker container ID
	DeviceName  string            // Device name for this instance
	Credentials map[string]string // App credentials
	Status      string            // Running, Stopped, etc.
	ProxyURL    string            // Proxy URL if using proxy
	SDKNodeID   string            // SDK node ID (for EarnApp etc.)
}

// InstanceManager manages all app instances
type InstanceManager struct {
	instances map[string]*AppInstance // instanceID -> AppInstance
	appMap    map[string][]string     // appID -> []instanceID
	proxyMap  map[string][]string     // proxyID -> []instanceID
	mu        sync.RWMutex
}

// NewInstanceManager creates a new instance manager
func NewInstanceManager() *InstanceManager {
	return &InstanceManager{
		instances: make(map[string]*AppInstance),
		appMap:    make(map[string][]string),
		proxyMap:  make(map[string][]string),
	}
}

// AddInstance adds a new app instance
func (im *InstanceManager) AddInstance(instance *AppInstance) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	// Add instance
	im.instances[instance.InstanceID] = instance

	// Update app map
	im.appMap[instance.AppID] = append(im.appMap[instance.AppID], instance.InstanceID)

	// Update proxy map
	if instance.ProxyID != "" {
		im.proxyMap[instance.ProxyID] = append(im.proxyMap[instance.ProxyID], instance.InstanceID)
	}

	return nil
}

// GetInstance retrieves an instance by ID
func (im *InstanceManager) GetInstance(instanceID string) (*AppInstance, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	instance, exists := im.instances[instanceID]
	if !exists {
		return nil, fmt.Errorf("instance not found: %s", instanceID)
	}

	return instance, nil
}

// GetAppInstances returns all instances for an app
func (im *InstanceManager) GetAppInstances(appID string) []*AppInstance {
	im.mu.RLock()
	defer im.mu.RUnlock()

	instanceIDs, exists := im.appMap[appID]
	if !exists {
		return []*AppInstance{}
	}

	result := make([]*AppInstance, 0, len(instanceIDs))
	for _, instanceID := range instanceIDs {
		if instance, ok := im.instances[instanceID]; ok {
			result = append(result, instance)
		}
	}

	return result
}

// GetProxyInstances returns all instances using a proxy
func (im *InstanceManager) GetProxyInstances(proxyID string) []*AppInstance {
	im.mu.RLock()
	defer im.mu.RUnlock()

	instanceIDs, exists := im.proxyMap[proxyID]
	if !exists {
		return []*AppInstance{}
	}

	result := make([]*AppInstance, 0, len(instanceIDs))
	for _, instanceID := range instanceIDs {
		if instance, ok := im.instances[instanceID]; ok {
			result = append(result, instance)
		}
	}

	return result
}

// RemoveInstance removes an instance
func (im *InstanceManager) RemoveInstance(instanceID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	instance, exists := im.instances[instanceID]
	if !exists {
		return fmt.Errorf("instance not found: %s", instanceID)
	}

	// Remove from app map
	instanceIDs := im.appMap[instance.AppID]
	for i, id := range instanceIDs {
		if id == instanceID {
			im.appMap[instance.AppID] = append(instanceIDs[:i], instanceIDs[i+1:]...)
			break
		}
	}

	// Remove from proxy map
	if instance.ProxyID != "" {
		instanceIDs = im.proxyMap[instance.ProxyID]
		for i, id := range instanceIDs {
			if id == instanceID {
				im.proxyMap[instance.ProxyID] = append(instanceIDs[:i], instanceIDs[i+1:]...)
				break
			}
		}
	}

	// Remove instance
	delete(im.instances, instanceID)

	return nil
}

// GetAllInstances returns all instances
func (im *InstanceManager) GetAllInstances() []*AppInstance {
	im.mu.RLock()
	defer im.mu.RUnlock()

	result := make([]*AppInstance, 0, len(im.instances))
	for _, instance := range im.instances {
		result = append(result, instance)
	}

	return result
}

// UpdateInstanceStatus updates an instance's status
func (im *InstanceManager) UpdateInstanceStatus(instanceID, status string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	instance, exists := im.instances[instanceID]
	if !exists {
		return fmt.Errorf("instance not found: %s", instanceID)
	}

	instance.Status = status
	return nil
}

// UpdateInstanceContainerID updates an instance's container ID
func (im *InstanceManager) UpdateInstanceContainerID(instanceID, containerID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	instance, exists := im.instances[instanceID]
	if !exists {
		return fmt.Errorf("instance not found: %s", instanceID)
	}

	instance.ContainerID = containerID
	return nil
}
