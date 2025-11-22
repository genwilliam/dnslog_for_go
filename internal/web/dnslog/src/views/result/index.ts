import { ref } from 'vue';
import { submitDns } from '@/api/submit';
import { domain, submitDomain } from '@/components/dnslog/input/use-domain';

export interface DnsResult {
  ip: string;
  address: string;
}

export const results = ref<DnsResult[]>([]);
export const error = ref<string>('');

export async function fetchDns() {
  try {
    const d = submitDomain();

    const res = await submitDns(d);

    if (res.code !== 200) {
      error.value = res.message;
      results.value = [];
      return;
    }

    domain.value = res.data.domain;
    results.value = res.data.results;

    error.value = '';
  } catch (e) {
    error.value = '请求失败，请检查服务器是否启动';
  }
}
