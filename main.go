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
	router := gin.Default()
	router.Static("/", "./")
	router.StaticFS("/avatar", http.Dir("user_img"))
	// router.StaticFile("/default_avatar", "default_avatar.jpg")

	r.Use(HuanxingTokenMiddleWare())
	{
		r.POST("/login", control.GetUserInfo) // 登录接口 ,  登录后返回相关数据
		// the create new user
		r.PUT("/user", control.RegisterListen) // 注册接口
		r.Run(":8080")
	}
	r.Run(":8080")
	router.Run(":8080")
}
