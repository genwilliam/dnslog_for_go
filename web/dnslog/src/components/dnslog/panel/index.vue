<template>
	<div class="btn-group">
		<generateBtn />
	</div>
</template>

<script setup lang="ts">
import generateBtn from '@/components/dnslog/button/generate.vue';
</script>

<style scoped>
.btn-group {
	display: flex;
	align-items: center;
	gap: 8px;
}
</style>
