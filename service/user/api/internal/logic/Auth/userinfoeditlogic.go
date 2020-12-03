package logic

import (
	"context"

	"hey-go-zero/service/user/api/internal/svc"
	"hey-go-zero/service/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UserInfoEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoEditLogic {
	return UserInfoEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoEditLogic) UserInfoEdit(req types.UserInfoReq) error {
	// todo: add your logic here and delete this line

	return nil
}
