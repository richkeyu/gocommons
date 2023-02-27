package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/richkeyu/gocommons/config"
	"github.com/richkeyu/gocommons/plog"

	"github.com/gin-gonic/gin"
)

const GatewaySecretKey = "h7G93M1Swa" // 网关秘钥标识 用于访问后端访问时判断来源是否合法
const GatewayHeaderKey = "x-gateway-auth"

var gatewaySecretKeyHistory = []string{
	GatewaySecretKey,
} // 修改key的时候需要把旧秘钥加到此列表中防止上线过程中的请求被拒绝

type CustomerHandler func(c *gin.Context, isPass bool) (newIsPass bool)

// IpAuthMiddleWare 限制访问IP
func IpAuthMiddleWare(ipList []string, handler ...CustomerHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		isPass := false
		// 获取客户端IP
		// 预发/生产
		//   google云header X-Forwarded-For 取倒数第二个
		//   201.83.6.2, 111.206.233.194, 34.111.82.191
		//   201.83.6.2, 111.206.233.194, 34.111.82.191, 169.254.1.1,0.0.0.0
		//   为防止生产环境误拦截 如header包含网关标识不拦截只记日志
		// 测试在nginx层面做限制
		// 开发不做限制
		env := os.Getenv(config.AppEnvName)
		if env == config.ProdEnv || env == config.PreEnv {
			ip := getGoogleClientIp(c)
			if ip == "" {
				ip = c.ClientIP()
			}
			for i := range ipList {
				if ip == ipList[i] {
					isPass = true
					break
				}
			}
		} else {
			isPass = true
		}

		// 根据header标识判断防止网络环境变化导致误拦截业务系统
		if !isPass {
			for i, s := range gatewaySecretKeyHistory {
				if c.Request.Header.Get(GatewayHeaderKey) == s {
					plog.GetDefaultFieldEntry(nil).WithField("header", c.Request.Header).Infof("IpAuthMiddleWare from gateway ip deny: %d", i)
					isPass = true
					break
				}
			}
		}

		// 自定义的处理逻辑
		if len(handler) > 0 {
			for i := range handler {
				isPass = handler[i](c, isPass)
			}
		}

		if !isPass {
			c.String(http.StatusUnauthorized, "Gateway 401 Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 根据google网络环境获取倒数第二个进行判断
func getGoogleClientIp(c *gin.Context) string {
	header := c.Request.Header.Get("X-Forwarded-For")
	if len(header) == 0 {
		return ""
	}
	ips := strings.Split(header, ",")
	if len(ips) < 2 {
		return ""
	}
	return strings.TrimSpace(ips[len(ips)-2])
}
