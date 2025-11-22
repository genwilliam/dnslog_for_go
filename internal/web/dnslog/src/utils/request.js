import axios from 'axios';
import { baseURL } from './index';
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
      return Promise.resolve(res);
    } else {
      return Promise.reject(new Error(res?.msg || '请求失败'));
    }
  },
  (error) => {
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
