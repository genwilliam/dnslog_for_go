import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueDevTools from 'vite-plugin-vue-devtools';

export default defineConfig({
  plugins: [vue(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: true,
    open: true,
    proxy: {
      '/dnslog': {
        target: 'http://localhost:8080/', // 后端地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''), // 去掉前缀
      },
    },
  },
});
