## 项目接口返回逻辑说明


### 登录
- 判断参数错误
- 检查错误参数
- 返回用户所需要信息
```base
➜  ~ http http://localhost:8080/login mobile='' passWord="123"
HTTP/1.1 400 Bad Request
Content-Length: 52
Content-Type: text/plain; charset=utf-8
Date: Sat, 04 Nov 2017 03:10:01 GMT

{
    "code": -1,
    "data": "",
    "msg": "手机号不能为空!"
}
```
