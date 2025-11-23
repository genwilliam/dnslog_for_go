import http from '@/utils/request';

export function dnsRecords(domain) {
	return http.get('/api/records');
}
