import { ref } from 'vue';
import { submitDomain } from '@/components/dnslog/input/use-domain.ts';

export interface DnsResult {
  ip: string;
  address: string;
}

export const results = ref<DnsResult[]>([]);
export const error = ref<string>('');

export default async function fetchDns() {
  try {
    const currentDomain = submitDomain();
    console.log('提交的域名:', currentDomain);

    const response = await fetch('http://localhost:8080/submit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ domain_name: currentDomain }),
    });

    const data = await response.json();

    if (data.error) {
      error.value = data.error;
      results.value = [];
    } else {
      error.value = '';
      results.value = data.results || [];
    }
  } catch (err: any) {
    console.error('请求失败:', err);
    error.value = '请求失败，请检查服务器是否启动';
  }
}
