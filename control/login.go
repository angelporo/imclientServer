package control

import (
	// "syscall"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	. "imClientServer/config"
	"github.com/xormplus/xorm"
	// "context"
	"errors"
	// "reflect"
)


// 登录接口参数字段
type Login struct {
	PassWord string ` json:"passWord" binding:"required"`
	UserMobile string `json:"mobile" binding:"required"`
}

// 用户信息类型
type UserInfoData struct {
	Id int64
	Name string
	Mobile string
	Avatar string
	Money int
	Uuid string
}

// 好用信息类型
type userItem struct {
	Age int `json:"age"`
	Sex int `json:"sex"`
	Mobile string `json:"mobile"`
	Avatar string `json:"avatar"`
	Name string `json:"name"`
}

type LoginData struct {
	User User `json:"user"`
	Friend []userItem `json:"friend"`
	RecentConcat []ChatRoomItem `json:"recentConcat"`
}

// 获取好友userName
func GetUserFriendById (userName string, engine *xorm.Engine) ([]UserRelationShip , error ){
	// 协程中出错
	friendNames := make([]UserRelationShip, 1)
	var friend  []UserRelationShip
	// findErr := make(chan error, 1)
	// defer func () {
	//	if err := recover(); err != nil {
	//		panic(err)
	//	}
	// }()

	// go func () {
		err := engine.Where("user_name = ?", userName ).Find(&friend)
		if err != nil {
			return friendNames,errors.New(err.Error())
		}else{
			friendNames = friend
		}
	// }()

	// select{
	// case e := <-findErr:
	//	close(findErr)
	//	close(friendNames)
	//	return nil,e
	// case s := <-friendNames:
	//	close(findErr)
	//	close(friendNames)
	//	return s,nil
	// }
	return friendNames, nil
}


// 1.检测用户名和密码
// 2.查询用户好友列表 得到id
// 3.查询好友信息
// 4.查询最近联系人
// 5.查询聊天内容
// http http://localhost:8080/login mobile='18303403737' passWord="angel"
func GetUserInfo (c *gin.Context) {
	// defer func () {
	//	if err := recover(); err != nil {
	//		c.JSON(200, gin.H{
	//			"code": -3,
	//			"msg": err,
	//			"content": "",
	//		})
	//	}
	// }()

	var (
		// password string
		json Login
	)
	err := c.Bind(&json)
	// 检查参数
	if err != nil {
		if json.UserMobile == "" {
			c.JSON(200, gin.H{
				"code": -1,
				"data": json.UserMobile,
				"msg": "手机号不能为空!",
			})
			return
		}
	}
	// TODO: 用户手机号检测
	// TODO: 添加用户二维码用作转账和收款
	engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
	// engine.ShowSQL(true)
	engine.ShowExecTime(true)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "系统发生错误",
			"content": "",
		})
		return
	}
	// 检测登录
	var res LoginData
	has, err := engine.Where("mobile = ?", json.UserMobile).Get(&res.User)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "查询出错",
			"content": res,
		})
		return
	}
	if !has {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "没有这个手机号!",
			"content": "",
		})
		return
	}
	if res.User.PassWord != json.PassWord {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "密码不正确!",
			"content": "",
		})
		return
	}
	userName := res.User.Name
	friendUserNames, err := GetUserFriendById(userName, engine)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "获取用户好友用户名出错",
			"content": err.Error(),
		})
		return
	}

	// 获取好友全部信息
	friendInfo, err := GetFriendInfoByUserName(friendUserNames, engine)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "登录获取全部好友信息失败",
			"content": err.Error(),
		})
		return
	}
	// 获取用户好友信息
	res.Friend = friendInfo
	// 获取最近联系人详情信息
	recentContent, err := GetRecentConcatById(userName, engine)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "获取最近联系人信息失败",
			"content": err.Error(),
		})
		return
	}
	res.RecentConcat = recentContent
	// 登录成功
	c.JSON(200, gin.H{
		"code": 0,
		"msg": "登录成功!",
		"content": res,
	})
}
