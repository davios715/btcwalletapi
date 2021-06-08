package response

type responseError = string

const (
	ErrInvalidPath = "INVALID_PATH"
	ErrInvalidInput = "INVALID_INPUT"
	ErrSeed = "INVALID_SEED"
	ErrInternal = "INTERNAL"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func GetResponse(e responseError) ErrorResponse {

	var msg = ""

	switch e {
	case ErrInvalidInput:
		msg = "Invalid input"
	case ErrInvalidPath:
		msg = "Invalid path"
	case ErrSeed:
		msg = "Invalid path"
	default:
		msg = "Internal server error"
	}

	return ErrorResponse{
		Code:    e,
		Message: msg,
	}
}
