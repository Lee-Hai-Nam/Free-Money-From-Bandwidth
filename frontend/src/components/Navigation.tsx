import { LayoutDashboard, Grid3x3, Store, Settings, Globe, ExternalLink } from 'lucide-react'

interface NavigationProps {
  activeTab: string
  onTabChange: (tab: string) => void
}

export default function Navigation({ activeTab, onTabChange }: NavigationProps) {
  const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { id: 'apps', label: 'My Apps', icon: Grid3x3 },
    { id: 'appstore', label: 'App Store', icon: Store },
    { id: 'proxy', label: 'Proxy', icon: Globe },
    { id: 'settings', label: 'Settings', icon: Settings },
    { id: 'support', label: 'Support', icon: ExternalLink },
  ]

  return (
    <nav className="bg-neutral-950 border-r border-neutral-800 w-64 p-4 flex flex-col">
      <div className="mb-8">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-brand-400 to-brand-600 bg-clip-text text-transparent">
          Free Income From Bandwidth
        </h1>
      </div>

      <ul className="space-y-2 flex-grow">
        {navItems.map((item) => {
          const isActive = activeTab === item.id
          
          return (
            <li key={item.id}>
              <button
                onClick={() => onTabChange(item.id)}
                className={`w-full flex items-center gap-3 px-4 py-2 rounded-md transition-colors ${
                  isActive
                    ? 'bg-brand-600 text-white'
                    : 'text-neutral-400 hover:bg-neutral-800 hover:text-neutral-100'
                }`}
              >
                <item.icon className="h-5 w-5" />
                <span className="font-medium text-sm">{item.label}</span>
              </button>
            </li>
          )
        })}
      </ul>

      <div className="mt-auto">
        <p className="text-xs text-neutral-600 text-center">Version 0.1.0</p>
      </div>
    </nav>
  )
}
