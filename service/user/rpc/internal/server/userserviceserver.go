// Code generated by goctl. DO NOT EDIT!
// Source: user.proto

package server

import (
	"context"

	"hey-go-zero/service/user/rpc/internal/logic"
	"hey-go-zero/service/user/rpc/internal/svc"
	"hey-go-zero/service/user/rpc/user"
)

type UserServiceServer struct {
	svcCtx *svc.ServiceContext
}

func NewUserServiceServer(svcCtx *svc.ServiceContext) *UserServiceServer {
	return &UserServiceServer{
		svcCtx: svcCtx,
	}
}

//  findone
func (s *UserServiceServer) FindOne(ctx context.Context, in *user.UserReq) (*user.UserReply, error) {
	l := logic.NewFindOneLogic(ctx, s.svcCtx)
	return l.FindOne(in)
}

//  findByIds
func (s *UserServiceServer) FindByIds(ctx context.Context, in *user.IdsReq) (*user.UserListReply, error) {
	l := logic.NewFindByIdsLogic(ctx, s.svcCtx)
	return l.FindByIds(in)
}