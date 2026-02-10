import { createRouter, createWebHistory } from 'vue-router';
import Layouts from '@/layouts/index.vue';

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			component: Layouts,
			redirect: '/observe',
			children: [
				{
					path: 'observe',
					name: 'observe',
					component: () => import('@/views/dns-query/index.vue'),
				},
				{
					path: 'dns_log',
					name: 'dnslog',
					component: () => import('@/views/dnslog/index.vue'),
				},
				{
					path: 'tokens',
					name: 'tokens',
					component: () => import('@/views/tokens/index.vue'),
				},
				{
					path: 'security',
					name: 'security',
					component: () => import('@/views/security/index.vue'),
				},
			],
		},
	],
});

export default router;
