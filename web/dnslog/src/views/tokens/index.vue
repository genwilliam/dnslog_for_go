<template>
	<div class="page-wrapper">
		<n-space vertical size="large">
				<n-card size="small" class="card">
					<n-space vertical size="small">
						<n-h3 prefix="bar" align-text>Tokens</n-h3>
						<n-text depth="3">按状态与时间筛选 token，支持复制与跳转。</n-text>
					</n-space>
				</n-card>

				<n-card title="筛选条件" size="small" class="card" :segmented="true">
					<n-form inline label-placement="left" label-width="90" :show-require-mark="false">
						<n-form-item label="关键字">
							<n-input v-model:value="filters.keyword" placeholder="token 或 domain" clearable />
						</n-form-item>
						<n-form-item label="状态">
							<n-select
								v-model:value="filters.status"
								:options="statusOptions"
								placeholder="任意"
								clearable
								style="width: 140px"
							/>
						</n-form-item>
						<n-form-item label="创建时间">
							<n-date-picker
								v-model:value="filters.createdRange"
								type="datetimerange"
								clearable
								format="yyyy-MM-dd HH:mm:ss"
								style="width: 280px"
							/>
						</n-form-item>
						<n-form-item label="命中时间">
							<n-date-picker
								v-model:value="filters.lastRange"
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
				</n-card>

				<n-card title="Token 列表" size="small" class="card" :segmented="true">
					<n-space vertical size="small">
						<div v-if="error" class="error-text">错误：{{ error }}</div>
						<n-data-table
							:loading="loading"
							:columns="columns"
							:data="tokens"
							:max-height="520"
							bordered
							:row-key="(row) => row.token"
						/>
						<n-pagination
							:page="pagination.page"
							:page-size="pagination.pageSize"
							:page-count="pagination.pageCount"
							size="small"
							@update:page="handlePageChange"
						/>
					</n-space>
				</n-card>
		</n-space>
	</div>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import dayjs from 'dayjs';
import type { DataTableColumns } from 'naive-ui';
import { useMessage, NButton } from 'naive-ui';
import { listTokens } from '@/api/tokens/index.js';

interface TokenItem {
	token: string;
	domain: string;
	status: 'INIT' | 'HIT' | 'EXPIRED';
	first_seen: number;
	last_seen: number;
	hit_count: number;
	created_at: number;
	expires_at: number;
}

const router = useRouter();
const message = useMessage();
const tokens = ref<TokenItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

const filters = reactive({
	keyword: '',
	status: null as null | string,
	createdRange: null as null | [number, number],
	lastRange: null as null | [number, number],
});

const statusOptions = [
	{ label: 'INIT', value: 'INIT' },
	{ label: 'HIT', value: 'HIT' },
	{ label: 'EXPIRED', value: 'EXPIRED' },
];

const pagination = reactive({
	page: 1,
	pageSize: 20,
	pageCount: 1,
});

function formatTime(ts?: number) {
	if (!ts) return '-';
	return dayjs(ts).format('YYYY-MM-DD HH:mm:ss');
}

function copyText(text: string) {
	if (!text) return;
	navigator.clipboard
		.writeText(text)
		.then(() => message.success('已复制'))
		.catch(() => message.error('复制失败'));
}

function jumpToQuery(domain: string) {
	router.push({ path: '/dnsquery', query: { domain } });
}

const columns: DataTableColumns<TokenItem> = [
	{
		title: 'Token',
		key: 'token',
		width: 150,
		render(row) {
			return h(
				'div',
				{ style: { display: 'flex', gap: '6px', alignItems: 'center' } },
				[
					h('span', row.token),
					h(
						NButton,
						{ size: 'tiny', type: 'primary', onClick: () => copyText(row.token) },
						{ default: () => '复制' },
					),
				],
			);
		},
	},
	{
		title: '域名',
		key: 'domain',
		minWidth: 220,
		render(row) {
			return h(
				'div',
				{ style: { display: 'flex', gap: '6px', alignItems: 'center' } },
				[
					h('span', row.domain),
					h(
						NButton,
						{ size: 'tiny', onClick: () => copyText(row.domain) },
						{ default: () => '复制' },
					),
					h(
						NButton,
						{ size: 'tiny', type: 'success', onClick: () => jumpToQuery(row.domain) },
						{ default: () => '跳转' },
					),
				],
			);
		},
	},
	{
		title: '状态',
		key: 'status',
		width: 90,
	},
	{
		title: '命中',
		key: 'hit_count',
		width: 80,
	},
	{
		title: '创建时间',
		key: 'created_at',
		width: 170,
		render(row) {
			return formatTime(row.created_at);
		},
	},
	{
		title: '最后命中',
		key: 'last_seen',
		width: 170,
		render(row) {
			return formatTime(row.last_seen);
		},
	},
	{
		title: '过期时间',
		key: 'expires_at',
		width: 170,
		render(row) {
			return formatTime(row.expires_at);
		},
	},
];

async function fetchTokens() {
	loading.value = true;
	error.value = null;
	try {
		const params: any = {
			page: pagination.page,
			pageSize: pagination.pageSize,
		};
		if (filters.keyword) params.keyword = filters.keyword;
		if (filters.status) params.status = filters.status;
		if (filters.createdRange && filters.createdRange.length === 2) {
			params.created_start = String(filters.createdRange[0]);
			params.created_end = String(filters.createdRange[1]);
		}
		if (filters.lastRange && filters.lastRange.length === 2) {
			params.last_start = String(filters.lastRange[0]);
			params.last_end = String(filters.lastRange[1]);
		}

		const res = await listTokens(params);
		tokens.value = res.data.items || [];
		const total = Number(res.data.total || 0);
		pagination.pageCount = Math.max(1, Math.ceil(total / pagination.pageSize));
	} catch (e: any) {
		error.value = e?.message || '加载失败';
	} finally {
		loading.value = false;
	}
}

function handleSearch() {
	pagination.page = 1;
	fetchTokens();
}

function handleReset() {
	filters.keyword = '';
	filters.status = null;
	filters.createdRange = null;
	filters.lastRange = null;
	pagination.page = 1;
	fetchTokens();
}

function handlePageChange(p: number) {
	pagination.page = p;
	fetchTokens();
}

onMounted(() => {
	fetchTokens();
});
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
