package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSelectionLogic {
	return GetSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelectionLogic) GetSelection() (*types.SelectionIdReq, error) {
	// todo: add your logic here and delete this line

	return &types.SelectionIdReq{}, nil
}
