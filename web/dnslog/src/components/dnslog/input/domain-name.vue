<template>
  <n-space vertical>
    <n-input size="large" round placeholder="Please enter the domain name" v-model:value="domain" />
  </n-space>
</template>

<script setup lang="ts">
import { domain } from './use-domain.ts';
</script>
