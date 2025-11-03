import { useState, useEffect } from 'react';
import { IsDockerAvailable, GetDockerInstallationURL } from '../../wailsjs/go/api/AppsAPI';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { AlertTriangle, ExternalLink } from 'lucide-react';

export default function DockerStatus() {
  const [dockerStatus, setDockerStatus] = useState<string>('checking');
  const [dockerError, setDockerError] = useState<string>('');
  const [installURL, setInstallURL] = useState<string>('https://www.docker.com/products/docker-desktop');

  useEffect(() => {
    const checkDocker = async () => {
      try {
        const result = await IsDockerAvailable();
        setDockerStatus(result.status);
        if (result.status === 'not_installed') {
          const url = await GetDockerInstallationURL();
          setInstallURL(url);
        }
      } catch (error) {
        setDockerStatus('error');
        setDockerError(String(error));
      }
    };

    checkDocker();
  }, []);

  if (dockerStatus === 'running' || dockerStatus === 'checking') {
    return null;
  }

  const handleInstallClick = () => {
    BrowserOpenURL(installURL);
  };

  return (
    <div className="bg-yellow-900/50 border-b border-yellow-700 text-yellow-200 p-3 text-sm">
      <div className="container mx-auto flex items-center justify-center gap-3">
        <AlertTriangle className="h-5 w-5 text-yellow-400" />
        <div className="flex-grow">
          {dockerStatus === 'not_installed' && (
            <span>
              Docker is not installed. Please install Docker to use this application.
            </span>
          )}
          {dockerStatus === 'not_running' && (
            <span>
              Docker is not running. Please start Docker to manage your apps.
            </span>
          )}
          {dockerStatus === 'error' && (
            <span>
              Error checking Docker status: {dockerError}
            </span>
          )}
        </div>
        {dockerStatus === 'not_installed' && (
          <button 
            onClick={handleInstallClick} 
            className="flex items-center gap-1.5 bg-yellow-600/50 hover:bg-yellow-600 text-white font-semibold px-3 py-1 rounded-md transition-colors text-xs"
          >
            <ExternalLink className="h-4 w-4" />
            Installation Instructions
          </button>
        )}
      </div>
    </div>
  );
}
