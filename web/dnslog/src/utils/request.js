import axios from 'axios';
import store from '@/store';
import { baseURL } from './index';
import { getApiKey } from '@/stores/auth';
import { maskApiKey } from '@/utils/apiKey';
axios.defaults.headers['Content-Type'] = 'application/json;charset=utf-8';
// 创建axios实例
const request = axios.create({
  // axios中请求配置有baseURL选项，表示请求URL公共部分
  baseURL,
  // 超时
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8',
  },
});
// request拦截器
request.interceptors.request.use(
  (config) => {
    const apiKey = getApiKey();
    if (apiKey) {
      config.headers['X-API-Key'] = apiKey;
    }
    if (import.meta.env.DEV) {
      const method = (config.method || 'GET').toUpperCase();
      const url = config.url || '';
      const masked = apiKey ? maskApiKey(apiKey) : 'none';
      console.info(`[http] ${method} ${url} key=${masked} injected=${!!apiKey}`);
    }
    // // 获取用户状态
    // const token = localStorage.getItem('token')
    // // 添加统一的 token
    // if (token) {
    //   config.headers['Authorization'] = 'Bearer ' + token
    // }
    return config;
  },
  (error) => Promise.reject(error),
);
// response拦截器
request.interceptors.response.use(
  (response) => {
    const res = response?.data;
    if (res && res.code === 200) {
      if (store?.commit && store?.state?.system) {
        store.commit('system/clearApiError');
      }
      return Promise.resolve(res);
    } else {
      if (store?.commit && store?.state?.system) {
        store.commit('system/setApiError', {
          status: response?.status || 0,
          type: 'error',
          message: res?.message || '请求失败',
        });
      }
      return Promise.reject(new Error(res?.message || '请求失败'));
    }
  },
  (error) => {
    const status = error?.response?.status;
    const apiKeyRequired = store?.state?.runtimeConfig?.apiKeyRequired;
    if (store?.commit && store?.state?.system) {
      if (status === 401) {
        if (apiKeyRequired === false) {
          store.commit('system/setApiError', {
            status,
            type: 'unauthorized',
            message: '后端未开启鉴权但返回 401，可能请求路径错误或后端配置不一致',
          });
        } else {
          store.commit('system/setApiError', {
            status,
            type: 'unauthorized',
            message: '未授权：请配置 API Key',
          });
        }
      } else if (status === 403) {
        store.commit('system/setApiError', {
          status,
          type: 'forbidden',
          message: '访问被拒绝：可能被黑名单或无权限',
        });
      } else if (status === 404) {
        store.commit('system/setApiError', {
          status,
          type: 'not_found',
          message: `接口路径错误或资源不存在: ${error?.config?.url || ''}`,
        });
      } else if (status === 429) {
        store.commit('system/setApiError', {
          status,
          type: 'rate_limited',
          message: '请求过于频繁，请稍后再试',
        });
      } else {
        store.commit('system/setApiError', {
          status: status || 0,
          type: 'network',
          message: '网络异常或服务不可用',
        });
      }
    }
    return Promise.reject(error);
  },
);

// 封装请求方法
const http = {
  get(url, params = {}) {
    return request.get(url, { params });
  },
  post(url, data = {}) {
    return request.post(url, data);
  },
  put(url, data = {}) {
    return request.put(url, data);
  },
  delete(url, data = {}) {
    return request.delete(url, { data });
  },
};

export default http;
