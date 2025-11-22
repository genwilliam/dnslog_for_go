<template>
  <n-space vertical>
    <n-button type="info" size="large" @click="handleGenerate">Generate</n-button>
  </n-space>
</template>

<script setup lang="ts">
import { domain } from '@/components/dnslog/input/use-domain';
import { generateRandomDomain } from '@/api/generate';

async function handleGenerate() {
  try {
    const res = await generateRandomDomain();

    domain.value = res.data.domain;
    console.log('获取随机域名成功：', res.data.domain);
  } catch (e) {
    console.error('获取随机域名失败：', e);
  }
}
</script>
