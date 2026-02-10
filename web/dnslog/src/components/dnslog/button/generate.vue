<template>
	<BaseButton type="info" text="Generate" @click="handleGenerate" />
</template>
<script setup lang="ts">
import BaseButton from '@/components/base/BaseButton.vue';
import { generateRandomDomain } from '@/api/generate';
import { domain } from '@/components/dnslog/input/use-domain.ts';
import { useMessage } from 'naive-ui';

const message = useMessage();

async function handleGenerate() {
	try {
		const res = await generateRandomDomain();
		if (!res?.data?.domain) {
			message.error('生成失败，请检查 API Key 或后端状态');
			return;
		}
		domain.value = res.data.domain;
	} catch (err) {
		message.error('请求失败，请检查 API Key 或网络');
	}
}
</script>
