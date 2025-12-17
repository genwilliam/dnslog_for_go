<template>
	<div class="page-wrapper">
		<n-config-provider>
			<n-space vertical size="large">
				<!-- 页面抬头 -->
				<n-card size="small" class="card">
					<n-space vertical size="small">
						<n-h3 prefix="bar" align-text>DNS 查询</n-h3>
						<n-text depth="3">生成测试域名，发起 DNS 查询并查看结果。</n-text>
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
							生成的测试域名会落在上述根域名下，请将该域名的解析请求打到 DNS 监听地址。
						</n-text>
					</template>
				</n-card>

				<!-- 查询结果 -->
				<n-card title="查询结果" size="small" class="card" :segmented="true">
					<n-space vertical size="small">
						<n-space align="center" wrap>
							<span class="label">轮询间隔(ms)：</span>
							<n-input-number v-model:value="pollInterval" :min="500" :step="500" size="small" style="width: 160px" />
						</n-space>

						<div v-if="error" class="error-text">错误：{{ error }}</div>

						<n-data-table
							v-if="records.length > 0"
							:columns="columns"
							:data="records"
							:max-height="520"
							bordered
						/>

						<div v-else class="empty">暂无 DNS 查询记录。</div>
					</n-space>
				</n-card>
			</n-space>
		</n-config-provider>
	</div>
</template>

<script setup lang="ts">
import { h, ref, watch, onMounted, reactive } from 'vue';
import dayjs from 'dayjs';
import { records, error } from './index';
import type { DataTableColumns } from 'naive-ui';
import { setPollingInterval } from '@/utils/pulling.ts';
import { fetchRuntimeConfig } from '@/api/config/index.js';

const pollInterval = ref<number>(Number(import.meta.env.VITE_POLL_INTERVAL_MS || 2000));
watch(pollInterval, (val) => {
	if (!val) return;
	setPollingInterval(val);
});

const runtimeConfig = reactive<any>({});
const loadConfig = async () => {
	try {
		const res = await fetchRuntimeConfig();
		if (res.code === 200) {
			Object.assign(runtimeConfig, res.data);
		}
	} catch (e) {
		console.warn('加载配置失败', e);
	}
};

onMounted(() => {
	loadConfig();
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
		title: '耗时 (ms)',
		key: 'query_cost',
		width: 120,
		render(row) {
			const ms = row.query_cost;
			const color = ms > 1500 ? 'red' : ms > 500 ? 'orange' : 'green';

			return h('span', { style: { color, fontWeight: 'bold' } }, `${ms} ms`);
		},
	},
	{
		title: '解析结果',
		key: 'results',
		render(row) {
			const text = row.results.map((r, i) => `${i + 1}. ${r.ip}`).join('\n');

			return h(
				'pre',
				{
					style: {
						margin: 0,
						whiteSpace: 'pre-wrap',
						fontFamily: 'monospace',
					},
				},
				text,
			);
		},
	},
	{
		title: 'DNS 服务器',
		key: 'dns_server',
		width: 150,
		render(row) {
			const server = row.results?.[0]?.address || '未知';
			return h(
				'span',
				{
					style: {
						padding: '4px 8px',
						background: '#eee',
						borderRadius: '4px',
						fontSize: '12px',
					},
				},
				server,
			);
		},
	},
];
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

.empty {
	color: #888;
	font-size: 14px;
	padding: 12px 0;
}
</style>
