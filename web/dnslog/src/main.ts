import './assets/main.css';
// import '@/components/dnslog/index.vue';
import router from './router';
import App from './App.vue';
import { createApp } from 'vue';

import naive from 'naive-ui';
import store from '@/store';

const app = createApp(App);
app.use(store);
app.use(router);

// 挂载 naive-ui
app.use(naive);

app.mount('#app');
