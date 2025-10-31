import * as wailsApps from '../../wailsjs/go/api/AppsAPI';
import * as wailsProxy from '../../wailsjs/go/api/ProxyAPI';
import * as wailsRuntime from '../../wailsjs/runtime/runtime';
import * as wailsSettings from '../../wailsjs/go/api/SettingsAPI';

// This is a conditional API client that uses the Wails bindings when available,
// and falls back to HTTP requests when running in a browser.

declare global {
  interface Window {
    wails: any;
  }
}

export async function GetRunningApps(): Promise<any> {
  if (window.wails) {
    return wailsApps.GetRunningApps();
  } else {
    const response = await fetch('/api/apps/running');
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function StartApp(appID: string): Promise<void> {
  if (window.wails) {
    return wailsApps.StartApp(appID);
  } else {
    const response = await fetch(`/api/apps/start/${appID}`, { method: 'POST' });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function StopApp(appID: string): Promise<void> {
  if (window.wails) {
    return wailsApps.StopApp(appID);
  } else {
    const response = await fetch(`/api/apps/stop/${appID}`, { method: 'POST' });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function RestartApp(appID: string): Promise<void> {
  if (window.wails) {
    return wailsApps.RestartApp(appID);
  } else {
    const response = await fetch(`/api/apps/restart/${appID}`, { method: 'POST' });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function RemoveApp(containerID: string): Promise<void> {
  if (window.wails) {
    return wailsApps.RemoveApp(containerID);
  } else {
    const response = await fetch(`/api/apps/remove/${containerID}`, { method: 'POST' });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function GetConfiguredAppsForProxy(): Promise<any> {
  if (window.wails) {
    return wailsProxy.GetConfiguredAppsForProxy();
  } else {
    const response = await fetch('/api/proxies/configured-apps');
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function GetAppCredentials(appID: string): Promise<any> {
  if (window.wails) {
    return wailsApps.GetAppCredentials(appID);
  } else {
    const response = await fetch(`/api/apps/credentials/${appID}`);
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function DeployApp(appID: string, formData: any): Promise<void> {
  if (window.wails) {
    return wailsApps.DeployApp(appID, formData);
  } else {
    const response = await fetch(`/api/apps/deploy/${appID}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function GetDashboardSummary(): Promise<any> {
  if (window.wails) {
    return wailsApps.GetDashboardSummary();
  } else {
    const response = await fetch('/api/apps/summary');
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export function BrowserOpenURL(url: string): void {
  if (window.wails) {
    wailsRuntime.BrowserOpenURL(url);
  } else {
    window.open(url, '_blank');
  }
}

export async function AddProxy(proxyStr: string, autoDeploy: boolean, selectedAppIDs: string[]): Promise<any> {
  if (window.wails) {
    return wailsProxy.AddProxy(proxyStr, autoDeploy, selectedAppIDs);
  } else {
    const response = await fetch('/api/proxies/add', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ proxyStr, autoDeploy, selectedAppIDs }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function ListProxies(): Promise<any> {
  if (window.wails) {
    return wailsProxy.ListProxies();
  } else {
    const response = await fetch('/api/proxies/list');
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function RemoveProxy(proxyID: string): Promise<any> {
  if (window.wails) {
    return wailsProxy.RemoveProxy(proxyID);
  } else {
    const response = await fetch(`/api/proxies/remove/${proxyID}`);
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function ConfirmRemoveProxy(proxyID: string): Promise<void> {
  if (window.wails) {
    return wailsProxy.ConfirmRemoveProxy(proxyID);
  } else {
    const response = await fetch(`/api/proxies/remove/confirm/${proxyID}`, { method: 'POST' });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function TestProxy(proxyID: string): Promise<any> {
  if (window.wails) {
    return wailsProxy.TestProxy(proxyID);
  } else {
    const response = await fetch(`/api/proxies/test/${proxyID}`);
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function GetAppsRunningOnProxies(proxyIDs: string[]): Promise<any> {
  if (window.wails) {
    return wailsProxy.GetAppsRunningOnProxies(proxyIDs);
  } else {
    const response = await fetch('/api/proxies/apps-running', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(proxyIDs),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function DeployAppWithProxy(appId: string, proxyUrl: string): Promise<any> {
  if (window.wails) {
    // The original wails function is DeployAppWithProxy(deploymentData map[string]interface{}) but the frontend calls it with (appId, proxyUrl)
    // I will assume the original implementation is wrong and the frontend is right.
    // I will create a new function in the backend to handle this.
    // For now, I will just call the http endpoint.
  }
  const response = await fetch('/api/apps/deploy-with-proxy', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ app_id: appId, proxy_url: proxyUrl }),
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }
  return await response.json();
}

export async function GetSettings(): Promise<any> {
  if (window.wails) {
    return wailsSettings.GetSettings();
  } else {
    const response = await fetch('/api/settings');
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    return await response.json();
  }
}

export async function SetAutoStart(enabled: boolean): Promise<void> {
  if (window.wails) {
    return wailsSettings.SetAutoStart(enabled);
  } else {
    const response = await fetch('/api/settings/autostart', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ enabled }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}

export async function SetShowInTray(enabled: boolean): Promise<void> {
  if (window.wails) {
    return wailsSettings.SetShowInTray(enabled);
  } else {
    const response = await fetch('/api/settings/showintray', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ enabled }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
  }
}
