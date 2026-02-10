import store from '@/store';
import {
  getEnvApiKey,
  getStoredApiKey,
  setStoredApiKey,
  clearStoredApiKey,
  isValidApiKeyFormat,
} from '@/utils/apiKey';

export type ApiKeySource = 'local' | 'env' | 'none';

const resolveApiKey = () => {
  const stored = getStoredApiKey();
  if (stored && isValidApiKeyFormat(stored)) {
    return { value: stored, source: 'local' as ApiKeySource };
  }
  const envKey = getEnvApiKey();
  if (envKey && isValidApiKeyFormat(envKey)) {
    return { value: envKey, source: 'env' as ApiKeySource };
  }
  return { value: '', source: 'none' as ApiKeySource };
};

export const loadApiKey = () => {
  if (store?.dispatch) {
    return store.dispatch('apiKey/load');
  }
  const resolved = resolveApiKey();
  if (resolved.value && resolved.source === 'local') {
    setStoredApiKey(resolved.value);
  }
  return resolved.value;
};

export const getApiKey = (): string => {
  const current = store?.state?.apiKey?.value;
  if (current) {
    return current;
  }
  const resolved = resolveApiKey();
  if (resolved.value && store?.commit) {
    store.commit('apiKey/setApiKey', resolved);
  }
  return resolved.value;
};

export const getApiKeySource = (): ApiKeySource => {
  return store?.state?.apiKey?.source || 'none';
};

export const setApiKey = (key: string) => {
  const trimmed = key.trim();
  if (!trimmed) {
    clearStoredApiKey();
    if (store?.commit) {
      store.commit('apiKey/clearApiKey');
    }
    return '';
  }
  if (!isValidApiKeyFormat(trimmed)) {
    clearStoredApiKey();
    if (store?.commit) {
      store.commit('apiKey/clearApiKey');
    }
    return '';
  }
  setStoredApiKey(trimmed);
  if (store?.commit) {
    store.commit('apiKey/setApiKey', { value: trimmed, source: 'local' as ApiKeySource });
  }
  return trimmed;
};

export const clearApiKey = () => {
  clearStoredApiKey();
  if (store?.commit) {
    store.commit('apiKey/clearApiKey');
  }
};
