import { createRouter, createWebHistory } from 'vue-router';
import Layouts from '@/layouts/index.vue';

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			component: Layouts,
			redirect: '/dnsquery',
			children: [
				{
					path: 'dnsquery',
					name: 'dnsquery',
					component: () => import('@/views/dns-query/index.vue'),
				},
				{
					path: 'dns_log',
					name: 'dnslog',
					component: () => import('@/views/dnslog/index.vue'),
				},
			],
		},
	],
});

export default router;
