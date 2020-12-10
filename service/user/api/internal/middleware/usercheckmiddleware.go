package middleware

import (
	"fmt"
	"net/http"

	"hey-go-zero/common/errorx"
	"hey-go-zero/common/jwtx"

	"github.com/tal-tech/go-zero/rest/httpx"
)

type UserCheckMiddleware struct {
}

func NewUserCheckMiddleware() *UserCheckMiddleware {
	return &UserCheckMiddleware{}
}

func (m *UserCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(jwtx.JwtWithUserKey)
		xUserId := r.Header.Get("x-user-id")
		if len(xUserId) == 0 {
			httpx.Error(w, errorx.NewDescriptionError("x-user-id不能为空"))
			return
		}

		if xUserId != fmt.Sprintf("%v", v) {
			httpx.Error(w, errorx.NewDescriptionError("用户信息不一致"))
			return
		}
		next(w, r)
	}
}
