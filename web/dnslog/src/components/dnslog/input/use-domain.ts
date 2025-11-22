// import { ref, watch } from 'vue';

// export const KEY = 'DOMAIN_KEY';

// export const domain = ref('');

// // 初始化 domain
// if (typeof window !== 'undefined') {
//   domain.value = localStorage.getItem(KEY) || '';
// }

// // watch 自动同步 localStorage
// watch(domain, (newValue) => {
//   if (typeof window !== 'undefined') {
//     localStorage.setItem(KEY, newValue);
//   }
// });

// export function submitDomain() {
//   return domain.value;
// }

// export default {
//   domain,
//   submitDomain,
//   KEY,
// };

import { ref, watch } from 'vue';
import store from '@/store';

const initValue = store.state.text?.domain ?? ''; // 加了 ?. 安全访问

export const domain = ref(initValue);

watch(domain, (newVal) => {
  if (store.state.text) {
    // 确保模块已初始化
    store.commit('text/setDomain', newVal);
  }
});

export function submitDomain() {
  return domain.value;
}
