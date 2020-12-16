package handler

import (
	"net/http"

	"hey-go-zero/common/jwtx"
	"hey-go-zero/service/selection/api/internal/logic"
	"hey-go-zero/service/selection/api/internal/svc"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func mineSelectionsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := jwtx.GetUserId(w, r)
		if !ok {
			return
		}

		l := logic.NewMineSelectionsLogic(r.Context(), ctx)
		resp, err := l.MineSelections(userId)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
