import { useState, useEffect } from 'react'
import { ArrowLeft, CheckCircle, AlertCircle, Loader } from 'lucide-react'
import { DeployApp, GetAppCredentials } from '../../wailsjs/go/api/AppsAPI';
import { useError } from '../components/ErrorContext';
import { apps } from '../../wailsjs/go/models';

interface AppSetupProps {
  appId: string
  appManifest: apps.AppManifest
  onBack: () => void
  onComplete: () => void
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

export default function AppSetup({ appId, appManifest, onBack, onComplete }: AppSetupProps) {
  const [formData, setFormData] = useState<Record<string, string>>({
    DEVICE_NAME: generateDeviceName()
  })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(false)
  const [completed, setCompleted] = useState(false)
  const [claimURL, setClaimURL] = useState<string | null>(null)
  
  const { showError } = useError();

  const handleChange = (key: string, value: string) => {
    setFormData({ ...formData, [key]: value })
    if (errors[key]) {
      delete errors[key]
      setErrors({ ...errors })
    }
  }

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {}
    
    if (appManifest.RequiredFields) {
      for (const field in appManifest.RequiredFields) {
        if (appManifest.RequiredFields[field] && !formData[field]) {
          newErrors[field] = `${field} is required`
        }
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  // Preload existing credentials if editing
  useEffect(() => {
    const loadExisting = async () => {
      try {
        const creds = await GetAppCredentials(appId)
        if (creds) {
          setFormData(prev => ({ DEVICE_NAME: prev.DEVICE_NAME, ...creds }))
        }
      } catch (e) {
        // Silently ignore if not configured yet
      }
    }
    loadExisting()
  }, [appId])

  const handleSubmit = async () => {
    if (!validate()) return

    setLoading(true)
    
    try {
      // Call backend API to deploy app
      await DeployApp(appId, formData)
      
      setLoading(false)
      setCompleted(true)
      
      setTimeout(() => {
        onComplete()
      }, 1500)
    } catch (error) {
      setLoading(false)
      console.error('Failed to deploy app:', error)
      const errorMessage = String(error)
      // If it's a Docker-related error, show specific Docker error
      if (errorMessage.includes('Please start Docker') || errorMessage.includes('Please install Docker')) {
        showError(errorMessage, 'Docker Issue')
      } else {
        showError(`Failed to deploy app: ${errorMessage}`, 'Deployment Error')
      }
    }
  }

  if (completed) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center max-w-md">
          <CheckCircle className="h-16 w-16 text-success mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-neutral-100 mb-2">App Setup Complete!</h2>
          <p className="text-neutral-400 mb-4">Your app is being deployed...</p>
          
          {/* Show claim URL for EarnApp */}
          {claimURL && (
            <div className="bg-yellow-900/20 border border-yellow-700 rounded-lg p-4 mt-6">
              <p className="text-yellow-200 text-sm font-medium mb-2">
                🔗 Claim Your Device:
              </p>
              <a 
                href={claimURL} 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-yellow-400 hover:text-yellow-300 text-sm break-all block mb-2"
              >
                {claimURL}
              </a>
              <p className="text-yellow-100 text-xs">
                Copy this link and visit it to claim your node on the EarnApp dashboard.
              </p>
            </div>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <button
        onClick={onBack}
        className="flex items-center gap-2 text-neutral-400 hover:text-white mb-6"
      >
        <ArrowLeft className="h-5 w-5" />
        Back to App Store
      </button>

      <div className="mb-6">
        <h1 className="text-3xl font-bold text-neutral-100 mb-2">Setup {appManifest.Name}</h1>
        <p className="text-neutral-400">Enter your credentials to get started</p>
      </div>

      <div className="bg-neutral-900 rounded-lg p-6 space-y-4 border border-neutral-800">
        {/* Device Name Field - Always present but optional */}
        <div>
          <label className="block text-neutral-300 mb-2">
            Device Name (Optional)
            <span className="text-neutral-500 text-xs ml-2">You can customize this or use the auto-generated name</span>
          </label>
          <input
            type="text"
            value={formData.DEVICE_NAME || ''}
            onChange={(e) => handleChange('DEVICE_NAME', e.target.value)}
            className="w-full bg-neutral-800 text-white px-4 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-brand-500 border border-neutral-700"
            placeholder="Auto-generated device name"
          />
          <p className="text-neutral-500 text-sm mt-1">
            This name will be used in the container and device identification
          </p>
        </div>

        {/* App-specific fields */}
        {appManifest.RequiredFields && Object.keys(appManifest.RequiredFields).filter(field => field !== 'DEVICE_NAME').map((field) => (
          <div key={field}>
            <label className="block text-neutral-300 mb-2">
              {field}
              {appManifest.RequiredFields[field] && <span className="text-danger ml-1">*</span>}
            </label>
            <input
              type="text"
              value={formData[field] || ''}
              onChange={(e) => handleChange(field, e.target.value)}
              className="w-full bg-neutral-800 text-white px-4 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-brand-500 border border-neutral-700"
              placeholder={`Your ${field}`}
            />
            {errors[field] && (
              <p className="text-danger text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="h-4 w-4" />
                {errors[field]}
              </p>
            )}
          </div>
        ))}

        <div className="pt-4">
          <button
            onClick={handleSubmit}
            disabled={loading}
            className="w-full bg-brand-600 hover:bg-brand-700 disabled:bg-brand-600/50 disabled:cursor-not-allowed text-white font-semibold px-6 py-3 rounded-lg flex items-center justify-center gap-2"
          >
            {loading ? (
              <>
                <Loader className="h-5 w-5 animate-spin" />
                Deploying...
              </>
            ) : (
              <>
                <CheckCircle className="h-5 w-5" />
                Deploy App
              </>
            )}
          </button>
        </div>
      </div>

      {/* Info Box */}
      <div className="bg-neutral-900 rounded-lg p-4 mt-6 border border-neutral-800">
        <p className="text-neutral-300 text-sm">
          🔒 Your credentials are encrypted and stored locally. They are only used to authenticate with {appManifest.Name}.
        </p>
      </div>
    </div>
  )
}

