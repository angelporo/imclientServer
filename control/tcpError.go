// 处理环信服务器返回错误!
package control

import (
	"github.com/gin-gonic/gin"
)

func ReturnErr (c *gin.Context, err string) {
	switch err {
	case "invalid_grant":
		c.JSON(200, gin.H{
			"code": -2,
			"msg": "用户名或者密码输入格式错误",
			"content": "",
		})
	case "organization_application_not_found":
		c.JSON(200, gin.H{
			"code": -2,
			"msg": "找不到aachatdemoui对应的APP，可能是URL写错了",
			"content": "",
		})
	case "json_parse":
		c.JSON(200, gin.H{
			"code": -2,
			"msg": "发送请求时请求体不符合标准的JSON格式，服务器无法正确解析",
			"content": "",
		})
	case "duplicate_unique_property_exists":
		c.JSON(200, gin.H{
			"code": -2,
			"msg": "用户名已存在",
			"content": "",
		})
	case "unauthorized":
		c.JSON(200, gin.H{
			"code": -2,
			"msg": "APP的用户注册模式为授权注册，但是注册用户时请求头没带token",
			"content": "",
		})
	case "auth_bad_access_token":
		c.JSON(200, gin.H{
			"code": -2,
			"msg": "发送请求时使用的token错误。注意：不是token过期",
			"content": "",
		})
	case "service_resource_not_found":
		c.JSON(200, gin.H{
			"code": -2,
			"msg":"URL指定的资源不存在",
			"content": "",
		})
	case  "reach_limit":
		c.JSON(200, gin.H{
			"code": -2,
			"msg":"	超过接口每秒调用次数，加大调用间隔或者联系商务调整限流大小",
			"content": "",
		})
	case "no_full_text_index":
		c.JSON(200, gin.H{
			"code": -2,
			"msg":"username不支持全文索引，不可以对该字段进行contains操作",
			"content": "",
		})
	case  "unsupported_service_operation":
		c.JSON(200, gin.H{
			"code": -2,
			"msg":"请求方式不被发送请求的URL支持",
			"content": "",
		})
	case "web_application":
		c.JSON(200, gin.H{
			"code": -2,
			"msg":"	错误的请求，给一个未提供的API发送了请求",
			"content": "",
		})
	default:
		c.JSON(200, gin.H{
			"code": -2,
			"msg":"环信服务器未知错误!",
			"content": "",
		})
	}
}
