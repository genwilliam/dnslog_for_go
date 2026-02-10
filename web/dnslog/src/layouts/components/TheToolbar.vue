<template>
	<div class="toolbar">
		<div class="top-row">
			<Inp class="input" />
			<Btn />
		</div>
	</div>
</template>

<script setup lang="ts">
import Inp from '@/components/dnslog/input/index.vue';
import Btn from '@/components/dnslog/panel/index.vue';
</script>

<style scoped>
.toolbar {
	padding: 16px 24px;
	background: #fff;
	border-bottom: 1px solid #e8e8e8;

	display: flex;
	flex-direction: column;
	gap: 16px;

	flex-shrink: 0;
	width: 100%;
	box-sizing: border-box;
}

.top-row {
	display: flex;
	align-items: center;
	width: 100%;
	box-sizing: border-box;
}

.top-row {
	gap: 12px;
}

.input {
	flex: 1;
}

</style>
