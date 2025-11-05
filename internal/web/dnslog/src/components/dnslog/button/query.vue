<template>
  <n-space>
    <n-button type="success" size="large" @click="fetchDns">Record DNS query</n-button>
  </n-space>
</template>

<script setup lang="ts">
import { submitDomain, domain } from '@/components/dnslog/input/use-domain.ts';

function query() {
  var demo = submitDomain();
  console.log(demo);
}

function fetchDns(domain) {
  fetch('http://localhost:8080/submit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ domain_name: domain }),
  })
    .then((response) => response.json())
    .then((data) => {
      const resultDiv = document.getElementById('result');
      if (!resultDiv) {
        throw new Error('Result div not found');
      }

      if (data.error) {
        resultDiv.innerHTML = `<p style="color:red;">错误: ${data.error}</p>`;
      } else {
        let tableHtml = `
                <table border="1" style="border-collapse: collapse; width: 100%; margin-top: 20px;">
                    <thead>
                        <tr style="background-color: #f2f2f2;">
                            <th>域名</th>
                            <th>IP 地址</th>
                            <th>DNS 服务器</th>
                        </tr>
                    </thead>
                    <tbody>
            `;

        data.results.forEach((result) => {
          tableHtml += `
                    <tr>
                        <td>${data.domain}</td>
                        <td>${result.ip}</td>
                        <td>${result.address}</td>
                    </tr>
                `;
        });

        tableHtml += `</tbody></table>`;
        resultDiv.innerHTML = tableHtml;
      }
    })
    .catch((error) => {
      console.error('请求失败:', error);
    });
}
</script>
