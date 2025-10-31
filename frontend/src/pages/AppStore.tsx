import { useState, useEffect } from 'react'
import { Plus, ChevronRight, ExternalLink, Edit, Power, PowerOff, CheckCircle, Copy } from 'lucide-react'
import { GetConfiguredAppsForProxy, BrowserOpenURL, GetAppCredentials, DeployApp } from '../services/api';

interface AvailableApp {
  id: string
  name: string
  description: string
  icon?: string
  referralLink: string
}

interface ConfiguredApp {
  app_id: string
  app_name: string
  device_name: string
  status: string
}

export default function AppStore({ onSetup, onRefresh }: { onSetup: (appId: string) => void, onRefresh?: () => void }) {
  const [configuredApps, setConfiguredApps] = useState<ConfiguredApp[]>([])
  const [apps] = useState<AvailableApp[]>([
    {
      id: 'earnapp',
      name: 'EarnApp',
      description: 'Make money by sharing your unused internet bandwidth',
      referralLink: 'https://earnapp.com/i/vjLzkum6'
    },
    {
      id: 'honeygain',
      name: 'Honeygain',
      description: 'Share your unused internet bandwidth and earn passive income',
      referralLink: 'https://join.honeygain.com/VIETB3A80B'
    },
    {
      id: 'iproyalpawns',
      name: 'IPRoyal Pawns',
      description: 'Share your unused bandwidth and earn passive income',
      referralLink: 'https://pawns.app/?r=2651145'
    },
    {
      id: 'packetstream',
      name: 'PacketStream',
      description: 'Share your bandwidth with PacketStream and earn income',
      referralLink: 'https://packetstream.io/?psr=5pob'
    },
    {
      id: 'traffmonetizer',
      name: 'TraffMonetizer',
      description: 'Monetize your unused bandwidth and earn passive income',
      referralLink: 'https://traffmonetizer.com/?aff=1861891'
    },
    {
      id: 'repocket',
      name: 'Repocket',
      description: 'Earn by sharing your unused network resources',
      referralLink: 'https://link.repocket.com/U2fc'
    },
    {
      id: 'earnfm',
      name: 'EarnFM',
      description: 'Monetize your bandwidth and earn passive income',
      referralLink: 'https://earn.fm/ref/NAMLNPT0'
    },
    {
      id: 'proxyrack',
      name: 'ProxyRack',
      description: 'Share your network and earn income',
      referralLink: 'https://peer.proxyrack.com/ref/tazpyxdqauaa81inkwaug1qte5zkocvdt70g7syi'
    },
    {
      id: 'proxylite',
      name: 'ProxyLite',
      description: 'Earn money by sharing your bandwidth',
      referralLink: 'https://proxylite.ru/?r=VAAUXWXS'
    },
    {
      id: 'bitping',
      name: 'BitPing',
      description: 'Earn by participating in their network',
      referralLink: 'https://app.bitping.com/'
    },
    {
      id: 'packetshare',
      name: 'PacketShare',
      description: 'Share your bandwidth and earn income',
      referralLink: 'https://www.packetshare.io/?code=B6F61AADDD273776'
    },
    {
      id: 'proxybase',
      name: 'ProxyBase',
      description: 'Earn by sharing your network',
      referralLink: 'https://peer.proxybase.org?referral=ba5Zu4tnRH'
    },
    {
      id: 'wipter',
      name: 'Wipter',
      description: 'Share bandwidth and get paid',
      referralLink: 'https://wipter.com/register?via=09CED22425'
    }
    ,
    {
      id: 'mystnode',
      name: 'MystNode',
      description: 'Run a Myst node and earn',
      referralLink: 'https://mystnodes.co/?referral_code=Tc7RaS7Fm12K3Xun6mlU9q9hbnjojjl9aRBW8ZA9'
    }
  ])

  useEffect(() => {
    loadConfiguredApps()
  }, [])

  const loadConfiguredApps = async () => {
    try {
      const apps = await GetConfiguredAppsForProxy()
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
        {apps.map((app) => {
          const isConfigured = isAppConfigured(app.id)
          const configuredApp = getConfiguredApp(app.id)
          return (
            <div key={app.id} className={`bg-neutral-900 rounded-lg p-4 transition border border-neutral-800 ${isConfigured ? 'border-l-4 border-success' : 'hover:bg-neutral-800'}`}>
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold text-neutral-100">{app.name}</h3>
                    {isConfigured && (
                      <div className="flex items-center gap-2">
                        <CheckCircle className="h-5 w-5 text-success" />
                        <span className="text-success text-sm font-medium">Configured</span>
                      </div>
                    )}
                  </div>
                  <p className="text-neutral-400 text-sm mt-1">{app.description}</p>
                </div>
                <div className="flex items-center gap-2">
                  {/* Register Button */}
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleRegister(app.referralLink);
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
                          const creds = await GetAppCredentials(app.id);
                          formData = { ...creds };
                        } catch {}
                        formData.DEVICE_NAME = `${Math.random().toString(36).slice(2, 10)}-${Date.now().toString().slice(-5)}`;
                        try {
                          await DeployApp(app.id, formData);
                          if (typeof onRefresh === 'function') onRefresh();
                          alert('Started new random container');
                        } catch (err) {
                          alert('Error starting random container: ' + err);
                        }
                      }}
                    >
                      <Plus className="h-4 w-4" />
                    </button>
                  )}
                  {/* Setup/Edit Button */}
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleEdit(app.id);
                    }}
                    className={`px-4 py-2 rounded-lg flex items-center gap-2 text-sm ${isConfigured ? 'bg-neutral-700 hover:bg-neutral-600 text-white' : 'bg-brand-600 hover:bg-brand-700 text-white'}`}
                  >
                    {isConfigured ? (<><Edit className="h-4 w-4" />Edit</>) : (<>Setup<ChevronRight className="h-4 w-4" /></>)}
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
