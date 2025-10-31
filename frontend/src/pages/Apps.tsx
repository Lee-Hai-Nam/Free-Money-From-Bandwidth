import { useState, useEffect, useRef } from 'react'
import { Play, Square, RefreshCw, Plus, Globe, Server, Trash2, Copy, Check, Loader2 } from 'lucide-react'
import { GetRunningApps, StartApp, StopApp, RestartApp, RemoveApp } from '../services/api';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

interface AppInstance {
  instance_id: string
  app_id: string
  proxy_id: string
  container_id: string
  device_name: string
  status: string
  proxy_url: string
  dashboardLink?: string
}

interface App {
  app_id: string
  name: string
  instances: AppInstance[]
  dashboardLink?: string
}

interface AppsProps {
  onAddNew?: () => void
}

export default function Apps({ onAddNew }: AppsProps) {
  const [apps, setApps] = useState<App[]>([])
  const [loading, setLoading] = useState(true)
  const [copiedLink, setCopiedLink] = useState<string | null>(null)
  const [actionLoadingId, setActionLoadingId] = useState<string | null>(null)


  useEffect(() => {
    loadApps()
  }, [])

  const loadApps = async () => {
    try {
      setLoading(true)
      const data = await GetRunningApps()
      const grouped = await groupByApp(data)
      setApps(grouped)
    } catch (error) {
      console.error('Failed to load apps:', error)
    } finally {
      setLoading(false)
    }
  }

  const groupByApp = async (containers: any[]): Promise<App[]> => {
    const appMap = new Map<string, App>()
    
    for (const container of containers) {
      const parts = container.name.split('_')
      if (parts.length < 3) continue
      
      const appId = parts[parts.length - 2]
      const instanceType = parts[parts.length - 1]
      const isLocal = instanceType === 'local'
      
      // Build dashboard link for EarnApp instances
      let instanceDashboardLink: string | undefined
      if (appId === 'earnapp' && container.sdkNodeID) {
        // Format: earnapp.com/r/sdk-node-{uuid}
        instanceDashboardLink = `https://earnapp.com/r/${container.sdkNodeID}`
      }

      const instance: AppInstance = {
        instance_id: container.id,
        app_id: appId,
        proxy_id: isLocal ? '' : instanceType,
        container_id: container.id,
        device_name: container.name,
        status: container.state === 'running' ? 'running' : 'stopped',
        proxy_url: isLocal ? 'Local' : `Proxy: ${instanceType.replace('proxy', '').substring(0, 8)}`,
        dashboardLink: instanceDashboardLink
      }

      if (!appMap.has(appId)) {
        appMap.set(appId, {
          app_id: appId,
          name: appId,
          instances: [],
          dashboardLink: instanceDashboardLink
        })
      }

      appMap.get(appId)!.instances.push(instance)
    }

    const result = Array.from(appMap.values())
    result.forEach(app => {
      app.instances.sort((a, b) => {
        if (a.proxy_id === '' && b.proxy_id !== '') return -1
        if (a.proxy_id !== '' && b.proxy_id === '') return 1
        if (a.status === 'running' && b.status !== 'running') return -1
        if (a.status !== 'running' && b.status === 'running') return 1
        return 0
      })
    })

    return result
  }

  const handleInstanceAction = async (action: string, instance: AppInstance) => {
    try {
      setActionLoadingId(`${action}:${instance.instance_id}`)
      if (action === 'start') {
        await StartApp(instance.container_id)
      } else if (action === 'stop') {
        await StopApp(instance.container_id)
      } else if (action === 'restart') {
        await RestartApp(instance.container_id)
      } else if (action === 'remove') {
        if (!confirm('Remove this container?')) return
        await RemoveApp(instance.container_id)
      }
      // Refresh list after action
      await loadApps()
    } catch (error) {
      console.error(`Failed to ${action}:`, error)
      alert(`Failed to ${action}: ${error}`)
    } finally {
      setActionLoadingId(null)
    }
  }

  const handleCopyLink = async (link: string, instance: AppInstance) => {
    try {
      await navigator.clipboard.writeText(link)
      setCopiedLink(instance.instance_id)
      setTimeout(() => setCopiedLink(null), 2000)
    } catch (error) {
      console.error('Failed to copy:', error)
    }
  }

  const handleOpenDashboard = (link: string) => {
    BrowserOpenURL(link)
  }

  // Logs feature removed per request

  if (loading) {
    return (
      <div className="p-6 flex items-center justify-center h-64">
        <div className="text-neutral-400">Loading apps...</div>
      </div>
    )
  }

  return (
    <div className="p-6 space-y-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-neutral-100">My Apps</h1>
          <p className="text-neutral-400 mt-2">Manage your bandwidth sharing applications</p>
        </div>
        {onAddNew && (
          <button
            onClick={onAddNew}
            className="flex items-center gap-2 px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg font-medium"
          >
            <Plus className="h-5 w-5" />
            Add New App
          </button>
        )}
      </div>

      <div className="space-y-6">
        {apps.length === 0 ? (
          <div className="bg-neutral-900 rounded-lg p-12 text-center border border-neutral-800">
            <p className="text-neutral-400 text-lg mb-4">No apps deployed yet</p>
            <p className="text-neutral-500 text-sm mb-6">
              Configure apps from the app store to start earning
            </p>
            {onAddNew && (
              <button
                onClick={onAddNew}
                className="inline-flex items-center gap-2 px-6 py-3 bg-brand-600 hover:bg-brand-700 text-white rounded-lg font-medium"
              >
                <Plus className="h-5 w-5" />
                Add Your First App
              </button>
            )}
          </div>
        ) : (
          apps.map((app) => (
            <div key={app.app_id} className="bg-neutral-900 rounded-lg overflow-hidden border border-neutral-800">
              <div className="p-6 border-b border-neutral-800">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-xl font-semibold text-neutral-100 capitalize">
                      {app.app_id.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                    </h3>
                    <p className="text-neutral-400 text-sm mt-1">
                      {app.instances.length} {app.instances.length === 1 ? 'instance' : 'instances'}
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="px-3 py-1 bg-brand-600/20 text-brand-400 rounded-full text-xs font-medium">
                      {app.instances.filter(i => i.status === 'running').length} running
                    </span>
                  </div>
                </div>
              </div>

              <div className="divide-y divide-neutral-800">
                {app.instances.map((instance) => (
                  <div key={instance.instance_id} className="p-6 hover:bg-neutral-800/50 transition">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4 flex-1">
                        {instance.proxy_id ? (
                          <div className="flex items-center gap-2 px-3 py-1 bg-purple-600/20 text-purple-400 rounded-lg text-xs font-medium border border-purple-500/30">
                            <Globe className="h-4 w-4" />
                            Proxy
                          </div>
                        ) : (
                          <div className="flex items-center gap-2 px-3 py-1 bg-blue-600/20 text-blue-400 rounded-lg text-xs font-medium border border-blue-500/30">
                            <Server className="h-4 w-4" />
                            Local (Primary)
                          </div>
                        )}

                        <div>
                          <p className="text-neutral-100 font-medium">{instance.device_name}</p>
                          <p className="text-neutral-400 text-sm">Container: {instance.container_id.substring(0, 12)}...</p>
                          {instance.dashboardLink && (
                            <div className="flex items-center gap-2 mt-1 text-xs">
                              <button
                                onClick={() => handleOpenDashboard(instance.dashboardLink!)}
                                className="text-brand-500 hover:underline"
                                title="Open in browser"
                              >
                                Link: {instance.dashboardLink.substring(instance.dashboardLink.lastIndexOf('/') + 1)}
                              </button>
                              <button
                                onClick={() => handleCopyLink(instance.dashboardLink!, instance)}
                                className="p-1 rounded hover:bg-neutral-700 text-neutral-300"
                                title="Copy"
                              >
                                {copiedLink === instance.instance_id ? (
                                  <Check className="h-3 w-3 text-success" />
                                ) : (
                                  <Copy className="h-3 w-3" />
                                )}
                              </button>
                            </div>
                          )}
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <span className={`px-3 py-1 rounded-full text-xs font-medium ${instance.status === 'running' ? 'bg-success text-white' : 'bg-neutral-600 text-neutral-200'}`}>
                          {instance.status}
                        </span>

                        <div className="flex gap-1">
                          {instance.status === 'running' ? (
                            <button
                              onClick={() => handleInstanceAction('stop', instance)}
                              className="p-2 text-danger hover:bg-danger/20 rounded transition disabled:opacity-50"
                              disabled={actionLoadingId === `stop:${instance.instance_id}`}
                              title="Stop"
                            >
                              {actionLoadingId === `stop:${instance.instance_id}` ? <Loader2 className="h-4 w-4 animate-spin" /> : <Square className="h-4 w-4" />}
                            </button>
                          ) : (
                            <button
                              onClick={() => handleInstanceAction('start', instance)}
                              className="p-2 text-success hover:bg-success/20 rounded transition disabled:opacity-50"
                              disabled={actionLoadingId === `start:${instance.instance_id}`}
                              title="Start"
                            >
                              {actionLoadingId === `start:${instance.instance_id}` ? <Loader2 className="h-4 w-4 animate-spin" /> : <Play className="h-4 w-4" />}
                            </button>
                          )}
                          
                          <button
                            onClick={() => handleInstanceAction('restart', instance)}
                            className="p-2 text-brand-500 hover:bg-brand-500/20 rounded transition disabled:opacity-50"
                            disabled={actionLoadingId === `restart:${instance.instance_id}`}
                            title="Restart"
                          >
                            {actionLoadingId === `restart:${instance.instance_id}` ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
                          </button>

                          {/* Logs button removed */}

                          <button
                            onClick={() => handleInstanceAction('remove', instance)}
                            className="p-2 text-danger hover:bg-danger/20 rounded transition disabled:opacity-50"
                            disabled={actionLoadingId === `remove:${instance.instance_id}`}
                            title="Remove"
                          >
                            {actionLoadingId === `remove:${instance.instance_id}` ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))
        )}
      </div>

      {/* Logs modal removed */}
    </div>
  )
}

// Logs Modal
// Render at end of component return