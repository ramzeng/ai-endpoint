package error

import (
	"net/http"

	"github.com/spf13/cast"
)

type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Parameters string `json:"param"`
}

func TooManyRequests() Error {
	return Error{
		Code:    cast.ToString(http.StatusTooManyRequests),
		Message: "too many requests",
	}
}

func InternalServerError() Error {
	return Error{
		Code:    cast.ToString(http.StatusInternalServerError),
		Message: "internal server error",
	}
}

func Unauthorized() Error {
	return Error{
		Code:    cast.ToString(http.StatusUnauthorized),
		Message: "unauthorized",
	}
}

func Forbidden(message string) Error {
	return Error{
		Code:    cast.ToString(http.StatusForbidden),
		Message: message,
	}
}

func BadRequest(message string) Error {
	return Error{
		Code:    cast.ToString(http.StatusBadRequest),
		Message: message,
	}
}

func GatewayTimeout(message string) Error {
	return Error{
		Code:    cast.ToString(http.StatusGatewayTimeout),
		Message: message,
	}
}
