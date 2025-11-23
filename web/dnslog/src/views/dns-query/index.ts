import { ref, type Ref } from 'vue';
import { submitDns } from '@/api/submit';
import { submitDomain } from '@/components/dnslog/input/use-domain';

// DNS 单条解析结果
export interface DnsResultItem {
	ip: string;
	address: string;
}

// 完整的 DNS 查询事件
export interface DnsRecord {
	domain: string;
	client_ip: string;
	query_cost: number;
	results: DnsResultItem[];
	timestamp: number;
	// trace_id: string;
}

// 所有查询记录
export const records: Ref<DnsRecord[]> = ref([]);

// 错误信息
export const error: Ref<string> = ref('');

// clearResults 清空所有历史 DNS 记录
export function clearResults() {
	records.value = [];
}

// fetchDns 请求一次 DNS 解析并记录事件
export async function fetchDns(shouldClear = false) {
	try {
		const domain = submitDomain();
		const res = await submitDns(domain);

		if (res.code !== 200) {
			error.value = res.message;
			return;
		}

		const data = res.data;

		const record: DnsRecord = {
			domain: data.domain,
			client_ip: data.client_ip,
			query_cost: data.query_cost,
			results: data.results,
			timestamp: data.timestamp,
			// trace_id: data.trace_id,
		};

		if (shouldClear) {
			// 覆盖旧数据
			records.value = [record];
		} else {
			// 追加新事件
			records.value.push(record);
		}

		error.value = '';
	} catch (e) {
		error.value = '请求失败: ' + (e instanceof Error ? e.message : String(e));
	}
}
