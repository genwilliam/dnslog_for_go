<template>
	<div class="page-wrapper">
		<n-config-provider>
			<n-space vertical>
				<div v-if="error" style="color: red; margin: 10px 0">错误：{{ error }}</div>

				<div style="display: flex; justify-content: flex-end; padding: 0 10px">
					<n-button size="small" @click="fetchLogs" :loading="loading">刷新日志</n-button>
				</div>

				<n-data-table
					:loading="loading"
					:columns="columns"
					:data="logs"
					:max-height="520"
					bordered
					:row-key="(row) => row.timestamp + row.domain"
				/>
			</n-space>
		</n-config-provider>
	</div>
</template>

<script setup lang="ts">
import { h, ref, onMounted } from 'vue';
import { NTag, NButton } from 'naive-ui';
import dayjs from 'dayjs';
import type { DataTableColumns } from 'naive-ui';

// --- 类型定义 (根据后端 JSON) ---
interface DnsLogItem {
	domain: string;
	client_ip: string;
	protocol: 'udp' | 'tcp';
	qtype: string;
	timestamp: number;
	server: string;
}

interface ApiResponse {
	code: number;
	message: string;
	data: {
		items: DnsLogItem[];
		total: number;
	};
}

// --- 状态管理 ---
const logs = ref<DnsLogItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

// --- 表格列定义 ---
const columns: DataTableColumns<DnsLogItem> = [
	{
		title: '时间',
		key: 'timestamp',
		width: 180,
		render(row) {
			return dayjs(row.timestamp).format('YYYY-MM-DD HH:mm:ss');
		},
	},
	{
		title: '请求域名',
		key: 'domain',
		minWidth: 180,
		ellipsis: {
			tooltip: true,
		},
	},
	{
		title: '类型',
		key: 'qtype',
		width: 80,
		render(row) {
			// 根据类型显示不同颜色的标签
			const type = row.qtype.toUpperCase();
			let typeType: 'default' | 'success' | 'info' | 'warning' | 'error' = 'default';

			if (type === 'A') typeType = 'success';
			else if (type === 'AAAA') typeType = 'info';
			else if (type === 'CNAME') typeType = 'warning';

			return h(NTag, { type: typeType, size: 'small', bordered: false }, { default: () => type });
		},
	},
	{
		title: '协议',
		key: 'protocol',
		width: 80,
		render(row) {
			return h(
				NTag,
				{
					type: row.protocol === 'udp' ? 'info' : 'primary',
					size: 'small',
					variant: 'outline', // 轮廓样式，以此区分 QType
				},
				{ default: () => row.protocol.toUpperCase() },
			);
		},
	},
	{
		title: '来源 IP',
		key: 'client_ip',
		width: 140,
	},
	{
		title: '监听端口',
		key: 'server',
		width: 100,
		render(row) {
			return h(
				'span',
				{
					style: {
						fontFamily: 'monospace',
						background: '#f4f4f5',
						padding: '2px 6px',
						borderRadius: '4px',
						fontSize: '12px',
					},
				},
				row.server,
			);
		},
	},
];

// --- 数据获取 ---
const fetchLogs = async () => {
	loading.value = true;
	error.value = null;
	try {
		const response = await fetch('/api/records');
		if (!response.ok) throw new Error('网络请求失败');

		const res: ApiResponse = await response.json();

		if (res.code === 200) {
			logs.value = res.data.items;
		} else {
			error.value = res.message || '获取数据失败';
		}
	} catch (err: any) {
		error.value = err.message || '未知错误';
		console.error(err);
	} finally {
		loading.value = false;
	}
};

// 挂载时加载
onMounted(() => {
	fetchLogs();
});
</script>

<style scoped>
.page-wrapper {
	width: 100%;
	height: 100%;
	padding: 0; /* 保持和你原来的布局一致 */
	box-sizing: border-box;
}
</style>
