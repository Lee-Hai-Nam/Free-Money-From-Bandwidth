import { useState } from 'react'
import React from 'react'
import Navigation from './components/Navigation'
import Dashboard from './pages/Dashboard'
import Apps from './pages/Apps'
import AppStore from './pages/AppStore'
import AppSetup from './pages/AppSetup'
import Proxy from './pages/Proxy'
import Settings from './pages/Settings'
import { BrowserOpenURL } from '../wailsjs/runtime/runtime'

interface AppConfig {
  appId: string
  name: string
  fields: Array<{
    key: string
    label: string
    required: boolean
    description: string
    type?: string
  }>
}

// Generate random device name
function generateDeviceName(): string {
  const adjectives = ['swift', 'mighty', 'quick', 'bright', 'strong', 'smart', 'cool', 'fast', 'sharp', 'bold']
  const nouns = ['eagle', 'tiger', 'hawk', 'wolf', 'bear', 'lion', 'fox', 'falcon', 'shark', 'dragon']
  const randAdj = adjectives[Math.floor(Math.random() * adjectives.length)]
  const randNoun = nouns[Math.floor(Math.random() * nouns.length)]
  const num = Math.floor(Math.random() * 10000)
  return `${randAdj}-${randNoun}-${num}`
}

const appConfigs: Record<string, AppConfig> = {
  earnapp: {
    appId: 'earnapp',
    name: 'EarnApp',
    fields: []
  },
  honeygain: {
    appId: 'honeygain',
    name: 'Honeygain',
    fields: [
      {
        key: 'HONEYGAIN_EMAIL',
        label: 'Email',
        required: true,
        description: 'Your Honeygain account email',
        type: 'email'
      },
      {
        key: 'HONEYGAIN_PASSWORD',
        label: 'Password',
        required: true,
        description: 'Your Honeygain account password',
        type: 'password'
      }
    ]
  },
  iproyalpawns: {
    appId: 'iproyalpawns',
    name: 'IPRoyal Pawns',
    fields: [
      {
        key: 'IPROYALPAWNS_EMAIL',
        label: 'Email',
        required: true,
        description: 'Your IPRoyal Pawns account email',
        type: 'email'
      },
      {
        key: 'IPROYALPAWNS_PASSWORD',
        label: 'Password',
        required: true,
        description: 'Your IPRoyal Pawns account password',
        type: 'password'
      }
    ]
  },
  packetstream: {
    appId: 'packetstream',
    name: 'PacketStream',
    fields: [
      {
        key: 'PACKETSTREAM_CID',
        label: 'CID',
        required: true,
        description: 'Your PacketStream CID',
        type: 'text'
      }
    ]
  },
  traffmonetizer: {
    appId: 'traffmonetizer',
    name: 'TraffMonetizer',
    fields: [
      {
        key: 'TRAFFMONETIZER_TOKEN',
        label: 'Token',
        required: true,
        description: 'Your TraffMonetizer token',
        type: 'text'
      }
    ]
  },
  repocket: {
    appId: 'repocket',
    name: 'Repocket',
    fields: [
      {
        key: 'REPOCKET_EMAIL',
        label: 'Email',
        required: true,
        description: 'Your Repocket account email',
        type: 'email'
      },
      {
        key: 'REPOCKET_APIKEY',
        label: 'API Key',
        required: true,
        description: 'Your Repocket API key',
        type: 'text'
      }
    ]
  },
  earnfm: {
    appId: 'earnfm',
    name: 'EarnFM',
    fields: [
      {
        key: 'EARNFM_APIKEY',
        label: 'API Key',
        required: true,
        description: 'Your EarnFM API key',
        type: 'text'
      }
    ]
  },
  proxyrack: {
    appId: 'proxyrack',
    name: 'ProxyRack',
    fields: [
      {
        key: 'PROXYRACK_APIKEY',
        label: 'API Key',
        required: true,
        description: 'Your ProxyRack API key',
        type: 'text'
      }
    ]
  },
  proxylite: {
    appId: 'proxylite',
    name: 'ProxyLite',
    fields: [
      {
        key: 'PROXYLITE_USERID',
        label: 'User ID',
        required: true,
        description: 'Your ProxyLite user ID',
        type: 'text'
      }
    ]
  },
  bitping: {
    appId: 'bitping',
    name: 'BitPing',
    fields: [
      {
        key: 'BITPING_EMAIL',
        label: 'Email',
        required: true,
        description: 'Your BitPing account email',
        type: 'email'
      },
      {
        key: 'BITPING_PASSWORD',
        label: 'Password',
        required: true,
        description: 'Your BitPing account password',
        type: 'password'
      }
    ]
  },
  packetshare: {
    appId: 'packetshare',
    name: 'PacketShare',
    fields: [
      {
        key: 'PACKETSHARE_EMAIL',
        label: 'Email',
        required: true,
        description: 'Your PacketShare account email',
        type: 'email'
      },
      {
        key: 'PACKETSHARE_PASSWORD',
        label: 'Password',
        required: true,
        description: 'Your PacketShare account password',
        type: 'password'
      }
    ]
  },
  proxybase: {
    appId: 'proxybase',
    name: 'ProxyBase',
    fields: [
      {
        key: 'PROXYBASE_USERID',
        label: 'User ID',
        required: true,
        description: 'Your ProxyBase user ID',
        type: 'text'
      }
    ]
  },
  wipter: {
    appId: 'wipter',
    name: 'Wipter',
    fields: [
      {
        key: 'WIPTER_EMAIL',
        label: 'Email',
        required: true,
        description: 'Your Wipter account email',
        type: 'email'
      },
      {
        key: 'WIPTER_PASSWORD',
        label: 'Password',
        required: true,
        description: 'Your Wipter account password',
        type: 'password'
      }
    ]
  },
  mystnode: {
    appId: 'mystnode',
    name: 'MYSTNODE',
    fields: [
      {
        key: 'DEVICE_NAME',
        label: 'Device Name',
        required: true,
        description: 'Unique name for this MystNode instance',
        type: 'text'
      },
      {
        key: 'MYSTNODE_PORT',
        label: 'Host Port',
        required: true,
        description: 'The port on your computer to map to MystNode dashboard (default: 4449)',
        type: 'number'
      }
    ]
  }
}

function App() {
  const [activeTab, setActiveTab] = useState('dashboard')
  const [setupApp, setSetupApp] = useState<string | null>(null)
  const lastTabRef = React.useRef(activeTab)

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
    if (setupApp) {
      const config = appConfigs[setupApp]
      if (config) {
        return <AppSetup appConfig={config} onBack={handleBack} onComplete={handleComplete} />
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
    <div className="flex h-screen bg-neutral-950 text-neutral-200 overflow-hidden">
      <Navigation activeTab={activeTab} onTabChange={handleTabChange} />
      <main className="flex-1 overflow-y-auto bg-neutral-900">
        <ErrorBoundary>
          {renderContent()}
        </ErrorBoundary>
      </main>
    </div>
  )
}

// Simple error boundary to prevent blank screen on render errors
class ErrorBoundary extends React.Component<{ children: React.ReactNode }, { hasError: boolean; message?: string }>{
  constructor(props: { children: React.ReactNode }) {
    super(props)
    this.state = { hasError: false }
  }
  static getDerivedStateFromError(error: any) {
    return { hasError: true, message: String(error) }
  }
  componentDidCatch(error: any, info: any) {
    console.error('Render error:', error, info)
  }
  render() {
    if (this.state.hasError) {
      return (
        <div className="p-6">
          <h2 className="text-xl font-semibold text-red-400">Something went wrong.</h2>
          <p className="text-gray-400 text-sm mt-2">{this.state.message}</p>
        </div>
      )
    }
    return this.props.children as any
  }
}

export default App
