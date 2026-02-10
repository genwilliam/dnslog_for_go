<template>
	<div class="page-wrapper">
		<n-space vertical size="large">
				<n-card size="small" class="card">
					<n-space vertical size="small">
						<n-h3 prefix="bar" align-text>Security</n-h3>
						<n-text depth="3">管理 API Keys 与 IP 黑名单。</n-text>
					</n-space>
				</n-card>

				<n-card title="API Keys" size="small" class="card" :segmented="true">
					<n-form inline label-placement="left" label-width="90" :show-require-mark="false">
						<n-form-item label="名称">
							<n-input v-model:value="keyForm.name" placeholder="ops" />
						</n-form-item>
						<n-form-item label="备注">
							<n-input v-model:value="keyForm.comment" placeholder="rotation-2025-01" />
						</n-form-item>
						<n-form-item>
							<n-button size="small" type="primary" @click="handleCreateKey" :loading="keyLoading"
								>创建</n-button
							>
						</n-form-item>
					</n-form>
					<div v-if="createdKey" class="info-box">
						<span>明文 Key（仅展示一次）：</span>
						<code>{{ createdKey }}</code>
						<n-button size="tiny" @click="copyText(createdKey)">复制</n-button>
						<n-button size="tiny" type="primary" @click="useCreatedKey">设为当前 Key</n-button>
					</div>
					<n-data-table
						:loading="keyLoading"
						:columns="keyColumns"
						:data="apiKeys"
						:max-height="360"
						bordered
						:row-key="(row) => row.id"
					/>
					<n-pagination
						:page="keyPagination.page"
						:page-size="keyPagination.pageSize"
						:page-count="keyPagination.pageCount"
						size="small"
						@update:page="handleKeyPageChange"
					/>
				</n-card>

				<n-card title="IP 黑名单" size="small" class="card" :segmented="true">
					<n-form inline label-placement="left" label-width="90" :show-require-mark="false">
						<n-form-item label="IP">
							<n-input v-model:value="blackForm.ip" placeholder="1.2.3.4" />
						</n-form-item>
						<n-form-item label="原因">
							<n-input v-model:value="blackForm.reason" placeholder="abuse" />
						</n-form-item>
						<n-form-item>
							<n-button size="small" type="primary" @click="handleAddBlacklist" :loading="blackLoading"
								>添加</n-button
							>
						</n-form-item>
					</n-form>
					<n-data-table
						:loading="blackLoading"
						:columns="blackColumns"
						:data="blacklist"
						:max-height="360"
						bordered
						:row-key="(row) => row.id"
					/>
					<n-pagination
						:page="blackPagination.page"
						:page-size="blackPagination.pageSize"
						:page-count="blackPagination.pageCount"
						size="small"
						@update:page="handleBlackPageChange"
					/>
				</n-card>
		</n-space>
	</div>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { useMessage, NButton } from 'naive-ui';
import dayjs from 'dayjs';
import {
	createApiKey,
	listApiKeys,
	disableApiKey,
	addBlacklist,
	listBlacklist,
	disableBlacklist,
} from '@/api/security/index.js';
import { setApiKey as setGlobalApiKey } from '@/stores/auth';

interface ApiKeyItem {
	id: number;
	name: string;
	enabled: boolean;
	created_at: number;
	last_used_at: number;
	comment: string;
}

interface BlacklistItem {
	id: number;
	ip: string;
	reason: string;
	enabled: boolean;
	created_at: number;
}

const message = useMessage();
const createdKey = ref('');

const keyForm = reactive({ name: '', comment: '' });
const apiKeys = ref<ApiKeyItem[]>([]);
const keyLoading = ref(false);
const keyPagination = reactive({ page: 1, pageSize: 20, pageCount: 1 });

const blackForm = reactive({ ip: '', reason: '' });
const blacklist = ref<BlacklistItem[]>([]);
const blackLoading = ref(false);
const blackPagination = reactive({ page: 1, pageSize: 20, pageCount: 1 });

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

function useCreatedKey() {
	if (!createdKey.value) {
		return;
	}
	setGlobalApiKey(createdKey.value);
	message.success('已设置为当前 API Key');
}

const keyColumns: DataTableColumns<ApiKeyItem> = [
	{ title: 'ID', key: 'id', width: 80 },
	{ title: '名称', key: 'name', width: 160 },
	{
		title: '状态',
		key: 'enabled',
		width: 90,
		render(row) {
			return row.enabled ? 'enabled' : 'disabled';
		},
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
		title: '最近使用',
		key: 'last_used_at',
		width: 170,
		render(row) {
			return formatTime(row.last_used_at);
		},
	},
	{
		title: '备注',
		key: 'comment',
	},
	{
		title: '操作',
		key: 'action',
		width: 120,
		render(row) {
			return h(
				NButton,
				{ size: 'tiny', type: 'error', onClick: () => handleDisableKey(row.id) },
				{ default: () => '禁用' },
			);
		},
	},
];

const blackColumns: DataTableColumns<BlacklistItem> = [
	{ title: 'ID', key: 'id', width: 80 },
	{ title: 'IP', key: 'ip', width: 140 },
	{ title: '原因', key: 'reason' },
	{
		title: '状态',
		key: 'enabled',
		width: 90,
		render(row) {
			return row.enabled ? 'enabled' : 'disabled';
		},
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
		title: '操作',
		key: 'action',
		width: 120,
		render(row) {
			return h(
				NButton,
				{ size: 'tiny', type: 'error', onClick: () => handleDisableBlacklist(row.id) },
				{ default: () => '禁用' },
			);
		},
	},
];

async function fetchKeys() {
	keyLoading.value = true;
	try {
		const res = await listApiKeys({
			page: keyPagination.page,
			pageSize: keyPagination.pageSize,
		});
		apiKeys.value = res.data.items || [];
		const total = Number(res.data.total || 0);
		keyPagination.pageCount = Math.max(1, Math.ceil(total / keyPagination.pageSize));
	} finally {
		keyLoading.value = false;
	}
}

async function fetchBlacklist() {
	blackLoading.value = true;
	try {
		const res = await listBlacklist({
			page: blackPagination.page,
			pageSize: blackPagination.pageSize,
		});
		blacklist.value = res.data.items || [];
		const total = Number(res.data.total || 0);
		blackPagination.pageCount = Math.max(1, Math.ceil(total / blackPagination.pageSize));
	} finally {
		blackLoading.value = false;
	}
}

async function handleCreateKey() {
	if (!keyForm.name) {
		message.error('名称不能为空');
		return;
	}
	keyLoading.value = true;
	try {
		const res = await createApiKey({ name: keyForm.name, comment: keyForm.comment });
		createdKey.value = res.data.key;
		keyForm.name = '';
		keyForm.comment = '';
		await fetchKeys();
	} finally {
		keyLoading.value = false;
	}
}

async function handleDisableKey(id: number) {
	await disableApiKey(id);
	await fetchKeys();
}

async function handleAddBlacklist() {
	if (!blackForm.ip) {
		message.error('IP 不能为空');
		return;
	}
	blackLoading.value = true;
	try {
		await addBlacklist({ ip: blackForm.ip, reason: blackForm.reason });
		blackForm.ip = '';
		blackForm.reason = '';
		await fetchBlacklist();
	} finally {
		blackLoading.value = false;
	}
}

async function handleDisableBlacklist(id: number) {
	await disableBlacklist(id);
	await fetchBlacklist();
}

function handleKeyPageChange(p: number) {
	keyPagination.page = p;
	fetchKeys();
}

function handleBlackPageChange(p: number) {
	blackPagination.page = p;
	fetchBlacklist();
}

onMounted(() => {
	fetchKeys();
	fetchBlacklist();
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

.info-box {
	margin: 12px 0;
	padding: 8px 12px;
	background: #f6f6f6;
	border: 1px solid #e5e5e5;
	border-radius: 6px;
	display: flex;
	gap: 8px;
	align-items: center;
	flex-wrap: wrap;
}
</style>
