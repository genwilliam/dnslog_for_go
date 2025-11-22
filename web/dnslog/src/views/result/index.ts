import { ref, type Ref } from 'vue';
import { submitDns } from '@/api/submit';
import { domain, submitDomain } from '@/components/dnslog/input/use-domain';

export interface DnsResult {
	ip: string;
	address: string;
}

export const results: Ref<DnsResult[]> = ref([]);
export const error: Ref<string> = ref('');

/**
 * 清空所有已有的 DNS 查询结果。
 */
export function clearResults() {
	results.value = [];
	// 如果清空也要清空域名的话
	// domain.value = '';
}

export async function fetchDns(shouldClear: boolean = false) {
	try {
		const d = submitDomain();

		const res = await submitDns(d);

		if (res.code !== 200) {
			error.value = res.message;
			// 错误时不应该清空 results
			// 而是设置为空
			results.value = [];
			// return;
		}

		// 成功获取数据
		domain.value = res.data.domain;

		// 如果外部调用时传递了 true，则先清空
		if (shouldClear) {
			results.value = res.data.results || [];
		} else {
			// 将新的结果追加到现有数组中
			if (res.data.results && res.data.results.length > 0) {
				results.value.push(...res.data.results);
			}
		}

		error.value = '';
	} catch (e) {
		// 确保错误信息是一个字符串
		error.value = '请求失败，请检查服务器是否启动: ' + (e instanceof Error ? e.message : String(e));
	}
}
