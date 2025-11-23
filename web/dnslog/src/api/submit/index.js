import http from '@/utils/request';

export const submitDns = (domain) => {
	return http.post('/api/submit', {
		domain_name: domain,
	});
};
