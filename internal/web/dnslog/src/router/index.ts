import { createRouter, createWebHistory } from 'vue-router';
import panel from '@/components/dnslog/panel/index.vue';
const router = createRouter({
  // history: createWebHistory(import.meta.env.VITE_APP_BASE_URL),
  routes: [
    {
      path: '/',
      name: 'panel',
      component: panel,
      redirect: '/login',
      children: [
        {
          path: '/user',
          name: 'user',
          component: () => import('@/views/user/index.vue'),
          meta: {
            title: '用户管理',
          },
        },
      ],
    },
  ],
});

export default router;
