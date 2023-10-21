package app_error

const (
	ErrorUserNotFound      ErrorType = "user_not_found"
	ErrorUserExist         ErrorType = "user_exist"
	ErrorUpdateUserSetting ErrorType = "update_user_setting"
	ErrorGetUserSetting    ErrorType = "get_user_setting"
	ErrorTaskNotFound      ErrorType = "task_not_found"
	ErrorTaskExist         ErrorType = "task_exist"
	ErrorBadRequest        ErrorType = "bad_request"
	ErrorNotFound          ErrorType = "not_found"
)

type ErrorType string

type AppError struct {
	Code    ErrorType
	Message string
}

func (e AppError) Error() string {
	return string(e.Code) + ":" + e.Message
}

func New(errType ErrorType, message string) AppError {
	return AppError{
		Code:    errType,
		Message: message,
	}
}

func WithError(errType ErrorType, err error) AppError {
	return AppError{
		Code:    errType,
		Message: err.Error(),
	}
}
