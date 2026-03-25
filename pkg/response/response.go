package response

// Response is the standard API response envelope
type Response struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code"`
	Data       interface{} `json:"data"`
}

// Success returns a successful response
func Success(data interface{}) Response {
	return Response{
		Successful: true,
		ErrorCode:  "",
		Data:       data,
	}
}

// Error returns an error response
func Error(errorCode string) Response {
	return Response{
		Successful: false,
		ErrorCode:  errorCode,
		Data:       nil,
	}
}
