package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type CreateSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateSelectionLogic {
	return CreateSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSelectionLogic) CreateSelection(req types.CreateSelectionReq) error {
	// todo: add your logic here and delete this line

	return nil
}
