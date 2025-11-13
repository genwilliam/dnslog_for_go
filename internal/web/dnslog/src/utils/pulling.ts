import fetchDns from '@/views/result/index';

let pollingTimer: number | null = null;
export function startPolling() {
  console.log('Polling started...');
  pollingTimer = setInterval(() => {
    fetchDns();
  }, 2000);
}

export function stopPolling() {
  console.log('Polling stopped.');
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
}
