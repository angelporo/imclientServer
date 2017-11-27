// 发送消息接口
// 环信: 在调用程序中，如果返回429或503错误，说明接口被限流了，请稍微暂停一下并重试。请求体如果超过 5kb 会导致413错误，需要拆成几个更小的请求体重试，同时用户消息+扩展字段的长度在40k字节以内
package control

import (
	// "time"
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
)

// 消息属性类型
type SendMsgRequest struct {
	TargetType string `json:"target_type" binding:"required"` // users 给用户发消息。chatgroups: 给群发消息，chatrooms: 给聊天室发消息
	Target []string `json:"target" binding:"required"` // 注意这里需要用数组，数组长度建议不大于20，即使只有一个用户，也要用数组 ['u1']，给用户发送时数组元素是用户名，给群组发送时.数组元素是groupid
	From string `json:"from" binding:"required"`  //表示消息发送者。无此字段Server会默认设置为"from":"admin"，有from字段但值为空串("")时请求失败
}

// 消息内容类型
type MassageContent struct {
	Type string `json:"type" binding:"required"` // 消息内容类型
	Msg string `json:"msg" binding:"required"` // 消息内容
}

// 环信服务器接受类型
type SendRequest struct {
	SendMsgRequest
	Msg MassageContent `json:"msg"`
}


// 发送消息到环信服务器
func SendMsgToHuanxinServer (requestBody *SendRequest) ([]byte, error) {

	request := requestBody
	requestByte, _ := json.Marshal(request)

	body := strings.NewReader(string(requestByte))
	client := &http.Client{}
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME + "/jmessages"
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
	fmt.Println(resp.Header)
	res, _ := ioutil.ReadAll(resp.Body)
	return res, nil
}

type HXSendMsgSuccess struct {
	Action string `json:"action"`
	Application string `json:"application"`
	Param map[string]string `json:"param"`
	Uri string `json:"uri"`
	Entities []interface{} `json:"entities"`
	Data map[string]string `json:"data"`
	Timestamp int `json:"timestamp"`
	Duration int `json:"duration"`
	Organization string `json:"organization"`
	ApplicationName string `json:"applicationName"`
}

// 环信返回json类型映射类型
type HXsendMsgContent struct {
	DataFaild
	HXSendMsgSuccess
}


// 发送消息
// http POST http://localhost:8080/sendmsg target_type=users target:='["liyuan2", "liyuan1"]' from=liyuan msg:='{"type":"tx9t","msg":"这才是消息内容"}'
func ListenSendMsg (c *gin.Context) {
	var (
		cInfo  SendRequest
	)
	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}
	err := c.Bind(&cInfo)
	if err != nil {
		rData(-1, "参数出错!检查数据", err.Error())
		return
	}
retrySendMsg:
	// 发送消息到环信服务器
	resp, resp_err := SendMsgToHuanxinServer(&cInfo)
	if resp_err != nil {
		CheckHXErr(c, resp_err.Error(), "消息发送失败")
		return
	}
	// 发送成功后写入自己数据库
	// token 错误重新获取Token
	// 序列化response
	var HXresponse HXsendMsgContent
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
	// 返回结果给用户,
	rData(0, "", HXresponse)


	// 并发写入相关数据
	go func () {
		// 写入数据表并发执行
		engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
		if err != nil {
			rData(-1, "写入聊天记录打开数据库出错!", err.Error())
			return
		}
		roomTarget := cInfo.Target
		var roomType int64 = 1
		if cInfo.TargetType == "chatgroups" {
			roomType = 2
		}

		// 不管单聊还是群聊, 写入最近聊天关系表
		rencent := &RecentConcat{
			RoomType: roomType,
			UserName: cInfo.From,
			RoomId: roomTarget[0],
		}
		// 检测是否已存在
		has, isExistRencentErr := engine.Exist(rencent)
		if has {
			return
		}
		if isExistRencentErr != nil {
			fmt.Printf("检测最近聊天是否存在失败,失败原因:%v",isExistRencentErr.Error())
			return
		}

		// 写入最近聊天关系表
		_ , inserErr := engine.Insert(rencent)
		if inserErr != nil{
			rData(-1, "写入最近联系关系表错误!", inserErr.Error())
			return
		}
		// 写入最近联系人详情表信息
	}()
}
