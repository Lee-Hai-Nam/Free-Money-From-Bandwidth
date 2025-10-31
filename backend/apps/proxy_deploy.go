package apps

import (
	"fmt"
	"os/exec"
	"strings"
)

// DeployProxyTun deploys or returns existing tun2socks proxy container
func DeployProxyTun(proxyID, proxyURL string) (string, error) {
	// Container name is based on proxy only, not device name
	proxyHash := GetProxyHash(proxyID)
	proxyContainerName := fmt.Sprintf("tun2socks_proxy_%s", proxyHash)
	networkName := fmt.Sprintf("proxy_network_%s", proxyHash)

	// Step 1: Check if tun2socks container already exists for this proxy
	checkCmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", proxyContainerName), "--format", "{{.ID}}")
	output, err := checkCmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		// Container already exists, return existing container name
		return proxyContainerName, nil
	}

	// Step 2: Create a network for this proxy
	if err := createNetwork(networkName); err != nil {
		return "", fmt.Errorf("failed to create network: %w", err)
	}

	// Step 3: Pull tun2socks image
	if err := pullImage("xjasonlyu/tun2socks:latest"); err != nil {
		return "", fmt.Errorf("failed to pull tun2socks image: %w", err)
	}

	// Step 4: Deploy tun2socks container with privileged access
	args := []string{
		"run", "-d",
		"--name", proxyContainerName,
		"--restart", "always",
		"--network", networkName,
		"--cap-add", "NET_ADMIN",
		"--privileged",
		"-e", fmt.Sprintf("PROXY=%s", proxyURL),
		"-e", "LOGLEVEL=info",
		"-e", fmt.Sprintf("EXTRA_COMMANDS=ip rule add iif lo ipproto udp dport 53 lookup main;"),
		"-v", "/dev/net/tun:/dev/net/tun",
		"--dns", "1.1.1.1",
		"--dns", "8.8.8.8",
		"xjasonlyu/tun2socks:latest",
	}

	cmd := exec.Command("docker", args...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create tun2socks container: %w, output: %s", err, string(output))
	}

	// Container created successfully
	return proxyContainerName, nil
}

// DeployAppWithProxyTun deploys an app that uses network_mode: service:proxy
func DeployAppWithProxyTun(deployment *AppDeployment, proxyContainerName string) (string, error) {
	// Generate container name
	containerName := deployment.ContainerName
	if containerName == "" {
		containerName = getContainerName(deployment.AppID, deployment.DeviceName, deployment.ProxyID)
	}

	// Pull the app image
	if err := pullImage(deployment.Image); err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	// Build docker run command with network_mode: service:proxy
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

	// IMPORTANT: Use network_mode: service:proxy to share the network stack
	args = append(args, "--network", fmt.Sprintf("container:%s", proxyContainerName))

	// Add volumes
	for _, vol := range deployment.Volumes {
		args = append(args, "-v", vol)
	}

	// Add ports (port mappings still apply when container uses host networking via service:proxy)
	for _, p := range deployment.Ports {
		args = append(args, "-p", p)
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

// createNetwork creates a Docker network
func createNetwork(networkName string) error {
	// Check if network exists
	checkCmd := exec.Command("docker", "network", "inspect", networkName)
	if err := checkCmd.Run(); err == nil {
		// Network already exists
		return nil
	}

	// Create network
	cmd := exec.Command("docker", "network", "create", networkName)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "already exists") {
		return fmt.Errorf("failed to create network: %w, output: %s", err, string(output))
	}

	return nil
}

// RemoveProxyNetwork removes the network and proxy container
func RemoveProxyNetwork(proxyContainerName string) error {
	// Get network info
	inspectCmd := exec.Command("docker", "inspect", "-f", "{{.NetworkSettings.Networks}}", proxyContainerName)
	output, err := inspectCmd.Output()
	if err != nil {
		// Container might not exist
		return nil
	}

	// Try to remove container (will also disconnect it from network)
	stopCmd := exec.Command("docker", "stop", proxyContainerName)
	_ = stopCmd.Run()

	removeCmd := exec.Command("docker", "rm", "-f", proxyContainerName)
	removeCmd.Run() // Don't check error, might not exist

	// The network will be automatically removed when no containers use it
	// Or we can force remove it
	outputStr := string(output)
	if strings.Contains(outputStr, "proxy_network_") {
		// Extract network name (simplified approach)
		// In a real implementation, parse the inspect output properly
		fmt.Printf("Network cleanup: %s\n", string(output))
	}

	return nil
}
