<template>
  <div class="background">
    <slot></slot>
  </div>
</template>

<script lang="ts" setup></script>

<style scoped>
.background {
  position: fixed; /* 固定在整个视口上 */
  top: 0;
  left: 0;
  width: 100vw; /* 占满整个宽度 */
  height: 100vh; /* 占满整个高度 */
  background: linear-gradient(135deg, #74abe2, #5563de); /* 一个柔和的渐变 */
  display: flex; /* 居中子元素 */
  justify-content: center;
  align-items: center;
  overflow: hidden; /* 防止内容溢出 */
}
</style>
