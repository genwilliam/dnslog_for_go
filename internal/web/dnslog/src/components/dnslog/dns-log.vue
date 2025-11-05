<template>
  <Background>
    <div class="form-group">
      <panel></panel>
    </div>
  </Background>
</template>

<script lang="ts" setup>
import Panel from '@/components/dnslog/panel/index.vue'
import Background from '@/layouts/background.vue'
</script>

<style scoped>
.form-group {
  /* position:让元素（.form-group）脱离普通文档流，
  并且相对于最近的有定位的父元素进行定位*/
  position: absolute;

  /* 把这个元素的左上角，放在父容器的宽度 50%、高度 50% 的位置 */
  top: 50%;
  left: 50%;

  /* translate(-50%, -50%) 表示让元素向左移动自身宽度的一半，
  并向上移动自身高度的一半 */
  transform: translate(-50%, -50%);
}
</style>
