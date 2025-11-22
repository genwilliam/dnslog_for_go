function generateRandomDomain() {
    // 发送请求获取随机域名
    fetch('http://localhost:8080/random-domain')
        .then(response => response.json())  // 后端返回 JSON 数据
        .then(data => {
            const randomDomain = data.domain;  // 返回的数据结构是 { domain: 'example.com' }
            console.log(randomDomain)
        })
        .catch(error => {
            console.error('Error fetching domain:', error);
        });
}