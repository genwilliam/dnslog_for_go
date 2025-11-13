<template>
  <div class="panel">
    <inp class="inp" />
    <btn class="btn" />
    <result />
  </div>
</template>

<script lang="ts" setup>
import Inp from '@/components/dnslog/input/index.vue';
import Btn from '@/components/dnslog/button/index.vue';
import result from '@/views/result/index.vue';
</script>

<style scoped>
.panel {
  display: flex;
  flex-direction: column;
  gap: 16px; /* 控制上下间距 */

  width: 668px;
  min-height: 325px;
  border: 1px solid #ccc;
  padding: 16px;
  border-radius: 8px;
  box-sizing: border-box;
  resize: vertical;
  overflow: auto;
}
</style>
