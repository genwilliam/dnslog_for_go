import http from '@/utils/request';

export const submitDns = (domain) => {
	return http.post('/submit', {
		domain_name: domain,
	});
};
