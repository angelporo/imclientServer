package control

import (
	"github.com/gin-gonic/gin"
	. "imClientServer/config"
	"github.com/xormplus/xorm"
	// "encoding/json"
)

type (
	MembersItem struct {
		Id int64 `json:"-"`
		Index int `xorm:"createuniques -> notnull" json:"-"`
		Owner string `xorm:"varchar(255)" json:"owner"`
		Member string `xorm:"varchar(255)" json:"member"`
	}

	// 最近聊天群信息 包括成员用户名列表
	GroupMembersContent struct {
		Id string `json:"id"`
		Index int `xorm:"pk autoincr unique notnull" json:"-"`
		Name string `xorm:"varchar(255)  notnull" json:"name"`
		Description string `xorm:"varchar(255) " json:"description"`
		//群组类型：true：公开群，false：私有群。
		Public bool `xorm:"bool " json:"public"`
		//加入群组是否需要群主或者群管理员审批。true：是，false：否。
		Membersonly bool `xorm:"bool " json:"membersonly"`
		//是否允许群成员邀请别人加入此群。 true：允许群成员邀请人加入此群，false：只有群主才可以往群里加人。
		Allowinvites bool `xorm:"bool " json:"allowinvites"`
		// 最大成员数
		Maxusers int64 `xorm:"bigint notnull " json:"maxusers"`
		// 现有成员数
		Affiliations_count int64 `xorm:"binint notnull " json:"affiliations_count"`
		Affiliations []MembersItem `json:"affiliations"`
		// 群主id
		Owner string `xorm:"varchar(255) notnull " json:"owner"`
		//群成员的环信 ID
		Member string `xorm:"varchar(255) "`
		Invite_need_confirm bool `xorm:"bool" json:"invite_need_confirm"`
		GroupAvatar string `xorm:"varchar(255) notnull" json:"groupAvatar"`
	}

	// 获取群组详情 包括成员列表
	GetGroupInfo struct {
		Data GroupMembersContent `json:"data"`
		// 未知属性用interface{}代替内部内容
		Entities []interface{} `json:"entities"`
		Uri string `json:"uri"`
		Params map[string]string `json:"params"`
		Application string `json:"application"`
		Action string `json:"action"`
		TimeStamp int `json:"timestamp"`
		Duration int `json:"duration"`
		Organization string `json:"organization"`
		AppLicationName string `json:"applicationName"`
	}

	GetGroupMembers struct {
		UserNames []map[string]string `json:"userName" binding:"required"`
	}

)


// 获取用户列表信息
// http POST http://localhost:8080/getusersInfo
func GetUsersInfo (c *gin.Context) {
	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}
	var (
		members GetGroupMembers
	)
	err := c.Bind(&members)
	if err != nil {
		rData(-1, "参数type 错误", err.Error())
		return
	}
	engine, openDbErr := xorm.NewEngine("mysql", DATABASE_LOGIN)
	if openDbErr != nil {
		rData(-1, "打开数据库错误", openDbErr.Error())
		return
	}
	resultEntity := make([]User, len(members.UserNames))
	for index, member := range members.UserNames {
		var res User
		_, finErr := engine.Where("name = ?", member["userName"]).Cols("mobile", "name", "sex", "age", "avatar").Get(&res)
		if finErr != nil {
			rData(-1, "查询失败", finErr.Error())
			return
		}
		resultEntity[index] = res
	}
	rData(0, "", resultEntity)
}
