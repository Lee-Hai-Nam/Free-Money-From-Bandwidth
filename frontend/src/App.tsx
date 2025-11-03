import { useState, useEffect } from 'react'
import React from 'react'
import Navigation from './components/Navigation'
import Dashboard from './pages/Dashboard'
import Apps from './pages/Apps'
import AppStore from './pages/AppStore'
import AppSetup from './pages/AppSetup'
import Proxy from './pages/Proxy'
import Settings from './pages/Settings'
import { BrowserOpenURL } from '../wailsjs/runtime/runtime'

import { GetAllAppManifests } from '../wailsjs/go/api/AppsAPI';
import { apps } from '../wailsjs/go/models';

import DockerStatus from './components/DockerStatus';

function App() {
  const [activeTab, setActiveTab] = useState('dashboard')
  const [setupApp, setSetupApp] = useState<string | null>(null)
  const [appManifests, setAppManifests] = useState<Record<string, apps.AppManifest>>({});
  const lastTabRef = React.useRef(activeTab)

  useEffect(() => {
    GetAllAppManifests().then(setAppManifests);
  }, []);

  const handleSetup = (appId: string) => {
    setSetupApp(appId)
  }

  const handleBack = () => {
    setSetupApp(null)
  }

  const handleComplete = () => {
    setSetupApp(null)
    setActiveTab('apps')
  }

  const handleRefresh = () => {
    // Force re-render of AppStore to refresh configured apps
    setActiveTab(activeTab)
  }

  const handleTabChange = (tab: string) => {
    if (tab === 'support') {
      BrowserOpenURL('https://discord.gg/h4cvRHCwRF')
      setActiveTab(lastTabRef.current)
      return
    }
    lastTabRef.current = tab
    if (tab !== 'appstore' && setupApp) setSetupApp(null)
    setActiveTab(tab)
  }

  const renderContent = () => {
    console.log("setupApp:", setupApp);
    console.log("appManifests:", appManifests);
    if (setupApp) {
      const manifest = appManifests[setupApp]
      if (manifest) {
        return <AppSetup appId={setupApp} appManifest={manifest} onBack={handleBack} onComplete={handleComplete} />
      }
      return <AppStore onSetup={handleSetup} onRefresh={handleRefresh} />
    }

    switch (activeTab) {
      case 'dashboard':
        return <Dashboard />
      case 'apps':
        return <Apps onAddNew={() => setActiveTab('appstore')} />
      case 'appstore':
        return <AppStore onSetup={handleSetup} onRefresh={handleRefresh} />
      case 'proxy':
        return <Proxy />
      case 'settings':
        return <Settings />
      default:
        return <Dashboard />
    }
  }

  return (
    <div className="flex h-screen bg-neutral-950 text-white">
      <Navigation activeTab={activeTab} onTabChange={handleTabChange} />
      <main className="flex-1 flex flex-col overflow-hidden">
        <DockerStatus />
        <div className="flex-1 overflow-y-auto">
          {renderContent()}
        </div>
      </main>
    </div>
  )
}

export default App
