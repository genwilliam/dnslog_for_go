import http from '@/utils/request';

export function dnsRecords(params = {}) {
	return http.get('/records', params);
}
