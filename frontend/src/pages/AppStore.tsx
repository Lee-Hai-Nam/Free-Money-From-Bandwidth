import { useState, useEffect } from 'react'
import { Plus, ChevronRight, ExternalLink, Edit, Power, PowerOff, CheckCircle, Copy } from 'lucide-react'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { GetConfiguredApps, GetAppCredentials, DeployApp, GetAllAppManifests } from '../../wailsjs/go/api/AppsAPI';
import { useError } from '../components/ErrorContext';
import { apps } from '../../wailsjs/go/models';

interface ConfiguredApp {
  app_id: string
  app_name: string
  device_name: string
  status: string
}

export default function AppStore({ onSetup, onRefresh }: { onSetup: (appId: string) => void, onRefresh?: () => void }) {
  const [configuredApps, setConfiguredApps] = useState<ConfiguredApp[]>([])
  const { showError } = useError();
  const [apps, setApps] = useState<Record<string, apps.AppManifest>>({});

  useEffect(() => {
    loadConfiguredApps()
    GetAllAppManifests().then(setApps);
  }, [])

  const loadConfiguredApps = async () => {
    try {
      const apps = await GetConfiguredApps()
      setConfiguredApps(apps as ConfiguredApp[])
    } catch (error) {
      console.error('Failed to load configured apps:', error)
    }
  }

  const isAppConfigured = (appId: string) => {
    return configuredApps.some(app => app.app_id === appId)
  }

  const getConfiguredApp = (appId: string) => {
    return configuredApps.find(app => app.app_id === appId)
  }

  const handleRegister = (referralLink: string) => {
    BrowserOpenURL(referralLink)
  }

  const handleEdit = (appId: string) => {
    onSetup(appId)
    // Refresh will be called when returning from setup
  }

  // Refresh configured apps when component becomes visible again
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (!document.hidden) {
        loadConfiguredApps()
      }
    }
    
    document.addEventListener('visibilitychange', handleVisibilityChange)
    return () => document.removeEventListener('visibilitychange', handleVisibilityChange)
  }, [])

  return (
    <div className="p-6 space-y-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-neutral-100">App Store</h1>
        <p className="text-neutral-400 mt-2">Select and configure bandwidth sharing apps</p>
      </div>

      {/* App List */}
      <div className="space-y-3">
        {Object.entries(apps).map(([id, app]) => {
          const isConfigured = isAppConfigured(id)
          const configuredApp = getConfiguredApp(id)
          return (
            <div key={id} className={`bg-neutral-900 rounded-lg p-4 transition border border-neutral-800 ${isConfigured ? 'border-l-4 border-success' : 'hover:bg-neutral-800'}`}>
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold text-neutral-100">{app.Name}</h3>
                    {isConfigured && (
                      <div className="flex items-center gap-2">
                        <CheckCircle className="h-5 w-5 text-success" />
                        <span className="text-success text-sm font-medium">Configured</span>
                      </div>
                    )}
                  </div>
                  <p className="text-neutral-400 text-sm mt-1">{app.Name}</p>
                </div>
                <div className="flex items-center gap-2">
                  {/* Register Button */}
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleRegister(app.Link);
                    }}
                    className={`${isConfigured ? 'w-8 h-8 p-0 bg-neutral-700 hover:bg-neutral-600' : 'px-3 py-2 bg-green-600 hover:bg-green-700 text-sm'} text-white rounded-lg flex items-center justify-center`}
                    title={isConfigured ? 'Open referral link' : 'Register'}
                  >
                    <ExternalLink className="h-4 w-4" />
                    {!isConfigured && <span className="ml-1">Register</span>}
                  </button>
                  {/* Start Random Container Button */}
                  {isConfigured && (
                    <button
                      className="w-8 h-8 p-0 bg-brand-600 hover:bg-brand-700 text-white rounded-lg flex items-center justify-center"
                      title="Start new random container"
                      onClick={async (e) => {
                        e.stopPropagation();
                        // Load saved credentials
                        let formData = {};
                        try {
                          const creds = await GetAppCredentials(id);
                          formData = { ...creds };
                        } catch {}
                        formData.DEVICE_NAME = `${Math.random().toString(36).slice(2, 10)}-${Date.now().toString().slice(-5)}`;
                        try {
                          await DeployApp(id, formData);
                          if (typeof onRefresh === 'function') onRefresh();
                          // Showing success is optional; can be handled differently if desired
                        } catch (err) {
                          showError('Error starting random container: ' + err, 'Start Container Error');
                        }
                      }}
                    >
                      <Plus className="h-4 w-4" />
                    </button>
                  )}
                  {/* Setup/Edit Button */}
                  <button
                    onClick={(e) => {
                      if (!app.Disabled) {
                        e.stopPropagation();
                        handleEdit(id);
                      }
                    }}
                    disabled={app.Disabled}
                    className={`px-4 py-2 rounded-lg flex items-center gap-2 text-sm ${
                      app.Disabled 
                        ? 'bg-neutral-800 text-gray-500 cursor-not-allowed' 
                        : isConfigured 
                          ? 'bg-neutral-700 hover:bg-neutral-600 text-white' 
                          : 'bg-brand-600 hover:bg-brand-700 text-white'
                    }`}
                  >
                    {isConfigured ? (
                      <>
                        <Edit className="h-4 w-4" /> Edit
                      </>
                    ) : app.Disabled ? (
                      <>
                        Setup <ChevronRight className="h-4 w-4" /> <span className="ml-1 text-xs">[Not Working]</span>
                      </>
                    ) : (
                      <>
                        Setup <ChevronRight className="h-4 w-4" />
                      </>
                    )}
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Tip removed per request */}
    </div>
  )
}
