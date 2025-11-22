<template>
  <n-space vertical>
    <div v-if="error" style="color: red; margin-top: 10px">错误: {{ error }}</div>

    <n-table v-if="results.length > 0" :bordered="true" style="margin-top: 20px" width="100%">
      <thead>
        <tr>
          <th>域名</th>
          <th>IP 地址</th>
          <th>DNS 服务器</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(result, index) in results" :key="index">
          <td>{{ domain }}</td>
          <td>{{ result.ip }}</td>
          <td>{{ result.address }}</td>
        </tr>
      </tbody>
    </n-table>
  </n-space>
</template>

<script lang="ts" setup>
import { results, error } from './index';
import { domain } from '@/components/dnslog/input/use-domain';
</script>
