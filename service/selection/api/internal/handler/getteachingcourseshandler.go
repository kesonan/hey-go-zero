package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"hey-go-zero/common/jwtx"
	"hey-go-zero/service/selection/api/internal/logic"
	"hey-go-zero/service/selection/api/internal/svc"
)

func getTeachingCoursesHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := jwtx.GetUserId(w, r)
		if !ok {
			return
		}

		l := logic.NewGetTeachingCoursesLogic(r.Context(), ctx)
		resp, err := l.GetTeachingCourses(userId)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
