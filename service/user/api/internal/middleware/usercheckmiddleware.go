package middleware

import (
	"net/http"

	"hey-go-zero/common/middleware"
)

type UserCheckMiddleware struct {
}

func NewUserCheckMiddleware() *UserCheckMiddleware {
	return &UserCheckMiddleware{}
}

func (m *UserCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return middleware.UserCheck(next)
}
