<template>
	<BaseButton type="success" size="large" text="query" @click="startPolling" />
</template>

<script setup lang="ts">
import { startPolling, stopPolling } from '@/utils/pulling.ts';
import BaseButton from '@/components/base/BaseButton.vue';
</script>
