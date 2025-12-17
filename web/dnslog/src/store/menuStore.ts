// 文件路径: src/stores/menuStore.ts

import { defineStore } from 'pinia';

export const useMenuStore = defineStore('menuStore', () => {
	// 定义新的菜单数据
	const menuList = [
		{ title: 'DNS query', path: '/dnsquery' }, // 用于发起新的 DNS 查询的页面
		{ title: 'DNS log', path: '/dnslog' }, // 用于展示历史/当前查询记录的页面 (您的表格组件)
	];

	return { menuList };
});
