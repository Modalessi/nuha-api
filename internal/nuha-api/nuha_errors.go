package nuha

import "fmt"

// const (
// 	INVALID_JSON = fmt.Errorf()
// )

type NuhaError struct {
	Code    int
	Message string
}

func (e NuhaError) Error() string {
	return fmt.Sprintf("error %v: %s", e.Code, e.Message)
}

var (
	INVALID_JSON_ERROR         = NuhaError{Code: 400, Message: "invalid request, please check your payload"}
	SERVER_ERROR               = NuhaError{Code: 500, Message: "something went wrong on the server side, sorry"}
	INVALID_CREDINTALS_ERROR   = NuhaError{Code: 400, Message: "please follow the credntals guidlines"}
	WRONG_CREDINTALS_ERROR     = NuhaError{Code: 400, Message: "email or password, one of them is wrong"}
	USER_ALREADY_EXIST_ERROR   = NuhaError{Code: 400, Message: "this user already exist"}
	AUTHORIZATION_HEADER_ERROR = NuhaError{Code: 401, Message: "no authorization header"}
	INVALID_TOKEN_ERROR        = NuhaError{Code: 401, Message: "invalid token, please check"}
)

func EntityDoesNotExistError(enitity string) NuhaError {
	return NuhaError{
		Code:    404,
		Message: enitity + " does not exist",
	}
}
