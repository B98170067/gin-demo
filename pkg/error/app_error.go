package errno

// AppError 是整個專案唯一允許往外丟的錯誤型別
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, msg string) error {
	return &AppError{Code: code, Message: msg}
}
