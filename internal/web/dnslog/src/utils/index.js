import axios from 'axios';
import { logInfo, redLog, blackLog } from '@/utils';
// 创建一个实例，此后都在此实例上改造
const http = axios.create({
  // timeout: 1000 * 4,
  withCredentials: true,
});
// 请求拦截器
const beforeRequest = (config) => {
  // 设置 token
  const token = localStorage.getItem('token');
  token && (config.headers.Authorization = token);
  // NOTE  添加自定义头部
  config.headers['my-header'] = 'jack';
  return config;
};

http.interceptors.request.use(beforeRequest);

// 响应拦截器
const responseSuccess = (response) => {
  // eslint-disable-next-line yoda
  // 这里没有必要进行判断，axios 内部已经判断
  // const isOk = 200 <= response.status && response.status < 300
  return Promise.resolve(response.data);
};

const responseFailed = (error) => {
  const { response } = error;
  if (response) {
    // handleError(response)
    logInfo(response);
    // cons error = new Error(response.data.msg)
    return Promise.reject();
  } else if (!window.navigator.onLine) {
    redLog('没有网络');
    return Promise.reject(new Error('请检查网络连接'));
  }
  return Promise.reject(error);
};
http.interceptors.response.use(responseSuccess, responseFailed);

export default {
  get,
  post,
  delete: del,
};
