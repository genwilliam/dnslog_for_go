// 全局变量用于存储轮询的定时器ID
let pollingInterval = null;

function init() {
  bindFormSubmit();
  setupGenerateDomainButton();
  ChangeDNSServer();
  ChangePact();
  stopPolling();
}

/**
 * 绑定表单提交事件
 * @returns {void}
 * @throws {Error} - 如果表单元素未找到，则抛出错误
 */
function bindFormSubmit() {
  const form = document.getElementById("dnslog-form");
  if (!form) {
    throw new Error("Form element not found");
  }

  form.addEventListener("submit", function (event) {
    event.preventDefault();

    const domainInput = document.getElementById("domain_name");
    if (domainInput && domainInput instanceof HTMLInputElement) {
      const domain = domainInput.value;

      // 停止之前的轮询
      if (pollingInterval) clearInterval(pollingInterval);

      // 开始新的轮询
      fetchDns(domain);

      pollingInterval = setInterval(() => {
        fetchDns(domain);
      }, 2000); // 每2秒请求一次
    } else {
      throw new Error(
        'Element with id "domain_name" is not an HTMLInputElement'
      );
    }
  });
}

/**
 * 从服务器获取DNS日志数据并更新UI
 * @param {string} domain - 要查询的域名
 */
function fetchDns(domain) {
  fetch("http://localhost:8080/submit", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ domain_name: domain }),
  })
    .then((response) => response.json())
    .then((data) => {
      const resultDiv = document.getElementById("result");
      if (!resultDiv) {
        throw new Error("Result div not found");
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
      console.error("请求失败:", error);
    });
}

/**
 * 设置生成域名按钮的点击事件
 */
function setupGenerateDomainButton() {
  const generateButton = document.getElementById("generate-domain-btn");
  if (!generateButton) {
    throw new Error("Generate domain button not found");
  }

  generateButton.addEventListener("click", function (event) {
    event.preventDefault();

    fetch("http://localhost:8080/random-domain", {
      method: "POST",
    })
      .then((response) => response.json())
      .then((data) => {
        const domainInput = document.getElementById("domain_name");
        if (!domainInput) {
          throw new Error("Domain input field not found");
        }
        domainInput.value = data.domain;
        fetchDns(data.domain);
      })
      .catch((error) => {
        console.error("Error fetching domain:", error);
      });
  });
}

/**
 * 更改DNS服务器
 */
function ChangeDNSServer() {
  const dnsSelect = document.getElementById("dns-select");
  if (!dnsSelect) {
    throw new Error("DNS select element not found");
  }

  dnsSelect.addEventListener("change", function () {
    const selectedValue = this.value;
    fetch("http://localhost:8080/change", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ num: parseInt(selectedValue) }),
    })
      .then((res) => res.text())
      .then((msg) => alert(msg))
      .catch((err) => console.error("请求失败:", err));
  });
}

/**
 * 更改协议
 */
function ChangePact() {
  const pactSelect = document.getElementById("pact");
  if (!pactSelect) {
    throw new Error("Pact select element not found");
  }

  pactSelect.addEventListener("change", function () {
    const selectPact = this.value.toLowerCase();
    fetch("http://localhost:8080/change-pact", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ pact: selectPact }),
    })
      .then((res) => res.json())
      .then((msg) => {
        if (msg.message) {
          alert(msg.message);
        } else {
          alert("发生错误：" + msg.error);
        }
      })
      .catch((err) => console.error("请求失败:", err));
  });
}

/**
 * 暂停轮询
 * @returns {void}
 * **/

function stopPolling() {
  const statusBtn = document.getElementById("pause");
  if (!statusBtn) {
    throw new Error("Pause button not found");
  }

  statusBtn.addEventListener("click", async function () {
    try {
      let action = statusBtn.innerText.toLowerCase(); // "start" or "pause"
      let url = `http://localhost:8080/${action}`;

      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: action }),
      });

      if (!response.ok) {
        throw new Error("Server error: " + response.status);
      }

      const msg = await response.text();
      alert(msg); // 弹窗显示服务器返回的信息

      // 切换按钮文字
      statusBtn.innerText = action === "pause" ? "start" : "pause";
    } catch (err) {
      console.error("请求失败:", err);
      alert("请求失败，请检查后端服务");
    }
  });
}

window.onload = init;
