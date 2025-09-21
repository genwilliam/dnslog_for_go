package domain

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var commonTLDs = []string{
	".com", ".net", ".org", ".cn", ".io", ".edu", ".gov", ".co", ".xyz",
}

// 基于uuid生成的域名，长度限制在10个字母之中
func TestGeneratingDomain(t *testing.T) {
	id := uuid.New().String()

	// 去掉 - 号
	cleaned := strings.ReplaceAll(id, "-", "")

	// 截取前 10 个字符作为域名主体
	if len(cleaned) < 10 {
		t.Fatalf("UUID 过短，不足10字符: %s", cleaned)
	}
	shortDomain := cleaned[:10]

	// 输出生成的短域名
	fmt.Println("生成的短域名为:", shortDomain)

	// 可选断言：校验长度
	if len(shortDomain) > 10 {
		t.Errorf("域名长度超过限制: %s", shortDomain)
	}

	// 添加后缀，在commonTLDs中随机选择
	i := rand.Intn(9)

	tld := commonTLDs[i]
	domainTest := fmt.Sprintf("%s%s", shortDomain, tld)
	// 输出完整的域名
	fmt.Println("完整的域名为:", domainTest)

	route(domainTest)

}
func route(domain string) {
	r := gin.Default()
	r.POST("/random-domain", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"domain": domain,
		})
	})
	fmt.Println(domain)
	r.Run(":8080")
}
