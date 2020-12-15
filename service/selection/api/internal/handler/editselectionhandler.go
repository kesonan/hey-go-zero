package handler

import (
	"net/http"

	"hey-go-zero/service/selection/api/internal/logic"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func editSelectionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EditSelectionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewEditSelectionLogic(r.Context(), ctx)
		err := l.EditSelection(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
