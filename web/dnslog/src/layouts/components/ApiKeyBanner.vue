<template>
	<div v-if="visible" class="api-banner">
		<n-alert :type="alertType" :show-icon="true">
			<template #header>
				{{ title }}
			</template>
			<n-space vertical size="small">
				<n-text>{{ message }}</n-text>
				<n-space v-if="showEditor" align="center" wrap>
					<n-input
						v-model:value="draft"
						placeholder="请输入 X-API-Key"
						size="small"
						style="min-width: 260px"
					/>
					<n-button size="small" type="primary" @click="saveKey">保存</n-button>
					<n-button v-if="hasKey" size="small" @click="cancelEdit">取消</n-button>
					<n-button size="small" @click="clearKey">清除</n-button>
				</n-space>
				<n-space v-else align="center" wrap>
					<n-tag type="success" round>已配置：{{ maskedKey }}</n-tag>
					<n-button size="small" @click="startEdit">修改</n-button>
					<n-button size="small" @click="clearKey">清除</n-button>
				</n-space>
				<n-text depth="3">
					{{ helperText }}
				</n-text>
			</n-space>
		</n-alert>
	</div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useStore } from 'vuex';
import { useMessage } from 'naive-ui';
import {
	ApiKeyLegacyStorageKey,
	ApiKeyStorageKey,
	getStoredApiKey,
	isValidApiKeyFormat,
	maskApiKey,
} from '@/utils/apiKey';
import { fetchRuntimeConfig } from '@/api/config/index.js';
import { clearApiKey, getApiKeySource, loadApiKey, setApiKey } from '@/stores/auth';

const store = useStore();
const messageBox = useMessage();
const apiError = computed(() => store.state.system?.apiError || null);

const apiKey = computed(() => store.state.apiKey?.value || '');
const apiKeySource = computed(() => getApiKeySource());
const draft = ref('');
const editing = ref(false);

const hasKey = computed(() => apiKey.value.length > 0);
const apiKeyRequired = computed(() => store.state.runtimeConfig?.apiKeyRequired);
const visible = computed(() => {
	if (apiError.value) return true;
	if (apiKeyRequired.value === false) return false;
	return !hasKey.value;
});
const showEditor = computed(() => !hasKey.value || editing.value);
const maskedKey = computed(() => maskApiKey(apiKey.value));

const alertType = computed(() => {
	if (!hasKey.value || apiError.value?.type === 'unauthorized') {
		return 'error';
	}
	if (apiError.value?.type === 'rate_limited') {
		return 'warning';
	}
	if (!apiError.value) {
		return 'success';
	}
	return 'error';
});

const title = computed(() => {
	if (!hasKey.value) {
		return '未配置 API Key';
	}
	if (apiError.value?.type === 'unauthorized') {
		return '未授权';
	}
	if (apiError.value?.type === 'forbidden') {
		return '访问被拒绝';
	}
	if (apiError.value?.type === 'rate_limited') {
		return '请求过于频繁';
	}
	if (!apiError.value) {
		return '已配置 API Key';
	}
	return '请求失败';
});

const message = computed(() => {
	if (!hasKey.value) {
		return '请配置 API Key 以访问受保护的接口。';
	}
	if (!apiError.value) {
		return '当前页面将自动携带 X-API-Key 访问后端。';
	}
	return apiError.value?.message || '请求失败，请检查后端与网络。';
});

const helperText = computed(() => {
	if (apiKeyRequired.value === false) {
		return '后端未开启鉴权，API Key 为可选。若需要访问受保护环境，可在此设置本地 DNSLOG_API_KEY。';
	}
	if (apiKeySource.value === 'env') {
		return '当前使用环境变量 VITE_API_KEY，可在此处覆盖并写入本地 DNSLOG_API_KEY。';
	}
	return 'API Key 会保存到本地 DNSLOG_API_KEY。';
});

const logMaskedKey = (label: string, value: string) => {
	if (!import.meta.env.DEV) {
		return;
	}
	const masked = maskApiKey(value || '');
	console.info(`[auth] ${label} key=${masked || 'none'}`);
};

async function saveKey() {
	const key = draft.value.trim();
	if (!key) {
		return;
	}
	if (!isValidApiKeyFormat(key)) {
		messageBox.error('API Key 格式不正确，应为 64 位十六进制字符串');
		return;
	}
	setApiKey(key);
	logMaskedKey('saved', getStoredApiKey());
	editing.value = false;
	store.commit('system/clearApiError');
	try {
		await fetchRuntimeConfig();
	} catch (error) {
		if (import.meta.env.DEV) {
			console.warn('[auth] verify config failed', error);
		}
	}
}

function clearKey() {
	clearApiKey();
	draft.value = '';
	editing.value = false;
	store.commit('system/clearApiError');
}

function startEdit() {
	editing.value = true;
	draft.value = apiKey.value;
}

function cancelEdit() {
	editing.value = false;
	draft.value = apiKey.value;
}

const onStorage = (event: StorageEvent) => {
	if (!event) {
		return;
	}
	if (
		event.key === ApiKeyStorageKey ||
		event.key === ApiKeyLegacyStorageKey ||
		event.key === null
	) {
		loadApiKey();
	}
};

watch(
	apiKey,
	(value, oldValue) => {
		if (!draft.value || draft.value === oldValue) {
			draft.value = value || '';
		}
	},
	{ immediate: true },
);

onMounted(() => {
	loadApiKey();
	if (typeof window !== 'undefined') {
		window.addEventListener('storage', onStorage);
	}
});

onBeforeUnmount(() => {
	if (typeof window !== 'undefined') {
		window.removeEventListener('storage', onStorage);
	}
});
</script>

<style scoped>
.api-banner {
	padding: 12px 20px 0;
	background: transparent;
}
</style>
