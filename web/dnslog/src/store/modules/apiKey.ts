import {
  getEnvApiKey,
  getStoredApiKey,
  setStoredApiKey,
  clearStoredApiKey,
  isValidApiKeyFormat,
} from '@/utils/apiKey';

export type ApiKeySource = 'local' | 'env' | 'none';

export default {
  namespaced: true,
  state: () => ({
    value: '',
    source: 'none' as ApiKeySource,
  }),
  mutations: {
    setApiKey(state, payload: { value: string; source: ApiKeySource }) {
      state.value = payload.value;
      state.source = payload.source;
    },
    clearApiKey(state) {
      state.value = '';
      state.source = 'none';
    },
  },
  actions: {
    load({ commit }) {
      const stored = getStoredApiKey();
      if (stored) {
        if (isValidApiKeyFormat(stored)) {
          commit('setApiKey', { value: stored, source: 'local' });
          return stored;
        }
        clearStoredApiKey();
      }
      const envKey = getEnvApiKey();
      if (envKey && isValidApiKeyFormat(envKey)) {
        commit('setApiKey', { value: envKey, source: 'env' });
        return envKey;
      }
      commit('clearApiKey');
      return '';
    },
    set({ commit }, key: string) {
      const trimmed = key.trim();
      if (!trimmed) {
        clearStoredApiKey();
        commit('clearApiKey');
        return '';
      }
      if (!isValidApiKeyFormat(trimmed)) {
        clearStoredApiKey();
        commit('clearApiKey');
        return '';
      }
      setStoredApiKey(trimmed);
      commit('setApiKey', { value: trimmed, source: 'local' });
      return trimmed;
    },
    clear({ commit }) {
      clearStoredApiKey();
      commit('clearApiKey');
    },
  },
};
