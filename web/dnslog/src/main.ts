import './assets/main.css';

import router from './router';
import App from './App.vue';
import { createApp } from 'vue';

import naive from 'naive-ui';
import store from '@/store';
import { loadApiKey } from '@/stores/auth';

const app = createApp(App);
app.use(store);
app.use(router);

// 挂载 naive-ui
app.use(naive);

loadApiKey();
store.dispatch('runtimeConfig/load');

app.mount('#app');
