const STORAGE_KEY = 'DNSLOG_API_KEY';
const LEGACY_KEY = 'api_key';

export function getEnvApiKey(): string {
  return import.meta.env.VITE_API_KEY || '';
}

export function getStoredApiKey(): string {
  if (typeof window === 'undefined') {
    return '';
  }
  return localStorage.getItem(STORAGE_KEY) || localStorage.getItem(LEGACY_KEY) || '';
}

export function getApiKey(): string {
  const stored = getStoredApiKey();
  if (stored) {
    return stored;
  }
  return getEnvApiKey();
}

export function maskApiKey(key: string, prefix = 6): string {
  if (!key) {
    return '';
  }
  const head = key.slice(0, prefix);
  return `${head}...(len=${key.length})`;
}

export function isValidApiKeyFormat(key: string): boolean {
  if (!key || key.length != 64) {
    return false;
  }
  for (let i = 0; i < key.length; i += 1) {
    const ch = key[i];
    const isHex =
      (ch >= '0' && ch <= '9') ||
      (ch >= 'a' && ch <= 'f') ||
      (ch >= 'A' && ch <= 'F');
    if (!isHex) {
      return false;
    }
  }
  return true;
}

export function setStoredApiKey(key: string) {
  if (typeof window === 'undefined') {
    return;
  }
  localStorage.setItem(STORAGE_KEY, key);
}

export function clearStoredApiKey() {
  if (typeof window === 'undefined') {
    return;
  }
  localStorage.removeItem(STORAGE_KEY);
}

export function setApiKey(key: string) {
  setStoredApiKey(key);
}

export function clearApiKey() {
  clearStoredApiKey();
}

export const ApiKeyStorageKey = STORAGE_KEY;
export const ApiKeyLegacyStorageKey = LEGACY_KEY;
