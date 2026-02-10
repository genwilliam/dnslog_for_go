<template>
	<n-config-provider>
		<n-message-provider>
			<RouterView />
		</n-message-provider>
	</n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider } from 'naive-ui';
</script>
