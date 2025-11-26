<template>
	<div class="layout-container">
		<TheHeader />
		<TheToolbar />
		<main class="layout-main">
			<RouterView />
		</main>
	</div>
</template>

<script setup lang="ts">
import TheHeader from './components/TheHeader.vue';
import TheToolbar from './components/TheToolbar.vue';
import { RouterView } from 'vue-router';
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
	padding: 20px;
}
</style>
