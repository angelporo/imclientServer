// 获取最近联系人信息以及聊天内容 通过用户id
package control
import (
	"github.com/xormplus/xorm"
	"errors"
	"time"
	"fmt"
)

// 创建最近联系人表结构
type RecentConcat struct {
	Id int64
	RoomType int64 `xorm:"notnull int"` // 聊天房间类型(1:单聊, 2群聊)
	UserName string `xorm:"varchar(255) notnull"`
	RoomId string `xorm:"varchar(255) notnull"` //单聊(用户名称字段) 群聊(群聊id)
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

// 最近联系聊天室内容以及聊天记录
type ChatRoomItem struct {
	ChatRoomInfo ChatGroupInfo `json:"chatRoomInfo"` // 聊天室内容
	ChatRoomHistory []ChatHistory `json:"chatRoomHistory"` // 聊天室聊天记录
	GroupMembers []GroupMembersContent `json:"members"` // 群组成员列表
	ChatRoomMembers []interface{} `json:"chatRoomMembers"` // 聊天室成员列表
}

// 获取最近联系人数据 通过用户名称
func GetRecentConcatById (userName string, engine *xorm.Engine) ([]ChatRoomItem, error) {

	var (
		RecentConcatInfo []RecentConcat
		r []ChatRoomItem
	)

	// 获取用户最近聊天室id
	err1 := engine.Where("user_name = ?", userName).Cols("room_id").Find(&RecentConcatInfo)
	if err1 != nil {
		return r, errors.New(err1.Error())
	}

	RecentSum := len(RecentConcatInfo)

	res := make([]ChatRoomItem, RecentSum)


	errD := engine.Sync2(new(ChatHistory))
	if errD != nil {

		return r, errors.New(errD.Error())
	}
	// 遍历用户最近联系用户id
	for i := 0; i < RecentSum; i++ {
		// TODO: RoomId是1 查询RoomId(userName)用户信息
		// TODO: RoomId是2 查询群详情
		roomId := RecentConcatInfo[i].RoomId
		fmt.Println(roomId)
		// 通过roomid查找聊天室信息
		has, err2 := engine.Table("chat_group_info").Where("id = ?", roomId).Get(&res[i].ChatRoomInfo)
		if err2 != nil {
			return r, errors.New(err2.Error())
		}
		if !has {
			// return r, errors.New("没有查到聊天室信息")
		}
		// 获取聊天记录
		err3 := engine.Table("chat_history").Where("room_id = ?", roomId).Find(&res[i].ChatRoomHistory)
		if err3 != nil {
			return r, errors.New(err3.Error())
		}

	}
	// TODO: 根据聊天房间类型 获取用户聊天内容
	return res, nil
}
