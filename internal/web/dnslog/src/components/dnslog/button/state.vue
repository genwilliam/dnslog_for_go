<template>
  <n-space>
    <!-- <n-button type="info" size="large">Generate</n-button> -->
    <!-- <n-button type="success" size="large">Record DNS query</n-button> -->
    <n-button type="error" size="large" @click="togglePolling">{{ state }}</n-button>
  </n-space>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

// true 表示正在 Polling，即显示 "Pause"
// false 显示 "Start"
const isPolling = ref(true)
const state = ref('Pause')

function togglePolling() {
  isPolling.value = !isPolling.value
  state.value = isPolling.value ? 'Pause' : 'Start'
}
</script>
