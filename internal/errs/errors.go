package errs

type ErrorType string

const (
	BadRequest          ErrorType = "BAD_REQUEST"
	NotFound            ErrorType = "NOT_FOUND"
	InternalServerError ErrorType = "INTERNAL"
	Unauthorized        ErrorType = "UNAUTHORIZED"
	Forbidden           ErrorType = "FORBIDDEN"
	Conflict            ErrorType = "CONFLICT"

	MultiBadRequest ErrorType = "MULTI_BAD_REQUEST"
)

type AppError struct {
	Type    ErrorType
	Message string
	Fields  map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(t ErrorType, msg string, fields map[string]string) *AppError {
	if fields == nil {
		return &AppError{
			Type:    t,
			Message: msg,
		}
	}
	return &AppError{
		Type:    t,
		Message: msg,
		Fields:  fields,
	}
}
