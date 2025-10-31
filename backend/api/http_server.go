package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func StartHeadlessServer(port int, appsAPI *AppsAPI, proxyAPI *ProxyAPI, settingsAPI *SettingsAPI, assets embed.FS) {
	mux := http.NewServeMux()

	// API handlers
	mux.HandleFunc("/api/apps/summary", func(w http.ResponseWriter, r *http.Request) {
		summary, err := appsAPI.GetDashboardSummary()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, summary, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/available", func(w http.ResponseWriter, r *http.Request) {
		apps, err := appsAPI.GetAvailableApps()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, apps, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/running", func(w http.ResponseWriter, r *http.Request) {
		running, err := appsAPI.GetRunningApps()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, running, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/start/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/start/")
		if err := appsAPI.StartApp(appID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "started"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/stop/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/stop/")
		if err := appsAPI.StopApp(appID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "stopped"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/restart/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/restart/")
		if err := appsAPI.RestartApp(appID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "restarted"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/logs/", func(w http.ResponseWriter, r *http.Request) {
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/logs/")
		tailStr := r.URL.Query().Get("tail")
		tail := 0
		if tailStr != "" {
			var err error
			tail, err = strconv.Atoi(tailStr)
			if err != nil {
				jsonResponse(w, map[string]string{"error": "Invalid tail value"}, http.StatusBadRequest)
				return
			}
		}
		logs, err := appsAPI.GetAppLogs(appID, tail)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"logs": logs}, http.StatusOK)
	})

	mux.HandleFunc("/api/container/logs/all/", func(w http.ResponseWriter, r *http.Request) {
		containerID := strings.TrimPrefix(r.URL.Path, "/api/container/logs/all/")
		logs, err := appsAPI.GetContainerLogsAll(containerID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"logs": logs}, http.StatusOK)
	})

	mux.HandleFunc("/api/container/logs/", func(w http.ResponseWriter, r *http.Request) {
		containerID := strings.TrimPrefix(r.URL.Path, "/api/container/logs/")
		tailStr := r.URL.Query().Get("tail")
		if tailStr != "" {
			tail, err := strconv.Atoi(tailStr)
			if err != nil {
				jsonResponse(w, map[string]string{"error": "Invalid tail value"}, http.StatusBadRequest)
				return
			}
			logs, err := appsAPI.GetContainerLogsTail(containerID, tail)
			if err != nil {
				jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
			jsonResponse(w, map[string]string{"logs": logs}, http.StatusOK)
		} else {
			logs, err := appsAPI.GetContainerLogs(containerID)
			if err != nil {
				jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
			jsonResponse(w, map[string]string{"logs": logs}, http.StatusOK)
		}
	})

	mux.HandleFunc("/api/apps/stats/", func(w http.ResponseWriter, r *http.Request) {
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/stats/")
		stats, err := appsAPI.GetAppStats(appID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, stats, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/instances/", func(w http.ResponseWriter, r *http.Request) {
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/instances/")
		instances, err := appsAPI.GetAppInstances(appID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, instances, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/configured", func(w http.ResponseWriter, r *http.Request) {
		configured, err := appsAPI.GetConfiguredApps()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, configured, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/credentials/", func(w http.ResponseWriter, r *http.Request) {
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/credentials/")
		credentials, err := appsAPI.GetAppCredentials(appID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, credentials, http.StatusOK)
	})

	mux.HandleFunc("/api/container/env/", func(w http.ResponseWriter, r *http.Request) {
		containerID := strings.TrimPrefix(r.URL.Path, "/api/container/env/")
		env, err := appsAPI.GetContainerEnvironmentVars(containerID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, env, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/deploy/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/deploy/")
		var formData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		if err := appsAPI.DeployApp(appID, formData); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "deployed"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/deploy-with-proxy/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/deploy-with-proxy/")
		var data struct {
			FormData map[string]string `json:"formData"`
			ProxyID  string            `json:"proxyID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		if err := appsAPI.DeployAppWithProxyId(appID, data.FormData, data.ProxyID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "deployed"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/deploy-with-proxies/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		appID := strings.TrimPrefix(r.URL.Path, "/api/apps/deploy-with-proxies/")
		var data struct {
			FormData map[string]string `json:"formData"`
			ProxyIDs []string          `json:"proxyIDs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		result, err := appsAPI.DeployAppWithProxies(appID, data.FormData, data.ProxyIDs)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/remove-instance/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		instanceID := strings.TrimPrefix(r.URL.Path, "/api/apps/remove-instance/")
		if err := appsAPI.RemoveAppInstance(instanceID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "removed"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/remove/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		containerID := strings.TrimPrefix(r.URL.Path, "/api/apps/remove/")
		if err := appsAPI.RemoveApp(containerID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "removed"}, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/deploy-selective", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			ProxyID  string   `json:"proxyID"`
			ProxyURL string   `json:"proxyURL"`
			AppIDs   []string `json:"appIDs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		result, err := appsAPI.DeployAppWithProxiesSelective(data.ProxyID, data.ProxyURL, data.AppIDs)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result, http.StatusOK)
	})

	mux.HandleFunc("/api/apps/deploy-with-proxy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		result, err := appsAPI.DeployAppWithProxy(data)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result, http.StatusOK)
	})

	// ProxyAPI Handlers
	mux.HandleFunc("/api/proxies/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			ProxyStr       string   `json:"proxyStr"`
			AutoDeploy     bool     `json:"autoDeploy"`
			SelectedAppIDs []string `json:"selectedAppIDs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		result, err := proxyAPI.AddProxy(data.ProxyStr, data.AutoDeploy, data.SelectedAppIDs)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/remove/", func(w http.ResponseWriter, r *http.Request) {
		proxyID := strings.TrimPrefix(r.URL.Path, "/api/proxies/remove/")
		result, err := proxyAPI.RemoveProxy(proxyID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/remove/confirm/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		proxyID := strings.TrimPrefix(r.URL.Path, "/api/proxies/remove/confirm/")
		if err := proxyAPI.ConfirmRemoveProxy(proxyID); err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "removed"}, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/list", func(w http.ResponseWriter, r *http.Request) {
		proxies, err := proxyAPI.ListProxies()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, proxies, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/test/", func(w http.ResponseWriter, r *http.Request) {
		proxyID := strings.TrimPrefix(r.URL.Path, "/api/proxies/test/")
		result, err := proxyAPI.TestProxy(proxyID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/containers/", func(w http.ResponseWriter, r *http.Request) {
		proxyID := strings.TrimPrefix(r.URL.Path, "/api/proxies/containers/")
		containers, err := proxyAPI.GetProxyContainers(proxyID)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, containers, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/apps-running", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		var proxyIDs []string
		if err := json.NewDecoder(r.Body).Decode(&proxyIDs); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		apps, err := proxyAPI.GetAppsRunningOnProxies(proxyIDs)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, apps, http.StatusOK)
	})

	mux.HandleFunc("/api/proxies/configured-apps", func(w http.ResponseWriter, r *http.Request) {
		apps, err := proxyAPI.GetConfiguredAppsForProxy()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, apps, http.StatusOK)
	})

	// SettingsAPI Handlers
	mux.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		settings, err := settingsAPI.GetSettings()
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, settings, http.StatusOK)
	})

	mux.HandleFunc("/api/settings/autostart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		_, err := settingsAPI.SetAutoStart(data.Enabled)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
	})

	mux.HandleFunc("/api/settings/showintray", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		_, err := settingsAPI.SetShowInTray(data.Enabled)
		if err != nil {
			jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Printf("Starting headless server on port %d\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting headless server: %v\n", err)
	}
}
