<template>
  <n-space vertical>
    <!-- 提交按钮 -->
    <n-button type="success" size="large" @click="startPolling"> Record DNS query </n-button>
  </n-space>
</template>

<script setup lang="ts">
import { startPolling, stopPolling } from '@/utils/pulling';
</script>
