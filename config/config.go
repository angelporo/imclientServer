// 项目全局config配置
// 包括数据库登录配置
// 错误码对照表
package config



// 系统错误码对照表
// 0 一切正常
// 10001 数据库查询出错
// 10002 参数与数据库不相符
// 10000 系统发生错误


// token 是否过期, 过期(1), 没有过期(0)

const (
	HUANXIN_DOMAIN = "https://a1.easemob.com/"
	TOKEN_FILE_NAME = "token.txt"
	DATABASE_LOGIN = "root:angel@tcp(localhost:3306)/xx?charset=utf8"
	CLIENT_ID = "YXA6Y9d0YJRtEeeDeXnjDm9J3g"
	CLIENT_SECRET = "YXA6aWIn9ZJ6G_CfLcqxSLox1jeWq7I"
	ORG_NAME = "1132170823115353"
	APP_NAME = "xinxin-imclient"
	MOBILE_ERR_MSG = "手机号码不正确!"
	USERNAME_ERR_MSG = "用户名称不能为空!"
)
