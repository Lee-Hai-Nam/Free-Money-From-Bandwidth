import { useState, useEffect } from 'react'
import { Activity, TrendingUp, Server, Clock } from 'lucide-react'
import { GetDashboardSummary } from '../services/api';

export default function Dashboard() {
  const [activeApps, setActiveApps] = useState(0)
  const [bandwidth, setBandwidth] = useState('0 MB')
  const [uptime, setUptime] = useState('0h 0m')
  const [recent, setRecent] = useState<string[]>([])

  useEffect(() => {
    const load = async () => {
      try {
        const data: any = await (GetDashboardSummary as any)()
        setActiveApps(data?.active_apps ?? 0)
        const bw = data?.bandwidth_used ?? 0
        setBandwidth(`${bw} MB`)
        const secs = data?.uptime_seconds ?? 0
        const h = Math.floor(secs / 3600)
        const m = Math.floor((secs % 3600) / 60)
        setUptime(`${h}h ${m}m`)
        setRecent(Array.isArray(data?.recent_activity) ? data.recent_activity.slice().reverse().slice(0, 10) : [])
      } catch {}
    }
    load()
    const id = setInterval(load, 5000)
    return () => clearInterval(id)
  }, [])

  return (
    <div className="p-6 space-y-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-neutral-100">Dashboard</h1>
        <p className="text-neutral-400 mt-2">Monitor your passive income</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-neutral-400 text-sm font-medium">Active Apps</p>
              <p className="text-brand-500 text-3xl font-bold mt-2">{activeApps}</p>
            </div>
            <Activity className="h-12 w-12 text-brand-500 opacity-30" />
          </div>
        </div>

        <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-neutral-400 text-sm font-medium">Bandwidth Used</p>
              <p className="text-brand-500 text-2xl font-bold mt-2">{bandwidth}</p>
            </div>
            <TrendingUp className="h-12 w-12 text-brand-500 opacity-30" />
          </div>
        </div>

        <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-neutral-400 text-sm font-medium">Uptime</p>
              <p className="text-brand-500 text-2xl font-bold mt-2">{uptime}</p>
            </div>
            <Server className="h-12 w-12 text-brand-500 opacity-30" />
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-neutral-900 rounded-lg p-6 border border-neutral-800">
        <h2 className="text-xl font-semibold text-neutral-100 mb-4">Recent Activity</h2>
        <div className="space-y-2">
          {recent.length === 0 ? (
            <div className="text-neutral-400 text-sm">No recent activity yet.</div>
          ) : (
            recent.map((line, idx) => (
              <div key={idx} className="text-neutral-300 text-sm flex items-center gap-2">
                <Clock className="h-4 w-4 text-neutral-500" />
                <span>{line}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}

