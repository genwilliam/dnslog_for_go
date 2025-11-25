<template>
	<!-- <n-space>
    <n-button type="error" size="large" @click="togglePolling">{{ state }}</n-button>
  </n-space> -->
	<BaseButton type="error" size="large" :text="state" @click="togglePolling" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import fetchDns from '@/views/dns-query/index.ts';
import { startPolling, stopPolling } from '@/utils/pulling.ts';
import BaseButton from '@/components/base/BaseButton.vue';
// 轮询状态控制
const isPolling = ref(true);
const state = ref('Pause');

// 切换按钮状态
function togglePolling() {
	isPolling.value = !isPolling.value;
	state.value = isPolling.value ? 'Pause' : 'Start';

	if (isPolling.value) {
		startPolling();
	} else {
		stopPolling();
	}
}
</script>
