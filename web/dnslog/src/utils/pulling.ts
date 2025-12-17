import { fetchDns } from '@/views/dns-query/index';

let pollingTimer: number | null = null;
let pollingInterval = Number(import.meta.env.VITE_POLL_INTERVAL_MS || 2000);

export function setPollingInterval(ms: number) {
  if (Number.isFinite(ms) && ms >= 500) {
    pollingInterval = ms;
  }
}

export function startPolling() {
  console.log('Polling started...');
  if (pollingTimer) {
    clearInterval(pollingTimer);
  }
  pollingTimer = setInterval(() => {
    fetchDns();
  }, pollingInterval);
}

export function stopPolling() {
  console.log('Polling stopped.');
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
}
