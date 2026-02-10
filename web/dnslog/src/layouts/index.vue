<template>
	<div class="layout-container">
		<TheHeader />
		<ApiKeyBanner />
		<!-- 仅在 DNS 查询页面显示工具栏，日志页不显示查询控件 -->
		<TheToolbar v-if="showToolbar" />
		<main class="layout-main">
			<RouterView />
		</main>
	</div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import TheHeader from './components/TheHeader.vue';
import TheToolbar from './components/TheToolbar.vue';
import ApiKeyBanner from './components/ApiKeyBanner.vue';
import { RouterView } from 'vue-router';

const route = useRoute();
const showToolbar = computed(() => route.name === 'dnsquery' || route.path === '/dnsquery');
</script>

<style scoped>
.layout-container {
	display: flex;
	flex-direction: column;
	height: 100vh;
	width: 100vw;
	overflow: hidden;
	background: #f0f2f5;
}

/* RouterView 显示区域 */
.layout-main {
	flex: 1;
	overflow: auto;
}
</style>
