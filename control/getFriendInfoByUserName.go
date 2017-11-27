// 通过好友id获取好友详情数据
package control
import (
	"github.com/xormplus/xorm"
	"errors"
)

// 通过好用id获取好友详情数据
func GetFriendInfoByUserName (arr []UserRelationShip, engine *xorm.Engine) ([]userItem, error) {
	friendInfo := make([]userItem, len(arr))
	for i:= 0; i < len(arr);i++ {
		var friendItem userItem
		_, err := engine.Table("user").Where("name = ?", arr[i].Friend_userName).Cols("avatar", "name", "mobile", "Sex").Get(&friendItem)
		if err != nil {
			return friendInfo, errors.New(err.Error())
		}
		friendInfo[i] = friendItem
	}
	return friendInfo, nil
}
