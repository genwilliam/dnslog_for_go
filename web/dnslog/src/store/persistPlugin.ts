const KEY = 'STORE';
// // 环境变量：若未定义，则默认启用持久化
// const enablePersistence = import.meta.env.VITE_ENABLE_PERSISTENCE !== 'false';
export default (store) => {
  window.addEventListener('beforeunload', () => {
    localStorage.setItem(KEY, JSON.stringify(store.state));
  });

  const item = localStorage.getItem(KEY);
  if (!item) {
    return;
  }

  try {
    const originState = JSON.parse(item);
    store.replaceState(originState);
  } catch (error) {
    console.log('error');
  }
};
