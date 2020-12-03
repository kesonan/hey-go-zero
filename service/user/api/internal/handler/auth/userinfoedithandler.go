package handler

import (
	"net/http"

	"hey-go-zero/service/user/api/internal/logic/auth"
	"hey-go-zero/service/user/api/internal/svc"
	"hey-go-zero/service/user/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UserInfoEditHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUserInfoEditLogic(r.Context(), ctx)
		err := l.UserInfoEdit(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
