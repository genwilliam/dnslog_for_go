import http from '@/utils/request';

export function getTokenStatus(token) {
	return http.get(`/tokens/${encodeURIComponent(token)}`);
}

export function getTokenRecords(token, params = {}) {
	return http.get(`/tokens/${encodeURIComponent(token)}/records`, params);
}

export function listTokens(params = {}) {
	return http.get('/tokens', params);
}
