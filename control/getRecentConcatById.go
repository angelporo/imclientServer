// 获取最近联系人信息以及聊天内容 通过用户id
package control
import (
	"github.com/xormplus/xorm"
	"errors"
	"time"
)

// 创建最近联系人表结构
type RecentConcat struct {
	Id int64
	RoomType int64 `xorm:"notnull int"` // 聊天房间类型(1:单聊, 2群聊)
	UserName string `xorm:"varchar(255) notnull" json:"userName"`
	TargetUserName string `xorm:"varchar(255) notnull" json:"targetUserName"` //单聊(用户名称字段) 群聊(群聊id)
	IsTop bool `xomr:"bool notnull" json:"isTop"`
	LastMessage string `xorm:"varchar(255)" json:"lastMessage"`
	LastMsgUpdated time.Time `xorm:"datetime"`
}

// 聊天室和群组id
// type ChatRoomGroupAffiliations struct {
//	Id int64
//	RoomId string `xorm:"notnull"` // 房间id
//	UserId int64 `xorm:"notnull bigint"` // 房间内成员id
// }

// 群组信息 结构 以及返回json字段配置
type ChatGroupInfo struct {
	Id int64 `json:"id"`
	// 是否顶置聊天内容
	IsTop bool `xorm:"notnull bool" json:"isTop"`
	// 聊天名称
	Name string `xorm:"notnull varchar(255)" json:"name"`
	// 描述
	Description string `xorm:"varchar(255)" json:"description"`
	// 最大成员
	Maxusers int `xorm:"notnull int" json:"maxusers"`
	// 现有成员总数
	AffiliationsCount int `xorm:"notnull int" json:"affiliationsCount"`
	// 创建聊天室name
	Owner string `xorm:"notnull varchar(255)" json:"owner"`
	// 单聊为1, 多聊为 n++
	Affiliations int `xorm:"notnull int" json:"affiliations"`
	//单聊成员
	Member string `xorm:"notnull varchar(255)" json:"member"`
	// 最后聊天内容
	LastMessage string `xorm:"varchar(255)" json:"lastMessage"`
	 //最后聊天内容发送时间
	LastMessageTime time.Time `xorm:"datetime" json:"lastMessageTime"`
	Avatar string `xorm:"varchar(255) notnull" json:"avatar"`
	// 聊天室(1) 群组(2)
	Type int `xorm:"int notnull" json:"type"`
}

// 创建聊天记录表结构
type ChatHistory struct {
	Id int64 `json:"chatId"`
	RoomId int64 `xorm:"bigint notnull" json:"roomId"`
	TargetId int64 `xorm:"bigint notnull" json:"targetId"`
	Type string `xorm:"Varchar(16)" json:"type"`
	Created time.Time `xorm:"created notnull" json:"createdTime"`
	Content string `xorm:"Varchar(255)" json:"chatContent"`
	Avatar string `xorm:"varchar(255)" json:"avatar"`
}

type ApiRecentChatInfo struct{
	IsTop bool `json:"isTop"`
	ChatType *string `json:"chatType"`
	ChatGroup GroupMembersContent `json:"chatGroup"`
	User User `json:"user"`
	RecentKey int64 `json:"recentKey"`
	LastMessage string `json:"lastMsg"`
	LastMsgUpdated time.Time `json:"lastMsgUpdated"`
}


// 最近联系聊天室内容以及聊天记录
type ChatRoomItem struct {
	ChatRoomHistory []ChatHistory `json:"chatRoomHistory"` // 聊天室聊天记录
	GroupMembers ApiRecentChatInfo `json:"members"` // 群聊和单聊具体详情
}

// 获取单聊最近联系人数据 通过用户名称
func GetRecentConcatById (userName string, engine *xorm.Engine) ([]ChatRoomItem, error) {

	var (
		RecentConcatInfo []RecentConcat
		r []ChatRoomItem
	)

	// 获取用户最近聊天室id
	err1 := engine.Where("user_name = ?", userName).Cols("target_user_name","room_type","is_top","last_message","last_msg_updated", "id").Find(&RecentConcatInfo)
	if err1 != nil {
		return r, errors.New(err1.Error())
	}

	RecentSum := len(RecentConcatInfo)

	res := make([]ChatRoomItem, RecentSum)

	for i := 0; i < RecentSum; i++ {
		// NOTE: RoomId是1 查询RoomId(userName)用户信息
		// NOTE: RoomId是2 查询群详情
		roomIdOrUserName := RecentConcatInfo[i].TargetUserName
		roomType := RecentConcatInfo[i].RoomType
		isTop := RecentConcatInfo[i].IsTop
		lastMsg := RecentConcatInfo[i].LastMessage
		lastMsgUpdaated := RecentConcatInfo[i].LastMsgUpdated
		key := RecentConcatInfo[i].Id

		var userType string = "users"
		var chatGroupType string = "chatgroups"
		var users int64 = 1
		var chatGroups int64 = 2
		if roomType == chatGroups{
			// 查询群聊详情
			// 查询群聊详情 group_members_content
			// 通过roomid查找聊天室信息
			res[i].GroupMembers.ChatType = &chatGroupType
			res[i].GroupMembers.IsTop = isTop
			res[i].GroupMembers.LastMessage = lastMsg
			res[i].GroupMembers.LastMsgUpdated = lastMsgUpdaated
			res[i].GroupMembers.RecentKey = key
			_, err2 := engine.Table("group_members_content").Where("id = ?", roomIdOrUserName).Get(&res[i].GroupMembers.ChatGroup)
			if err2 != nil {
				return r, errors.New(err2.Error())
			}
		}
		if roomType == users {
			// 查询用户信息 user表
			res[i].GroupMembers.ChatType = &userType
			res[i].GroupMembers.IsTop = isTop
			res[i].GroupMembers.LastMessage = lastMsg
			res[i].GroupMembers.LastMsgUpdated = lastMsgUpdaated
			res[i].GroupMembers.RecentKey = key
			has, err2 := engine.Where("name = ?", roomIdOrUserName).Cols("name","mobile","avatar","activated","uuid", "id").Get(&res[i].GroupMembers.User)
			if err2 != nil {
				return r, errors.New(err2.Error())
			}
			if !has {
				return r, errors.New("没有查到聊天室信息")
			}
		}
		// 获取聊天记录
	}
	// TODO: 根据聊天房间类型 获取用户聊天内容
	return res, nil
}
