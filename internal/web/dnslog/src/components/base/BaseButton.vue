<template>
  <n-button :type="type" :size="size" @click="$emit('click')">
    <slot>{{ text }}</slot>
  </n-button>
</template>

<script setup lang="ts">
const props = defineProps({
  type: { type: String, default: 'default' },
  size: { type: String, default: 'large' },
  text: { type: String, default: '' },
});
</script>
