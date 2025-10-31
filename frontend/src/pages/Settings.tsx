import { useEffect, useState } from 'react'
import { GetSettings, SetAutoStart, SetShowInTray } from '../services/api';

export default function Settings() {
  const [autoStart, setAutoStart] = useState(false)
  const [showInTray, setShowInTray] = useState(true)

  useEffect(() => {
    const load = async () => {
      try {
        const cfg: any = await (GetSettings as any)()
        setAutoStart(!!cfg?.auto_start)
        setShowInTray(cfg?.show_in_tray !== false)
      } catch {}
    }
    load()
  }, [])

  const onToggleAutoStart = async (val: boolean) => {
    setAutoStart(val)
    try { await (SetAutoStart as any)(val) } catch {}
  }

  const onToggleShowInTray = async (val: boolean) => {
    setShowInTray(val)
    try { await (SetShowInTray as any)(val) } catch {}
  }

  return (
    <div className="p-6 space-y-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-neutral-100">Settings</h1>
        <p className="text-neutral-400 mt-2">Configure your application preferences</p>
      </div>

      <div className="space-y-6">
        {/* General Settings */}
        <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
          <h2 className="text-xl font-semibold text-neutral-100 mb-4">General</h2>
          
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <label className="block text-neutral-300">Auto-start on boot</label>
              <input type="checkbox" className="w-5 h-5 rounded-md bg-neutral-800 border-neutral-700 text-brand-600 focus:ring-brand-500" checked={autoStart} onChange={e => onToggleAutoStart(e.target.checked)} />
            </div>
            
            <div className="flex items-center justify-between">
              <label className="block text-neutral-300">Show in system tray</label>
              <input type="checkbox" className="w-5 h-5 rounded-md bg-neutral-800 border-neutral-700 text-brand-600 focus:ring-brand-500" checked={showInTray} onChange={e => onToggleShowInTray(e.target.checked)} />
            </div>
          </div>
        </div>

        {/* Docker Settings */}
        <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
          <h2 className="text-xl font-semibold text-neutral-100 mb-4">Docker</h2>
          
          <div className="space-y-4">
            <div>
              <label className="block text-neutral-300 mb-2">Docker Host</label>
              <input 
                type="text" 
                defaultValue="unix:///var/run/docker.sock"
                className="w-full bg-neutral-800 text-neutral-300 px-4 py-2 rounded-lg border border-neutral-700"
                disabled
              />
              <p className="text-neutral-500 text-sm mt-1">Using local Docker daemon</p>
            </div>
          </div>
        </div>

        {/* Proxy Settings - simplified (toggle removed) */}
        <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
          <h2 className="text-xl font-semibold text-neutral-100 mb-4">Proxy Configuration</h2>
          <p className="text-neutral-400">Proxies are managed from the Proxy tab.</p>
        </div>
      </div>
    </div>
  )
}

