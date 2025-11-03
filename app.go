package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"free-income-from-bandwidth/backend/api"
	"free-income-from-bandwidth/backend/apps"
	"free-income-from-bandwidth/backend/config"
	"free-income-from-bandwidth/backend/docker"
	"free-income-from-bandwidth/backend/monitor"
	"free-income-from-bandwidth/backend/notifications"
	"free-income-from-bandwidth/backend/orchestrator"
	"free-income-from-bandwidth/backend/proxy"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Define command-line flags
	headless := flag.Bool("headless", false, "Run in headless mode")
	port := flag.Int("port", 8080, "Port for headless server")
	flag.Parse()

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}

	configsDir := filepath.Join(wd, "configs")
	dbPath := filepath.Join(wd, "data", "monitor.db")

	// Create data directory if it doesn't exist
	_ = os.MkdirAll(filepath.Join(wd, "data"), 0755)

	// Initialize Docker client
	dockerClient, err := docker.NewDockerClient("")
	if err != nil {
		fmt.Printf("Warning: Failed to initialize Docker client: %v\n", err)
		// Continue without Docker for now
	}

	// Initialize config loader
	configLoader := config.NewLoader(configsDir)
	if err := configLoader.LoadAppConfigs(); err != nil {
		fmt.Printf("Warning: Failed to load app configs: %v\n", err)
	}



	// Initialize monitor
	monitorCollector, err := monitor.NewCollector(dbPath)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize monitor: %v\n", err)
	}

	// Initialize proxy manager
	proxyManager := proxy.NewManager()

	// Initialize instance manager
	instanceManager := apps.NewInstanceManager()

	// Initialize credential store
	credentialStore := config.NewCredentialStore()

	// Initialize orchestrator (not used yet)
	_ = orchestrator.NewManager()

	// Initialize notifications (not used yet)
	notifConfig := &notifications.Config{
		Enabled: true,
	}
	_ = notifications.NewHandler(notifConfig, monitorCollector)

	// Initialize API
	appsAPI := api.NewAppsAPI(dockerClient, configLoader, monitorCollector, instanceManager, credentialStore, proxyManager)

	// Initialize Proxy API
	proxyAPI := api.NewProxyAPI(proxyManager, instanceManager, credentialStore, appsAPI)

	// Settings API
	settingsAPI := api.NewSettingsAPI(wd)

	if *headless {
		// Start headless server
		api.StartHeadlessServer(*port, appsAPI, proxyAPI, settingsAPI, assets)
	} else {
		// Create application with options
		err = wails.Run(&options.App{
			Title:  "Free Income From Bandwidth",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
			OnStartup: func(ctx context.Context) {
				appsAPI.OnStartup(ctx)
				proxyAPI.OnStartup(ctx)
				settingsAPI.OnStartup(ctx)
			},
			Bind: []interface{}{
				appsAPI,
				proxyAPI,
				settingsAPI,
			},
		})

		if err != nil {
			println("Error:", err.Error())
		}
	}
}
