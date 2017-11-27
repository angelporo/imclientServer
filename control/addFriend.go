// 添加好友接口
// 环信: 404（此IM用户或被添加的好友不存在）、401（未授权[无token、token错误、token过期]）
package control

import (
	"github.com/gin-gonic/gin"
	"errors"
	. "imClientServer/config"
	. "imClientServer/util"
	"github.com/xormplus/xorm"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

// 给用户添加好友到环信服务器
func HXAddFriendGiveUsers (friendUserName string, targetUserName string) ([]byte, error, string) {
	client := &http.Client{}
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME + "/users/" + targetUserName + "/contacts/users/" + friendUserName
	req, err_req := http.NewRequest("POST", path, nil)

	token, getTokenErr := GetToken()
	if getTokenErr != nil {
		return []byte(""), errors.New(ToStr(getTokenErr)), ""
	}
	tokenStr := "Bearer " + ToStr(token)
	req.Header.Set("Authorization", tokenStr)

	if err_req != nil {
		return []byte(""), errors.New(err_req.Error()), ""
	}
	resp, err_resp := client.Do(req)
	if err_resp != nil {
		return []byte(""),  errors.New("获取token错误"), ""
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	return res, nil, resp.Status
}

type AddFriendRequestBody struct {
	UserId string `jaon:"userName"`
	UserName string `json:"userName"`
	FriendName string `json:"friendName"`
}

// 环信添加好友返回response主体内容
type HXresponseContent struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
	UserName string `json:"userName"`
	Activated bool `json:"activated"`
}

type AddSuccessResponse struct {
	HxSuccessData
	DataFaild
	Entities []HXresponseContent
}


// http POST http://localhost:8080/addfriend userName="angelporo" friendName="angelporo1" userId="68"
func AddFriendToUser (c *gin.Context) {
	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}
	// 绑定用户request
	var requestBody AddFriendRequestBody
	bindErr := c.Bind(&requestBody)
	if bindErr != nil {
		rData(-1, "添加好友绑定request参数错误!", bindErr.Error())
		return
	}
	if requestBody.UserName == "" {
		rData(-1, "参数错误", "用户名不能为空!")
		return
	}
	if requestBody.FriendName == "" {
		rData(-1, "参数错误", "添加的好友用户名不能为空")
		return
	}
	resp , addErr, statusStr := HXAddFriendGiveUsers(requestBody.FriendName, requestBody.UserName)
	status := Substr(statusStr, 0, 3)
	if addErr != nil {
		CheckHXErr(c, addErr.Error(), "添加好友发送环信服务器失败")
		return
	}
	if status == "404" {
		rData(404, "添加好友出错", "被添加的好友不存在")
		return
	}
	if status == "401" {
		rData(401, "添加好友出错", "token错误或者过期或者未授权!")
		return
	}
	if Substr(statusStr, 0, 1) == "5" {
		rData(500, "环信服务器发生错误, 稍后重试", "")
		return
	}
	var response AddSuccessResponse
	marshallErr := json.Unmarshal(resp, &response)
	if marshallErr != nil {
		rData(-1, "添加好友序列化response错误!", marshallErr.Error())
		return
	}

	getUserId64, _ := strconv.ParseInt(requestBody.UserId, 10, 64)
	// 环信添加好友成功, 写入用户好友关系库
	userFriend := &UserRelationShip{
		UserId: getUserId64,
		UserName: requestBody.UserName,
		Friend_userName: requestBody.FriendName,
	}
	engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
	if err != nil {
		rData(-1, "添加好友打开数据库出错, 请重试!", err.Error())
		return
	}

	// 插入之前先检测是否成为好友
	has, err := engine.Table("user_relation_ship").Where("user_name = ?", requestBody.UserName).And("friend_user_name = ? ", requestBody.FriendName).Exist()
	if err != nil {
		rData(-1, "检测好友是否存在失败", err.Error())
		return
	}
	if has {
		rData(-1, "被添加用户已经是您好友", "")
		return
	}

	// 写入数据库
	_ , insertFriendErr := engine.Insert(userFriend)
	if insertFriendErr != nil {
		rData(-1, "写出好友数据库出错,请重试!", insertFriendErr.Error())
	}
	rData(0, "", response)
}
