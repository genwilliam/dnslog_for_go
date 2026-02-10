import http from '@/utils/request';

export function createApiKey(data = {}) {
	return http.post('/keys', data);
}

export function listApiKeys(params = {}) {
	return http.get('/keys', params);
}

export function disableApiKey(id) {
	return http.post(`/keys/${id}/disable`);
}

export function addBlacklist(data = {}) {
	return http.post('/blacklist', data);
}

export function listBlacklist(params = {}) {
	return http.get('/blacklist', params);
}

export function disableBlacklist(id) {
	return http.post(`/blacklist/${id}/disable`);
}
