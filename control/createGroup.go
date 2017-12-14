package control

import (
	. "imClientServer/config"
	. "imClientServer/util"
	"github.com/gin-gonic/gin"
	"github.com/xormplus/xorm"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"errors"
	"strings"
	"strconv"
	"fmt"
)

type CreateGoupRequest struct {
	UserId string `json:"userId" binding:"required"`
	Groupname string `json:"groupname" binding:"required"` // 群名称
	Desc string `json:"desc" binding:"required"` // 群组描述
	Maxusers int `json:"maxusers"` //群组成员最大数（包括群主），值为数值类型，默认值200，
	Members_only bool `json:"membersOnly"` //加入群是否需要群主或者群管理员审批，默认是false
	Public bool `json:"public" binding:"required"` // 是否为公开群
	Allowinvites bool `json:"allowinvites"` //是否允许群成员邀请别人加入此群
	Owner string `json:"owner" binding:"required"` // 群组管理员
	Members []string `json:"members" binding:"required"` // 群组成员
}

type CreateGroupHXresponse struct {
	HxSuccessData
	UserId string `json:"userId"`
	Data map[string]string `json:"data"`
}

// 用户参与的群聊
type GroupRelationShip struct {
	Index int64 `xorm:"notnull bigint pk autoincr"`
	UserName string `xorm:"varchar(255) notnull"`
	UserId int64 `xorm:"bigint notnull"`
	GroupRoomId string `xorm:"varchar(14) notnull"`
	GroupAvatar string `xorm:"varshar(255) notnull"`
	MyNickName string `xormj:"notnull varchar(255)"`
}


// 发送创建群聊到环信服务器
func CreatGroupToHX (requestBody CreateGoupRequest) ([]byte, error, string) {
	client := &http.Client{}
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME + "/chatgroups"
	b, _ := json.Marshal(requestBody)
	body := strings.NewReader(string(b))
	req, err_req := http.NewRequest("POST", path, body)
	if err_req != nil {
		return []byte(""), errors.New("发送创建群聊失败"), ""
	}
	token , getTokenErr := GetToken()
	if getTokenErr != nil {
		return []byte(""), errors.New("创建群聊Token出错"), ""
	}
	req.Header.Add("Authorization", "Bearer " + ToStr(token))
	resp, resp_err := client.Do(req)
	if resp_err != nil {
		return []byte(""), errors.New("创建群聊发送request失败"), ""
	}
	defer resp.Body.Close()
	r, _ := ioutil.ReadAll(resp.Body)
	return r, nil, resp.Status
}


// // 获取群组详情
func GetGroupMembersListById (groupId string) ([]byte, error, string){
	client := &http.Client{}
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME + "/chatgroups/" + groupId
	req, err_req := http.NewRequest("GET", path, nil)
	if err_req != nil {
		return []byte(""), errors.New("发送获取群组详情失败"), ""
	}
	token , getTokenErr := GetToken()
	if getTokenErr != nil {
		return []byte(""), errors.New("获取群聊详情Token出错"), ""
	}
	req.Header.Add("Authorization", "Bearer " + ToStr(token))
	resp, resp_err := client.Do(req)
	if resp_err != nil {
		return []byte(""), errors.New("获取群聊详情发送request失败"), ""
	}
	defer resp.Body.Close()
	byteB, _ := ioutil.ReadAll(resp.Body)

	return byteB, nil, resp.Status
}


// 群组聊天详情类型
type GroupInfoContent struct {
	HxSuccessData
	Data []GroupMembersContent `json:"data"`
}


// 写入群组聊天详情群组详情数据表
func WriteGroupInfo (groupId string)  {

	var tokenType GroupInfoContent

	// 获取群组详情信息
	groupInfo, getGroupErr, statusStr := GetGroupMembersListById(groupId)
	if getGroupErr != nil {
		fmt.Printf("获取群组聊天详情错误! %v%v", Substr(statusStr, 0, 3), getGroupErr.Error())
		return
	}

	// 获取环信群组聊天详情服务器出错, 暂不考虑groupRoomId错误情况
	if Substr(statusStr, 0, 1) == "5" {
		// 写入日志
		fmt.Printf("获取群组聊天详情错误! %v%v",Substr(statusStr, 0, 3) , getGroupErr.Error())
		return
	}
	// 状态码没有出错, 暂不考虑环形返回状态码
	marshErr := json.Unmarshal(groupInfo, &tokenType)
	if marshErr != nil {
		fmt.Printf("序列化群组详情出错:%v \n groupInfo内容:%v", marshErr.Error(), groupInfo)
		return
	}

	// 提取群消息数据
	tokenType.Data[0].GroupAvatar = "/files/group/groupavatar.png"
	engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
	if err != nil {
		fmt.Printf("打开数据库出错: %v", err.Error())
		return
	}

	_ , writeErr := engine.Insert(&tokenType.Data)
	if writeErr != nil {
		fmt.Printf("写入群组聊天详情出错: %v", writeErr.Error())
		return
	}

}


// 用户创建好友群组聊天
// http POST http://localhost:8080/creategroup groupname=test desc=测试群组聊天 public:=true maxusers:=400 allowinvites:=true owner=liyuan members:='["liyuan1","liyuan2"]' membersOnly:=false userId=68
func CreateGoup (c *gin.Context) {
	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}

	var creReqBody CreateGoupRequest
	bindErr := c.Bind(&creReqBody)
	if bindErr != nil{
		rData(-1, "创建用户群聊绑定request参数错误!", bindErr.Error())
		return
	}
	if creReqBody.Groupname == "" {
		rData(-1, "用户名称必填!", "")
		return
	}
	if creReqBody.Desc == "" {
		rData(-1, "群描述!", "")
		return
	}
	if creReqBody.Owner == "" {
		rData(-1, "群管理员不可少!", "")
		return
	}

	// 获取环信返回response结果
	resp, respErr, statusStr := CreatGroupToHX(creReqBody)
	if respErr != nil {
		rData(-1, "发送创建群聊信息失败, 请重试!", respErr.Error())
		return
	}
	status := Substr(statusStr, 0, 3)
	if status == "400" {
		rData(-1, "管理员不存在", "")
		return
	}
	if status == "401" {
		rData(-1, "token出错", "")
		return
	}
	var  HXresponse CreateGroupHXresponse
	marshallErr := json.Unmarshal(resp, &HXresponse)
	if marshallErr != nil {
		rData(-1, "创建群聊序列化response错误!", marshallErr.Error())
		return
	}
	// 环信创建成功
	// 写入用户参与群聊数据表中
	userId64 , _ := strconv.ParseInt(creReqBody.UserId, 10, 64)
	userGroups := &GroupRelationShip{
		UserName: creReqBody.Owner,
		UserId: userId64,
		GroupRoomId: HXresponse.Data["groupid"],
		GroupAvatar: "/files/group/groupavatar.png",
		MyNickName: creReqBody.Owner,
	}
	// 写入群组关系数据库
	engine, openDbErr := xorm.NewEngine("mysql", DATABASE_LOGIN)
	if openDbErr != nil {
		rData(-1, "创建群组聊天打开数据库失败", openDbErr.Error())
		return
	}
	_, insertErr := engine.Insert(userGroups)
	if insertErr != nil {
		rData(-1, "写入聊天群组失败", insertErr.Error())
		return
	}
	go WriteGroupInfo(userGroups.GroupRoomId)
	rData(0, "", userGroups)
}
