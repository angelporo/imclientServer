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
	UserId int64 `xorm:"int notnull"`
	RoomId string `xorm:"varchar(255) notnull"`
}

// 群聊用户id 和房间id
type ChatRoomGroupAffiliations struct {
	Index int64 `xorm:"notnull bigint autoincr pk"`
	RoomId string `xorm:"notnull"` // 房间id
	UserId int64 `xorm:"notnull bigint"` // 房间内成员id
}
// 聊天室信息 结构 以及返回json字段配置
type ChatRoom struct {
	Id int64 `json:"id"`
	IsTop bool `xorm:"notnull bool" json:"isTop"` // 是否顶置聊天内容
	Name string `xorm:"notnull varchar(255)" json:"name"` // 聊天名称
	Description string `xorm:"varchar(255)" json:"description"` // 描述
	Maxusers int `xorm:"notnull int" json:"maxusers"` // 最大成员
	AffiliationsCount int `xorm:"notnull int" json:"affiliationsCount"` // 现有成员总数
	Owner string `xorm:"notnull varchar(255)" json:"owner"` // 创建聊天室name
	Affiliations int `xorm:"notnull int" json:"affiliations"` // 单聊为1, 多聊为 n++
	Member string `xorm:"notnull varchar(255)" json:"member"` //单聊成员
	LastMessage string `xorm:"varchar(255)" json:"lastMessage"` // 最后聊天内容
	LastMessageTime time.Time `xorm:"datetime" json:"lastMessageTime"` //最后聊天内容发送时间
}

// 创建聊天记录表结构
type ChatHistory struct {
	Id int64 `json:"chatId"`
	RoomId int64 `xorm:"bigint notnull" json:"roomId"`
	TargetId int64 `xorm:"bigint notnull" json:"targetId"`
	Type string `xorm:"Varchar(16)" json:"type"`
	Created time.Time `xorm:"created notnull" json:"createdTime"`
	Content string `xorm:"Varchar(255)" json:"chatContent"`
}

// 最近联系聊天室内容以及聊天记录
type ChatRoomItem struct {
	ChatRoom ChatRoom `json:"chatRoomInfo"` // 聊天室内容
	ChatRoomHistory []ChatHistory `json:"chatRoomHistory"` // 聊天室聊天记录
}

// 获取最近联系人数据 通过最近联系人id
func GetRecentConcatById (id int64, engine *xorm.Engine) ([]ChatRoomItem, error) {

	var (
		RecentConcatInfo []RecentConcat
		r []ChatRoomItem
	)

	// 同步最近联系人表结构
	err := engine.Sync2(new(RecentConcat))
	if err != nil {
		return r, errors.New(err.Error())
	}
	// 获取用户最近聊天室id号
	err1 := engine.Table("recent_concat").Where("user_id = ?", id).Cols("room_id").Find(&RecentConcatInfo)
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

		id := RecentConcatInfo[i].RoomId

		// 通过roomid查找聊天室信息
		has, err2 := engine.Table("chat_room").Where("id = ?", id).Get(&res[i].ChatRoom)
		if err2 != nil {
			return r, errors.New(err2.Error())
		}
		if !has {
			return r, errors.New("没有查到聊天室信息")
		}
		// 获取聊天记录
		err3 := engine.Table("chat_history").Where("room_id = ?", id).Find(&res[i].ChatRoomHistory)
		if err3 != nil {
			return r, errors.New(err3.Error())
		}

	}
	// TODO: 通过聊天室id获取第一页聊天记录
	return res, nil
}
