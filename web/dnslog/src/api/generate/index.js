import http from '@/utils/request';

// 生成随机域名
export const generateRandomDomain = () => {
	return http.get('/random-domain');
};
