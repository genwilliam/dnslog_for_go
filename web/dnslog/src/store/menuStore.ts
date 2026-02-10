// 文件路径: src/stores/menuStore.ts

import { defineStore } from 'pinia';

export const useMenuStore = defineStore('menuStore', () => {
	// 定义新的菜单数据
	const menuList = [
		{ title: 'Observe', path: '/observe' }, // 观测/令牌详情
		{ title: 'DNS log', path: '/dns_log' }, // 历史/当前记录
		{ title: 'Tokens', path: '/tokens' }, // token 列表
		{ title: 'Security', path: '/security' }, // API Keys 与黑名单
	];

	return { menuList };
});
