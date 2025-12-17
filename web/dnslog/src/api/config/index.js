import http from '@/utils/request';

export const fetchRuntimeConfig = () => http.get('/api/config');

