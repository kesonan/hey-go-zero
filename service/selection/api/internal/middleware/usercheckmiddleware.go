package middleware

import (
	"net/http"

	"hey-go-zero/common/errorx"
	"hey-go-zero/common/jwtx"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserCheckMiddleware struct {
	userRpcClient userservice.UserService
	role          string
}

func NewManagerCheckMiddleware(role string, userRpcClient userservice.UserService) *UserCheckMiddleware {
	return &UserCheckMiddleware{
		role:          role,
		userRpcClient: userRpcClient,
	}
}

func (m *UserCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := jwtx.GetUserId(w, r)
		if !ok {
			return
		}

		data, err := m.userRpcClient.FindOne(r.Context(), &userservice.UserReq{
			Id: userId,
		})
		if err != nil {
			st := status.Convert(err)
			if st.Code() == codes.NotFound {
				httpx.Error(w, errorx.NewDescriptionError("用户不存在"))
				return
			}

			httpx.Error(w, errorx.NewDescriptionError("用户信息获取失败"))
			return
		}

		if data.Role != m.role {
			httpx.Error(w, errorx.NewDescriptionError("无权限访问"))
			return
		}

		next(w, r)
	}
}
