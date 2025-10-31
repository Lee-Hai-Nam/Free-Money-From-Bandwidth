package apps

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

// AppDeployment represents an app deployment configuration
type AppDeployment struct {
	AppID         string
	ProxyID       string
	ProxyURL      string
	DeviceName    string
	Image         string
	Environment   []string
	Volumes       []string
	Ports         []string
	Command       string
	RestartPolicy string
	NetworkMode   string
	ContainerName string
}

// DeployApp deploys an app using Docker CLI
func DeployApp(deployment *AppDeployment) (string, error) {
	// Generate container name if not provided
	containerName := deployment.ContainerName
	if containerName == "" {
		containerName = getContainerName(deployment.AppID, deployment.DeviceName, deployment.ProxyID)
	}

	// First, pull the image
	if err := pullImage(deployment.Image); err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	// Build docker run command
	args := []string{"run", "-d", "--name", containerName}

	// Add restart policy
	if deployment.RestartPolicy != "" {
		args = append(args, "--restart", deployment.RestartPolicy)
	} else {
		args = append(args, "--restart", "always")
	}

	// Add environment variables
	for _, env := range deployment.Environment {
		args = append(args, "-e", env)
	}

	// Add proxy environment variables if proxy is configured
	if deployment.ProxyURL != "" {
		args = append(args, "-e", fmt.Sprintf("HTTP_PROXY=%s", deployment.ProxyURL))
		args = append(args, "-e", fmt.Sprintf("HTTPS_PROXY=%s", deployment.ProxyURL))
		args = append(args, "-e", fmt.Sprintf("ALL_PROXY=%s", deployment.ProxyURL))
	}

	// Add volumes
	for _, vol := range deployment.Volumes {
		args = append(args, "-v", vol)
	}

	// Add ports
	for _, p := range deployment.Ports {
		args = append(args, "-p", p)
	}

	// Add network mode
	if deployment.NetworkMode != "" {
		args = append(args, "--network", deployment.NetworkMode)
	}

	// Add command
	if deployment.Command != "" {
		args = append(args, deployment.Image)
		args = append(args, strings.Split(deployment.Command, " ")...)
	} else {
		args = append(args, deployment.Image)
	}

	// Run docker command
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w, output: %s", err, string(output))
	}

	// Extract container ID from output
	containerID := strings.TrimSpace(string(output))

	return containerID, nil
}

// pullImage pulls a Docker image
func pullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull image %s: %w", image, err)
	}
	return nil
}

// getContainerName generates container name from app ID, device name, and proxy ID
func getContainerName(appID, deviceName, proxyID string) string {
	if proxyID == "" {
		// Local instance
		return fmt.Sprintf("%s_%s_local", deviceName, appID)
	}

	// Generate short hash for proxy
	hash := sha256.Sum256([]byte(proxyID))
	proxyHash := fmt.Sprintf("%x", hash[:4])[:8] // First 8 chars of hash

	return fmt.Sprintf("%s_%s_proxy%s", deviceName, appID, proxyHash)
}

// GetProxyHash wrapper to keep compatibility
func GetProxyHash(proxyID string) string {
	hash := sha256.Sum256([]byte(proxyID))
	return fmt.Sprintf("%x", hash[:4])[:8]
}
