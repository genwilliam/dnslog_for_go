<template>
  <n-space vertical>
    <n-button type="info" size="large" @click="generateDomain">Generate</n-button>
  </n-space>
</template>

<script setup lang="ts">
import { domain } from '@/components/dnslog/input/use-domain.ts';

function generateDomain() {
  fetch('http://localhost:8080/random-domain', {
    method: 'POST',
  })
    .then((res) => res.json())
    .then((data) => {
      domain.value = data.domain; // 直接更新 ref，输入框自动刷新
      fetchDns(data.domain); // 如果你有其他逻辑
    })
    .catch((err) => {
      console.error('Error fetching domain:', err);
    });
}

// 处理生成的域名
function fetchDns(domainStr) {
  if (!domainStr) {
    console.warn('No domain provided');
    return;
  }
  console.log('Processing domain:', domainStr);
  // TODO: 其他逻辑，更新 Vuex 或发送请求
}
</script>
