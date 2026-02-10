import http from '@/utils/request';

export const fetchRuntimeConfig = () => http.get('/config');
