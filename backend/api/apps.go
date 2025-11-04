package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"free-income-from-bandwidth/backend/apps"
	"free-income-from-bandwidth/backend/config"
	"free-income-from-bandwidth/backend/docker"
	"free-income-from-bandwidth/backend/monitor"
	"free-income-from-bandwidth/backend/proxy"
)

type wailsRuntime interface {
	WindowShow()
	WindowHide()
}

// AppsAPI provides the API for app management
type AppsAPI struct {
	ctx             context.Context
	docker          *docker.Client
	config          *config.Loader
	monitor         *monitor.Collector
	instanceManager *apps.InstanceManager
	credentialStore *config.CredentialStore
	proxyManager    *proxy.Manager
	startTime       time.Time
	recentActivity  []string
}

// NewAppsAPI creates a new AppsAPI
func NewAppsAPI(dockerClient *docker.Client, configLoader *config.Loader, monitorCollector *monitor.Collector, instanceManager *apps.InstanceManager, credentialStore *config.CredentialStore, proxyManager *proxy.Manager) *AppsAPI {
	return &AppsAPI{
		docker:          dockerClient,
		config:          configLoader,
		monitor:         monitorCollector,
		instanceManager: instanceManager,
		credentialStore: credentialStore,
		proxyManager:    proxyManager,
		startTime:       time.Now(),
		recentActivity:  make([]string, 0, 50),
	}
}

func (a *AppsAPI) addActivity(entry string) {
	// keep last 50 entries
	a.recentActivity = append(a.recentActivity, time.Now().Format(time.RFC3339)+" "+entry)
	if len(a.recentActivity) > 50 {
		a.recentActivity = a.recentActivity[len(a.recentActivity)-50:]
	}
}

	// GetDashboardSummary returns active containers, uptime and recent activity
func (a *AppsAPI) GetDashboardSummary() (map[string]interface{}, error) {
	containers, err := a.docker.ListContainersByLabel("com.free-income-from-bandwidth.managed=true")
	if err != nil {
		return nil, err
	}

	var runningIDs []string
	for _, c := range containers {
		if strings.ToLower(c.State) == "running" {
			runningIDs = append(runningIDs, c.ID)
				}
			}
		
			stats, _ := a.docker.GetContainersStats(runningIDs)
			startTimes, _ := a.docker.GetContainersStartTimes(runningIDs)
			totalCPU := 0.0
			totalRAM := int64(0)
			var oldestStart *time.Time
			for cid, stat := range stats {
				totalCPU += stat.CPUUsage
				totalRAM += stat.MemoryUsage
				if t, ok := startTimes[cid]; ok {
					if oldestStart == nil || t.Before(*oldestStart) {
						oldestStart = &t
					}
				}
			}
		
			uptimeSec := int64(0)
			if oldestStart != nil {
				uptimeSec = int64(time.Since(*oldestStart).Seconds())
			}
		
			summary := map[string]interface{}{
				"active_apps":     len(runningIDs),
				"cpu_usage":       totalCPU,
				"ram_usage":       totalRAM,
				"uptime_seconds":  uptimeSec,
				"recent_activity": a.recentActivity,
			}
	fmt.Printf("Dashboard Summary: %+v\n", summary)
	return summary, nil
}

// GetAvailableApps returns all available apps from config
func (a *AppsAPI) GetAvailableApps() (map[string]interface{}, error) {
	apps := a.config.GetApps()

	result := make(map[string]interface{})
	for id, appConfig := range apps {
		result[id] = map[string]interface{}{
			"app_id":           appConfig.AppID,
			"name":             appConfig.Name,
			"docker_image":     appConfig.DockerImage,
			"description":      appConfig.Description,
			"proxy_support":    appConfig.ProxySupport,
			"environment_vars": appConfig.EnvironmentVars,
		}
	}

	return result, nil
}

// GetRunningApps returns all running apps
func (a *AppsAPI) GetRunningApps() ([]map[string]interface{}, error) {
	containers, err := a.docker.ListContainers()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// Get SDK node IDs from instances
	instanceMap := make(map[string]string) // containerID -> sdkNodeID
	allInstances := a.instanceManager.GetAllInstances()
	for _, inst := range allInstances {
		if inst.SDKNodeID != "" {
			instanceMap[inst.ContainerID] = inst.SDKNodeID
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, container := range containers {
		// Filter only our app containers
		// TODO: Add filtering logic based on container name patterns
		containerData := map[string]interface{}{
			"id":     container.ID,
			"name":   container.Name,
			"image":  container.Image,
			"status": container.Status,
			"state":  container.State,
		}

		// Add published ports information
		if len(container.PublishedPorts) > 0 {
			containerData["published_ports"] = container.PublishedPorts
		}
		// Add detailed ports information
		if len(container.PortsDetail) > 0 {
			containerData["ports_detail"] = container.PortsDetail
		}

		// Get SDK node ID from instance manager or from container environment
		if sdkNodeID, exists := instanceMap[container.ID]; exists {
			containerData["sdkNodeID"] = sdkNodeID
		} else {
			// Try to get from container environment variables
			env, err := a.GetContainerEnvironmentVars(container.ID)
			if err == nil {
				if earnAppUUID, ok := env["EARNAPP_UUID"]; ok && earnAppUUID != "" {
					// Check if it has the prefix
					if strings.HasPrefix(earnAppUUID, "sdk-node-") {
						containerData["sdkNodeID"] = earnAppUUID
					}
				}
			}
		}

		result = append(result, containerData)
	}

	return result, nil
}

// StartApp starts an app by ID
func (a *AppsAPI) StartApp(appID string) error {
	err := a.docker.StartContainer(appID)
	if err == nil {
		a.addActivity("Started container " + appID)
	}
	return err
}

// StopApp stops an app by ID
func (a *AppsAPI) StopApp(appID string) error {
	err := a.docker.StopContainer(appID)
	if err == nil {
		a.addActivity("Stopped container " + appID)
	}
	return err
}

// RestartApp restarts an app by ID
func (a *AppsAPI) RestartApp(appID string) error {
	err := a.docker.RestartContainer(appID)
	if err == nil {
		a.addActivity("Restarted container " + appID)
	}
	return err
}

// GetAppLogs gets logs for an app
func (a *AppsAPI) GetAppLogs(appID string, tail int) (string, error) {
	return a.docker.GetContainerLogs(appID, tail)
}

// GetContainerLogs returns recent logs for a container by ID or name
func (a *AppsAPI) GetContainerLogs(containerID string) (string, error) {
	// Default to last 300 lines
	return a.docker.GetContainerLogs(containerID, 300)
}

// GetContainerLogsTail returns logs with a specified tail count
func (a *AppsAPI) GetContainerLogsTail(containerID string, tail int) (string, error) {
	if tail <= 0 {
		tail = 300
	}
	if tail > 10000 {
		tail = 10000
	}
	return a.docker.GetContainerLogs(containerID, tail)
}

// GetContainerLogsAll returns full logs for a container
func (a *AppsAPI) GetContainerLogsAll(containerID string) (string, error) {
	return a.docker.GetContainerLogsAll(containerID)
}

// GetAppStats gets statistics for an app
func (a *AppsAPI) GetAppStats(appID string) (map[string]interface{}, error) {
	stats, err := a.monitor.GetContainerStats(appID, 100)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"stats": stats,
	}

	return result, nil
}

// IsDockerAvailableResult represents the result of checking Docker availability
type IsDockerAvailableResult struct {
	Available bool   `json:"available"`
	Error     string `json:"error,omitempty"`
	Status    string `json:"status"` // "not_installed", "not_running", or "running"
}

// IsDockerAvailable checks if Docker is available and running
func (a *AppsAPI) IsDockerAvailable() (*IsDockerAvailableResult, error) {
	// First check if Docker is installed
	isInstalled := IsDockerInstalled()
	
	if !isInstalled {
		return &IsDockerAvailableResult{
			Available: false, 
			Error:     "Docker is not installed", 
			Status:    "not_installed",
		}, nil
	}
	
	// Docker is installed, now check if the daemon is running
	// If a.docker is nil, try to initialize it
	if a.docker == nil {
		// Try to create a new Docker client to test if daemon is running
		client, err := docker.NewDockerClient("")
		if err != nil {
			// Could not initialize Docker client - daemon likely not running
			return &IsDockerAvailableResult{
				Available: false, 
				Error:     "Docker daemon is not running", 
				Status:    "not_running",
			}, nil
		}
		
		// Test if the daemon is responsive
		testErr := client.TestConnection()
		if testErr != nil {
			return &IsDockerAvailableResult{
				Available: false, 
				Error:     "Docker daemon is not running", 
				Status:    "not_running",
			}, nil
		}
		
		// If we get here, Docker is running but wasn't initialized during app startup
		// For this check we just need to confirm it's available, so return success
		return &IsDockerAvailableResult{
			Available: true, 
			Error:     "", 
			Status:    "running",
		}, nil
	}
	
	// Docker client exists, test the connection
	err := a.docker.TestConnection()
	if err != nil {
		// Docker daemon is not responsive
		return &IsDockerAvailableResult{
			Available: false, 
			Error:     "Docker daemon is not running", 
			Status:    "not_running",
		}, nil
	}
	
	// Docker is running and available
	return &IsDockerAvailableResult{
		Available: true, 
		Error:     "", 
		Status:    "running",
	}, nil
}

// DeployApp deploys a new app with proper Docker configuration
func (a *AppsAPI) DeployApp(appID string, formData map[string]string) error {
	// Check if Docker is available before attempting deployment
	result, err := a.IsDockerAvailable()
	if err != nil {
		return fmt.Errorf("docker availability check failed: %w", err)
	}
	if !result.Available {
		if result.Status == "not_installed" {
			return fmt.Errorf("1. Please install Docker")
		} else if result.Status == "not_running" {
			return fmt.Errorf("1. Please start Docker")
		} else if result.Error != "" {
			// Check if the error contains specific Docker messages to reformat
			if strings.Contains(result.Error, "Docker daemon is not running") {
				return fmt.Errorf("1. Please start Docker")
			} else if strings.Contains(result.Error, "Docker is not installed") {
				return fmt.Errorf("1. Please install Docker")
			} else {
				return fmt.Errorf("docker is not available: %s", result.Error)
			}
		}
		return fmt.Errorf("1. Please install and start Docker")
	}
	
	// If Docker is available but a.docker is nil (not initialized during startup), 
	// try to initialize it now if needed for deployment
	if a.docker == nil {
		dockerClient, err := docker.NewDockerClient("")
		if err != nil {
			return fmt.Errorf("failed to initialize Docker client: %w", err)
		}
		a.docker = dockerClient
	}
	
	return a.DeployAppWithProxyId(appID, formData, "")
}

// DeployAppWithProxyId deploys an app with a specific proxy (or without if proxyID is empty)
func (a *AppsAPI) DeployAppWithProxyId(appID string, formData map[string]string, proxyID string) error {
	// Check if deploying local instance (proxyID empty) and if one already exists
	if proxyID == "" {
		instances := a.instanceManager.GetAppInstances(appID)
		// Check if there's already a local instance
		for _, instance := range instances {
			if instance.ProxyID == "" {
				return fmt.Errorf("a local instance already exists for this app. Only one local instance is allowed per app. Use proxy instances for additional connections")
			}
		}
	}

	// Get the app manifest
	manifest := apps.GetAppManifest(appID)
	if manifest == nil {
		return fmt.Errorf("app not found: %s", appID)
	}

	// Get device name from form data
	deviceName, ok := formData["DEVICE_NAME"]
	if !ok || deviceName == "" {
		return fmt.Errorf("DEVICE_NAME is required")
	}

	// Get proxy info if proxyID is provided
	var proxyURL string
	if proxyID != "" {
		prox, err := a.proxyManager.GetProxy(proxyID)
		if err != nil {
			return fmt.Errorf("proxy not found: %w", err)
		}
		proxyURL = prox.FormatProxy()
	}

	// Build environment variables
	env := []string{}

	// Add environment variables from manifest
	for key, value := range manifest.Environment {
		// Add all manifest environment variables (including those with placeholders)
		if value != "" {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Handle auto-generated fields (like EARNAPP_UUID)
	if manifest.AutoGenerateFields != nil {
		for fieldName, genConfig := range manifest.AutoGenerateFields {
			// Generate UUID
			uuid := generateUUID(genConfig.Length, genConfig.Charset)
			// Combine prefix with UUID
			fullValue := genConfig.Prefix + uuid
			env = append(env, fmt.Sprintf("%s=%s", fieldName, fullValue))

			// Store for later claim URL generation
			if fieldName == "EARNAPP_UUID" {
				formData["claimURL"] = fullValue
			}
		}
	}

	// Add regular environment variables from form data
	for key, value := range formData {
		// Skip fields that are handled elsewhere or shouldn't be added
		if key == "DEVICE_NAME" || key == "claimURL" || key == "EARNAPP_UUID" {
			continue
		}
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Also add device name to environment
	env = append(env, fmt.Sprintf("DEVICE_NAME=%s", deviceName))

	// Resolve port mappings: replace ${VAR} with user-specified or default to container port
	ports := []string{}
	if len(manifest.Ports) > 0 {
		for _, pm := range manifest.Ports {
			parts := strings.Split(pm, ":")
			if len(parts) == 2 {
				host := parts[0]
				container := parts[1]
				// If host is a placeholder like ${VAR}, use formData[VAR] or default to container port
				if strings.HasPrefix(host, "${") && strings.HasSuffix(host, "}") {
					varName := strings.TrimSuffix(strings.TrimPrefix(host, "${"), "}")
					val := formData[varName]
					if val == "" {
						val = container // default to same as container port
					}
					host = val
				}
				ports = append(ports, fmt.Sprintf("%s:%s", host, container))
			} else {
				ports = append(ports, pm)
			}
		}
	}

	// Create deployment config
	deployment := &apps.AppDeployment{
		AppID:         appID,
		ProxyID:       proxyID,
		ProxyURL:      proxyURL,
		DeviceName:    deviceName,
		Image:         manifest.Image,
		Environment:   env,
		Volumes:       manifest.Volumes,
		Ports:         ports,
		Command:       manifest.Command,
		RestartPolicy: "always",
		Labels:        map[string]string{"com.free-income-from-bandwidth.managed": "true"},
	}

	// Deploy the app
	var containerID string
	if proxyID != "" && proxyURL != "" {
		// Deploy using TUN proxy approach
		// Note: One tun2socks container per proxy, shared by all apps
		proxyContainerName, err := apps.DeployProxyTun(proxyID, proxyURL)
		if err != nil {
			return fmt.Errorf("failed to deploy proxy tun: %w", err)
		}

		// Ensure per-instance data volume directory for EarnApp (proxy instance)
		if appID == "earnapp" {
			// Set deterministic container name matching deploy helper
			deployment.ContainerName = fmt.Sprintf("%s_%s_proxy%s", deviceName, appID, apps.GetProxyHash(proxyID))
			// Rewrite volumes to use .data/<container>_earnapp
			newVolumes := make([]string, 0, len(deployment.Volumes))
			for _, v := range deployment.Volumes {
				if strings.Contains(v, ".data/.earnapp:") {
					newVolumes = append(newVolumes, strings.Replace(v, ".data/.earnapp:", ".data/"+deployment.ContainerName+"_earnapp:", 1))
				} else {
					newVolumes = append(newVolumes, v)
				}
			}
			deployment.Volumes = newVolumes
		}

		// Offset ports for proxy instances to avoid conflicts
		if len(deployment.Ports) > 0 {
			proxyOffset := int(apps.GetProxyHash(proxyID)[0]) % 50 // simple small offset
			newPorts := make([]string, 0, len(deployment.Ports))
			for _, pm := range deployment.Ports {
				// format "HOST:CONTAINER"
				parts := strings.Split(pm, ":")
				if len(parts) == 2 {
					host := parts[0]
					container := parts[1]
					// try parse int and offset
					if h, err := strconv.Atoi(host); err == nil {
						newPorts = append(newPorts, fmt.Sprintf("%d:%s", h+proxyOffset, container))
						continue
					}
				}
				newPorts = append(newPorts, pm)
			}
			deployment.Ports = newPorts
		}

		// Deploy app with network_mode: service:proxy
		containerID, err = apps.DeployAppWithProxyTun(deployment, proxyContainerName)
		if err != nil {
			return fmt.Errorf("failed to deploy app with proxy: %w", err)
		}
	} else {
		// Deploy without proxy
		var err error
		// Ensure per-instance data volume directory for EarnApp
		if appID == "earnapp" {
			// Generate container name to derive unique volume path (local instance)
			containerName := deployment.ContainerName
			if containerName == "" {
				containerName = fmt.Sprintf("%s_%s_local", deviceName, appID)
			}
			// Rewrite volumes to use .data/<container>_earnapp
			newVolumes := make([]string, 0, len(deployment.Volumes))
			for _, v := range deployment.Volumes {
				if strings.Contains(v, ".data/.earnapp:") {
					newVolumes = append(newVolumes, strings.Replace(v, ".data/.earnapp:", ".data/"+containerName+"_earnapp:", 1))
				} else {
					newVolumes = append(newVolumes, v)
				}
			}
			deployment.Volumes = newVolumes
		}
		containerID, err = apps.DeployApp(deployment)
		if err != nil {
			return err
		}
	}

	// Generate instance ID
	instanceID := fmt.Sprintf("%s_%s_%d", appID, deviceName, time.Now().Unix())

	// Extract SDK node ID for EarnApp (if present)
	sdkNodeID := ""
	if appID == "earnapp" && formData["claimURL"] != "" {
		// claimURL format is "sdk-node-{uuid}"
		claimURL := formData["claimURL"]
		// Extract UUID after "sdk-node-"
		if strings.HasPrefix(claimURL, "sdk-node-") {
			sdkNodeID = claimURL
		}
	}

	// Create instance record
	instance := &apps.AppInstance{
		InstanceID:  instanceID,
		AppID:       appID,
		ProxyID:     proxyID,
		ContainerID: containerID,
		DeviceName:  deviceName,
		Credentials: formData,
		Status:      "running",
		ProxyURL:    proxyURL,
		SDKNodeID:   sdkNodeID,
	}

	// Add instance to manager
	if err := a.instanceManager.AddInstance(instance); err != nil {
		return fmt.Errorf("failed to add instance: %w", err)
	}

	// Save credentials
	creds := &config.AppCredentials{
		AppID:       appID,
		DeviceName:  deviceName,
		Credentials: formData,
	}
	if err := a.credentialStore.SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	a.addActivity("Deployed app " + appID + " container " + containerID)
	return nil
}

// DeployAppWithProxy deploys an app with proxy configuration from API
func (a *AppsAPI) DeployAppWithProxy(deploymentData map[string]interface{}) (map[string]interface{}, error) {
	appID, _ := deploymentData["app_id"].(string)
	deviceName, _ := deploymentData["device_name"].(string)
	proxyID, _ := deploymentData["proxy_id"].(string)
	proxyURL, _ := deploymentData["proxy_url"].(string)
	credentials, _ := deploymentData["credentials"].(map[string]string)

	// Get manifest
	manifest := apps.GetAppManifest(appID)
	if manifest == nil {
		return nil, fmt.Errorf("app not found: %s", appID)
	}

	// Build environment
	env := []string{}
	for key, value := range credentials {
		if key != "DEVICE_NAME" && key != "claimURL" {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}
	env = append(env, fmt.Sprintf("DEVICE_NAME=%s", deviceName))

	// Create deployment
	deployment := &apps.AppDeployment{
		AppID:         appID,
		ProxyID:       proxyID,
		ProxyURL:      proxyURL,
		DeviceName:    deviceName,
		Image:         manifest.Image,
		Environment:   env,
		Volumes:       manifest.Volumes,
		Command:       manifest.Command,
		RestartPolicy: "always",
	}

	// Deploy
	containerID, err := apps.DeployApp(deployment)
	if err != nil {
		return nil, err
	}

	// Generate instance ID
	instanceID := fmt.Sprintf("%s_%s_%d", appID, deviceName, time.Now().Unix())

	// Add instance
	instance := &apps.AppInstance{
		InstanceID:  instanceID,
		AppID:       appID,
		ProxyID:     proxyID,
		ContainerID: containerID,
		DeviceName:  deviceName,
		Credentials: credentials,
		Status:      "running",
		ProxyURL:    proxyURL,
	}

	if err := a.instanceManager.AddInstance(instance); err != nil {
		return nil, fmt.Errorf("failed to add instance: %w", err)
	}

	return map[string]interface{}{
		"instance_id":  instanceID,
		"container_id": containerID,
		"device_name":  deviceName,
	}, nil
}

// DeployAppWithProxies deploys an app with multiple proxies
func (a *AppsAPI) DeployAppWithProxies(appID string, formData map[string]string, proxyIDs []string) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	usedPorts := map[int]bool{}
	// 1. Gather all ports used by running containers
	allContainers, _ := a.docker.ListContainers()
	for _, cont := range allContainers {
		for _, portstr := range cont.PublishedPorts {
			h, err := strconv.Atoi(portstr)
			if err == nil {
				usedPorts[h] = true
			}
		}
	}

	// 2. Track ports we assign in this batch
	assigned := map[int]bool{}
	getFreePort := func() int {
		for i := 0; i < 100; i++ {
			p := rand.Intn(50000) + 10000
			if !usedPorts[p] && !assigned[p] {
				assigned[p] = true
				return p
			}
		}
		return rand.Intn(50000) + 10000 // fallback (possible collision)
	}

	// Check if local instance already exists
	instances := a.instanceManager.GetAppInstances(appID)
	localExists := false
	for _, instance := range instances {
		if instance.ProxyID == "" {
			localExists = true
			break
		}
	}

	// Deploy local instance first (only if it doesn't exist)
	if !localExists {
		if err := a.DeployAppWithProxyId(appID, formData, ""); err != nil {
			return nil, fmt.Errorf("failed to deploy local instance: %w", err)
		}
		results = append(results, map[string]interface{}{
			"proxy_id": "",
			"status":   "deployed",
		})
	} else {
		results = append(results, map[string]interface{}{
			"proxy_id": "",
			"status":   "skipped (already exists)",
		})
	}

	// Deploy with each proxy
	for _, proxyID := range proxyIDs {
		manifest := apps.GetAppManifest(appID)
		if manifest == nil {
			continue
		}
		useFormData := map[string]string{}
		for k, v := range formData {
			useFormData[k] = v
		}
		// For each manifest.Ports entry, auto-assign free port if needed
		ports := []string{}
		for _, pm := range manifest.Ports {
			parts := strings.Split(pm, ":")
			if len(parts) == 2 {
				userHost := useFormData["HOSTPORT"]
				containerPort := parts[1]
				// If port blank or already assigned, pick random
				var hostPort int
				if userHost != "" {
					h, err := strconv.Atoi(userHost)
					if err != nil || usedPorts[h] || assigned[h] {
						hostPort = getFreePort()
					} else {
						hostPort = h
						assigned[h] = true
					}
				} else {
					hostPort = getFreePort()
				}
				ports = append(ports, fmt.Sprintf("%d:%s", hostPort, containerPort))
				usedPorts[hostPort] = true
			}
		}
		// Use assigned ports for this deployment
		useFormData["PORTS"] = strings.Join(ports, ",")
		if err := a.DeployAppWithProxyId(appID, useFormData, proxyID); err != nil {
			fmt.Printf("failed to deploy with proxy %s: %v\n", proxyID, err)
			results = append(results, map[string]interface{}{
				"proxy_id": proxyID,
				"status":   "error",
				"error":    err.Error(),
			})
			continue
		}
		results = append(results, map[string]interface{}{
			"proxy_id": proxyID,
			"status":   "deployed",
			"ports":    ports,
		})
	}
	return results, nil
}

// GetAppInstances returns all instances for an app
func (a *AppsAPI) GetAppInstances(appID string) ([]map[string]interface{}, error) {
	instances := a.instanceManager.GetAppInstances(appID)
	result := make([]map[string]interface{}, 0, len(instances))

	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"instance_id":  instance.InstanceID,
			"app_id":       instance.AppID,
			"proxy_id":     instance.ProxyID,
			"container_id": instance.ContainerID,
			"device_name":  instance.DeviceName,
			"status":       instance.Status,
			"proxy_url":    instance.ProxyURL,
		})
	}

	return result, nil
}

// RemoveAppInstance removes a specific app instance
func (a *AppsAPI) RemoveAppInstance(instanceID string) error {
	instance, err := a.instanceManager.GetInstance(instanceID)
	if err != nil {
		return err
	}

	// Stop and remove container
	if instance.ContainerID != "" {
		if err := a.docker.StopContainer(instance.ContainerID); err != nil {
			fmt.Printf("failed to stop container: %v\n", err)
		}
		// TODO: Remove container
	}

	// Remove from instance manager
	return a.instanceManager.RemoveInstance(instanceID)
}

// RemoveApp removes a container by container ID (for compatibility)
func (a *AppsAPI) RemoveApp(containerID string) error {
	// Stop container first
	if err := a.docker.StopContainer(containerID); err != nil {
		fmt.Printf("failed to stop container: %v\n", err)
	}
	// Attempt to backup EarnApp data directory if present
	// Determine container info to derive name
	cont, contErr := a.docker.GetContainer(containerID)
	if contErr == nil {
		// Expected volume path: .data/<container>_earnapp
		dataDir := filepath.Join(".data", cont.Name+"_earnapp")
		if _, statErr := os.Stat(dataDir); statErr == nil {
			// Create backup dir
			backupRoot := filepath.Join(".data", "backup")
			_ = os.MkdirAll(backupRoot, 0755)
			backupDir := filepath.Join(backupRoot, cont.Name+"_earnapp_"+time.Now().Format("20060102_150405"))
			// Rename (move) directory to backup location
			if err := os.Rename(dataDir, backupDir); err != nil {
				fmt.Printf("failed to backup data dir %s: %v\n", dataDir, err)
			}
		}
	}
	// Remove container
	err := a.docker.RemoveContainer(containerID)
	if err == nil {
		a.addActivity("Removed container " + containerID)
	}
	
	// Also remove the instance from the instance manager if it exists
	// Find instance by container ID and remove it
	allInstances := a.instanceManager.GetAllInstances()
	for _, instance := range allInstances {
		if instance.ContainerID == containerID {
			// Remove the instance from manager
			removeErr := a.instanceManager.RemoveInstance(instance.InstanceID)
			if removeErr != nil {
				fmt.Printf("failed to remove instance %s from manager: %v\n", instance.InstanceID, removeErr)
			}
			break // Found and attempted to remove, break the loop
		}
	}
	
	return err
}

// GetConfiguredApps returns all configured apps with credentials
func (a *AppsAPI) GetConfiguredApps() ([]map[string]interface{}, error) {
	appIDs, err := a.credentialStore.GetAllConfiguredApps()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(appIDs))
	for _, appID := range appIDs {
		creds, err := a.credentialStore.LoadCredentials(appID)
		if err != nil {
			continue
		}

		// Get app config for display name
		config, err := a.config.GetApp(appID)
		appName := appID
		if err == nil {
			appName = config.Name
		}

		result = append(result, map[string]interface{}{
			"app_id":      appID,
			"app_name":    appName,
			"device_name": creds.DeviceName,
		})
	}

	return result, nil
}

// GetAppCredentials returns saved credentials for an app (for editing)
func (a *AppsAPI) GetAppCredentials(appID string) (map[string]string, error) {
	creds, err := a.credentialStore.LoadCredentials(appID)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}
	// Return only the credentials map; DEVICE_NAME is part of credentials as well
	return creds.Credentials, nil
}

// DeployAppWithProxiesSelective deploys app with selected proxies
func (a *AppsAPI) DeployAppWithProxiesSelective(proxyID string, proxyURL string, appIDs []string) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	for _, appID := range appIDs {
		// Load credentials for this app
		creds, err := a.credentialStore.LoadCredentials(appID)
		if err != nil {
			fmt.Printf("failed to load credentials for app %s: %v\n", appID, err)
			continue
		}

		// Deploy app with this proxy
		err = a.DeployAppWithProxyId(appID, creds.Credentials, proxyID)
		if err != nil {
			fmt.Printf("failed to deploy app %s with proxy: %v\n", appID, err)
			continue
		}

		results = append(results, map[string]interface{}{
			"app_id": appID,
			"status": "deployed",
		})
	}

	return results, nil
}

// GetDockerInstallationURL returns the appropriate Docker installation URL based on the OS
func (a *AppsAPI) GetDockerInstallationURL() (string, error) {
	// Determine the OS and return appropriate Docker installation URL
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}

	switch goos {
	case "windows":
		return "https://www.docker.com/products/docker-desktop", nil
	case "darwin": // macOS
		return "https://www.docker.com/products/docker-desktop", nil
	case "linux":
		return "https://docs.docker.com/engine/install/", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", goos)
	}
}

// IsDockerInstalled checks if Docker is installed on the system by looking for common installation locations
func IsDockerInstalled() bool {
	switch runtime.GOOS {
	case "windows":
		// Check common Windows Docker installation locations
		possiblePaths := []string{
			"C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe",
			"C:\\Program Files\\Docker\\Docker.exe",
			"C:\\Docker\\Docker.exe",
		}
		
		        for _, path := range possiblePaths {
		            if _, err := os.Stat(path); err == nil {
		                return true
		            }
		        }
		        
		        // Also check if docker command is in PATH
		        cmd := exec.Command("docker", "--version")
		        hideConsoleWindow(cmd)
		        if err := cmd.Run(); err == nil {
		            return true
		        }		
	case "darwin": // macOS
		// Check common macOS Docker installation locations
		possiblePaths := []string{
			"/Applications/Docker.app",
			"/usr/local/bin/docker",
			"/opt/homebrew/bin/docker", // For Apple Silicon Macs
		}
		
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				// For Docker.app directory, just checking if the directory exists is sufficient
				if strings.HasSuffix(path, ".app") {
					return true
				}
				// For executable files, try to run version command
				if strings.Contains(path, "docker") {
					if err := exec.Command(path, "--version").Run(); err == nil {
						return true
					}
				}
			}
		}
		
		// Also check if docker command is in PATH
		if err := exec.Command("docker", "--version").Run(); err == nil {
			return true
		}
		
	case "linux":
		// Check common Linux Docker installation locations
		possiblePaths := []string{
			"/usr/bin/docker",
			"/usr/local/bin/docker",
		}
		
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				// For executable files, try to run version command
				if err := exec.Command(path, "--version").Run(); err == nil {
					return true
				}
			}
		}
		
		// Also check if docker command is in PATH
		if err := exec.Command("docker", "--version").Run(); err == nil {
			return true
		}
	}
	
	return false
}

// GetContainerEnvironmentVars gets environment variables from a container
func (a *AppsAPI) GetContainerEnvironmentVars(containerID string) (map[string]string, error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{json .Config.Env}}", containerID)
	// Hide flashing console on Windows
	hideConsoleWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	var envArray []string
	if err := json.Unmarshal(output, &envArray); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	env := make(map[string]string)
	for _, envVar := range envArray {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	return env, nil
}

// OnStartup is called when the Wails app starts
func (a *AppsAPI) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

// SetContext sets the context for the API
func (a *AppsAPI) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// OnStartupContext wrapper for Wails
func (a *AppsAPI) OnStartupContext(ctx context.Context) {
	a.ctx = ctx
}

// generateUUID generates a random UUID with specified length and charset
func generateUUID(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune(charset)
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// GetAllAppManifests returns all app manifests
func (a *AppsAPI) GetAllAppManifests() (map[string]*apps.AppManifest, error) {
	return apps.GetAllManifests(), nil
}
