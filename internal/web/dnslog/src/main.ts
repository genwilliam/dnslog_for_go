import './assets/main.css';
import '@/components/dnslog/dns-log.vue';

import { createApp } from 'vue';
import App from './App.vue';

import naive from 'naive-ui';
import store from '@/store';

const app = createApp(App);

Vue.prototype.$http = http; // 组件中 使用 this.$http 使用

// 挂载 naive-ui
app.use(naive);

app.mount('#app');
