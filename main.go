package main

import (
	"github.com/gin-gonic/gin"
	"imClientServer/control"
	// "imClientServer/util"
	"net/http"
	// "log"
)



// 检查管理员token是否过期
// 过期重新获取
// 不过期直接Next()
func HuanxingTokenMiddleWare () gin.HandlerFunc {
	return func(c *gin.Context) {
		// expires_in := 234353253
		// nowTime := time.Time()
		// 获取环信服务器管理员token
		// 超时 重新获取
		// token , err := control.GetToken();
		// if err != nil {
		//	fmt.Println(err)
		// }
		// i := util.ToStr(token)
		c.Next()
	}
}

func main() {
	// Creates a router without any middleware by default
	r := gin.New()
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// the change user content
	// get user all info
	// authorized := r.Group("/")
	// 开启static file

	// router := gin.Default()
	r.StaticFS("/files", http.Dir("./static"))
	// r.Static("/files", "./assets")
	// 给指定文件定制路由
	// r.StaticFile("/image", "./assets/1.png")

	r.Use(HuanxingTokenMiddleWare())
	{
		r.POST("/login", control.GetUserInfo) // 登录接口 ,  登录后返回相关数据
		// the create new user
		r.PUT("/user", control.RegisterListen) // 注册接口
		r.POST("/sendmsg", control.ListenSendMsg) // 发送消息
		r.POST("/addfriend", control.AddFriendToUser) // 给用户添加好友
		r.POST("/creategroup", control.CreateGoup) // 用户创建群组聊天
		r.POST("/getgroupinfo", control.GetGroupInfoByGroupId) // 获取用户群组聊天详情
	}
	r.Run(":8080")
}
