declare global {
  interface Window {
    __CONFIG__?: {publicBaseUrl?: string}
  }
}

export const BaseURL = (): string => window.__CONFIG__?.publicBaseUrl ?? window.location.origin
