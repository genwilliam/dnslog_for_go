<template>
	<div class="page-wrapper">
		<n-space vertical size="large">
				<!-- 页面抬头 + 运行配置提示 -->
				<n-card size="small" class="card">
					<n-space vertical size="small">
						<n-h3 prefix="bar" align-text>DNS 日志</n-h3>
						<n-text depth="3">
							这里被动展示所有到达 DNS 服务器的查询记录，无需在页面手动触发，只要目标系统把 DNS
							请求指向本服务即可。
						</n-text>
						<n-space wrap size="small" style="margin-top: 4px">
							<n-tag type="primary" round
								>根域名: {{ runtimeConfig.root_domain || '加载中' }}</n-tag
							>
							<n-tag type="info" round
								>DNS监听: {{ runtimeConfig.dns_listen_addr || '加载中' }}</n-tag
							>
							<n-tag type="success" round>协议: {{ runtimeConfig.protocol || '加载中' }}</n-tag>
							<n-tag type="warning" round>
								上游DNS: {{ runtimeConfig.upstream_dns?.join(', ') || '加载中' }}
							</n-tag>
						</n-space>
					</n-space>
				</n-card>

				<!-- 筛选表单 -->
				<n-card title="筛选条件" size="small" class="card" :segmented="true">
					<n-form inline label-placement="left" label-width="80" :show-require-mark="false">
						<n-form-item label="域名">
							<n-input v-model:value="filters.domain" placeholder="子/根 域关键词" clearable />
						</n-form-item>
						<n-form-item label="Token">
							<n-input v-model:value="filters.token" placeholder="token" clearable />
						</n-form-item>
						<n-form-item label="来源 IP">
							<n-input v-model:value="filters.client_ip" placeholder="IP 关键字" clearable />
						</n-form-item>
						<n-form-item label="协议">
							<n-select
								v-model:value="filters.protocol"
								:options="protocolOptions"
								placeholder="任意"
								clearable
								style="width: 120px"
							/>
						</n-form-item>
						<n-form-item label="类型">
							<n-select
								v-model:value="filters.qtype"
								:options="qtypeOptions"
								placeholder="任意"
								clearable
								style="width: 120px"
							/>
						</n-form-item>
						<n-form-item label="时间范围">
							<n-date-picker
								v-model:value="filters.timerange"
								type="datetimerange"
								clearable
								format="yyyy-MM-dd HH:mm:ss"
								style="width: 280px"
							/>
						</n-form-item>
						<n-form-item>
							<n-space>
								<n-button size="small" type="primary" @click="handleSearch" :loading="loading"
									>查询</n-button
								>
								<n-button size="small" @click="handleReset">重置</n-button>
							</n-space>
						</n-form-item>
					</n-form>
					<template #footer>
						<n-text depth="3"
							>当前总数：{{ pagination.itemCount }}，每页：{{ pagination.pageSize }}</n-text
						>
					</template>
				</n-card>

				<!-- 日志表 -->
				<n-card title="日志列表" size="small" class="card" :segmented="true">
					<n-space vertical size="small">
						<div v-if="error" class="error-text">错误：{{ error }}</div>
						<n-data-table
							:loading="loading"
							:columns="columns"
							:data="logs"
							:max-height="520"
							:pagination="pagination"
							bordered
							:row-key="(row) => row.timestamp + row.domain"
						/>
					</n-space>
				</n-card>
		</n-space>
	</div>
</template>

<script setup lang="ts">
import { h, ref, onMounted, reactive } from 'vue';
import { NTag } from 'naive-ui';
import dayjs from 'dayjs';
import type { DataTableColumns } from 'naive-ui';
import { fetchRuntimeConfig } from '@/api/config/index.js';
import { dnsRecords } from '@/api/records/index.js';

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
		page?: number;
		size?: number;
	};
}

// --- 状态管理 ---
const logs = ref<DnsLogItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const runtimeConfig = reactive<any>({});
const pagination = reactive({
	page: 1,
	pageSize: 20,
	itemCount: 0,
	showSizePicker: true,
	pageSizes: [10, 20, 50, 100],
	onChange: (page: number) => {
		pagination.page = page;
		fetchLogs();
	},
	onUpdatePageSize: (size: number) => {
		pagination.pageSize = size;
		pagination.page = 1;
		fetchLogs();
	},
});

const filters = reactive({
	domain: '',
	token: '',
	client_ip: '',
	protocol: null as null | string,
	qtype: null as null | string,
	timerange: null as null | [number, number],
});

const protocolOptions = [
	{ label: 'UDP', value: 'udp' },
	{ label: 'TCP', value: 'tcp' },
];
const qtypeOptions = [
	{ label: 'A', value: 'A' },
	{ label: 'AAAA', value: 'AAAA' },
	{ label: 'CNAME', value: 'CNAME' },
	{ label: 'MX', value: 'MX' },
	{ label: 'TXT', value: 'TXT' },
];

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
		const params = new URLSearchParams({
			page: String(pagination.page),
			pageSize: String(pagination.pageSize),
		});
		if (filters.domain) params.append('domain', filters.domain);
		if (filters.token) params.append('token', filters.token);
		if (filters.client_ip) params.append('client_ip', filters.client_ip);
		if (filters.protocol) params.append('protocol', filters.protocol);
		if (filters.qtype) params.append('qtype', filters.qtype);
		if (filters.timerange && filters.timerange.length === 2) {
			params.append('start', String(filters.timerange[0]));
			params.append('end', String(filters.timerange[1]));
		}

		const res: ApiResponse = await dnsRecords(Object.fromEntries(params));

		if (res.code === 200) {
			logs.value = res.data.items;
			pagination.itemCount = res.data.total;
			if (res.data.page) pagination.page = res.data.page;
			if (res.data.size) pagination.pageSize = res.data.size;
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
	loadConfig();
	fetchLogs();
});

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

const handleSearch = () => {
	pagination.page = 1;
	fetchLogs();
};

const handleReset = () => {
	filters.domain = '';
	filters.token = '';
	filters.client_ip = '';
	filters.protocol = null;
	filters.qtype = null;
	filters.timerange = null;
	pagination.page = 1;
	fetchLogs();
};
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

.error-text {
	color: #d03050;
	font-size: 14px;
}
</style>
