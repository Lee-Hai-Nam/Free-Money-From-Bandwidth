package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Client manages Docker via CLI commands
type Client struct {
	dockerCmd string
	ctx       context.Context
}

// NewDockerClient creates a new Docker client instance
func NewDockerClient(host string) (*Client, error) {
	dockerCmd := "docker"
	if host != "" {
		dockerCmd = fmt.Sprintf("docker -H %s", host)
	}

	return &Client{
		dockerCmd: dockerCmd,
		ctx:       context.Background(),
	}, nil
}

// TestConnection tests the Docker daemon connection
func (c *Client) TestConnection() error {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	args := c.parseCommand("ps")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	return nil
}

// CreateContainer creates a new container
func (c *Client) CreateContainer(config *ContainerConfig) error {
	args := c.parseCommand("run", "-d", "--name", config.Name, config.Image)

	// Add environment variables
	for _, env := range config.Env {
		args = append(args, "-e", env)
	}

	// Add volumes
	for key := range config.Volumes {
		args = append(args, "-v", key)
	}

	// Network mode
	if config.NetworkMode != "" {
		args = append(args, "--network", config.NetworkMode)
	}

	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)
	return cmd.Run()
}

// GetContainer gets container by name or ID
func (c *Client) GetContainer(name string) (*ContainerInfo, error) {
	args := c.parseCommand("ps", "-a", "--filter", fmt.Sprintf("name=^%s$", name), "--format", "json")
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	containers := []ContainerInfo{}
	if len(output) > 0 {
		// Split by lines if multiple containers
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			var cont ContainerInfo
			if err := json.Unmarshal([]byte(line), &cont); err == nil {
				// Parse Names field to set Name
				names := strings.TrimPrefix(cont.Names, "/")
				cont.Name = names
				containers = append(containers, cont)
			}
		}
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("container not found: %s", name)
	}

	return &containers[0], nil
}

// StartContainer starts a container
func (c *Client) StartContainer(name string) error {
	args := c.parseCommand("start", name)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)
	return cmd.Run()
}

// StopContainer stops a container
func (c *Client) StopContainer(name string) error {
	args := c.parseCommand("stop", name)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)
	return cmd.Run()
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(name string) error {
	args := c.parseCommand("rm", "-f", name)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)
	return cmd.Run()
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(name string) error {
	args := c.parseCommand("restart", name)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)
	return cmd.Run()
}

func (c *Client) ListContainers() ([]ContainerInfo, error) {
	return c.listContainersWithFilter("")
}

// ListContainersByLabel lists all containers with a given label
func (c *Client) ListContainersByLabel(label string) ([]ContainerInfo, error) {
	return c.listContainersWithFilter(fmt.Sprintf("label=%s", label))
}

// listContainersWithFilter is a helper to list containers with an optional filter
func (c *Client) listContainersWithFilter(filter string) ([]ContainerInfo, error) {
	args := c.parseCommand("ps", "-a", "--format", "{{json .}}")
	if filter != "" {
		args = append(args, "--filter", filter)
	}
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var containers []ContainerInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var cont ContainerInfo
		if err := json.Unmarshal([]byte(line), &cont); err == nil {
			names := strings.TrimPrefix(cont.Names, "/")
			cont.Name = names

			// Parse Ports field for published ports of format "0.0.0.0:5902->5900/tcp, ..."
			cont.PublishedPorts = []string{}
			cont.PortsDetail = []PortMapping{}
			portsStr := cont.Ports
			for _, pfrag := range strings.Split(portsStr, ",") {
				pfrag = strings.TrimSpace(pfrag)
				if pfrag == "" {
					continue
				}
				if strings.Contains(pfrag, ":") && strings.Contains(pfrag, "->") {
					// Extract hostPort:containerPort format
					parts := strings.Split(pfrag, "->")
					if len(parts) == 2 {
						hostSide := parts[0]
						containerSide := parts[1] // e.g., "5900/tcp"
						
						hostPort := hostSide[strings.LastIndex(hostSide, ":")+1:]
						containerPort := containerSide
						if strings.Contains(containerSide, "/") {
							containerPort = containerSide[:strings.Index(containerSide, "/")]
						}
						
						portMapping := PortMapping{
							HostPort:      hostPort,
							ContainerPort: containerPort,
						}
						cont.PortsDetail = append(cont.PortsDetail, portMapping)
						
						// Still add for backward compatibility
						cont.PublishedPorts = append(cont.PublishedPorts, fmt.Sprintf("%s:%s", hostPort, containerPort))
					}
				}
			}
			containers = append(containers, cont)
		}
	}

	return containers, nil
}

// GetContainerLogs gets container logs
func (c *Client) GetContainerLogs(name string, tail int) (string, error) {
	args := c.parseCommand("logs", "--tail", fmt.Sprintf("%d", tail), name)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GetContainerLogsAll gets full container logs without tail limit
func (c *Client) GetContainerLogsAll(name string) (string, error) {
	args := c.parseCommand("logs", name)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// parseCommand parses the command string into executable and arguments
func (c *Client) parseCommand(cmd string, args ...string) []string {
	parts := strings.Fields(c.dockerCmd)
	parts = append(parts, cmd)
	parts = append(parts, args...)
	return parts
}

// PortMapping represents a port mapping
type PortMapping struct {
	HostPort      string `json:"host_port"`
	ContainerPort string `json:"container_port"`
}

// ContainerInfo represents container information (Docker CLI format)
type ContainerInfo struct {
	ID             string        `json:"ID"`
	Name           string        // populated from Names field
	Names          string        `json:"Names"`
	Image          string        `json:"Image"`
	Status         string        `json:"Status"`
	State          string        `json:"State"`
	Ports          string        `json:"Ports"`
	PublishedPorts []string      // host:container port mappings in format "host:container"
	PortsDetail    []PortMapping // detailed port mapping information
}

// ContainerConfig represents container configuration
type ContainerConfig struct {
	Name         string
	Image        string
	Env          []string
	PortBindings map[string]interface{}
	Volumes      map[string]interface{}
	NetworkMode  string
}

// Helper to parse ports
func ParsePort(port string) string {
	return port
}

type DockerContainerStats struct {
	RxBytes     int64
	TxBytes     int64
	CPUUsage    float64
	MemoryUsage int64
}

func (c *Client) GetContainersStats(containerIDs []string) (map[string]DockerContainerStats, error) {
	result := map[string]DockerContainerStats{}
	if len(containerIDs) == 0 {
		return result, nil
	}
	args := c.parseCommand("stats", "--no-stream", "--format",
		"{{.Container}}|{{.NetIO}}|{{.CPUPerc}}|{{.MemUsage}}")
	args = append(args, containerIDs...)
	cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
	hideConsoleWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}
		cid, netio, cpu, mem := parts[0], parts[1], parts[2], parts[3]

		// netio is "2.14kB / 1.23kB"
		nets := strings.Split(netio, "/")
		if len(nets) != 2 {
			continue
		}
		parseBytes := func(s string) int64 {
			s = strings.TrimSpace(s)
			mult := int64(1)
			if strings.HasSuffix(s, "KiB") {
				mult = 1024
				s = strings.TrimSuffix(s, "KiB")
			} else if strings.HasSuffix(s, "MiB") {
				mult = 1024 * 1024
				s = strings.TrimSuffix(s, "MiB")
			} else if strings.HasSuffix(s, "GiB") {
				mult = 1024 * 1024 * 1024
				s = strings.TrimSuffix(s, "GiB")
			} else if strings.HasSuffix(s, "B") {
				mult = 1
				s = strings.TrimSuffix(s, "B")
			}
			val, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
			return int64(val * float64(mult))
		}

		cpuUsage, _ := strconv.ParseFloat(strings.TrimRight(strings.TrimSpace(cpu), "%"), 64)

		memParts := strings.Split(mem, "/")
		memUsed := int64(0)
		if len(memParts) > 0 {
			memUsed = parseBytes(memParts[0])
		}

		result[cid] = DockerContainerStats{
			RxBytes:     parseBytes(nets[0]),
			TxBytes:     parseBytes(nets[1]),
			CPUUsage:    cpuUsage,
			MemoryUsage: memUsed,
		}
	}
	return result, nil
}

func (c *Client) GetContainersStartTimes(containerIDs []string) (map[string]time.Time, error) {
	result := map[string]time.Time{}
	for _, id := range containerIDs {
		args := c.parseCommand("inspect", "-f", "{{.State.StartedAt}}", id)
		cmd := exec.CommandContext(c.ctx, args[0], args[1:]...)
		hideConsoleWindow(cmd)
		output, err := cmd.Output()
		if err != nil {
			continue
		}
		text := strings.TrimSpace(string(output))
		if text == "" {
			continue
		}
		parsed, err := time.Parse(time.RFC3339Nano, text)
		if err != nil {
			continue
		}
		result[id] = parsed
	}
	return result, nil
}
