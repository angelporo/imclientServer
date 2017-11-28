package control

type MembersItem struct {
	Id int64 `json:"-"`
	Index int `xorm:"createuniques -> notnull" json:"-"`
	Owner string `xorm:"varchar(255)" json:"owner"`
	Member string `xorm:"varchar(255)" json:"member"`
}

// 最近聊天群信息 包括成员用户名列表
type GroupMembersContent struct {
	Id string
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
	Affiliations []MembersItem
	 // 群主id
	Owner string `xorm:"varchar(255) notnull " json:"owner"`
	//群成员的环信 ID
	Member string `xorm:"varchar(255) "`
	Invite_need_confirm bool `xorm:"bool "`
}

// 获取群组详情 包括成员列表
type GetGroupInfo struct {
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
