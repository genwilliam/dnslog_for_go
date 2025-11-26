<template>
	<div class="toolbar">
		<div class="top-elements">
			<Inp class="input" />
			<queryBtn />
			<button @click="dnsRecords()">ddd</button>
		</div>

		<div class="bottom-elements">
			<Btn />
		</div>
	</div>
</template>

<script setup lang="ts">
import Inp from '@/components/dnslog/input/index.vue';
import Btn from '@/components/dnslog/panel/index.vue';
import queryBtn from '@/components/dnslog/button/query.vue';
import { dnsRecords } from '@/api/records/index.js';
</script>

<style scoped>
.toolbar {
	padding: 16px 24px;
	background: #fff;
	border-bottom: 1px solid #e8e8e8;

	display: flex;
	flex-direction: column; /* 垂直堆叠顶部和底部元素 */
	gap: 16px;

	flex-shrink: 0;
	width: 100%; /* 确保工具栏本身撑满父容器宽度 */
	box-sizing: border-box; /* 保证 padding 不会超出 100% 宽度 */
}

/* 关键：确保顶部和底部元素能够横向撑满整个工具栏 */
.top-elements,
.bottom-elements {
	display: flex;
	width: 100%;
	box-sizing: border-box;
	/* 确保它们内部的元素能正确排列 */
}

.top-elements {
	align-items: center;
	gap: 12px;
}

.input {
	/* 确保输入框占据所有可用空间 */
	flex: 1;
}

.bottom-elements {
	justify-content: flex-start;
}
</style>
