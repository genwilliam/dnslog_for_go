import http from '@/utils/request';

export const submitDns = (domain) => {
  return http.post('/dnslog/submit', {
    domain_name: domain,
  });
};
