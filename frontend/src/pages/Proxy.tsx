import { useState, useEffect } from 'react'
import { Plus, RefreshCw, Trash2, CheckCircle, XCircle, Upload, X } from 'lucide-react'
import { AddProxy, ListProxies, RemoveProxy, ConfirmRemoveProxy, TestProxy, GetConfiguredAppsForProxy, GetAppsRunningOnProxies, DeployAppWithProxy } from '../services/api';

interface Proxy {
  id: string
  url: string
  status: 'active' | 'failed' | 'testing'
  latency?: number
  lastChecked?: string
}

interface ConfiguredApp {
  app_id: string
  app_name: string
  device_name: string
}

export default function Proxy() {
  const [proxies, setProxies] = useState<Proxy[]>([])
  const [proxyInput, setProxyInput] = useState('')
  const [selectedFormat, setSelectedFormat] = useState<'single' | 'bulk'>('single')
  const [autoDeploy, setAutoDeploy] = useState(false)
  const [adding, setAdding] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [showAppSelector, setShowAppSelector] = useState(false)
  const [configuredApps, setConfiguredApps] = useState<ConfiguredApp[]>([])
  const [selectedApps, setSelectedApps] = useState<string[]>([])
  const [bulkProxyInput, setBulkProxyInput] = useState('')
  const [selectedProxies, setSelectedProxies] = useState<string[]>([])
  const [showBulkActions, setShowBulkActions] = useState(false)
  const [runningAppsOnSelectedProxies, setRunningAppsOnSelectedProxies] = useState<string[]>([])

  useEffect(() => {
    loadProxies()
    loadConfiguredApps()
  }, [])

  const loadProxies = async () => {
    try {
      const data = await ListProxies()
      const proxyList: Proxy[] = data.map((p: any) => ({
        id: p.id,
        url: p.url,
        status: 'active',
        latency: 0,
        lastChecked: new Date().toISOString()
      }))
      setProxies(proxyList)
    } catch (error) {
      console.error('Failed to load proxies:', error)
    }
  }

  const loadConfiguredApps = async () => {
    try {
      const apps = await GetConfiguredAppsForProxy()
      setConfiguredApps(apps as ConfiguredApp[])
    } catch (error) {
      console.error('Failed to load configured apps:', error)
    }
  }

  const handleAddProxy = async () => {
    if (!proxyInput.trim() || adding) return
    
    // Show app selector modal if there are configured apps
    if (configuredApps.length > 0 && !autoDeploy) {
      setShowAppSelector(true)
      return
    }
    
    // Proceed with adding
    await addProxyWithSelections()
  }

  const addProxyWithSelections = async () => {
    if (adding) return // Prevent multiple clicks
    
    setAdding(true)
    
    try {
      // Add proxy with selected apps
      const result = await AddProxy(proxyInput.trim(), autoDeploy, selectedApps)
      
      // Reload proxy list
      await loadProxies()
      
      // Clear input and reset
      setProxyInput('')
      setSelectedApps([])
      setShowAppSelector(false)
      
      // Show success message
      if (result.deployed_containers) {
        const count = result.deployed_containers.length
        alert(`Proxy added successfully! Deployed to ${count} app(s).`)
      } else {
        alert('Proxy added successfully!')
      }
    } catch (error: any) {
      console.error('Failed to add proxy:', error)
      alert(`Failed to add proxy: ${error.message || error}`)
    } finally {
      setAdding(false)
    }
  }

  const handleSelectAll = () => {
    // Get available apps (not running on selected proxies)
    const availableApps = configuredApps.filter(app => {
      if (selectedProxies.length > 0) {
        return !runningAppsOnSelectedProxies.includes(app.app_id)
      }
      return true
    })
    
    if (selectedApps.length === availableApps.length) {
      setSelectedApps([])
    } else {
      setSelectedApps(availableApps.map(app => app.app_id))
    }
  }

  const handleAppToggle = (appId: string) => {
    if (selectedApps.includes(appId)) {
      setSelectedApps(selectedApps.filter(id => id !== appId))
    } else {
      setSelectedApps([...selectedApps, appId])
    }
  }

  const handleImportFile = () => {
    // Create file input element
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.txt,.list'
    input.onchange = async (e: any) => {
      const file = e.target.files?.[0]
      if (!file) return

      try {
        const text = await file.text()
        setBulkProxyInput(text)
      } catch (error) {
        console.error('Failed to read file:', error)
        alert('Failed to read file')
      }
    }
    input.click()
  }

  const handleBulkAdd = async () => {
    const lines = bulkProxyInput.split('\n').filter(line => line.trim())
    if (lines.length === 0) return

    // Show app selection modal if there are configured apps
    if (configuredApps.length > 0) {
      setShowAppSelector(true)
    } else {
      // Proceed with bulk add without app selection
      await performBulkAdd(lines)
    }
  }

  const performBulkAdd = async (proxyLines: string[]) => {
    setAdding(true)
    
    try {
      let successCount = 0
      let failCount = 0

      for (const proxyStr of proxyLines) {
        try {
          await AddProxy(proxyStr.trim(), autoDeploy, selectedApps)
          successCount++
        } catch (error) {
          console.error(`Failed to add proxy ${proxyStr}:`, error)
          failCount++
        }
      }

      // Reload proxy list
      await loadProxies()
      
      // Clear input
      setBulkProxyInput('')
      setSelectedApps([])
      setShowAppSelector(false)
      
      // Show results
      alert(`Added ${successCount} proxy/proxies successfully${failCount > 0 ? `, ${failCount} failed` : ''}`)
    } catch (error: any) {
      alert(`Bulk import failed: ${error.message || error}`)
    } finally {
      setAdding(false)
    }
  }

  const handleBulkSelectAll = () => {
    if (selectedApps.length === configuredApps.length) {
      setSelectedApps([])
    } else {
      setSelectedApps(configuredApps.map(app => app.app_id))
    }
  }

  const handleProxySelect = (proxyId: string) => {
    if (selectedProxies.includes(proxyId)) {
      setSelectedProxies(selectedProxies.filter(id => id !== proxyId))
    } else {
      setSelectedProxies([...selectedProxies, proxyId])
    }
  }

  const handleSelectAllProxies = () => {
    if (selectedProxies.length === proxies.length) {
      setSelectedProxies([])
    } else {
      setSelectedProxies(proxies.map(p => p.id))
    }
  }

  const handleBulkDeleteProxies = async () => {
    if (selectedProxies.length === 0) return

    const confirmMessage = `Are you sure you want to delete ${selectedProxies.length} proxy/proxies? This will also stop and remove all containers using these proxies.`
    if (!confirm(confirmMessage)) return

    setAdding(true)
    try {
      let successCount = 0
      let failCount = 0

      for (const proxyId of selectedProxies) {
        try {
          await ConfirmRemoveProxy(proxyId)
          successCount++
        } catch (error) {
          console.error(`Failed to delete proxy ${proxyId}:`, error)
          failCount++
        }
      }

      await loadProxies()
      setSelectedProxies([])
      setShowBulkActions(false)
      
      alert(`Deleted ${successCount} proxy/proxies successfully${failCount > 0 ? `, ${failCount} failed` : ''}`)
    } catch (error: any) {
      alert(`Bulk delete failed: ${error.message || error}`)
    } finally {
      setAdding(false)
    }
  }

  const handleBulkAddAppsToProxies = async () => {
    if (selectedProxies.length === 0) return
    
    // Load running apps on selected proxies
    try {
      const runningApps = await GetAppsRunningOnProxies(selectedProxies)
      setRunningAppsOnSelectedProxies(runningApps)
    } catch (error) {
      console.error('Failed to load running apps:', error)
      setRunningAppsOnSelectedProxies([])
    }
    
    setShowAppSelector(true)
  }

  const handleAddAppsToSelectedProxies = async () => {
    if (selectedProxies.length === 0 || selectedApps.length === 0) return
    if (adding) return // Prevent multiple clicks

    setAdding(true)
    try {
      let successCount = 0
      let failCount = 0

      // For each selected proxy, deploy all selected apps
      for (const proxyId of selectedProxies) {
        try {
          // Get proxy details
          const proxy = proxies.find(p => p.id === proxyId)
          if (!proxy) continue

          // Deploy each selected app to this proxy
          for (const appId of selectedApps) {
            try {
              await DeployAppWithProxy(appId, proxy.url)
              successCount++
            } catch (error) {
              console.error(`Failed to deploy app ${appId} to proxy ${proxyId}:`, error)
              failCount++
            }
          }
        } catch (error) {
          console.error(`Failed to process proxy ${proxyId}:`, error)
          failCount++
        }
      }

      await loadProxies()
      setSelectedApps([])
      setSelectedProxies([])
      setShowAppSelector(false)
      setShowBulkActions(false)
      
      alert(`Deployed ${successCount} app instances successfully${failCount > 0 ? `, ${failCount} failed` : ''}`)
    } catch (error: any) {
      alert(`Deployment failed: ${error.message || error}`)
    } finally {
      setAdding(false)
    }
  }

  const handleTestAll = async () => {
    // Test each proxy
    const updatedProxies = await Promise.all(proxies.map(async (p) => {
      try {
        const result = await TestProxy(p.id)
        return {
          ...p,
          status: result.status === 'success' ? 'active' : 'failed',
          latency: result.latency
        }
      } catch (error) {
        return { ...p, status: 'failed' as const }
      }
    }))
    setProxies(updatedProxies)
  }

  const handleRemoveProxy = async (id: string) => {
    try {
      // First check affected containers
      const check = await RemoveProxy(id)
      
      if (check.affected_containers > 0) {
        const confirm = window.confirm(
          `This proxy is used by ${check.affected_containers} container(s). ` +
          `Removing it will stop and remove all these containers. Continue?`
        )
        if (!confirm) return
      }
      
      // Remove proxy and containers
      await ConfirmRemoveProxy(id)
      
      // Reload list
      await loadProxies()
    } catch (error: any) {
      console.error('Failed to remove proxy:', error)
      alert(`Failed to remove proxy: ${error.message || error}`)
    }
  }

  return (
    <div className="p-6 space-y-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-white">Proxy Configuration</h1>
        <p className="text-gray-400 mt-2">Manage proxies for your apps</p>
      </div>

      {/* Removed Enable Proxy Toggle */}

      {/* Add Proxy Section */}
      <div className="bg-gray-800 rounded-lg p-6 space-y-4">
        <h2 className="text-xl font-semibold text-white">Add Proxies</h2>

        {/* Format Selection */}
        <div className="flex gap-4">
          <button
            onClick={() => setSelectedFormat('single')}
            className={`px-4 py-2 rounded-lg font-medium transition-colors ${
              selectedFormat === 'single'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            }`}
          >
            Single Proxy
          </button>
          <button
            onClick={() => setSelectedFormat('bulk')}
            className={`px-4 py-2 rounded-lg font-medium transition-colors ${
              selectedFormat === 'bulk'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            }`}
          >
            Bulk Import
          </button>
        </div>

        {/* Input Section */}
        {selectedFormat === 'single' ? (
          <div className="flex gap-2">
            <input
              type="text"
              value={proxyInput}
              onChange={(e) => setProxyInput(e.target.value)}
              placeholder="http://user:pass@host:port or socks5://host:port"
              className="flex-1 bg-gray-700 text-white px-4 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              onKeyPress={(e) => e.key === 'Enter' && handleAddProxy()}
            />
            <button
              onClick={handleAddProxy}
              disabled={adding}
              className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {adding ? (
                <>
                  <RefreshCw className="h-5 w-5 animate-spin" />
                  Adding...
                </>
              ) : (
                <>
                  <Plus className="h-5 w-5" />
                  Add
                </>
              )}
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            <textarea
              value={bulkProxyInput}
              onChange={(e) => setBulkProxyInput(e.target.value)}
              placeholder="Paste proxies line by line or import from file&#10;Format:&#10;http://user:pass@host:port&#10;socks5://host:port"
              className="w-full bg-gray-700 text-white px-4 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 h-32 font-mono text-sm"
            />
            <div className="flex gap-2">
              <button
                onClick={handleImportFile}
                className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium flex items-center gap-2"
              >
                <Upload className="h-5 w-5" />
                Import File
              </button>
              <button
                onClick={handleBulkAdd}
                disabled={adding || !bulkProxyInput.trim()}
                className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {adding ? (
                  <>
                    <RefreshCw className="h-5 w-5 animate-spin" />
                    Adding...
                  </>
                ) : (
                  <>
                    <Plus className="h-5 w-5" />
                    Add All ({bulkProxyInput.split('\n').filter(l => l.trim()).length})
                  </>
                )}
              </button>
            </div>
          </div>
        )}

        {/* Auto-deploy Toggle */}
        <div className="flex items-center gap-3 bg-gray-700/50 rounded-lg p-3">
          <input
            type="checkbox"
            id="autoDeploy"
            checked={autoDeploy}
            onChange={(e) => setAutoDeploy(e.target.checked)}
            className="h-4 w-4 rounded border-gray-600 bg-gray-800 text-blue-600 focus:ring-blue-500"
          />
          <label htmlFor="autoDeploy" className="text-sm text-gray-300 cursor-pointer">
            Auto-deploy to configured apps when proxy is added
          </label>
        </div>

        {/* Info */}
        <p className="text-gray-400 text-sm">
          💡 Supported formats: HTTP, HTTPS, SOCKS4, SOCKS5
        </p>
      </div>

      {/* Proxy List */}
      <div className="bg-gray-800 rounded-lg p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-white">
            Configured Proxies ({proxies.length})
          </h2>
          <div className="flex gap-2">
            {selectedProxies.length > 0 && (
              <button
                onClick={() => {
                  setSelectedProxies([])
                  setShowBulkActions(false)
                }}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm"
              >
                Clear Selection ({selectedProxies.length})
              </button>
            )}
            {!showBulkActions && proxies.length > 0 && (
              <button
                onClick={() => setShowBulkActions(true)}
                className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-sm"
              >
                Bulk Select
              </button>
            )}
            {proxies.length > 0 && (
              <button
                onClick={handleTestAll}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg flex items-center gap-2"
              >
                <RefreshCw className="h-5 w-5" />
                Test All
              </button>
            )}
          </div>
        </div>

        {proxies.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-400">No proxies added yet</p>
          </div>
        ) : (
          <div className="space-y-2">
            {showBulkActions && (
              <div className="bg-purple-900/20 border border-purple-700 rounded-lg p-4 mb-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <input
                      type="checkbox"
                      checked={selectedProxies.length === proxies.length && proxies.length > 0}
                      onChange={handleSelectAllProxies}
                      className="w-4 h-4 text-purple-600 bg-gray-700 border-gray-600 rounded focus:ring-purple-500"
                    />
                    <span className="text-purple-200 font-medium">
                      Select All ({selectedProxies.length}/{proxies.length})
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={handleBulkAddAppsToProxies}
                      disabled={selectedProxies.length === 0}
                      className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg text-sm"
                    >
                      Add Apps to Selected
                    </button>
                    <button
                      onClick={handleBulkDeleteProxies}
                      disabled={selectedProxies.length === 0}
                      className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg text-sm"
                    >
                      Delete Selected ({selectedProxies.length})
                    </button>
                  </div>
                </div>
              </div>
            )}
            {proxies.map((proxy) => (
              <div
                key={proxy.id}
                className={`flex items-center justify-between p-4 bg-gray-700 rounded-lg hover:bg-gray-650 transition ${
                  selectedProxies.includes(proxy.id) ? 'ring-2 ring-purple-500' : ''
                }`}
              >
                <div className="flex items-center gap-3 flex-1">
                  {showBulkActions && (
                    <input
                      type="checkbox"
                      checked={selectedProxies.includes(proxy.id)}
                      onChange={() => handleProxySelect(proxy.id)}
                      className="w-4 h-4 text-purple-600 bg-gray-700 border-gray-600 rounded focus:ring-purple-500"
                    />
                  )}
                  {proxy.status === 'active' ? (
                    <CheckCircle className="h-5 w-5 text-green-500" />
                  ) : proxy.status === 'failed' ? (
                    <XCircle className="h-5 w-5 text-red-500" />
                  ) : (
                    <RefreshCw className="h-5 w-5 text-yellow-500 animate-spin" />
                  )}
                  <div className="flex-1">
                    <p className="text-white font-mono text-sm">{proxy.url}</p>
                    {proxy.latency && (
                      <p className="text-gray-400 text-xs">Latency: {proxy.latency}ms</p>
                    )}
                    {proxy.containerCount !== undefined && (
                      <p className="text-gray-400 text-xs">Used by {proxy.containerCount} container(s)</p>
                    )}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {!showBulkActions && (
                    <button
                      onClick={() => handleRemoveProxy(proxy.id)}
                      className="p-2 text-red-400 hover:text-red-300 hover:bg-red-900/20 rounded"
                    >
                      <Trash2 className="h-5 w-5" />
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Removed Proxy Usage Stats dependent on enable toggle */}

      {/* App Selection Modal */}
      {showAppSelector && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-lg max-w-2xl w-full max-h-[80vh] overflow-hidden flex flex-col">
            {/* Header */}
            <div className="p-6 border-b border-gray-700 flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-white">
                  {selectedFormat === 'bulk' ? 'Select Apps for Bulk Import' : 
                   selectedProxies.length > 0 ? 'Add Apps to Selected Proxies' : 
                   'Select Apps to Deploy'}
                </h2>
                <p className="text-gray-400 text-sm mt-1">
                  {selectedFormat === 'bulk' ? 'Choose which apps to deploy with the imported proxies' :
                   selectedProxies.length > 0 ? `Add selected apps to ${selectedProxies.length} proxy/proxies` :
                   'Choose which apps should use this proxy'}
                </p>
                {selectedProxies.length > 0 && (
                  <div className="mt-2 p-2 bg-purple-900/20 border border-purple-700 rounded text-sm">
                    <p className="text-purple-200">Selected proxies: {selectedProxies.length}</p>
                  </div>
                )}
              </div>
              <button
                onClick={() => {
                  setShowAppSelector(false)
                  setSelectedApps([])
                  setRunningAppsOnSelectedProxies([])
                }}
                className="p-2 hover:bg-gray-700 rounded-lg transition"
              >
                <X className="h-6 w-6 text-gray-400" />
              </button>
            </div>

            {/* App List */}
            <div className="flex-1 overflow-y-auto p-6">
              {configuredApps.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-gray-400">No configured apps found</p>
                  <p className="text-gray-500 text-sm mt-2">Configure apps from the App Store first</p>
                </div>
              ) : (
                <div className="space-y-2">
                  {/* Select All Button */}
                  {(() => {
                    const availableApps = configuredApps.filter(app => {
                      if (selectedProxies.length > 0) {
                        return !runningAppsOnSelectedProxies.includes(app.app_id)
                      }
                      return true
                    })
                    
                    return (
                      <button
                        onClick={handleSelectAll}
                        className="w-full px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm font-medium mb-4"
                      >
                        {selectedApps.length === availableApps.length ? 'Deselect All' : `Select All (${availableApps.length} available)`}
                      </button>
                    )
                  })()}

                  {/* App Checkboxes */}
                  {configuredApps
                    .filter(app => {
                      // When adding to selected proxies, hide apps already running on those proxies
                      if (selectedProxies.length > 0) {
                        return !runningAppsOnSelectedProxies.includes(app.app_id)
                      }
                      return true
                    })
                    .map((app) => {
                      const isRunning = runningAppsOnSelectedProxies.includes(app.app_id)
                      return (
                        <label
                          key={app.app_id}
                          className={`flex items-center gap-3 p-4 rounded-lg transition ${
                            isRunning 
                              ? 'bg-gray-800/50 opacity-50 cursor-not-allowed' 
                              : 'bg-gray-700/50 hover:bg-gray-700 cursor-pointer'
                          }`}
                        >
                          <input
                            type="checkbox"
                            checked={selectedApps.includes(app.app_id)}
                            onChange={() => handleAppToggle(app.app_id)}
                            disabled={isRunning}
                            className="h-5 w-5 rounded border-gray-600 bg-gray-800 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
                          />
                          <div className="flex-1">
                            <p className="text-white font-medium">
                              {app.app_name}
                              {isRunning && (
                                <span className="ml-2 text-xs bg-yellow-600 text-yellow-100 px-2 py-1 rounded">
                                  Already Running
                                </span>
                              )}
                            </p>
                            <p className="text-gray-400 text-sm">{app.device_name}</p>
                          </div>
                        </label>
                      )
                    })}
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="p-6 border-t border-gray-700 flex justify-end gap-3">
              <button
                onClick={() => {
                  setShowAppSelector(false)
                  setSelectedApps([])
                  setRunningAppsOnSelectedProxies([])
                }}
                className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium"
              >
                Cancel
              </button>
              <button
                onClick={
                  selectedFormat === 'bulk' 
                    ? () => performBulkAdd(bulkProxyInput.split('\n').filter(line => line.trim())) 
                    : selectedProxies.length > 0 
                      ? handleAddAppsToSelectedProxies
                      : addProxyWithSelections
                }
                disabled={adding || selectedApps.length === 0}
                className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {adding ? (
                  <>
                    <RefreshCw className="h-5 w-5 animate-spin inline mr-2" />
                    {selectedProxies.length > 0 ? 'Deploying...' : 'Adding...'}
                  </>
                ) : (
                  selectedProxies.length > 0 
                    ? `Deploy to ${selectedProxies.length} proxy/proxies`
                    : `Deploy to ${selectedApps.length} app(s)`
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}



