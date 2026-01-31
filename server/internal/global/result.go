package g

import (
	"fmt"
)

const (
	SUCCESS = 0   // 成功业务码
	FAIL    = 500 // 失败业务码
)

// Result 自定义业务结果类型
type Result struct {
	code int    // 业务码
	msg  string // 业务消息
}

func (r Result) Code() int {
	return r.code
}

func (r Result) Msg() string {
	return r.msg
}

var (
	_codes    = make(map[int]struct{}) // 注册过的错误码集合,防止重复
	_messages = make(map[int]string)   // 错误码对应的错误信息
)

// RegisterResult 注册一个响应码,不允许重复注册
func RegisterResult(code int, msg string) Result {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在,请更换一个", code))
	}
	if msg == "" {
		panic("错误码消息不能为空")
	}

	_codes[code] = struct{}{}
	_messages[code] = msg

	return Result{
		code: code,
		msg:  msg,
	}
}

// GetMsg 根据响应码获取响应信息
func GetMsg(code int) string {
	return _messages[code]
}

// 系统级响应码
var (
	Success = RegisterResult(1000, "success")
	Err     = RegisterResult(1001, "服务器端错误")
)

// 通用业务响应码
var (
	ErrRequest  = RegisterResult(9001, "请求参数格式错误")
	ErrDbOp     = RegisterResult(9004, "数据库操作异常")
	ErrRedisOp  = RegisterResult(9005, "Redis 操作异常")
	ErrUserAuth = RegisterResult(9006, "用户认证异常")
)

// Token相关响应码
var (
	ErrTokenNotExist = RegisterResult(1201, "TOKEN 不存在,请重新登陆")
	ErrTokenRuntime  = RegisterResult(1202, "TOKEN 已过期,请重新登陆")
	ErrTokenWrong    = RegisterResult(1203, "TOKEN 不正确,请重新登陆")
	ErrTokenType     = RegisterResult(1204, "TOKEN 格式错误,请重新登陆")
	ErrTokenCreate   = RegisterResult(1205, "TOKEN 生成失败")
)

// 用户相关响应码
var (
	ErrUserNotExist = RegisterResult(1003, "该用户不存在")
	ErrUserExist    = RegisterResult(1004, "该用户已存在,请更换用户名")
	ErrPassword     = RegisterResult(1002, "密码错误")
	ErrPermission   = RegisterResult(1005, "权限验证失败")
)

// 文件/资源相关
var (
	ErrFileNotExist = RegisterResult(1101, "文件/资源不存在")
)

// TODO相关响应码
var (
	ErrTodoListGet = RegisterResult(1301, "获取TODO列表失败")
	ErrNoTodoList  = RegisterResult(1302, "暂无TODO")
	ErrTodoListSet = RegisterResult(1303, "更新TODO失败")
)
