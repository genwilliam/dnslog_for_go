<template>
  <div class="layout-container">
    <header class="navbar">
      <div class="brand-logo">
        <span class="icon">📡</span>
        <span class="text">DNS Logger</span>
      </div>
    </header>

    <div class="toolbar">
      <div class="search-section">
        <Inp class="search-input" />
        <queryBtn />
      </div>

      <div class="action-section">
        <Btn />
      </div>
    </div>

    <main class="main-content">
      <div class="result-wrapper">
        <result class="result-component" />
      </div>
    </main>
  </div>
</template>

<script lang="ts" setup>
import Inp from '@/components/dnslog/input/index.vue';
import Btn from '@/components/dnslog/index.vue';
import result from '@/views/result/index.vue';
import queryBtn from '@/components/dnslog/button/query.vue';
</script>

<style scoped>
/* 整体布局架构 */
.layout-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background-color: #f0f2f5;
  overflow: hidden;
}

/* 顶部导航栏 */
.navbar {
  height: 56px;
  background-color: #001529;
  color: #fff;

  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: center;

  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
  z-index: 20;
  position: relative; /* 为了以后可能的绝对定位元素做准备 */
}

.brand-logo {
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 10px;
  letter-spacing: 1px;

  /* 文字本身不可被选中 */
  user-select: none;
  cursor: default;
}

/* 操作工具栏 */
.toolbar {
  background-color: #ffffff;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
  gap: 20px;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  max-width: 800px;
}

.search-input {
  flex: 1;
}

.action-section {
  display: flex;
  align-items: center;
}

/* 主内容区域 */
.main-content {
  flex: 1;
  padding: 20px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.result-wrapper {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  padding: 20px;
  overflow-y: auto;
  position: relative;
}

.result-component {
  min-height: 100%;
  width: 100%;
}
</style>
