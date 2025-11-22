<template>
	<n-config-provider>
		<n-space vertical>
			<div v-if="error" style="color: red; margin-top: 10px">错误: {{ error }}</div>

			<n-data-table
				v-if="results.length > 0"
				:columns="columns"
				:data="tableData"
				:max-height="tableMaxHeight"
				:bordered="true"
				style="margin-top: 20px"
			/>

			<div v-else>
				<p>暂无 DNS 查询结果。</p>
			</div>
		</n-space>
	</n-config-provider>
</template>

<script lang="ts" setup>
import { computed, h } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { results, error } from './index';
import { domain } from '@/components/dnslog/input/use-domain';

interface DnsResult {
	ip: string;
	address: string;
}

// 定义表格的列配置
// 我们需要使用 render 函数来确保 '域名' 列显示外部的 domain 变量
const createColumns = ({
	domainValue, // 接收 domain 的响应式值
}: {
	domainValue: string;
}): DataTableColumns<DnsResult> => {
	return [
		{
			title: '域名',
			key: 'domain',
			// 渲染函数：返回一个 span 元素，显示固定的 domain 值
			render() {
				// h 是 Vue 的渲染函数，用于创建 VNode
				return h('span', null, domainValue);
			},
		},
		{
			title: 'IP 地址',
			key: 'ip',
		},
		{
			title: 'DNS 服务器',
			key: 'address',
		},
	];
};

// 处理数据格式
// n-data-table 要求数据是一个对象数组，并且每一项最好有一个唯一的 key
const tableData = computed(() => {
	return results.value.map((item, index) => ({
		...item,
		key: index, // 必须为每一行添加一个唯一的 key
	}));
});

// 响应式地生成列配置
const columns = computed(() => {
	return createColumns({ domainValue: domain.value });
});

// 设置最大高度
const tableMaxHeight = '400px';
</script>
