import { ref, type Ref } from 'vue';
import { getTokenRecords, getTokenStatus } from '@/api/tokens';
import { submitDomain } from '@/components/dnslog/input/use-domain';

// DNSLog 观测记录
export interface DnsRecord {
	domain: string;
	client_ip: string;
	protocol: string;
	qtype: string;
	timestamp: number;
	server: string;
	token: string;
}

export interface TokenStatus {
	token: string;
	domain: string;
	status: 'INIT' | 'HIT' | 'EXPIRED';
	first_seen: number;
	last_seen: number;
	hit_count: number;
	expires_at: number;
	expired: boolean;
}

// 所有查询记录
export const records: Ref<DnsRecord[]> = ref([]);
export const tokenStatus: Ref<TokenStatus | null> = ref(null);
export const recordsTotal: Ref<number> = ref(0);

// 错误信息
export const error: Ref<string> = ref('');

const tokenRegex = /^[a-f0-9]{10}$/i;

function parseToken(input: string) {
	const trimmed = input.trim();
	if (!trimmed) {
		return { token: '', reason: '请输入 token 或使用 Generate 生成域名' };
	}
	const firstPart = trimmed.split('.')[0] || '';
	if (!tokenRegex.test(firstPart)) {
		return { token: '', reason: '请输入 token 或使用 Generate 生成域名' };
	}
	return { token: firstPart.toLowerCase(), reason: '' };
}

export async function fetchTokenStatus() {
	try {
		const domain = submitDomain();
		const { token, reason } = parseToken(domain);
		if (!token) {
			error.value = reason;
			tokenStatus.value = null;
			return;
		}

		const res = await getTokenStatus(token);
		if (res.code !== 200 || !res.data) {
			error.value = res.message || '请求失败';
			return;
		}

		tokenStatus.value = res.data as TokenStatus;
		error.value = '';
	} catch (e) {
		error.value = '请求失败: ' + (e instanceof Error ? e.message : String(e));
	}
}

export async function fetchTokenRecords(page = 1, pageSize = 20, order: 'asc' | 'desc' = 'desc') {
	try {
		const domain = submitDomain();
		const { token, reason } = parseToken(domain);
		if (!token) {
			error.value = reason;
			records.value = [];
			recordsTotal.value = 0;
			return;
		}

		const res = await getTokenRecords(token, { page, pageSize, order });
		if (res.code !== 200 || !res.data) {
			error.value = res.message || '请求失败';
			return;
		}

		const items = (res.data.items || []) as DnsRecord[];
		records.value = items;
		recordsTotal.value = Number(res.data.total || 0);

		error.value = '';
	} catch (e) {
		error.value = '请求失败: ' + (e instanceof Error ? e.message : String(e));
	}
}
