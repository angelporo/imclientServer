// im client register api
package control

import (
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"regexp"
	. "imClientServer/config"
	"imClientServer/util"
	"time"
	"encoding/json"
	"os"
	"errors"
)

var engine *xorm.Engine
const (
	// 判断手机号正则
	regular = `^1([38][0-9]|14[57]|5[^4])\d{8}$`
)

// 客户端传参数 模型绑定
type Register struct {
	UserName string `json:"userName" binding:"required"`
	NickName string `json:"nickName" binding:"required"`
	PassWord string `json:"passWord" binding:"required"`
	UserMobile string `json:"mobile" binding:"required"`
}

// 数据表映射
// User 对应数据表name
// 属性对应数据表内容字段
type User struct {
	Id int64
	Name string `xorm:"notnull"`
	Mobile string `xorm:"notnull"`
	Age int `xorm:"int notnull"`
	Sex int `xorm:"int notnull"`
	Money int `xorm:"int notnull default 0"`
	Avatar string `xorm:"notnull"`
	PassWord string `xorm:"notnull"`
	Created time.Time `xorm:"created notnull"`
	Updated time.Time `xorm:"updated"`
	Uuid string `xorm:"notnull"`
	Activated bool `xorm:"notnull bool"` // 是否在线
}
// 用户好友表
type UserRelationShip struct {
	Index int64 `xorm:"pk notnull autoincr unique"`
	UserId int64  `xorm:"notnull bigint"`
	UserName string `xorm:"notnull varchar(255)"`
	Friend_userName string `xorm:"notnull varchar(255)"`
}

// 环信返回结果type
type Huanxin struct {
	HxSuccessData
	Entities []RegisterSuccessContent
}

type HxSuccessData struct {
	ApplicationName string `json:"applicationName"`
	Duration int `json:"duration"`
	Path string `json:"path"`
	Action string `json:"action"`
	Organization string `json:"organization"`
	Uri string `json:"uri"`
	Timestamp int `json:"timestamp"`
	Application string `json:"application"`
	DataFaild
}

type DataFaild struct {
	Error string `json:"error"`
	Error_description string `json:"error_description"`
	Exception string `json:"exception"`
}

// 注册成功环信返回主体类型
type RegisterSuccessContent struct {
	Activated bool
	Created int
	Modified int
	Nickname string
	Type string
	UserName string
	Uuid string
}


// 判断文件是否存在
func checkFileIsExist(filename string) (bool) {
	var exist = true;
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false;
	}
	return exist;
}

// 从环信服务器 获取管理员Token
func GetTokenFromHuanxinServer () ([]byte, error) {
	// 环信注册用户
	// 构建获取环信token结构
	type GetHuanxinTokenType struct {
		Grant_type  string `json:"grant_type"`
		Client_id   string `json:"client_id"`
		Client_secret string `json:"client_secret"`
	}
	client := &http.Client{}
	bodyStr := &GetHuanxinTokenType{
		Grant_type:     "client_credentials",
		Client_id:   CLIENT_ID,
		Client_secret: CLIENT_SECRET,
	}
	b, _ := json.Marshal(bodyStr)
	body := strings.NewReader(string(b))
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME+ "/token"
	req, err_req := http.NewRequest("POST", path, body)
	if err_req != nil {
		return []byte(""),  errors.New("获取token错误")
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err_resp := client.Do(req)
	if err_resp != nil {
		return []byte("") , errors.New("获取token错误")
	}
	defer resp.Body.Close()
	r, _ := ioutil.ReadAll(resp.Body)
	return r, nil
}


// 在环信服务器注册用户
func RegisterUserByHuanxing (name string, passWord string, nickname string) ([]byte, string,  error) {
	type RegiserData struct {
		UserName string `json:"username"`
		PassWord string `json:"password"`
		Nickname string `json:"nickname"`
	}
	requestBody := &RegiserData{
		UserName:name,
		PassWord:passWord,
		Nickname: nickname,
	}
	client := &http.Client{}
	b, _ := json.Marshal(requestBody)
	body := strings.NewReader(string(b))
	path := HUANXIN_DOMAIN + ORG_NAME + "/" + APP_NAME + "/users"
	req, err_req := http.NewRequest("POST", path, body)
	if err_req != nil {
		return []byte(""), "600" , errors.New("环信注册失败")
	}
	req.Header.Add("Content-Type", "application/json")
	token, err := GetToken()
	if err != nil {
		return []byte(""), "600", errors.New("环信注册用户Token出错!")
	}
	req.Header.Add("Authorization", "Bearer" + util.ToStr(token))
	resp, err_resp := client.Do(req)
	if err_resp != nil {
		return []byte("") ,"600",  errors.New("获取token错误")
	}
	defer resp.Body.Close()
	r, _ := ioutil.ReadAll(resp.Body)
	return r, resp.Status, nil
}


func GetToken () (interface{}, error) {
	var data interface{}
retryOpenToken:
	f, err := os.OpenFile(TOKEN_FILE_NAME, os.O_APPEND, 0666)
	if err != nil {
		// err := errors.New("打开token文件失败!")
		_, getErr := GetAbleToken()
		if getErr != nil {
			goto retryOpenToken
		}
	}
	fileInfo, err := os.Stat(TOKEN_FILE_NAME)
	if err != nil {
		err := errors.New("打开token文件失败!")
		return  data, err
	}
	fileSize := fileInfo.Size() //获取size
	tokenSlice := make([]byte, fileSize)
	_, err5 := f.Read(tokenSlice)
	if err5 != nil {
		err := errors.New("读取Token文件失败!")
		return data, err
	}
	// 去除Token字符串中不符合json格式符号
	str := strings.Replace(string(tokenSlice), "'", "\"", -1)
	str = strings.Replace(str, "\n", "", -1)
	var tokenType map[string]interface{}
	err6 := json.Unmarshal([]byte(str), &tokenType)
	if err6 != nil {
		err := errors.New("读取Token文件失败!")
		return data, err
	}
	token := tokenType["access_token"]
	return  token, nil
}
// 获取有效的Token
func GetAbleToken () (bool, error) {
	token, err := GetTokenFromHuanxinServer()
	if err != nil {
		return false, errors.New("err")
	}
	f, err1 := os.Create(TOKEN_FILE_NAME)
	if err1 != nil {
		return false, errors.New("创建token文件失败!")
	}
	f.Write(token)
	return true, nil
}

/** 用户注册 put 请求
 * 功能介绍 手机号检测(重号, 格式不对), 用户名(重名, 格式(225byte, 英文开头, 不能汉子))
 * http PUT http://localhost:8080/user mobile='18303403747' userName="liyuan" passWord="angel" NickName="会上树的猪"
**/
func RegisterListen(c *gin.Context){
	var (
		paramJson Register
	)
	// isTokenExpires, ok := c.Get("isExpires")
	// 程序错误没有正确的isTokenExpires value
	// if !ok {
	//	c.JSON(200, gin.H{
	//		"code": -1,
	//		"msg": "没有isTokenExpires这个key!",
	//		"content": "",
	//	})
	//	return
	// }
	// if isTokenExpires != nil {
	//	c.JSON(200, gin.H{
	//		"code": -1,
	//		"msg": isTokenExpires,
	//		"content": "",
	//	})
	//	return
	// }
	err := c.Bind(&paramJson)
	// 检测参数
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "参数错误",
			"content": "",
		})
		return
	}

	if ok,  _ := regexp.MatchString(regular, paramJson.UserMobile); !ok {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "填写正确的手机号!",
			"content": "",
		})
		return
	}

	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", paramJson.UserName); !ok {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "用户名格式只能使用字母和数字组成, 长度在4~16!",
			"content": "",
		})
		return
	}

	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", paramJson.PassWord); !ok {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "密码格式错误!",
			"content": "",
		})
		return
	}

	// create xorm entity
	engine, err := xorm.NewEngine("mysql", DATABASE_LOGIN)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "系统发生错误",
			"content": "",
		})
		return
	}

	// 同步用户数据表结构
	errA := engine.Sync2(
		new(User),
		// new(ChatRoomGroupAffiliations), // 只有单聊了群聊,  没有聊天室
		new(ChatGroupInfo),
		new(ChatHistory),
		new(GroupMembersContent), // 群组聊天详情
		new(UserRelationShip),
		new(GroupRelationShip), // 群组聊天关系
		new(RecentConcat),
		new(MembersItem))
	if errA != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "同步数据表错误",
			"content": errA.Error(),
		})
		return
	}
	// 检查手机号是否使用
	mobile := &User{
		Mobile: paramJson.UserMobile,
	}
	has, _ := engine.Exist(mobile)
	if has {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "手机号码被注册, 请使用其他手机号!",
			"content": "",
		})
		return
	}
	// 检查用户名是否使用
	name := &User{
		Name: paramJson.UserName,
	}
	hasName, _ := engine.Exist(name)
	if hasName {
		c.JSON(200, gin.H{
			"code": "",
			"msg": "用户名已存在",
			"content": "",
		})
		return
	}
	// 环信服务器注册
	resp, status, err := RegisterUserByHuanxing(paramJson.UserName, paramJson.PassWord, paramJson.NickName)
	if err != nil {
		c.JSON(200, gin.H{
			"code": status,
			"msg": err,
			"content": "",
		})
		return
	}

	var HuanxinData Huanxin
	err7 := json.Unmarshal(resp, &HuanxinData)
	if err7 != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "反序列化环信注册结果失败!",
			"content": err7.Error(),
		})
		return
	}
	if HuanxinData.Entities == nil {
		CheckHXErr(c, HuanxinData.Error, "环信注册失败")
		return
	}
	HuanxinResContent := HuanxinData.Entities
	userUuid := HuanxinResContent[0].Uuid
	// 数据没有任何问题, 插入数据库
	user := &User{
		Name: paramJson.UserName,
		Age: 26,
		Mobile: paramJson.UserMobile,
		PassWord: paramJson.PassWord,
		Money: 0,
		Sex: 1, // 默认性别
		Avatar: "/static/default_avatar.png", // 默认头像为相对路径
		Uuid: userUuid,
		Activated: true,
	}
	_ , err1 := engine.Insert(user)
	if err1 != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg": "数据写入出错",
			"content": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg": "创建成功, 欢迎使用信信聊天",
		"content": HuanxinData,
	})
	return
}
