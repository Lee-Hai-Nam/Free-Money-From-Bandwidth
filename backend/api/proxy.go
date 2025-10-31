package api

import (
	"bandwidth-income-manager/backend/apps"
	"bandwidth-income-manager/backend/config"
	"bandwidth-income-manager/backend/proxy"
	"context"
	"fmt"
)

// ProxyAPI provides the API for proxy management
type ProxyAPI struct {
	ctx             context.Context
	proxyManager    *proxy.Manager
	instanceManager *apps.InstanceManager
	credentialStore *config.CredentialStore
	appsAPI         *AppsAPI
}

// NewProxyAPI creates a new ProxyAPI
func NewProxyAPI(proxyManager *proxy.Manager, instanceManager *apps.InstanceManager, credentialStore *config.CredentialStore, appsAPI *AppsAPI) *ProxyAPI {
	return &ProxyAPI{
		proxyManager:    proxyManager,
		instanceManager: instanceManager,
		credentialStore: credentialStore,
		appsAPI:         appsAPI,
	}
}

// AddProxy adds a new proxy with optional auto-deployment
func (p *ProxyAPI) AddProxy(proxyStr string, autoDeploy bool, selectedAppIDs []string) (map[string]interface{}, error) {
	// Validate proxy format
	_, err := proxy.ParseProxy(proxyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy format: %w", err)
	}

	// Add proxy to manager
	addedProxy, err := p.proxyManager.AddProxy(proxyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to add proxy: %w", err)
	}

	// Test connectivity
	testErr := p.proxyManager.TestConnectivity(addedProxy)
	isHealthy := testErr == nil

	result := map[string]interface{}{
		"proxy_id":  addedProxy.ID,
		"proxy_url": addedProxy.FormatProxy(),
		"healthy":   isHealthy,
	}

	// Deploy to selected apps if provided, or auto-deploy to all if requested
	if isHealthy && (autoDeploy || len(selectedAppIDs) > 0) {
		deployedContainers, err := p.deployToSelectedApps(addedProxy.ID, addedProxy.FormatProxy(), selectedAppIDs)
		if err != nil {
			result["deployment_error"] = err.Error()
		} else {
			result["deployed_containers"] = deployedContainers
		}
	} else if !isHealthy {
		result["error"] = "Proxy validation failed. Auto-deployment skipped."
	}

	return result, nil
}

// RemoveProxy checks which containers use this proxy (returns info, doesn't remove)
func (p *ProxyAPI) RemoveProxy(proxyID string) (map[string]interface{}, error) {
	// Check if proxy exists
	prox, err := p.proxyManager.GetProxy(proxyID)
	if err != nil {
		return nil, fmt.Errorf("proxy not found: %s", proxyID)
	}

	// Get all instances using this proxy
	instances := p.instanceManager.GetProxyInstances(proxyID)

	containers := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		containers = append(containers, map[string]interface{}{
			"instance_id":  instance.InstanceID,
			"app_id":       instance.AppID,
			"container_id": instance.ContainerID,
			"device_name":  instance.DeviceName,
		})
	}

	return map[string]interface{}{
		"proxy_id":            proxyID,
		"proxy_url":           prox.FormatProxy(),
		"affected_containers": len(containers),
		"containers":          containers,
	}, nil
}

// ConfirmRemoveProxy actually removes the proxy and all its containers
func (p *ProxyAPI) ConfirmRemoveProxy(proxyID string) error {
	// Get all instances using this proxy
	instances := p.instanceManager.GetProxyInstances(proxyID)

	// Track proxy containers to remove
	proxyContainers := make(map[string]bool)

	// Stop and remove all containers
	for _, instance := range instances {
		if instance.ContainerID != "" {
			// Stop container
			if err := p.appsAPI.StopApp(instance.ContainerID); err != nil {
				// Log error but continue
				fmt.Printf("failed to stop container %s: %v\n", instance.ContainerID, err)
			}

			// Remove container
			if err := p.appsAPI.RemoveApp(instance.ContainerID); err != nil {
				// Log error but continue
				fmt.Printf("failed to remove container %s: %v\n", instance.ContainerID, err)
			}

			// Track proxy container (one per proxy, shared by all apps)
			proxyContainerName := fmt.Sprintf("tun2socks_proxy_%s", p.hashProxyID(proxyID))
			proxyContainers[proxyContainerName] = true
		}

		// Remove from instance manager
		p.instanceManager.RemoveInstance(instance.InstanceID)
	}

	// Remove proxy containers
	for proxyContainerName := range proxyContainers {
		if err := apps.RemoveProxyNetwork(proxyContainerName); err != nil {
			fmt.Printf("failed to remove proxy container %s: %v\n", proxyContainerName, err)
		}
	}

	// Remove proxy
	return p.proxyManager.RemoveProxy(proxyID)
}

func (p *ProxyAPI) hashProxyID(proxyID string) string {
	return apps.GetProxyHash(proxyID)
}

// ListProxies returns all proxies with usage stats
func (p *ProxyAPI) ListProxies() ([]map[string]interface{}, error) {
	proxies := p.proxyManager.ListProxies()
	result := make([]map[string]interface{}, 0, len(proxies))

	for _, prox := range proxies {
		instances := p.instanceManager.GetProxyInstances(prox.ID)

		proxyMap := map[string]interface{}{
			"id":               prox.ID,
			"url":              prox.FormatProxy(),
			"protocol":         prox.Protocol,
			"host":             prox.Host,
			"port":             prox.Port,
			"containers_count": len(instances),
		}

		result = append(result, proxyMap)
	}

	return result, nil
}

// TestProxy tests proxy connectivity
func (p *ProxyAPI) TestProxy(proxyID string) (map[string]interface{}, error) {
	err := p.proxyManager.ValidateProxy(proxyID)
	if err != nil {
		return map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}, nil
	}

	health, _ := p.proxyManager.GetProxyHealth(proxyID)
	return map[string]interface{}{
		"status":     "success",
		"health":     health.Status,
		"last_check": health.LastCheck,
		"latency":    health.Latency.String(),
	}, nil
}

// GetProxyContainers returns all containers using a proxy
func (p *ProxyAPI) GetProxyContainers(proxyID string) ([]map[string]interface{}, error) {
	instances := p.instanceManager.GetProxyInstances(proxyID)
	result := make([]map[string]interface{}, 0, len(instances))

	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"instance_id":  instance.InstanceID,
			"app_id":       instance.AppID,
			"container_id": instance.ContainerID,
			"device_name":  instance.DeviceName,
			"status":       instance.Status,
		})
	}

	return result, nil
}

// GetAppsRunningOnProxies returns app IDs that are already running on the specified proxies
func (p *ProxyAPI) GetAppsRunningOnProxies(proxyIDs []string) ([]string, error) {
	runningApps := make(map[string]bool)

	for _, proxyID := range proxyIDs {
		instances := p.instanceManager.GetProxyInstances(proxyID)
		for _, instance := range instances {
			if instance.Status == "running" {
				runningApps[instance.AppID] = true
			}
		}
	}

	result := make([]string, 0, len(runningApps))
	for appID := range runningApps {
		result = append(result, appID)
	}

	return result, nil
}

// deployToSelectedApps deploys containers for selected apps using the new proxy
func (p *ProxyAPI) deployToSelectedApps(proxyID, proxyURL string, selectedAppIDs []string) ([]map[string]interface{}, error) {
	deployedContainers := make([]map[string]interface{}, 0)

	// Use selected apps if provided, otherwise get all configured apps
	appIDs := selectedAppIDs
	if len(appIDs) == 0 {
		configuredAppIDs, err := p.credentialStore.GetAllConfiguredApps()
		if err != nil {
			return nil, fmt.Errorf("failed to get configured apps: %w", err)
		}
		appIDs = configuredAppIDs
	}

	// For each app, load credentials and deploy
	for _, appID := range appIDs {
		creds, err := p.credentialStore.LoadCredentials(appID)
		if err != nil {
			// Skip apps without credentials
			fmt.Printf("skipping app %s: credentials not found\n", appID)
			continue
		}

		// Deploy app with this proxy
		err = p.appsAPI.DeployAppWithProxyId(appID, creds.Credentials, proxyID)
		if err != nil {
			// Log error but continue
			fmt.Printf("failed to deploy app %s with proxy: %v\n", appID, err)
			continue
		}

		deployedContainers = append(deployedContainers, map[string]interface{}{
			"app_id": appID,
			"status": "deployed",
		})
	}

	return deployedContainers, nil
}

// GetConfiguredAppsForProxy returns all configured apps
func (p *ProxyAPI) GetConfiguredAppsForProxy() ([]map[string]interface{}, error) {
	return p.appsAPI.GetConfiguredApps()
}

// SetContext sets the context
func (p *ProxyAPI) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// OnStartup is called when Wails starts
func (p *ProxyAPI) OnStartup(ctx context.Context) {
	p.ctx = ctx
}
