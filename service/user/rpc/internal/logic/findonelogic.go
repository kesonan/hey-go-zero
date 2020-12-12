package logic

import (
	"context"

	"hey-go-zero/service/user/model"
	"hey-go-zero/service/user/rpc/internal/svc"
	"hey-go-zero/service/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FindOneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneLogic {
	return &FindOneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  findone
func (l *FindOneLogic) FindOne(in *user.UserReq) (*user.UserReply, error) {
	data, err := l.svcCtx.UserModel.FindOne(in.Id)
	switch err {
	case nil:
		return convertUserFromDbToPb(data), nil
	case model.ErrNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Unknown, err.Error())
	}
}
