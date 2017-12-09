// 发送消息接口
// 环信: 在调用程序中，如果返回429或503错误，说明接口被限流了，请稍微暂停一下并重试。请求体如果超过 5kb 会导致413错误，需要拆成几个更小的请求体重试，同时用户消息+扩展字段的长度在40k字节以内
package control

import (
	"github.com/xormplus/xorm"
	"github.com/gin-gonic/gin"
	"errors"
	. "imClientServer/config"
	. "imClientServer/util"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
)
// 消息属性类型
type (
	// 发送消息目标对象详情
	TargetInfo struct {
		ChatInfo User `json:"chatInfo"`
		GroupsInfo GroupMembersContent `json:"groupInfo"`
	}
	SendContent struct {
		Created int `json:"created"`
		From string `json:"from"`
		Modified int `json:"modified"`
		Msg SendMsg `json:"msg"`
		Target []string `json:"target"`
		Target_type string `json:"target_type"`
		Type string `json:"type"`
		FromAvatar string `json:"fromAvatar"`
		// FromAvatar string `json:"fromAvatar"`
		// TargetAvatar string `json:"targetAvatar"`
	}

	SendMsg struct{
		Type string `json:"type"`
		Msg string `json:"msg"`
	}

	HXSendMsgSuccess struct {
		Action string `json:"action"`
		Application string `json:"application"`
		Param map[string]string `json:"param"`
		Uri string `json:"uri"`
		Entities []SendContent `json:"entities"`
		Data map[string]string `json:"data"`
		Timestamp int `json:"timestamp"`
		Duration int `json:"duration"`
		Organization string `json:"organization"`
		ApplicationName string `json:"applicationName"`
	}

	// 环信返回json类型映射类型
	HXsendMsgContent struct {
		DataFaild
		HXSendMsgSuccess
	}

	// 消息内容类型
	MassageContent struct {
		Type string `json:"type" binding:"required"` // 消息内容类型
		Msg string `json:"msg" binding:"required"` // 消息内容
		FileName string `json:"fileName" binding:"required"` // 文件名称
		Size map[string]int `json:"size" binding:"required"` // 图片大小
		Secret string `json:"secret" binding:"required"` // 成功上传文件后返回的secret
	}

	ExtMsg struct {
		FromAvatar string `json:"fromAvatar" binding:"required"` // 发送者头像
		SendTime int `json:"sendTime" binding:"required"` // 发送者本机当前时间
		SendMsgUri string `json:"sendMsgUri"`
		// 添加发送数据类型
		TargetInfo TargetInfo `json:"targetInfo"`
	}

	// 环信服务器接受类型
	SendRequest struct {
		Ext ExtMsg `json:"ext"`
		TargetType string `json:"target_type" binding:"required"` // users 给用户发消息。chatgroups: 给群发消息，chatrooms: 给聊天室发消息
		Target []string `json:"target" binding:"required"` // 注意这里需要用数组，数组长度建议不大于20，即使只有一个用户，也要用数组 ['u1']，给用户发送时数组元素是用户名，给群组发送时.数组元素是groupid
		From string `json:"from" binding:"required"`  //表示消息发送者。无此字段Server会默认设置为"from":"admin"，有from字段但值为空串("")时请求失败
		Msg MassageContent `json:"msg"`
		RecentId string `json:"recentId" binding:"required"`
	}
)


// 发送图片消息 通过消息类型分发
func(requestBody *SendRequest) SendMsgToHXServer () ([]byte, error){
	msgType := requestBody.Msg.Type
	var(
		res []byte
		err error
	)
	switch msgType{
	case "txt":
		resp, resp_err := requestBody.sendTxtMsgToHuanxinServer()
		res = resp
		err = resp_err
	case "img":
		resp, resp_err := requestBody.sendTxtMsgToHuanxinServer()
		res = resp
		err = resp_err
	default:
		return []byte(""), errors.New("消息类型不确定")

	}
	return res, err
}


// 发送图片消息到环信服务器
// func (requestBody *SendRequest) sendImgMsgHXServer () ([]byte, error) {

// }


// // 接受客户端发送来的图片消息内容
// func(requestBody *SendRequest) getClientImg () ([]byte, error){

// }


// 发送消息到环信服务器
func (requestBody *SendRequest) sendTxtMsgToHuanxinServer () ([]byte, error) {

	request := requestBody
	requestByte, _ := json.Marshal(request)

	body := strings.NewReader(string(requestByte))
	client := &http.Client{}
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME + "/messages"
	req, err_req := http.NewRequest("POST", path, body)

	token, getTokenErr := GetToken()
	if getTokenErr != nil {
		return []byte(""), errors.New(ToStr(getTokenErr))
	}
	req.Header.Add("Authorization", "Bearer " + ToStr(token))

	if err_req != nil {
		return []byte(""), errors.New(err_req.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err_resp := client.Do(req)
	if err_resp != nil {
		return []byte(""), errors.New(err_resp.Error())
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	return res, nil
}


/*发送消息
 * http POST http://localhost:8080/sendmsg target_type=users target:='["liyuan"]' from=liyuan1 msg:='{"type":"txt","msg":"这才是消息内容", "size":='{"width":0, "height":0}', "fileName":"", "scret":""}' recentId=3 ext:='{"fromAvatar":"/files/default_avatar.png","sendTime":1512099413}'
 *
 * arget_type == "chatgroups" 群聊 target_type == "users" 单聊
*/
func ListenSendMsg (c *gin.Context) {
	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}
	var (
		cInfo SendRequest
		HXresponse HXsendMsgContent

	)
	engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
	if err != nil {
		rData(-1, "写入聊天记录打开数据库出错! 错误原因", err.Error())
		return
	}

	bindErr := c.Bind(&cInfo)
	if bindErr != nil {
		rData(-1, "参数出错!检查数据" + bindErr.Error(), bindErr.Error())
		return
	}

	// 获取目标对象信息
	if cInfo.TargetType == "chatgroups" {
		// 群聊详情信息
		_, getErr := engine.Where("id = ?", cInfo.Target[0]).Get(&cInfo.Ext.TargetInfo.GroupsInfo)
		if getErr != nil {
			rData(-1, "查询群聊详情失败!", getErr.Error())
			return
		}
	}else {
		// 单聊用户详情信息
		_, getErr := engine.Where("name = ?", cInfo.Target[0]).Cols("mobile", "age", "sex", "activated", "name").Get(&cInfo.Ext.TargetInfo.ChatInfo)
		if getErr != nil {
			rData(-1, "查询目标用户详情失败!", getErr.Error())
			return
		}
	}

retrySendMsg:
	// 发送消息到环信服务器
	resp, resp_err := cInfo.SendMsgToHXServer()
	if resp_err != nil {
		CheckHXErr(c, resp_err.Error(), "消息发送失败")
		return
	}
	// 发送成功后写入自己数据库
	// token 错误重新获取Token
	// 序列化response
	marshalErr := json.Unmarshal(resp, &HXresponse)
	if marshalErr != nil {
		rData(-1, "发送消息序列化失败!", marshalErr.Error())
		return
	}
	// 检测Token是都过期
	if HXresponse.Error_description == "Unable to authenticate due to expired access token" {
		isSuccess, err := GetAbleToken()
		if isSuccess {
			goto retrySendMsg
		}else {
			rData(-1, "获取Token错误!", err.Error())
			return
		}
	}

	if HXresponse.Data == nil {
		rData(-1, "发送失败!", "")
		return
	}
	// 查询发送目标对象信息
	rData(0, "", cInfo)

	// 并发写入相关数据
	go func () {
		// 写入数据表并发执行
		roomTarget := cInfo.Target
		var roomType int64 = 1
		if cInfo.TargetType == "chatgroups" {
			roomType = 2
		}
		timeInt64 := int64(cInfo.Ext.SendTime)
		t := time.Unix(timeInt64, 0)
		lastMsg := func () string{
			if cInfo.Msg.Type == "img" {
				return "[图片]"
			} else {
				return cInfo.Msg.Msg
			}
		}
		rencent := &RecentConcat{
			RoomType: roomType,
			UserName: cInfo.From,
			TargetUserName: roomTarget[0],
			IsTop: false,
			LastMessage: lastMsg(),
			LastMsgUpdated: t,
		}
		// 单聊写入单聊关系表  群聊忽略
		// 单聊需要检测是否存在
		// 检测是否已存在
		// 已存在更新发送内容和发送内容时间
		has, isExistRencentErr := engine.Table("recent_concat").Where("target_user_name = ?", roomTarget[0]).And("user_name = ?", cInfo.From).Exist()
		if has {
			fmt.Println("存在")
			// 更新更新发送内容和发送内容时间
			_, updateErr := engine.Id(cInfo.RecentId).Cols("last_message","last_msg_updated").Update(rencent)
			if updateErr != nil {
				// TODO: 出错写入日志
				fmt.Printf("更新最后聊天记录失败,失败原因:%v",updateErr.Error())
				return
			}
			return
		}
		// 如果不存在,rencent.isTop属性默认为false
		if isExistRencentErr != nil {
			// TODO: 出错写入日志
			fmt.Printf("检测最近聊天是否存在失败,失败原因:%v",isExistRencentErr.Error())
			return
		}

		// 写入最近聊天关系表
		_ , inserErr := engine.Insert(rencent)
		if inserErr != nil{
			rData(-1, "写入最近联系关系表错误!", inserErr.Error())
			return
		}
	}()
}
