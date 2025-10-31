import { useState, useEffect } from 'react'
import { Play, Square, RefreshCw, Trash2, FileText, Search } from 'lucide-react'

interface Container {
  id: string
  name: string
  image: string
  status: string
  created: string
  state: string
}

export default function Docker() {
  const [containers, setContainers] = useState<Container[]>([])
  const [searchTerm, setSearchTerm] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // TODO: Load containers from backend
    loadContainers()
  }, [])

  const loadContainers = async () => {
    setLoading(true)
    // TODO: Call backend API
    setLoading(false)
  }

  const filteredContainers = containers.filter(container =>
    container.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    container.image.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div className="p-6 space-y-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-neutral-100">Docker Containers</h1>
          <p className="text-neutral-400 mt-2">Manage your running containers</p>
        </div>
        <button
          onClick={loadContainers}
          className="flex items-center gap-2 px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg font-medium"
        >
          <RefreshCw className="h-5 w-5" />
          Refresh
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-neutral-400" />
        <input
          type="text"
          placeholder="Search containers by name or image..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full pl-10 pr-4 py-3 bg-neutral-800 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-brand-500 border border-neutral-700"
        />
      </div>

      {/* Containers List */}
      <div className="bg-neutral-900 rounded-lg overflow-hidden border border-neutral-800">
        {loading ? (
          <div className="p-12 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-500 mx-auto mb-4"></div>
            <p className="text-neutral-400">Loading containers...</p>
          </div>
        ) : filteredContainers.length === 0 ? (
          <div className="p-12 text-center">
            <p className="text-neutral-400 text-lg mb-2">No containers found</p>
            <p className="text-neutral-500 text-sm">
              {searchTerm ? 'Try adjusting your search terms' : 'Deploy apps from the App Store to get started'}
            </p>
          </div>
        ) : (
          <div className="divide-y divide-neutral-800">
            {filteredContainers.map((container) => (
              <div key={container.id} className="p-4 hover:bg-neutral-800/50 transition">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <h3 className="text-lg font-semibold text-neutral-100">{container.name}</h3>
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${container.state === 'running' ? 'bg-success text-white' : container.state === 'exited' ? 'bg-neutral-600 text-neutral-200' : 'bg-warning text-white'}`}>
                        {container.state}
                      </span>
                    </div>
                    <p className="text-neutral-400 text-sm mb-1">{container.image}</p>
                    <p className="text-neutral-500 text-xs">ID: {container.id.substring(0, 12)} • Created: {container.created}</p>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2 ml-4">
                    {container.state === 'running' ? (
                      <>
                        <button className="px-4 py-2 bg-neutral-700 hover:bg-neutral-600 text-white rounded-lg flex items-center gap-2">
                          <FileText className="h-4 w-4" />
                          Logs
                        </button>
                        <button className="px-4 py-2 bg-neutral-700 hover:bg-neutral-600 text-white rounded-lg flex items-center gap-2">
                          <RefreshCw className="h-4 w-4" />
                          Restart
                        </button>
                        <button className="px-4 py-2 bg-danger hover:opacity-80 text-white rounded-lg flex items-center gap-2">
                          <Square className="h-4 w-4" />
                          Stop
                        </button>
                      </>
                    ) : (
                      <button className="px-4 py-2 bg-success hover:opacity-80 text-white rounded-lg flex items-center gap-2">
                        <Play className="h-4 w-4" />
                        Start
                      </button>
                    )}
                    <button className="px-4 py-2 bg-danger hover:opacity-80 text-white rounded-lg flex items-center gap-2">
                      <Trash2 className="h-4 w-4" />
                      Remove
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Stats Summary */}
      {containers.length > 0 && (
        <div className="grid grid-cols-3 gap-4">
          <div className="bg-neutral-900 rounded-lg p-4 border border-neutral-800">
            <p className="text-neutral-400 text-sm mb-1">Total Containers</p>
            <p className="text-white text-2xl font-bold">{containers.length}</p>
          </div>
          <div className="bg-neutral-900 rounded-lg p-4 border border-neutral-800">
            <p className="text-neutral-400 text-sm mb-1">Running</p>
            <p className="text-success text-2xl font-bold">
              {containers.filter(c => c.state === 'running').length}
            </p>
          </div>
          <div className="bg-neutral-900 rounded-lg p-4 border border-neutral-800">
            <p className="text-neutral-400 text-sm mb-1">Stopped</p>
            <p className="text-danger text-2xl font-bold">
              {containers.filter(c => c.state !== 'running').length}
            </p>
          </div>
        </div>
      )}
    </div>
  )
}

