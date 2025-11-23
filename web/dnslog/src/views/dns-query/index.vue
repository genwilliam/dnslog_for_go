<template>
	<n-config-provider>
		<n-space vertical>
			<!-- 错误信息 -->
			<div v-if="error" style="color: red; margin: 10px 0">错误：{{ error }}</div>

			<!-- 有记录时显示表格 -->
			<n-data-table
				v-if="records.length > 0"
				:columns="columns"
				:data="records"
				:max-height="520"
				bordered
			/>

			<!-- 无结果 -->
			<div v-else>
				<p>暂无 DNS 查询记录。</p>
			</div>
		</n-space>
	</n-config-provider>
</template>

<script setup lang="ts">
import { h } from 'vue';
import dayjs from 'dayjs';
import { records, error } from './index';
import type { DataTableColumns } from 'naive-ui';

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
