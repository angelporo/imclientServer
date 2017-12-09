package control

import (
	"github.com/xormplus/xorm"
	"github.com/gin-gonic/gin"
	// "net/http"
	"errors"
	. "imClientServer/config"
	// . "imClientServer/util"
	// "net/http"
	// "strings"
	// "io/ioutil"
	"encoding/json"
	"fmt"
	// "time"
	"log"
	"io"
	"os"
	"strconv"
	"time"
)


// 返回数据结构
type SendImgResponse struct {
	Info SendRequest `json:"res"`
	Res HXsendMsgContent `json:"info"`
}


// 发送图片到环信服务器
// 发送成功不做任何事,
func (request *SendRequest) SendImgMsgToHXserver () error {
	return errors.New("hahah")
}



// http -f POST http://localhost:8080/sendimg file@/Users/angel/Downloads/th.jpeg -v
func ListenSendImgMsg (c *gin.Context) {
	const (
		upload_path string = "./static/chatImg/"
	)
	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}

	// 获取form/data中的file对象
	form, _ := c.MultipartForm()
	files := form.File["file"]
	// 获取append到form/data中的自定义数据
	from := c.PostForm("from")
	target := c.PostForm("target")
	sendTime := c.PostForm("sendTime")
	fromAvatar := c.PostForm("fromAvatar")
	target_type := c.PostForm("target_type")

	savePath := upload_path + from + "_TO_" + target + "/"
	if _, err := os.Stat(savePath); err == nil {
		log.Println("存在文件夹")
	} else {
		// 文件夹不存在则,  自动创建
		err := os.MkdirAll(savePath, 0711)

		if err != nil {
			rData(-1, "创建文件夹失败", err.Error())
			return
		}
	}
	// 创建切片数组来返回发送情况
	result := make([]SendImgResponse, len(files))

	for i, _ := range files {
		file, openErr := files[i].Open()
		var targetInfo User
		var groupMembersContent GroupMembersContent
		if openErr != nil {
			rData(-1, "openFile失败", openErr.Error())
			return
		}
		fileName := strconv.FormatInt(time.Now().Unix(), 10) + from + "_to_" + target + files[i].Filename
		Fw, CreateErr := os.Create(savePath + fileName)
		fmt.Println(savePath + fileName)
		if CreateErr != nil {
			rData(-1, "创建文件失败", CreateErr.Error())
			return
		}
		defer Fw.Close()
		_, saveErr := io.Copy(Fw, file)
		if saveErr != nil {
			rData(-1, "创建文件复制操作错误!", saveErr.Error())
			return
		}
		engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
		if err != nil {
			rData(-1, "写入聊天记录打开数据库出错! 错误原因", err.Error())
			return
		}
		if target_type == "chatgroups" {
			// 群聊详情信息
			_, getErr := engine.Where("id = ?", target).Get(&groupMembersContent)
			if getErr != nil {
				rData(-1, "查询群聊详情失败!", getErr.Error())
				return
			}
		}else {
			// 单聊用户详情信息
			_, getErr := engine.Where("name = ?", target).Cols("mobile", "age", "sex", "activated", "name").Get(&targetInfo)
			if getErr != nil {
				rData(-1, "查询目标用户详情失败!", getErr.Error())
				return
			}
		}
		// 并发发送图片到环信服务器
		sendTimeInt, _ := strconv.Atoi(sendTime)
		ExtMsg := ExtMsg{
			FromAvatar: fromAvatar,
			SendTime: sendTimeInt,
			SendMsgUri: fileName,
			TargetInfo: TargetInfo{
				ChatInfo: targetInfo,
				GroupsInfo:groupMembersContent,
			},
		}
		// 获取发送实际图片大小
		fileSize := map[string]int{
			"width": 123,
			"height":234,
		}
		msg := MassageContent {
			Type: "img",
			FileName: fileName,
			Secret: "",
			Size: fileSize,
		}
		// 获取发送消息必要属性
		request := SendRequest {
			Ext: ExtMsg,
			TargetType: target_type,
			Target: []string{target},
			From: from,
			Msg: msg,
		}
		res, sendErr := request.SendMsgToHXServer()
		if sendErr != nil {
			rData(0, "发送失败", res)
			return
		}
		var HXresponse HXsendMsgContent
		marshalErr := json.Unmarshal( res, &HXresponse )
		if marshalErr != nil {
			rData(-1, "发送消息序列化失败!", marshalErr.Error())
		}
		// 发送图片实际情况 装载到切片中
		result[i].Res = HXresponse
		result[i].Info = request
	}

	rData(0, "发送成功", result)
}
