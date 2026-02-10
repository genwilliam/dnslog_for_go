import { fetchRuntimeConfig } from '@/api/config/index.js';

export interface RuntimeConfigState {
  apiKeyRequired: boolean;
  dnsPort: string;
  rootDomain: string;
}

export default {
  namespaced: true,
  state: (): RuntimeConfigState => ({
    apiKeyRequired: true,
    dnsPort: '15353',
    rootDomain: '',
  }),
  mutations: {
    setConfig(state, payload: Partial<RuntimeConfigState>) {
      state.apiKeyRequired = payload.apiKeyRequired ?? state.apiKeyRequired;
      state.dnsPort = payload.dnsPort ?? state.dnsPort;
      state.rootDomain = payload.rootDomain ?? state.rootDomain;
    },
  },
  actions: {
    async load({ commit }) {
      try {
        const res = await fetchRuntimeConfig();
        const data = res?.data || {};
        let apiKeyRequired = true;
        if (data.apiKeyRequired !== undefined) {
          apiKeyRequired = Boolean(data.apiKeyRequired);
        } else if (data.api_key_required !== undefined) {
          apiKeyRequired = Boolean(data.api_key_required);
        }
        const dnsPort = data.dns_port || inferPort(data.dns_listen_addr) || '15353';
        commit('setConfig', {
          apiKeyRequired: Boolean(apiKeyRequired),
          dnsPort: String(dnsPort),
          rootDomain: data.root_domain || '',
        });
      } catch (err) {
        // 若失败，不影响功能，但保持默认需要 key
      }
    },
  },
};

function inferPort(listen: string | undefined): string {
  if (!listen) return '';
  const idx = listen.lastIndexOf(':');
  if (idx >= 0 && idx < listen.length - 1) {
    return listen.slice(idx + 1);
  }
  return '';
}
