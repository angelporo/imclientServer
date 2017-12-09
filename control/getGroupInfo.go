package control
import (
	. "imClientServer/util"
	"github.com/gin-gonic/gin"
	"encoding/json"
)

type getGroupRequest struct {
	GroupId string `json:"groupId"`
}

// 获取群组聊天详情
// http POST http://localhost:8080/getgroupinfo groupId=33739229954049
func GetGroupInfoByGroupId (c *gin.Context) {

	var(
		request getGroupRequest
		tokenType GroupInfoContent
	)

	rData := func(code interface{}, msg interface{}, desc interface{}) {
		c.JSON(200, gin.H{
			"code": code,
			"msg": msg,
			"content": desc,
		})
	}

	bindErr := c.Bind(&request)
	if bindErr != nil {
		rData(-1, "获取群组详情绑定request参数错误", bindErr.Error())
		return
	}

	// 获取群组详情信息
	groupInfo, getGroupErr, statusStr := GetGroupMembersListById(request.GroupId)
	if getGroupErr != nil {
		rData(-1, "获取群组详情失败", getGroupErr.Error())
		return
	}

	status := Substr(statusStr, 0, 3)
	// 获取环信群组聊天详情服务器出错, 暂不考虑groupRoomId错误情况
	if Substr(statusStr, 0, 1) == "5" {
		// TODO: 写入日志
		rData(status, "获取群组详情环信服务器发生" + status + "错误", getGroupErr.Error())
		return
	}
	if status == "404" {
		rData(status, "群组id不存在", "")
		return
	}
	marshErr := json.Unmarshal(groupInfo, &tokenType)
	if marshErr != nil {
		rData(status, "解析群组详情错误", marshErr.Error())
		return
	}
	rData(0, "", tokenType)
}
