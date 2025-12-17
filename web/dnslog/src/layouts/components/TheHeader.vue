<template>
	<header class="layout-header">
		<div class="brand-logo" @click="goHome">
			<span class="icon">📡</span>
			<span class="text">DNS Logger</span>
		</div>

		<nav class="nav">
			<button
				v-for="item in menus"
				:key="item.path"
				class="nav-item"
				:class="{ active: isActive(item) }"
				type="button"
				@click="switchTo(item)"
			>
				{{ item.label }}
			</button>
		</nav>
	</header>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';

const router = useRouter();
const route = useRoute();

const menus = [
	{ label: 'DNS query', path: '/dnsquery' },
	{ label: 'DNS log', path: '/dns_log' },
];

const isActive = (item: { path: string }) => {
	// 当前路径以菜单 path 开头则视为激活
	return route.path === item.path;
};

const switchTo = (item: { path: string }) => {
	if (route.path !== item.path) {
		router.push(item.path);
	}
};

const goHome = () => {
	router.push('/dnsquery');
};
</script>

<style scoped>
.layout-header {
	height: 56px;
	background-color: #001529;
	color: #fff;
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0 24px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
	flex-shrink: 0;
	z-index: 20;
}

.brand-logo {
	font-size: 18px;
	font-weight: 600;
	display: flex;
	align-items: center;
	gap: 10px;
	letter-spacing: 1px;
	user-select: none;
	cursor: pointer;
}

.nav {
	display: flex;
	align-items: center;
	gap: 12px;
}

.nav-item {
	border: none;
	outline: none;
	padding: 6px 14px;
	border-radius: 16px;
	background: transparent;
	color: #d9d9d9;
	font-size: 14px;
	cursor: pointer;
	transition:
		background-color 0.2s ease,
		color 0.2s ease;
}

.nav-item:hover {
	background: rgba(255, 255, 255, 0.12);
	color: #ffffff;
}

.nav-item.active {
	background: #1890ff;
	color: #ffffff;
}
</style>
