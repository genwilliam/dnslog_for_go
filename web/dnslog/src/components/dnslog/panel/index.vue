<template>
	<n-config-provider>
		<n-loading-bar-provider>
			<n-message-provider>
				<n-notification-provider>
					<n-modal-provider>
						<n-dialog-provider>
							<div class="btn-group">
								<generateBtn />
								<stateBtn />
								<dataPersisBtn />
								<clearBtn />
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
import stateBtn from '@/components/dnslog/button/state.vue';
import dataPersisBtn from '@/components/dnslog/button/persis.vue';
import clearBtn from '@/components/dnslog/button/clear.vue';
</script>

<style scoped>
.btn-group {
	display: flex;
	align-items: center;

	gap: 8px;
}
</style>
