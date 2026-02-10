const normalizeBaseURL = (val: string | undefined) => {
  if (!val) {
    return '/api';
  }
  const trimmed = val.replace(/\/+$/, '');
  if (trimmed === '' || trimmed === '/') {
    return '/api';
  }
  if (/^https?:\/\//.test(trimmed) || trimmed.startsWith('//')) {
    if (trimmed.endsWith('/api')) {
      return trimmed;
    }
    return `${trimmed}/api`;
  }
  const withSlash = trimmed.startsWith('/') ? trimmed : `/${trimmed}`;
  if (withSlash.endsWith('/api')) {
    return withSlash;
  }
  return `${withSlash}/api`;
};

const baseURL = normalizeBaseURL(import.meta.env.VITE_API_BASE_URL);
const PULL_REQUEST = import.meta.env.VITE_PULL_REQUEST;
export { baseURL, PULL_REQUEST };
