<template>
	<div class="page-wrapper">
		<n-space vertical size="large">
			<!-- 页面抬头 -->
			<n-card size="small" class="card">
				<n-space vertical size="small">
					<n-h3 prefix="bar" align-text>域名生成</n-h3>
					<n-text depth="3">生成测试域名，将目标系统的 DNS 请求打到本服务后进行观测。</n-text>
				</n-space>
			</n-card>

			<!-- 运行配置 -->
			<n-card title="运行配置" size="small" class="card" :segmented="true">
				<n-space wrap size="small">
					<n-tag type="primary" round>根域名: {{ runtimeConfig.root_domain || '加载中' }}</n-tag>
					<n-tag type="info" round>DNS监听: {{ runtimeConfig.dns_listen_addr || '加载中' }}</n-tag>
					<n-tag type="info" round>HTTP: {{ runtimeConfig.http_listen || '加载中' }}</n-tag>
					<n-tag type="success" round>协议: {{ runtimeConfig.protocol || '加载中' }}</n-tag>
					<n-tag type="warning" round>
						上游DNS: {{ runtimeConfig.upstream_dns?.join(', ') || '加载中' }}
					</n-tag>
				</n-space>
				<template #footer>
					<n-text depth="3">
						生成的测试域名会落在上述根域名下，请将该域名的解析请求打到 DNS 监听地址后轮询记录。
					</n-text>
				</template>
			</n-card>

			<!-- 查询结果 -->
			<n-card title="查询结果" size="small" class="card" :segmented="true">
				<n-space vertical size="small">
					<n-space align="center" wrap>
						<n-text depth="3">本地验证请复制并执行：</n-text>
						<n-input size="small" readonly style="min-width: 420px" :value="digCommand" />
						<n-button size="small" @click="copyDig">复制 dig 命令</n-button>
					</n-space>
					<n-space align="center" wrap>
						<span class="label">轮询间隔(ms)：</span>
						<n-input-number
							v-model:value="pollInterval"
							:min="500"
							:step="500"
							size="small"
							style="width: 160px"
						/>
					</n-space>

					<div v-if="error" class="error-text">错误：{{ error }}</div>
					<div v-if="tokenStatus" class="status-row">
						<n-tag :type="statusTagType" round>状态: {{ tokenStatus.status }}</n-tag>
						<n-tag type="info" round>命中次数: {{ tokenStatus.hit_count }}</n-tag>
						<n-tag type="warning" round>过期时间: {{ formatTime(tokenStatus.expires_at) }}</n-tag>
					</div>

					<n-data-table
						v-if="records.length > 0"
						:columns="columns"
						:data="records"
						:max-height="520"
						bordered
					/>
					<n-pagination
						v-if="records.length > 0"
						:page="page"
						:page-size="pageSize"
						:page-count="pageCount"
						size="small"
						@update:page="handlePageChange"
					/>

					<div v-else class="empty">暂无 DNS 查询记录。</div>
				</n-space>
			</n-card>
		</n-space>
	</div>
</template>

<script setup lang="ts">
import { h, ref, watch, onMounted, onBeforeUnmount, reactive, computed } from 'vue';
import { useRoute } from 'vue-router';
import dayjs from 'dayjs';
import { records, error, tokenStatus, recordsTotal, fetchTokenRecords } from './index';
import type { DataTableColumns } from 'naive-ui';
import { useMessage } from 'naive-ui';
import { setPollingInterval, stopPolling, startPolling } from '@/utils/pulling.ts';
import { domain as domainRef } from '@/components/dnslog/input/use-domain';
import { fetchRuntimeConfig } from '@/api/config/index.js';
import { useStore } from 'vuex';

const pollInterval = ref<number>(Number(import.meta.env.VITE_POLL_INTERVAL_MS || 2000));
watch(pollInterval, (val) => {
	if (!val) return;
	setPollingInterval(val);
});

const message = useMessage();
const lastStatus = ref<string | null>(null);
const route = useRoute();
const store = useStore();

const runtimeConfig = reactive<any>({});
const loadConfig = async () => {
	try {
		const res = await fetchRuntimeConfig();
		if (res.code === 200) {
			Object.assign(runtimeConfig, res.data);
			let apiKeyRequired = true;
			if (res.data.apiKeyRequired !== undefined) {
				apiKeyRequired = Boolean(res.data.apiKeyRequired);
			} else if (res.data.api_key_required !== undefined) {
				apiKeyRequired = Boolean(res.data.api_key_required);
			}
			store.commit('runtimeConfig/setConfig', {
				apiKeyRequired: Boolean(apiKeyRequired),
				dnsPort: res.data.dns_port || inferPort(res.data.dns_listen_addr) || '15353',
				rootDomain: res.data.root_domain || '',
			});
		}
	} catch (e) {
		console.warn('加载配置失败', e);
	}
};

onMounted(() => {
	loadConfig();
	startPolling();
});

onBeforeUnmount(() => {
	stopPolling();
});

watch(
	() => route.query.domain,
	(val) => {
		if (typeof val === 'string' && val) {
			domainRef.value = val;
		}
	},
	{ immediate: true },
);

const page = ref(1);
const pageSize = ref(20);
const pageCount = ref(1);

function handlePageChange(p: number) {
	page.value = p;
	fetchTokenRecords(page.value, pageSize.value, 'desc');
}

const statusTagType = computed(() => {
	if (!tokenStatus.value) return 'default';
	if (tokenStatus.value.status === 'HIT') return 'success';
	if (tokenStatus.value.status === 'EXPIRED') return 'error';
	return 'warning';
});

function formatTime(ts?: number) {
	if (!ts) return '-';
	return dayjs(ts).format('YYYY-MM-DD HH:mm:ss');
}

watch(tokenStatus, (val) => {
	if (!val) return;
	if (lastStatus.value !== val.status) {
		lastStatus.value = val.status;
		if (val.status === 'HIT') {
			message.success('首次命中已捕获，已降频轮询');
			setPollingInterval(10000);
			page.value = 1;
			fetchTokenRecords(page.value, pageSize.value, 'desc');
		} else if (val.status === 'EXPIRED') {
			message.warning('Token 已过期，停止轮询');
			stopPolling();
		} else if (val.status === 'INIT') {
			setPollingInterval(2000);
		}
	}
	pageCount.value = Math.max(1, Math.ceil((recordsTotal.value || 0) / pageSize.value));
});

// 表格列
const columns: DataTableColumns = [
	{
		title: '时间',
		key: 'timestamp',
		width: 180,
		render(row) {
			return dayjs(row.timestamp).format('YYYY-MM-DD HH:mm:ss');
		},
	},
	{
		title: '子域名',
		key: 'domain',
		width: 180,
	},
	{
		title: '来源 IP',
		key: 'client_ip',
		width: 120,
	},
	{
		title: '协议',
		key: 'protocol',
		width: 90,
		render(row) {
			return h(
				'span',
				{
					style: {
						padding: '4px 8px',
						background: '#eef',
						borderRadius: '4px',
						fontSize: '12px',
					},
				},
				row.protocol?.toUpperCase?.() || 'UNKNOWN',
			);
		},
	},
	{
		title: '类型',
		key: 'qtype',
		width: 80,
	},
	{
		title: 'Token',
		key: 'token',
		width: 140,
	},
	{
		title: 'DNS 服务器',
		key: 'server',
		width: 150,
	},
];

const digCommand = computed(() => {
	const dnsPort =
		store.state.runtimeConfig?.dnsPort ||
		runtimeConfig.dns_port ||
		inferPort(runtimeConfig.dns_listen_addr) ||
		'15353';
	const domain = domainRef.value || '<token>.demo.com';
	return `dig @127.0.0.1 -p ${dnsPort} ${domain}`;
});

function inferPort(addr?: string): string {
	if (!addr) return '';
	const idx = addr.lastIndexOf(':');
	if (idx >= 0 && idx < addr.length - 1) {
		return addr.slice(idx + 1);
	}
	return '';
}

function copyDig() {
	navigator.clipboard
		.writeText(digCommand.value)
		.then(() => message.success('已复制 dig 命令'))
		.catch(() => message.error('复制失败'));
}
</script>
<style scoped>
.page-wrapper {
	width: 100%;
	height: 100%;
	padding: 16px 20px 24px;
	box-sizing: border-box;
	background: #f5f6f8;
}

.card {
	box-shadow: 0 6px 18px rgba(0, 0, 0, 0.06);
}

.label {
	color: #555;
	font-size: 14px;
}

.error-text {
	color: #d03050;
	font-size: 14px;
}

.status-row {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
}

.empty {
	color: #888;
	font-size: 14px;
	padding: 12px 0;
}
</style>
