package app_code

const (
	Success                CodeType = "success"
	ErrorUserNotFound      CodeType = "user_not_found"
	ErrorUserPassword      CodeType = "password_error"
	ErrorServerError       CodeType = "service_error"
	ErrorAuthFailed        CodeType = "auth_failed"
	ErrorTokenNotSet       CodeType = "token_not_set"
	ErrorTokenTimeout      CodeType = "token_timeout"
	ErrorForbidden         CodeType = "forbidden"
	ErrorUserExist         CodeType = "user_exist"
	ErrorUpdateUserSetting CodeType = "update_user_setting"
	ErrorGetUserSetting    CodeType = "get_user_setting"
	ErrorTaskNotFound      CodeType = "task_not_found"
	ErrorTaskExist         CodeType = "task_exist"
	ErrorBadRequest        CodeType = "bad_request"
	ErrorAuditParamInValid CodeType = "audit_param_invalid"
	ErrorNotFound          CodeType = "not_found"
	ErrorTaskRunNotFound   CodeType = "task_run_not_found"
	ErrorCreateTask        CodeType = "create_task"
)

type CodeType string

type AppError struct {
	Code    CodeType
	Message string
}

func (e AppError) Error() string {
	return string(e.Code) + ":" + e.Message
}

func New(errType CodeType, message string) AppError {
	return AppError{
		Code:    errType,
		Message: message,
	}
}

func WithError(errType CodeType, err error) AppError {
	return AppError{
		Code:    errType,
		Message: err.Error(),
	}
}
