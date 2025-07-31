package app_code

const (
	Success                CodeType = "success"
	ErrorUserNotFound      CodeType = "user_not_found"
	ErrorGetUserList       CodeType = "find_user_list"
	ErrorUserPassword      CodeType = "password_error"
	ErrorServerError       CodeType = "service_error"
	ErrorAuthFailed        CodeType = "auth_failed"
	ErrorWxAuthFailed      CodeType = "wx_auth_failed"
	ErrorTokenNotSet       CodeType = "token_not_set"
	ErrorRepeat            CodeType = "repeat"
	ErrorTokenTimeout      CodeType = "token_timeout"
	ErrorForbidden         CodeType = "forbidden"
	ErrorUserExist         CodeType = "user_exist"
	ErrorUpdateUserSetting CodeType = "update_user_setting"
	ErrorGetUserSetting    CodeType = "get_user_setting"
	ErrorTaskNotFound      CodeType = "task_not_found"
	ErrorUserTaskNotFound  CodeType = "user_task_not_found"
	ErrorTaskExist         CodeType = "task_exist"
	ErrorBadRequest        CodeType = "bad_request"
	ErrorNotImage          CodeType = "not_image"
	ErrorAuditParamInValid CodeType = "audit_param_invalid"
	ErrorNotFound          CodeType = "not_found"
	ErrorTaskRunNotFound   CodeType = "task_run_not_found"
	ErrorCreateTask        CodeType = "create_task"
)

// 将错误代码翻译成中文
var CodeTypeToChineseMap = map[CodeType]string{
	Success:                "成功",
	ErrorUserNotFound:      "用户未找到",
	ErrorGetUserList:       "查找用户列表失败",
	ErrorUserPassword:      "密码错误",
	ErrorServerError:       "服务器错误",
	ErrorAuthFailed:        "令牌解析失败，请重新登录",
	ErrorWxAuthFailed:      "微信身份验证失败",
	ErrorTokenNotSet:       "令牌未设置，请重新登录",
	ErrorRepeat:            "重复操作",
	ErrorTokenTimeout:      "令牌已过期，请重新登录",
	ErrorForbidden:         "权限不足",
	ErrorUserExist:         "用户已存在",
	ErrorUpdateUserSetting: "更新用户设置失败",
	ErrorGetUserSetting:    "获取用户设置失败",
	ErrorTaskNotFound:      "任务未找到",
	ErrorUserTaskNotFound:  "用户任务未找到",
	ErrorTaskExist:         "任务已存在",
	ErrorBadRequest:        "请求参数错误",
	ErrorNotImage:          "不是图片格式",
	ErrorAuditParamInValid: "审核参数无效",
	ErrorNotFound:          "记录未找到",
	ErrorTaskRunNotFound:   "任务运行记录未找到",
	ErrorCreateTask:        "创建任务失败",
}

var AppErrorRoot = New("root", "root error")

type CodeType string

type AppError struct {
	Code    CodeType
	Message string
}

func (e AppError) Error() string {
	return string(e.Code) + ":" + e.Message
}

func New(errType CodeType, message string) *AppError {
	return &AppError{
		Code:    errType,
		Message: message,
	}
}

func WithError(errType CodeType, err error) *AppError {
	return &AppError{
		Code:    errType,
		Message: err.Error(),
	}
}
func GetErrMsg(code CodeType) string {
	msg, ok := CodeTypeToChineseMap[code]
	if !ok {
		msg = "操作失败"
	}
	return msg
}
