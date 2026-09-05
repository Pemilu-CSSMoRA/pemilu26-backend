package dto

import (
	"net/http"

	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
)

const (
	// English
	MESSAGE_FAILED_TOKEN_NOT_VALID   = "token not valid"
	MESSAGE_FAILED_TOKEN_NOT_FOUND   = "token not found"
	MESSAGE_FAILED_TOKEN_EXPIRED     = "token expired"
	MESSAGE_FAILED_TOKEN_NOT_ALLOWED = "token not allowed"
	MESSAGE_API_IS_LOCKED            = "api is now locked"
)

var (
	ErrTokenInvalid   = MyError.New("token invalid", http.StatusBadRequest)
	ErrTokenExpired   = MyError.New("token expired", http.StatusUnauthorized)
	ErrRoleNotAllowed = MyError.New("role not allowed", http.StatusForbidden)
	ErrTokenNotFound  = MyError.New("token not found", http.StatusUnauthorized)
)
