package control

import (
	// "syscall"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	. "imClientServer/config"
	"github.com/xormplus/xorm"
	"errors"
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
	Id int64 `json:"id"`
	Avatar string `json:"avatar"`
	Name string `json:"name"`
}

type LoginData struct {
	User User
	Friend []userItem
	RecentConcat []ChatRoomItem
}

// 获取好友id
func GetUserFriendById (id int64, engine *xorm.Engine) ([]UserRelationShip , error){
	var friend []UserRelationShip
	err := engine.Where("id = ?", id).Find(&friend)
	if err != nil {
		return friend, errors.New(err.Error())
	}
	return friend, nil
}


// 1.检测用户名和密码
// 2.查询用户好友列表 得到id
// 3.查询好友信息
// 4.查询最近联系人
// 5.查询聊天内容
// http http://localhost:8080/login mobile='18303403747' passWord="angel"
func GetUserInfo (c *gin.Context) {
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
	engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
	engine.ShowSQL(true)
	engine.ShowExecTime(true)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "系统发生错误",
			"content": "",
		})
		return
	}
	// 同步数据表结构
	err = engine.Sync2(new(User))
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "同步数据表错误",
			"content": err.Error(),
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
			"content": res,
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
	id := res.User.Id
	friendIds, err := GetUserFriendById(id, engine)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -6,
			"msg": err.Error(),
			"content": friendIds,
		})
		return
	}
	// 获取好友全部信息
	friendInfo, err := GetFriendInfoById(friendIds, engine)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -6,
			"msg": err.Error(),
			"content": friendInfo,
		})
		return
	}
	// 获取用户好友信息
	res.Friend = friendInfo
	// 获取最近联系人信息
	recentContent, err := GetRecentConcatById(id, engine)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": err.Error(),
			"content": "",
		})
		return
	}
	res.RecentConcat = recentContent
	// 登录成功
	c.JSON(200, gin.H{
		"code": -1,
		"msg": "登录成功!",
		"content": res,
	})
}
