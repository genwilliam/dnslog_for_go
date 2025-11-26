<template>
	<header class="layout-header">
		<div class="brand-logo">
			<span class="icon">📡</span>
			<span class="text">DNS Logger</span>
		</div>
	</header>
</template>

<script setup lang="ts"></script>

<style scoped>
.layout-header {
	height: 56px;
	background-color: #001529;
	color: #fff;
	display: flex;
	align-items: center; /* 垂直居中 (Y轴) */
	justify-content: center; /* <-- 修改这里：水平居中 (X轴) */
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
	cursor: default;
}
</style>
