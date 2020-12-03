package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"hey-go-zero/service/user/api/internal/logic/auth"
	"hey-go-zero/service/user/api/internal/svc"
)

func UserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewUserInfoLogic(r.Context(), ctx)
		resp, err := l.UserInfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
