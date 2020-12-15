package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"hey-go-zero/service/selection/api/internal/logic"
	"hey-go-zero/service/selection/api/internal/svc"
)

func getSelectionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGetSelectionLogic(r.Context(), ctx)
		resp, err := l.GetSelection()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
