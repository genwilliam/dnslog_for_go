import { fetchTokenStatus } from '@/views/dns-query/index';

let pollingTimer: number | null = null;
let pollingInterval = Number(import.meta.env.VITE_POLL_INTERVAL_MS || 2000);
let pollingStartedAt: number | null = null;
const maxPollingMs = Number(import.meta.env.VITE_POLL_MAX_MS || 300000);

function tick() {
  if (pollingStartedAt && maxPollingMs > 0) {
    if (Date.now() - pollingStartedAt >= maxPollingMs) {
      stopPolling();
      return;
    }
  }
  fetchTokenStatus();
}

export function setPollingInterval(ms: number) {
  if (Number.isFinite(ms) && ms >= 500) {
    pollingInterval = ms;
    if (pollingTimer) {
      clearInterval(pollingTimer);
      pollingTimer = setInterval(() => {
        tick();
      }, pollingInterval);
    }
  }
}

export function startPolling() {
  console.log('Polling started...');
  if (pollingTimer) {
    clearInterval(pollingTimer);
  }
  pollingStartedAt = Date.now();
  tick();
  pollingTimer = setInterval(() => {
    tick();
  }, pollingInterval);
}

export function stopPolling() {
  console.log('Polling stopped.');
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
  pollingStartedAt = null;
}
