<template>
	<BaseButton type="info" text="Generate" @click="handleGenerate" />
</template>

<script setup lang="ts">
import BaseButton from '@/components/base/BaseButton.vue';
import { generateRandomDomain } from '@/api/generate';
import { domain } from '@/components/dnslog/input/use-domain.ts';

async function handleGenerate() {
	const res = await generateRandomDomain();
	domain.value = res.data.domain;
}
</script>
