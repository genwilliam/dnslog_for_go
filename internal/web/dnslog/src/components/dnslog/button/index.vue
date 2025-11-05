<template>
  <n-config-provider>
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-modal-provider>
            <n-dialog-provider>
              <div class="btn-group">
                <generateBtn />
                <queryBtn />
                <stateBtn />
                <dataPersis />
              </div>
            </n-dialog-provider>
          </n-modal-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import generateBtn from '@/components/dnslog/button/generate.vue';
import queryBtn from '@/components/dnslog/button/query.vue';
import stateBtn from '@/components/dnslog/button/state.vue';
import dataPersis from '@/components/dnslog/button/persis.vue';
</script>

<style scoped>
.btn-group {
  display: flex;
  gap: 12px;
}
</style>
