<template>
	<BaseButton type="tertiary" size="large" text="clear" @click="clearResults" />
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import BaseButton from '@/components/base/BaseButton.vue';
import { clearResults } from '@/views/result/index.ts';
</script>

<style scoped></style>
