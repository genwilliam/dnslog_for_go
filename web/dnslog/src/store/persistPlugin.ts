const KEY = 'STORE';

export default (store) => {
  if (typeof window === 'undefined') {
    return;
  }

  window.addEventListener('beforeunload', () => {
    localStorage.setItem(KEY, JSON.stringify(store.state));
  });

  const item = localStorage.getItem(KEY);
  if (!item) {
    return;
  }

  try {
    const originState = JSON.parse(item);
    const merged = { ...store.state, ...originState };

    if (originState?.text && typeof originState.text === 'object') {
      merged.text = { ...store.state.text, ...originState.text };
    }
    if (originState?.system && typeof originState.system === 'object') {
      merged.system = { ...store.state.system, ...originState.system };
    }
    if (originState?.runtimeConfig && typeof originState.runtimeConfig === 'object') {
      merged.runtimeConfig = { ...store.state.runtimeConfig, ...originState.runtimeConfig };
    }

    store.replaceState(merged);
  } catch (error) {
    console.log('restore persisted state failed');
  }
};
